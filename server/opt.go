package server

import (
	"fmt"
	"sync"
	"time"
)

// New 创建本地节点列表
func New(node Node, listenAddr string, listenPort int, amount int, cycle int64, buffer int, timeout int64) NodeList {

	var nodes sync.Map
	nodes.Store(node, time.Now().Unix())

	nodeList := NodeList{
		Nodes:      nodes,
		Amount:     amount,
		Cycle:      cycle,
		Buffer:     buffer,
		Timeout:    timeout,
		localNode:  node,
		ListenAddr: listenAddr,
		ListenPort: listenPort,
		status:     true,
	}

	return nodeList
}

// Join 加入集群
func (nodeList *NodeList) Join() {

	//定时发布本地节点列表信息
	go task(nodeList)

	//监听队列（UDP监听缓冲区）
	var mq = make(chan []byte, nodeList.Buffer)

	//监听其他节点的信息，并放入mq队列
	go listen(nodeList, mq)

	//消费mq队列中的信息
	go consume(nodeList, mq)
}

// Quit 退出集群
func (nodeList *NodeList) Quit() {
	fmt.Println("[Quit]: ", nodeList.ListenAddr+":", nodeList.ListenPort)
	nodeList.status = false
}

// Set 向本地节点列表中加入其他节点
func (nodeList *NodeList) Set(node Node) {
	nodeList.Nodes.Store(node, time.Now().Unix())
}

// Get 获取本地节点列表
func (nodeList *NodeList) Get() []Node {
	var nodes []Node
	// 遍历所有sync.Map中的键值对
	nodeList.Nodes.Range(func(k, v interface{}) bool {
		//如果该节点超过一段时间没有更新
		if v.(int64)+nodeList.Timeout < time.Now().Unix() {
			nodeList.Nodes.Delete(k)
		} else {
			nodes = append(nodes, k.(Node))
		}
		return true
	})
	return nodes
}
