package gossip

import (
	"sync"
	"sync/atomic"
)

// NodeList 节点列表
type NodeList struct {
	nodes   sync.Map //节点集合（key为Node结构体，value为节点最近更新的秒级时间戳）
	Amount  int      //每次给多少个节点发送同步信息
	Cycle   int64    //同步时间周期（每隔多少秒向其他节点发送一次列表同步信息）
	Buffer  int      //UDP接收缓冲区大小（决定UDP监听服务可以异步处理多少个请求）
	Size    int      //单个UDP心跳数据包的最大容量（单位：字节）
	Timeout int64    //单个节点的过期删除界限（多少秒后删除）

	localNode Node //本地节点信息

	ListenAddr string //本地UDP监听地址，用这个监听地址接收其他节点发来的心跳包（一般填0.0.0.0即可）

	status atomic.Value //本地节点列表更新状态（true：正常运行，false：停止发布心跳）

	IsPrint bool //是否打印列表同步信息到控制台

	metadata atomic.Value //元数据，集群中各个节点的元数据内容一致，相当于集群的公共数据（可存储一些公共配置信息），可以通过广播更新各个节点的元数据内容
}

// Node 节点
type Node struct {
	Addr string //节点IP地址（公网环境下填公网IP）
	Port int    //端口号
	Name string //节点名称（自定义）
	Tag  string //节点标签（自定义，可以写一些基本信息）
}

// 数据包
type packet struct {

	//节点信息
	Node       Node           //心跳数据包中的节点信息
	TargetAddr string         //发送目标的IP地址
	TargetPort int            //发送目标的端口号
	Infected   map[string]int //已被该数据包传染的节点列表，key为Addr:Port拼接的字符串，value为判定该节点是否已被传染的参数（1：是，0：否）

	//元数据信息
	Metadata string //新的元数据信息，如果该数据包是元数据更新数据包（isUpdate=true），则用newData覆盖掉原先的集群元数据metadata
	IsUpdate bool   //判定该数据包是否为元数据更新数据包（true：是，false：否）
}
