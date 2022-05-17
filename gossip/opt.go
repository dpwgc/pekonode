package gossip

import (
	"sync"
	"time"
)

// New 创建本地节点列表
func New(node Node, listenAddr string, listenPort int, amount int, cycle int64, buffer int, timeout int64, isPrint bool) NodeList {

	var nodes sync.Map
	nodes.Store(node, time.Now().Unix())

	var status = make(map[int]bool, 1)
	status[1] = true

	nodeList := NodeList{
		nodes:      nodes,
		amount:     amount,
		cycle:      cycle,
		buffer:     buffer,
		timeout:    timeout,
		localNode:  node,
		listenAddr: listenAddr,
		listenPort: listenPort,
		status:     status,
		isPrint:    isPrint,
	}

	return nodeList
}

// Join 加入集群
func (nodeList *NodeList) Join() {

	//定时发布本地节点列表信息
	go task(nodeList)

	//监听队列（UDP监听缓冲区）
	var mq = make(chan []byte, nodeList.buffer)

	//监听其他节点的信息，并放入mq队列
	go listen(nodeList, mq)

	//消费mq队列中的信息
	go consume(nodeList, mq)

	nodeList.println(time.Now().Format("2006-01-02 15:04:05"), "/ [Join]:", nodeList.localNode)
}

// Stop 停止同步
func (nodeList *NodeList) Stop() {
	nodeList.println(time.Now().Format("2006-01-02 15:04:05"), "/ [Stop]:", nodeList.listenAddr, nodeList.listenPort)
	nodeList.status[1] = false
}

// Set 向本地节点列表中加入其他节点
func (nodeList *NodeList) Set(node Node) {
	nodeList.nodes.Store(node, time.Now().Unix())
}

// Get 获取本地节点列表
func (nodeList *NodeList) Get() []Node {
	var nodes []Node
	// 遍历所有sync.Map中的键值对
	nodeList.nodes.Range(func(k, v interface{}) bool {
		//如果该节点超过一段时间没有更新
		if v.(int64)+nodeList.timeout < time.Now().Unix() {
			nodeList.nodes.Delete(k)
		} else {
			nodes = append(nodes, k.(Node))
		}
		return true
	})
	return nodes
}
