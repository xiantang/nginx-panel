package main

import (
	"bytes"
	"context"
	"golang/internal/etcd"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var filePath = "/etc/nginx/tests/test.conf"

type NginxReq struct {
	ConfigBody string `json:"config_body"`
	SererName  string `json:"server_name"`
}

// submit nginx config to etcd
func submit(c *gin.Context) {
	var nginxReq NginxReq
	if err := c.BindJSON(&nginxReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := handleSubmit(nginxReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return

	}
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

func handleSubmit(nginxReq NginxReq) error {
	err := NginxTest(nginxReq.ConfigBody)
	if err != nil {
		return err
	}

	// write success config into etcd
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{etcd.ETCDPath},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.WithError(err).Error("connect etcd error")
		return err
	}
	defer cli.Close()

	log.WithField("config", nginxReq.SererName).Info("write config to etcd")

	key := "/config/" + nginxReq.SererName
	// with timeout
	log.WithField("key", key).Info("put config to etcd")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = cli.Put(ctx, key, nginxReq.ConfigBody)
	if err != nil {
		log.WithError(err).Error("put config to etcd error")
		return err
	}

	return nil
}

// NginxTest do nginx test nginx -t
func NginxTest(body string) error {
	if os.Getenv("SKIP_NGINX_TEST") == "true" {
		return nil
	}
	log.WithField("config", body).Info("nginx test")
	// write file to /etc/nginx/test/ than do nginx -t

	if err := ioutil.WriteFile(filePath, []byte(body), 0644); err != nil {
		log.WithField("config", body).WithError(err).Error("write file error")
		return err
	}
	defer os.Remove(filePath)
	cmd := exec.Command("nginx", "-c", "/etc/nginx/nginx_test.conf", "-t")
	buf := bytes.NewBufferString("")
	cmd.Stderr = buf
	err := cmd.Run()
	if err != nil {
		log.WithField("output", buf.String()).WithError(err).Error("nginx -t error")
		return err
	}

	return nil
}
