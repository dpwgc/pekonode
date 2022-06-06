package pekonode

// write 发送数据
func write(addr string, port int, data []byte) {
	udpWrite(addr, port, data)
}

// listen 服务端监听
func listen(addr string, port int, size int, mq chan []byte) {
	udpListen(addr, port, size, mq)
}
