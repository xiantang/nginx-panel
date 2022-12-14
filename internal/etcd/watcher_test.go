package etcd

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func TestUpdateConfig(t *testing.T) {
	t.Log("TestUpdateConfig")
	t.Setenv("SKIP_NGINX_TEST", "true")
	t.Setenv("SKIP_RELOAD", "true")
	ETCDPath = "http://localhost:2379"
	client := GetClient()
	ctx := context.Background()
	contxt, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	client.Delete(contxt, "/", clientv3.WithPrefix())

	contxt, cancel = context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	key := "/config/" + "abc"
	value := "server {}"
	_, err := client.Put(contxt, key, value)
	assert.NoError(t, err)

	syncer := NewSyncer()
	err = syncer.sync()
	assert.NoError(t, err)
	assert.Len(t, syncer.schemas, 1)

	contxt, cancel = context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	value = "server {listen 80;}"
	_, err = client.Put(contxt, key, value)
	assert.NoError(t, err)

	err = syncer.sync()
	assert.NoError(t, err)
	assert.Len(t, syncer.schemas, 1)
	assert.Equal(t, value, syncer.schemas[key])
}

// TestBasicAddConfig tests the basic add config
// write config just `server {}` to etcd
func TestBasicAddConfig(t *testing.T) {
	t.Setenv("SKIP_NGINX_TEST", "true")
	t.Setenv("SKIP_RELOAD", "true")
	t.Log("TestBasicWatch")
	ETCDPath = "http://localhost:2379"
	client := GetClient()
	ctx := context.Background()
	contxt, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	client.Delete(contxt, "/", clientv3.WithPrefix())

	contxt, cancel = context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	key := "/config/" + "abc"
	value := "server {}"
	_, err := client.Put(contxt, key, value)
	assert.NoError(t, err)

	syncer := NewSyncer()
	err = syncer.sync()
	assert.NoError(t, err)
	assert.Len(t, syncer.schemas, 1)
}

func TestBasicWatch(t *testing.T) {
	t.Log("TestBasicWatch")
	t.Setenv("SKIP_NGINX_TEST", "true")
	t.Setenv("SKIP_RELOAD", "true")
	ETCDPath = "http://localhost:2379"
	client := GetClient()
	ctx := context.Background()
	contxt, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	client.Delete(contxt, "/", clientv3.WithPrefix())
	syncer := NewSyncer()
	err := syncer.sync()
	assert.NoError(t, err)
}
