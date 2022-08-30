package main

import (
	"context"
	"github.com/stretchr/testify/assert"
	clientv3 "go.etcd.io/etcd/client/v3"
	"testing"
)

func Test_handleSubmit(t *testing.T) {
	t.Setenv("SKIP_NGINX_TEST", "true")
	req := NginxReq{
		ConfigBody: "server {\\nlisten 127.0.0.1:8089;\\n# Additional server configuration\\nlocation /some/path/ {\\nreturn 200;\\n}\\n}",
		SererName:  "shopee.com",
	}

	path := t.TempDir() + "/test.conf"
	filePath = path
	etcdPATH = "http://localhost:2379"
	err := handleSubmit(req)
	assert.NoError(t, err)
	// get value from etcd
	cli, err := clientv3.New(clientv3.Config{
		Endpoints: []string{etcdPATH},
	})
	assert.NoError(t, err)
	defer cli.Close()
	resp, err := cli.Get(context.Background(), "/config/shopee.com")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(resp.Kvs))
	assert.Equal(t, req.ConfigBody, string(resp.Kvs[0].Value))
}
