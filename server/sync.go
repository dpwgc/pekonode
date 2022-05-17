package server

import (
	"encoding/json"
	"fmt"
	"time"
)

//定时同步任务
func task(nodeList *NodeList) {
	for {
		//停止同步
		if !nodeList.status[1] {
			break
		}

		//发送本地节点信息
		broadcast(nodeList, nodeList.localNode)
		fmt.Println("[Listen]:", nodeList.localNode, "/ [Node list]:", nodeList.Get())
		time.Sleep(time.Duration(nodeList.Cycle) * time.Second)
	}
}

//监听其他节点发来的同步信息
func listen(nodeList *NodeList, mq chan []byte) {
	//监听协程
	Listen(nodeList.ListenAddr, nodeList.ListenPort, mq)
}

//消费信息
func consume(nodeList *NodeList, mq chan []byte) {
	for {
		//从监听队列中取出消息
		data := <-mq
		var node Node
		err := json.Unmarshal(data, &node)
		if err != nil {
			fmt.Println(err)
		}

		//查看本地节点列表是否存在该节点
		_, ok := nodeList.Nodes.Load(node)

		//更新本地列表
		nodeList.Set(node)

		//如果存在，则跳过广播环节
		if ok {
			continue
		}

		//广播推送
		broadcast(nodeList, node)
	}
}

//广播推送信息
func broadcast(nodeList *NodeList, node Node) {

	//取出所有未过期的节点
	nodes := nodeList.Get()
	//如果当前节点数量小于Amount最大推送数量
	if len(nodes) <= nodeList.Amount {
		for _, n := range nodes {
			b, _ := json.Marshal(node)
			Write(n.Addr, n.Port, b)
		}
	} else {
		//向部分节点发送信息
		for i := 0; i < nodeList.Amount; i++ {
			b, _ := json.Marshal(node)
			Write(nodes[i].Addr, nodes[i].Port, b)
		}
	}
}
