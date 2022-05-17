package main

import (
	"pekonode/server"
	"time"
)

func main() {

	createNode(8000)
	time.Sleep(5 * time.Second)
	createNode(8001)
	time.Sleep(15 * time.Second)
	createNode(8002)
	time.Sleep(25 * time.Second)
	nl := createNode(8003)
	time.Sleep(40 * time.Second)
	stop(nl)

	for {
		time.Sleep(3000)
	}
}

func createNode(port int) server.NodeList {
	var node server.Node
	node.Addr = "0.0.0.0"
	node.Port = port

	nodeList := server.New(node, "0.0.0.0", port, 3, 10, 5, 60)

	nodeList.Set(server.Node{
		Addr: "0.0.0.0",
		Port: 8000,
	})

	nodeList.Join()
	return nodeList
}

func stop(nodeList server.NodeList) {
	nodeList.Quit()
}
