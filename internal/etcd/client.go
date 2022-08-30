package etcd

import (
	clientv3 "go.etcd.io/etcd/client/v3"
)

var c *clientv3.Client

func init() {
	initClient()
}

func initClient() error {
	// init etcd client
	client, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"http://etcd:2379"},
	})
	if err != nil {
		return err
	}
	c = client
	return nil
}

func GetClient() *clientv3.Client {
	return c
}
