package server

import (
	"fmt"
	"net"
	"os"
)

// Write 发送udp数据
func Write(addr string, port int, data []byte) {
	socket, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.ParseIP(addr),
		Port: port,
	})
	if err != nil {
		fmt.Println(err)
	}

	_, err = socket.Write(data) // 发送数据
	if err != nil {
		fmt.Println(err)
	}

	err = socket.Close()
	if err != nil {
		fmt.Println(err)
	}
}

// Listen udp服务端监听
func Listen(addr string, port int, mq chan []byte) {

	udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer func(conn *net.UDPConn) {
		err = conn.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(conn)

	for {
		bs := make([]byte, 4096)
		_, _, err := conn.ReadFromUDP(bs)

		var b []byte
		for _, v := range bs {
			if v == 0x0 {
				break
			}
			b = append(b, v)
		}
		mq <- b

		if err != nil {
			fmt.Println(err)
			continue
		}
	}
}
