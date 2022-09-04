package main

import (
	"fmt"
	"golang/internal/etcd"
	"log"
	"net"
)

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
func main() {
	fmt.Println(GetOutboundIP().To4().String())
	etcd.NewSyncer().Run()
}
