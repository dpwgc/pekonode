package pekonode

import (
	"fmt"
	"net"
	"os"
)

//向TCP服务端发送数据
func tcpWrite(nodeList *NodeList, addr string, port int, data []byte) {

	tcpAddr := fmt.Sprintf("%s:%v", addr, port)
	server, err := net.ResolveTCPAddr("tcp4", tcpAddr)

	if err != nil {
		nodeList.println("[Error]:", err)
		return
	}

	//建立服务器连接
	conn, err := net.DialTCP("tcp", nil, server)
	if err != nil {
		nodeList.println("[Error]:", err)
		return
	}

	_, err = conn.Write(data) //给服务器发信息
	if err != nil {
		nodeList.println("[Error]:", err)
	}

	defer func(conn *net.TCPConn) {
		err = conn.Close()
		if err != nil {
			nodeList.println("[Error]:", err)
		}
	}(conn)
}

func tcpListen(nodeList *NodeList, mq chan []byte) {
	server, err := net.Listen("tcp", fmt.Sprintf("%s:%v", nodeList.ListenAddr, nodeList.localNode.Port))
	if err != nil {
		nodeList.println("[Error]:", err)
		return
	}
	defer func(server net.Listener) {
		err = server.Close()
		if err != nil {
			nodeList.println("[Error]:", err)
		}
		os.Exit(1)
	}(server)

	for {
		conn, err := server.Accept()
		if err != nil {
			continue
		}
		go func() {

			//接收数组
			bs := make([]byte, nodeList.Size)
			n, err := conn.Read(bs)
			if err != nil {
				nodeList.println("[Error]:", err)
				return
			}

			//获取有效数据
			b := bs[:n]

			//将数据放入缓冲队列，异步处理数据
			mq <- b
		}()
	}
}
