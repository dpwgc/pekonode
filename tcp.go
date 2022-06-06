package pekonode

import (
	"fmt"
	"net"
)

//向TCP服务端发送数据
func tcpWrite(addr string, port int, data []byte) {

	tcpAddr := fmt.Sprintf("%s:%v", addr, port)
	server, err := net.ResolveTCPAddr("tcp4", tcpAddr)

	if err != nil {
		println("[Error]:", err)
	}

	//建立服务器连接
	conn, err := net.DialTCP("tcp", nil, server)

	defer func(conn *net.TCPConn) {
		err = conn.Close()
		if err != nil {
			println("[Error]:", err)
		}
	}(conn)

	_, err = conn.Write(data) //给服务器发信息
	if err != nil {
		println("[Error]:", err)
	}
}

func tcpListen(addr string, port int, size int, mq chan []byte) {
	server, err := net.Listen("tcp", fmt.Sprintf("%s:%v", addr, port))
	if err != nil {
		println("[Error]:", err)
		return
	}
	defer func(server net.Listener) {
		err = server.Close()
		if err != nil {
			println("[Error]:", err)
		}
	}(server)

	for {
		conn, err := server.Accept()
		if err != nil {
			continue
		}
		go func() {

			//接收数组
			bs := make([]byte, size)
			n, err := conn.Read(bs)
			if err != nil {
				println("[Error]:", err)
				return
			}

			//获取有效数据
			b := bs[:n]

			//将数据放入缓冲队列，异步处理数据
			mq <- b
		}()
	}
}
