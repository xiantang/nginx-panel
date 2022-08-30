package main

import (
	"golang/internal/etcd"
)

func main() {
	etcd.NewSyncer().Run()
}
