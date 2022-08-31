package etcd

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func TestBasicWatch(t *testing.T) {
	t.Log("TestBasicWatch")
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
