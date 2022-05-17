package server

import "sync"

// NodeList 节点列表
type NodeList struct {
	Nodes   sync.Map //节点集合（key为Node结构体，value为节点最近更新的秒级时间戳）
	Amount  int      //每次给多少个节点发送同步信息
	Cycle   int64    //同步时间周期（每隔多少秒向其他节点发送一次列表同步信息）
	Buffer  int      //接收缓冲区大小
	Timeout int64    //单个节点的过期删除界限（多少秒后删除）

	localNode Node //本地服务节点

	ListenAddr string //本地节点列表更新监听地址（这两一般与本地节点Node设置相同）
	ListenPort int    //本地节点列表更新监听端口

	status map[int]bool //本地节点列表更新状态（map[1] = true：正常运行，map[1] = false：停止同步更新）

	isPrint bool //是否打印信息到控制台
}

// Node 节点
type Node struct {
	Addr string //节点IP地址（公网环境下填公网IP）
	Port int    //端口号
}
