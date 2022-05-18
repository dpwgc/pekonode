package gossip

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

		//将本地节点加入已传染的节点列表sentList
		var sentList = make(map[Node]int)
		sentList[nodeList.localNode] = 1
		//更新本地节点信息
		nodeList.Set(nodeList.localNode)

		//设置心跳数据包
		sn := sendNode{
			Node:     nodeList.localNode,
			SentList: sentList,
		}

		//发送本地节点心跳数据包
		broadcast(nodeList, sn)
		nodeList.println(time.Now().Format("2006-01-02 15:04:05"), "/ [Listen]:", nodeList.localNode, "/ [Node list]:", nodeList.Get())
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
		bs := <-mq
		var sn sendNode
		err := json.Unmarshal(bs, &sn)
		if err != nil {
			fmt.Println(err)
		}

		node := sn.Node

		//更新本地列表
		nodeList.Set(node)

		//广播推送该节点
		broadcast(nodeList, sn)
	}
}

//广播推送信息
func broadcast(nodeList *NodeList, sn sendNode) {

	//取出所有未过期的节点
	nodes := nodeList.Get()
	var sendNodes []sendNode

	//选取部分未被传染的节点
	i := 0
	for _, n := range nodes {

		//如果超过Amount最大推送数量
		if i >= nodeList.Amount {
			//结束广播
			break
		}

		//如果该节点已经被传染过了
		if sn.SentList[n] == 1 {
			//跳过该节点
			continue
		}

		//将该节点添加进发送列表
		sn.SentList[n] = 1 //标记该节点为已传染状态
		sn.TargetNode = n  //设置发送目标节点
		sendNodes = append(sendNodes, sn)
		i++
	}

	//向这些未被传染的节点广播传染数据
	for _, n := range sendNodes {
		bs, _ := json.Marshal(sn)
		Write(n.TargetNode.Addr, n.TargetNode.Port, bs)
	}
}
