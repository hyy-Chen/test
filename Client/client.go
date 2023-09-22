package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

var tag string

const HAND_SHAKE_MSG = "keeplive"

func main() {
	fmt.Print("输入标识: ")
	fmt.Scan(&tag)
	var port int
	fmt.Print("输入端口: ")
	fmt.Scan(&port)

	srcAddr := &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: port,
	}

	dstAddr := &net.UDPAddr{
		IP:   net.ParseIP("8.134.141.213"),
		Port: 9001,
	}
	conn, err := net.DialUDP("udp", srcAddr, dstAddr)
	if err != nil {
		panic(err)
	}
	if _, err := conn.Write([]byte("hello, I`m new peer:" + tag)); err != nil {
		panic(err)
	}
	data := make([]byte, 1024)

	n, remoteAddr, err := conn.ReadFromUDP(data)
	if err != nil {
		fmt.Printf("error during read: %s\n", err)
	}
	conn.Close()
	anotherPeer := parseAddr(string(data[:n]))
	log.Printf("local:%s server %s another: %s\n", srcAddr, remoteAddr, anotherPeer.String())
	bidirectionHole(srcAddr, &anotherPeer)
}

func parseAddr(addr string) net.UDPAddr {
	t := strings.Split(addr, ":")
	port, _ := strconv.Atoi(t[1])
	return net.UDPAddr{
		IP:   net.ParseIP(t[0]),
		Port: port,
	}
}

func bidirectionHole(srcAddr *net.UDPAddr, antherAddr *net.UDPAddr) {

	conn, err := net.DialUDP("udp", srcAddr, antherAddr)
	if err != nil {
		log.Panic("send handshake:", err)
	}

	go func() {
		for {
			time.Sleep(time.Second * 10)
			if _, err := conn.Write([]byte("keepLive")); err != nil {
				log.Panic("send keepLive fail", err)
			}
		}
	}()

	go func() {
		for {
			time.Sleep(time.Second * 10)
			if _, err := conn.Write([]byte("from [" + tag + "]")); err != nil {
				log.Panic("send msg fail", err)
			}
		}
	}()

	for {
		conn.SetReadDeadline(time.Now().Add(15 * time.Second))
		data := make([]byte, 1024)
		n, _, err := conn.ReadFromUDP(data)
		if err != nil {
			log.Panic(err)
		} else {
			log.Printf("收到数据: %s\n", data[:n])
		}
	}
}

func keepAlive() {

}
