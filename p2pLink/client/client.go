package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

var serverAddr = getAddr("8.134.141.213:9001")
var localAddr *net.UDPAddr
var peerAddr *net.UDPAddr
var addrs = map[string]*net.UDPAddr{}

func main() {
	// 先读入自己端口号
	var port int
	fmt.Print("请输入端口号: ")
	fmt.Scan(&port)
	localAddr = &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: port,
	}
	// 开启端口监听
	conn, err := net.ListenUDP("udp", localAddr)
	if err != nil {
		log.Panic(err)
	}
	// 对指令做出判断以及对应输出
	go func() {
		data := make([]byte, 1024)
		for {
			n, addr, err2 := conn.ReadFromUDP(data)
			if err2 != nil {
				log.Println(err)
			} else {
				// 将地址添加进保活名单里
				if _, ok := addrs[addr.String()]; !ok {
					addrs[addr.String()] = addr
				}
				// 获取发送的信息
				msg := string(data[:n])
				log.Printf("监听到地址<%v> [%v]\n", addr.String(), string(data[:n]))
				if msg == "down" {
					delete(addrs, addr.String())
					time.Sleep(2)
					delete(addrs, addr.String())
				} else {
					strs := strings.Split(msg, " ")
					if strs[0] == "goLink" {
						if inAddrList(strs[1]) {
							conn.WriteToUDP([]byte("linkOK"), addr)
						} else {
							sendMsg(getAddr(strs[1]), "hello I`m "+localAddr.String())
							conn.WriteToUDP([]byte("sendOK"), addr)
							log.Printf("发送信息到<%s>\n", strs[1])
						}
						time.Sleep(2 * time.Second)
					}
				}
			}
		}
	}()

	// 与每一个端口保持通信
	go func() {
		for {
			// 循环记录表发出信号
			for addr := range addrs {
				conn.WriteToUDP([]byte("keepLive"), addrs[addr])
			}
			time.Sleep(time.Second * 10)
		}
	}()
	sendMsg(serverAddr, "link")

	ch := make(chan struct{})
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

func inAddrList(addrStr string) bool {
	_, ok := addrs[addrStr]
	return ok
}

func sendMsg(addr *net.UDPAddr, msg string) {
	conn, err := net.DialUDP("udp", localAddr, addr)
	if err != nil {
		log.Panic(err)
	}
	defer conn.Close()
	_, err = conn.Write([]byte(msg))
	if err != nil {
		log.Panic(err)
	}
}
