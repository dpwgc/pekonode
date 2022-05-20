package gossip

import (
	"encoding/json"
	"strconv"
	"time"
)

//定时心跳广播任务
func task(nodeList *NodeList) {
	for {
		//停止同步
		if !nodeList.status.Load().(bool) {
			break
		}

		//将本地节点加入已传染的节点列表infected
		var infected = make(map[string]int)
		infected[nodeList.localNode.Addr+":"+strconv.Itoa(nodeList.localNode.Port)] = 1

		//更新本地节点信息
		nodeList.Set(nodeList.localNode)

		//设置心跳数据包
		p := packet{
			Node:     nodeList.localNode,
			Infected: infected,
		}

		//广播心跳数据包
		broadcast(nodeList, p)

		//向集群中的某个节点发起数据交换请求
		swapRequest(nodeList)

		nodeList.println("[Listen]:", nodeList.ListenAddr+":"+strconv.Itoa(nodeList.localNode.Port), "/ [Node list]:", nodeList.Get(), "/ [Metadata]:", nodeList.Read())

		//间隔时间
		time.Sleep(time.Duration(nodeList.Cycle) * time.Second)
	}
}

//监听其他节点发来的同步信息
func listener(nodeList *NodeList, mq chan []byte) {
	//监听协程
	listen(nodeList.ListenAddr, nodeList.localNode.Port, nodeList.Size, mq)
}

//消费信息
func consume(nodeList *NodeList, mq chan []byte) {
	for {
		//从监听队列中取出消息
		bs := <-mq
		var p packet
		err := json.Unmarshal(bs, &p)
		//如果数据解析错误
		if err != nil {
			println("[error]:", err)
			//跳过
			continue
		}

		//从节点心跳数据包中取出节点信息
		node := p.Node

		//更新本地列表
		nodeList.Set(node)

		//如果该数据包是元数据更新数据包，且数据包中的元数据版本要比本地存储的元数据版本新
		if p.IsUpdate && p.Metadata.Update > nodeList.metadata.Load().(metadata).Update {
			//更新本地节点中存储的元数据信息
			nodeList.metadata.Store(p.Metadata)
		}

		//如果该数据包是两节点间的元数据交换数据包
		if p.IsSwap != 0 {
			//如果数据包中的元数据版本要比本地存储的元数据版本新
			if p.Metadata.Update > nodeList.metadata.Load().(metadata).Update {
				//更新本地节点中存储的元数据信息
				nodeList.metadata.Store(p.Metadata)
			}
			//如果是接收方接收到发起方的数据交换请求
			if p.IsSwap == 1 {
				//回应发送方，完成交换流程
				swapResponse(nodeList, p.Node)
			}
			//跳过，不广播
			continue
		}

		//广播推送该节点信息
		broadcast(nodeList, p)
	}
}

//广播推送信息
func broadcast(nodeList *NodeList, p packet) {

	//取出所有未过期的节点
	nodes := nodeList.Get()
	var packets []packet

	//选取部分未被传染的节点
	i := 0
	for _, v := range nodes {

		//如果超过Amount最大推送数量
		if i >= nodeList.Amount {
			//结束广播
			break
		}

		//如果该节点已经被传染过了
		if p.Infected[v.Addr+":"+strconv.Itoa(v.Port)] == 1 {
			//跳过该节点
			continue
		}

		p.Infected[v.Addr+":"+strconv.Itoa(v.Port)] = 1 //标记该节点为已传染状态
		p.TargetAddr = v.Addr                           //设置发送目标地址
		p.TargetPort = v.Port                           //设置发送目标端口

		//将该节点添加进广播列表
		packets = append(packets, p)
		i++
	}

	//向这些未被传染的节点广播传染数据
	for _, v := range packets {
		bs, err := json.Marshal(p)
		if err != nil {
			println("[error]:", err)
		}
		write(v.TargetAddr, v.TargetPort, bs)
	}
}

//发起两节点数据交换请求
func swapRequest(nodeList *NodeList) {

	//设置为数据交换数据包
	p := packet{
		//将本地节点信息存入数据包，接收方根据这个信息回复请求
		Node:     nodeList.localNode,
		Infected: make(map[string]int),
		IsSwap:   1,
		Metadata: nodeList.metadata.Load().(metadata),
	}

	//取出所有未过期的节点
	nodes := nodeList.Get()

	bs, err := json.Marshal(p)
	if err != nil {
		println("[error]:", err)
	}

	//在节点列表中随机选取一个节点，发起数据交换请求
	for i := 0; i < len(nodes); i++ {
		//如果遍历到自己，则跳过
		if nodes[i].Addr == nodeList.localNode.Addr && nodes[i].Port == nodeList.localNode.Port {
			continue
		}
		write(nodes[i].Addr, nodes[i].Port, bs)
		nodeList.println("[Swap Request]:", nodeList.localNode.Addr+":"+strconv.Itoa(nodeList.localNode.Port), "->", nodes[i].Addr+":"+strconv.Itoa(nodes[i].Port))
	}
}

//接收数据交换请求并回应发送方，完成交换工作
func swapResponse(nodeList *NodeList, node Node) {

	//设置为数据交换数据包
	p := packet{
		Node:     nodeList.localNode,
		Infected: make(map[string]int),
		IsSwap:   2,
		Metadata: nodeList.metadata.Load().(metadata),
	}

	bs, err := json.Marshal(p)
	if err != nil {
		println("[error]:", err)
	}

	//回应发起节点
	write(node.Addr, node.Port, bs)

	nodeList.println("[Swap Response]:", node.Addr+":"+strconv.Itoa(node.Port), "<-", nodeList.localNode.Addr+":"+strconv.Itoa(nodeList.localNode.Port))
}
