// 开启一个udp监听
// 当收到link时代表想要连接，保存对应的ip
// 当收到两个不同link时交换两边ip
// 当收到两个ip的linkOk时关闭服务

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
	// 使用切片保存不同的peer
	var peerA *net.UDPAddr
	data := make([]byte, 1024)
	// 开启监听
	for {
		n, remoteAddr, err := listener.ReadFromUDP(data)
		if err != nil {
			panic(err)
		}
		log.Printf("监听到地址<%v> %v\n", remoteAddr.String(), string(data[:n]))
		if string(data[:n]) == "link" {
			if peerA == nil {
				peerA = remoteAddr
			} else {
				// 开启双向连接以及新的检测
				go p2pLink(listener, peerA, remoteAddr)
				break
			}
		}
	}
	ch := make(chan string)
	<-ch
}

func p2pLink(conn *net.UDPConn, addr1 *net.UDPAddr, addr2 *net.UDPAddr) {
	// 向两边发送对应的ip以及提示词 goLink   也就是"goLink ip"格式
	// 向让a给b发送信息
	sendPeer(conn, addr1, "goLink "+addr2.String())
	flg1, flg2 := false, false
	data := make([]byte, 1024)
	// 开启监听
	// 收到a发送成功后让b给a发送信息   sendOK
	//
	//	// 收到b发送成功后再让a发送信息
	//
	//	// 直到两边都发送了 linkOK 为止
	for {
		if flg1 && flg2 {
			// 让两边都断开连接
			sendPeer(conn, addr1, "down")
			sendPeer(conn, addr2, "down")
			break
		}
		n, remoteAddr, err := conn.ReadFromUDP(data)
		if err != nil {
			panic(err)
		}
		log.Printf("监听到地址<%v> %v\n", remoteAddr.String(), string(data[:n]))
		if string(data[:n]) == "sendOK" {
			if remoteAddr.String() == addr1.String() {
				sendPeer(conn, addr2, "goLink "+addr1.String())
			} else if remoteAddr.String() == addr2.String() {
				sendPeer(conn, addr1, "goLink "+addr2.String())
			}
		} else if string(data[:n]) == "linkOK" {
			if remoteAddr.String() == addr1.String() {
				sendPeer(conn, addr2, "goLink "+addr1.String())
				flg1 = true
			} else if remoteAddr.String() == addr2.String() {
				sendPeer(conn, addr1, "goLink "+addr2.String())
				flg2 = true
			}
			time.Sleep(time.Second)
		}
	}

}

func sendPeer(conn *net.UDPConn, addr *net.UDPAddr, msg string) {
	defer func() {
		println(recover())
	}()
	_, err := conn.WriteToUDP([]byte(msg), addr)
	if err != nil {
		log.Panic(err)
	}
}
