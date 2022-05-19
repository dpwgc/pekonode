package gossip

import (
	"sync"
	"time"
)

// New 初始化本地节点列表
func (nodeList *NodeList) New(localNode Node) {

	//ListenAddr 缺省值：0.0.0.0
	if nodeList.ListenAddr == "" {
		nodeList.ListenAddr = "0.0.0.0"
	}

	//Amount 缺省值：3
	if nodeList.Amount == 0 {
		nodeList.Amount = 3
	}

	//Cycle 缺省值：6
	if nodeList.Cycle == 0 {
		nodeList.Cycle = 6
	}

	//Buffer 缺省值：不填则默认等于Amount乘3
	if nodeList.Buffer == 0 {
		nodeList.Buffer = nodeList.Amount * 3
	}

	//Size 缺省值：1024
	if nodeList.Size == 0 {
		nodeList.Size = 1024
	}

	//Timeout 缺省值：如果当前Timeout小于或等于Cycle，则自动扩大Timeout的值
	if nodeList.Timeout <= nodeList.Cycle {
		nodeList.Timeout = nodeList.Cycle*3 + 2
	}

	//节点集合
	var nodes sync.Map

	//将本地节点添加进本地节点集合
	nodes.Store(localNode, time.Now().Unix())

	//初始化节点服务状态
	var status = make(map[int]bool, 1)
	status[1] = true

	nodeList.nodes = nodes         //本地节点集合
	nodeList.localNode = localNode //本地节点信息
	nodeList.status = status       //节点服务状态
}

// Join 加入集群
func (nodeList *NodeList) Join() {

	//如果该节点的本地节点列表还未初始化
	if nodeList.status == nil {
		//直接返回
		return
	}

	//定时广播本地节点信息
	go task(nodeList)

	//监听队列（UDP监听缓冲区）
	var mq = make(chan []byte, nodeList.Buffer)

	//监听其他节点的信息，并放入mq队列
	go listener(nodeList, mq)

	//消费mq队列中的信息
	go consume(nodeList, mq)

	nodeList.println("[Join]:", nodeList.localNode)
}

// Stop 停止广播心跳
func (nodeList *NodeList) Stop() {

	//如果该节点的本地节点列表还未初始化
	if nodeList.status == nil {
		//直接返回
		return
	}

	nodeList.println("[Stop]:", nodeList.localNode)
	nodeList.status[1] = false
}

// Start 重新开始广播心跳
func (nodeList *NodeList) Start() {

	//如果该节点的本地节点列表还未初始化
	if nodeList.status == nil {
		//直接返回
		return
	}

	//如果当前心跳服务正常
	if nodeList.status[1] {
		//返回
		return
	}
	nodeList.println("[Start]:", nodeList.localNode)
	nodeList.status[1] = true
	//定时广播本地节点信息
	go task(nodeList)
}

// Set 向本地节点列表中加入其他节点
func (nodeList *NodeList) Set(node Node) {

	//如果该节点的本地节点列表还未初始化
	if nodeList.status == nil {
		//直接返回
		return
	}

	nodeList.nodes.Store(node, time.Now().Unix())
}

// Get 获取本地节点列表
func (nodeList *NodeList) Get() []Node {

	//如果该节点的本地节点列表还未初始化
	if nodeList.status == nil {
		//直接返回
		return nil
	}

	var nodes []Node
	// 遍历所有sync.Map中的键值对
	nodeList.nodes.Range(func(k, v interface{}) bool {
		//如果该节点超过一段时间没有更新
		if v.(int64)+nodeList.Timeout < time.Now().Unix() {
			nodeList.nodes.Delete(k)
		} else {
			nodes = append(nodes, k.(Node))
		}
		return true
	})
	return nodes
}
