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

		sn := sendNode{
			Node: nodeList.localNode,
			Num:  0,
		}

		//发送本地节点信息
		broadcast(nodeList, sn)
		nodeList.println(time.Now().Format("2006-01-02 15:04:05"), "/ [Listen]:", nodeList.localNode, "/ [Node list]:", nodeList.Get())
		time.Sleep(time.Duration(nodeList.cycle) * time.Second)
	}
}

//监听其他节点发来的同步信息
func listen(nodeList *NodeList, mq chan []byte) {
	//监听协程
	Listen(nodeList.listenAddr, nodeList.listenPort, mq)
}

//消费信息
func consume(nodeList *NodeList, mq chan []byte) {
	for {
		//从监听队列中取出消息
		bs := <-mq
		var sn sendNode
		err := json.Unmarshal(bs, &sn)
		if err != nil {
			fmt.Println(err)
		}

		//数据包转发次数+1
		sn.Num = sn.Num + 1

		node := sn.Node

		//查看本地节点列表是否存在该节点
		_, ok := nodeList.nodes.Load(node)

		//更新本地列表
		nodeList.Set(node)

		//如果存在，且该消息广播次数不超过3次
		if ok && sn.Num >= 3 {
			//跳过广播
			continue
		}
		//广播推送该节点
		broadcast(nodeList, sn)
	}
}

//广播推送信息
func broadcast(nodeList *NodeList, sn sendNode) {

	//取出所有未过期的节点
	nodes := nodeList.Get()
	//如果当前节点数量小于Amount最大推送数量
	if len(nodes) <= nodeList.amount {
		for _, n := range nodes {
			bs, _ := json.Marshal(sn)
			Write(n.Addr, n.Port, bs)
		}
	} else {
		//向部分节点发送信息
		for i := 0; i < nodeList.amount; i++ {
			bs, _ := json.Marshal(sn)
			Write(nodes[i].Addr, nodes[i].Port, bs)
		}
	}
}
