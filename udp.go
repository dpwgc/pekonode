package pekonode

import (
	"fmt"
	"net"
	"os"
)

// udpWrite 发送udp数据
func udpWrite(nodeList *NodeList, addr string, port int, data []byte) {
	socket, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.ParseIP(addr),
		Port: port,
	})
	if err != nil {
		nodeList.println("[Error]:", err)
		return
	}

	_, err = socket.Write(data) // 发送数据
	if err != nil {
		nodeList.println("[Error]:", err)
		return
	}

	defer func(socket *net.UDPConn) {
		err = socket.Close()
		if err != nil {
			nodeList.println("[Error]:", err)
		}
	}(socket)
}

// udpListen udp服务端监听
func udpListen(nodeList *NodeList, mq chan []byte) {

	udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", nodeList.ListenAddr, nodeList.localNode.Port))
	if err != nil {
		nodeList.println("[Error]:", err)
		os.Exit(1)
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		nodeList.println("[Error]:", err)
		os.Exit(1)
	}
	defer func(conn *net.UDPConn) {
		err = conn.Close()
		if err != nil {
			nodeList.println("[Error]:", err)
		}
	}(conn)

	for {
		//接收数组
		bs := make([]byte, nodeList.Size)

		//从UDP监听中接收数据
		n, _, err := conn.ReadFromUDP(bs)
		if err != nil {
			nodeList.println("[Error]:", err)
			continue
		}

		//获取有效数据
		b := bs[:n]

		//将数据放入缓冲队列，异步处理数据
		mq <- b
	}
}
