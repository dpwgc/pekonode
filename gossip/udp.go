package gossip

import (
	"fmt"
	"net"
	"os"
)

// write 发送udp数据
func write(addr string, port int, data []byte) {
	socket, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.ParseIP(addr),
		Port: port,
	})
	if err != nil {
		println("[Error]:", err)
	}

	_, err = socket.Write(data) // 发送数据
	if err != nil {
		println("[Error]:", err)
	}

	err = socket.Close()
	if err != nil {
		println("[Error]:", err)
	}
}

// listen udp服务端监听
func listen(addr string, port int, mq chan []byte) {

	udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		println("[Error]:", err)
		os.Exit(1)
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		println("[Error]:", err)
		os.Exit(1)
	}
	defer func(conn *net.UDPConn) {
		err = conn.Close()
		if err != nil {
			println("[Error]:", err)
		}
	}(conn)

	for {
		bs := make([]byte, 8192)
		_, _, err = conn.ReadFromUDP(bs)
		var b []byte
		for _, v := range bs {
			if v == 0x0 {
				break
			}
			b = append(b, v)
		}

		mq <- b

		if err != nil {
			println("[Error]:", err)
			continue
		}
	}
}
