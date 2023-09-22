package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

// 监听9091端口，向9001发送信息，读取信息发现有从9002到的进行输出，1min后没有信息接收显示down

var serverAddr = getAddr("8.134.141.213:9001")

func main() {
	var port int
	fmt.Print("输入开启监听端口: ")
	fmt.Scan(&port)
	// 开启本地监听地址9001
	listener, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: port,
	})
	if err != nil {
		panic(err)
	}
	localAddr := getAddr("0.0.0.0:9091")
	conn, err := net.DialUDP("udp", localAddr, serverAddr)
	if err != nil {
		log.Panic(err)
	}
	// 要发送的数据
	message := []byte("link")
	// 发送数据
	_, err = conn.Write(message)
	if err != nil {
		panic(err)
	}
	conn.Close()

	data := make([]byte, 1024)

	for {
		n, remoteAddr, err := listener.ReadFromUDP(data)
		if err != nil {
			panic(err)
		}
		log.Printf("监听到地址<%v> %s\n", remoteAddr.String(), data[:n])
		listener.WriteToUDP([]byte("I here you "+listener.LocalAddr().String()), remoteAddr)
	}

	// 先向服务器发起信息

	ch := make(chan struct{}, 1)
	<-ch

}

func getAddr(addr string) *net.UDPAddr {
	t := strings.Split(addr, ":")
	port, _ := strconv.Atoi(t[1])
	return &net.UDPAddr{
		IP:   net.ParseIP(t[0]),
		Port: port,
	}
}
