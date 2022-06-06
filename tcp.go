package pekonode

import (
	"fmt"
	"net"
)

//向TCP服务端发送数据
func tcpWrite(addr string, port string, data string) {

	tcpAddr := fmt.Sprintf("%s:%s", addr, port)
	server, err := net.ResolveTCPAddr("tcp4", tcpAddr)

	if err != nil {
		println("[Error]:", err)
	}

	//建立服务器连接
	conn, err := net.DialTCP("tcp", nil, server)

	defer func(conn *net.TCPConn) {
		err := conn.Close()
		if err != nil {
			println("[Error]:", err)
		}
	}(conn)

	_, err = conn.Write([]byte(data)) //给服务器发信息
	if err != nil {
		println("[Error]:", err)
	}
}
