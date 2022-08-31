package etcd

import (
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

var c *clientv3.Client
var ETCDPath = "http://etcd:2379"

func initClient() error {
	// init etcd client
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{ETCDPath},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return err
	}
	c = client
	return nil
}

func GetClient() *clientv3.Client {
	if err := initClient(); err != nil {
		return nil
	}
	return c
}
