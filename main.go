package main

import (
	"pekonode/gossip"
	"time"
)

func main() {

	createNode(8000)
	time.Sleep(5 * time.Second)
	createNode(8001)
	time.Sleep(15 * time.Second)
	createNode(8002)
	time.Sleep(25 * time.Second)
	createNode(8003)
	time.Sleep(40 * time.Second)
	//stop(nl)

	for {
		time.Sleep(10 * time.Second)
	}
}

func createNode(port int) gossip.NodeList {
	var node gossip.Node
	node.Addr = "0.0.0.0"
	node.Port = port
	node.Tag = "test"

	nodeList := gossip.New(node, "0.0.0.0", port, 3, 10, 5, 30, true)

	nodeList.Set(gossip.Node{
		Addr: "0.0.0.0",
		Port: 8000,
		Tag:  "test",
	})

	nodeList.Join()
	return nodeList
}

func stop(nodeList gossip.NodeList) {
	nodeList.Stop()
}
