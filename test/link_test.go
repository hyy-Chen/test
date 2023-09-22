/**
 @author : ikkk
 @date   : 2023/9/21
 @text   : nil
**/

package test

import (
	"fmt"
	"net"
	"os"
	"testing"
	"time"
)

func TestLink(t *testing.T) {
	// 目标地址
	targetIP := "120.236.146.57" // 请替换为实际的目标IP地址
	targetPort := 44448          // 请替换为实际的目标端口号

	// 创建UDP连接
	targetAddr := &net.UDPAddr{
		IP:   net.ParseIP(targetIP),
		Port: targetPort,
	}

	conn, err := net.DialUDP("udp", nil, targetAddr)
	if err != nil {
		fmt.Println("Error creating UDP connection:", err)
		os.Exit(1)
	}
	defer conn.Close()

	// 要发送的数据
	message := []byte("Hello, UDP!")

	for {
		// 发送数据
		_, err = conn.Write(message)
		if err != nil {
			fmt.Println("Error sending UDP message:", err)
			os.Exit(1)
		}
		time.Sleep(2 * time.Second)
	}
	ch := make(chan struct{})
	<-ch
}

func TestA(t *testing.T) {
	a := "ok"
	b := string([]byte("ok"))
	fmt.Println(a == b)
}
