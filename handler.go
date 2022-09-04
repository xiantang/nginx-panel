package main

import (
	"context"
	"errors"
	"golang/internal/etcd"
	"golang/internal/nginx"
	"net/http"
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

// list all nginx config with server name and body
func list(c *gin.Context) {
	bodys, err := handleList(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": bodys,
	})
}

// delete nginx config from etcd by server name
func del(c *gin.Context) {
	// read from url  	group.DELETE("/del/:server_name", del)
	serverName := c.Param("server_name")
	err := handleDelete(c, serverName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "delete success",
	})
}

// handleDelete delete nginx config from etcd by server name
func handleDelete(c context.Context, serverName string) error {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{etcd.ETCDPath},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.WithError(err).Error("connect etcd error")
		return err
	}
	defer cli.Close()
	// delete config from etcd
	key := "/config/" + serverName
	log.WithField("key", key).Info("delete config from etcd")
	// check config exist
	resp, err := cli.Get(c, key)
	if err != nil {
		log.WithError(err).Error("get etcd error")
		return err
	}
	if len(resp.Kvs) == 0 {
		return errors.New("config not exist")
	}

	_, err = cli.Delete(c, key)
	if err != nil {
		log.WithError(err).Error("delete config from etcd error")
		return err
	}
	return nil

}

// create nginx config to etcd
func create(c *gin.Context) {
	var nginxReq NginxReq
	if err := c.BindJSON(&nginxReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := handleCreate(nginxReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return

	}
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

// handleList list all nginx config with server name and body
func handleList(c context.Context) ([]NginxReq, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{etcd.ETCDPath},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.WithError(err).Error("connect etcd error")
		return nil, err
	}
	resp, err := cli.Get(c, "/config", clientv3.WithPrefix())
	if err != nil {
		log.WithError(err).Error("get etcd error")
		return nil, err
	}
	defer cli.Close()

	var nginxResp []NginxReq
	for _, ev := range resp.Kvs {
		log.WithField("key", string(ev.Key)).Info("get config from etcd")
		nginxReq := NginxReq{
			ConfigBody: string(ev.Value),
			SererName:  string(ev.Key),
		}
		nginxResp = append(nginxResp, nginxReq)
	}

	return nginxResp, nil
}

func handleCreate(nginxReq NginxReq) error {
	err := nginx.Test(nginxReq.ConfigBody)
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
