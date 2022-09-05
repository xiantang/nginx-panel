package etcd

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type Syncer struct {
	client  *clientv3.Client
	schemas map[string]string
}

func NewSyncer() *Syncer {
	client := GetClient()
	return &Syncer{
		client,
		map[string]string{},
	}
}

func (w *Syncer) Run() {
	for {
		// start sync
		err := w.sync()
		if err != nil {
			log.WithError(err).Error("sync error")
		}
		time.Sleep(10 * time.Second)
	}
}

func (w *Syncer) sync() error {
	// get config from etcd
	resp, err := w.client.Get(context.Background(), "/config", clientv3.WithPrefix())
	if err != nil {
		return err
	}
	// map
	configMap := make(map[string]string)
	for _, kv := range resp.Kvs {
		log.WithField("key", string(kv.Key)).Info("get config from etcd")
		configMap[string(kv.Key)] = string(kv.Value)
	}
	log.WithField("configMap", configMap).Info("configMap")

	// add map
	addedMap := make(map[string]string)

	// update map
	updatedMap := make(map[string]string)

	// compare from memory
	for k, v := range configMap {
		if _, ok := w.schemas[k]; !ok {
			// add into addedMap
			addedMap[k] = v
			continue
		}
		if old, ok := w.schemas[k]; ok {
			if old != v {
				// update into updatedMap
				updatedMap[k] = v
			}
		}
	}

	// mkdir which is /etc/nginx/test/ if not exist

	_ = os.MkdirAll("/etc/nginx/tests/", 0755)

	// write file
	for k, v := range addedMap {
		// write file
		s := strings.Split(k, "/")
		fileName := s[len(s)-1]
		path := fmt.Sprintf("%v%v.conf", "/etc/nginx/tests/", fileName)
		log.WithField("path", path).Info("write file")
		ioutil.WriteFile(path, []byte(v), 0644)
	}

	if len(addedMap) == 0 && len(updatedMap) == 0 {
		return nil
	}

	log.Info("file written")
	if os.Getenv("SKIP_NGINX_TEST") == "true" {
		log.Info("skip nginx test")
	} else {
		// nginx test

		cmd := exec.Command("nginx", "-c", "/etc/nginx/nginx_test.conf", "-t")

		buf := bytes.NewBufferString("")
		cmd.Stderr = buf
		err = cmd.Run()
		if err != nil {
			log.WithError(err).Error("nginx -t error")
			return err
		}
		log.Info("nginx -t success")
	}

	// copy files from nginx/tests/ to nginx/http-enabled/

	for k, v := range addedMap {
		// write file
		s := strings.Split(k, "/")
		fileName := s[len(s)-1]
		path := fmt.Sprintf("%v%v.conf", "/etc/nginx/http-enabled/", fileName)
		ioutil.WriteFile(path, []byte(v), 0644)
	}

	for k, v := range addedMap {
		w.schemas[k] = v
	}

	log.Info("schemas updated")
	if os.Getenv("SKIP_RELOAD") == "true" {
		log.Info("skip nginx reload")
	} else {
		cmd := exec.Command("nginx", "-s", "reload")
		buf := bytes.NewBufferString("")
		cmd.Stderr = buf
		err = cmd.Run()
		if err != nil {
			log.WithField("output", buf.String()).WithError(err).Error("nginx -s reload error")
			return err
		}
		log.Info("nginx reload success")
	}
	return nil
}
