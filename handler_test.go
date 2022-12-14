package main

import (
	"context"
	"golang/internal/etcd"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func Test_handleCURD(t *testing.T) {
	t.Setenv("SKIP_NGINX_TEST", "true")
	req := NginxReq{
		ConfigBody: "server {\\nlisten 127.0.0.1:8089;\\n# Additional server configuration\\nlocation /some/path/ {\\nreturn 200;\\n}\\n}",
		SererName:  "shopee.com",
	}

	path := t.TempDir() + "/test.conf"
	etcd.ETCDPath = "http://localhost:2379"
	filePath = path
	err := handleCreate(req)
	assert.NoError(t, err)
	// get value from etcd
	t.Log("get value from etcd")
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{etcd.ETCDPath},
		DialTimeout: 5 * time.Second,
	})
	assert.NoError(t, err)
	defer cli.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	resp, err := cli.Get(ctx, "/config/shopee.com")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(resp.Kvs))
	assert.Equal(t, req.ConfigBody, string(resp.Kvs[0].Value))

	// list config
	t.Log("list config")
	bodys, err := handleList(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 1, len(bodys))

	// get config
	t.Log("get config")
	body, err := handleGet(context.Background(), req.SererName)
	assert.NoError(t, err)
	assert.Equal(t, req.ConfigBody, body)

	// delete config
	t.Log("delete config")
	err = handleDelete(context.Background(), req.SererName)
	assert.NoError(t, err)
	resp, err = cli.Get(ctx, "/config/shopee.com")
	assert.NoError(t, err)
	assert.Equal(t, 0, len(resp.Kvs))

}
