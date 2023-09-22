package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

var ip net.IP

func main() {
	link1()
}

func link1() {
	// 开启本地监听地址9001
	listener, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 9001,
	})
	if err != nil {
		panic(err)
	}
	log.Printf("本地监听[%v]\n", listener.LocalAddr().String())
	data := make([]byte, 1024)
	for {
		n, remoteAddr, err := listener.ReadFromUDP(data)
		if err != nil {
			panic(err)
		}
		log.Printf("监听到地址<%v> %s\n", remoteAddr.String(), data[:n])
		if string(data[:n]) == "link" {
			log.Printf("link addr %v", remoteAddr.String())
			go func() {
				for {
					listener.WriteToUDP([]byte("I here you "+listener.LocalAddr().String()), remoteAddr)
					fmt.Println("发送成功")
					time.Sleep(5 * time.Second)
				}
			}()
		}
	}
	ch := make(chan struct{})
	<-ch
}
