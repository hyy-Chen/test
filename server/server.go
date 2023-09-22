package main

import (
	"log"
	"net"
	"time"
)

func main() {
	// 开启本地监听地址9001
	listener, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 9001,
	})
	if err != nil {
		panic(err)
	}
	log.Printf("本地监听[%v]\n", listener.LocalAddr().String())
	// 存储一下想要开启p2p的地址
	peers := make([]net.UDPAddr, 0, 2)
	data := make([]byte, 1024)
	for {
		n, remoteAddr, err := listener.ReadFromUDP(data)
		if err != nil {
			panic(err)
		}
		log.Printf("监听到地址<%v> %s\n", remoteAddr.String(), data[:n])
		peers = append(peers, *remoteAddr)
		if len(peers) == 2 {
			log.Printf("进行UDP打洞， 建立 %s<---->%s 的链接\n", peers[0].String(), peers[1].String())
			listener.WriteToUDP([]byte(peers[0].String()), &peers[1])
			listener.WriteToUDP([]byte(peers[1].String()), &peers[0])
			time.Sleep(time.Second * 10)
			log.Println("中转服务器退出")
			return
		}
	}
}
