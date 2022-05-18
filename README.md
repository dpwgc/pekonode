# PekoNode

## 一个基于Golang整合UDP实现的Gossip协议工具包

***

### 实现原理
![avatar](./img/gossip.jpg)
* 每个节点都有一个本地节点列表NodeList。
* 每个节点的后台同步协程定时将节点信息封装成心跳数据包，并广播给集群中的3个节点（可设置数量）。
* 其他节点接收到心跳数据包后，更新自己的本地节点列表NodeList，并将该心跳数据包广播给集群中的3个节点。
* 心跳数据包中存在Num计数器字段，每广播一次，Num+1，当广播次数大于3时，终止广播。
* 如果本地节点列表NodeList中，存在超时未发送心跳更新的节点，则删除该超时节点数据。

***

### 导入包
```
go get github.com/dpwgc/pekonode
```

***

### 简单示例
* Gossip集群的第一个服务节点，8000节点启动。
```
package main

import (
	"github.com/dpwgc/pekonode/gossip"
	"time"
)

func main() {

	//设置本地节点信息
	var node gossip.Node
	node.Addr = "0.0.0.0"   //节点的IP地址，公网环境下请填公网IP
	node.Port = 8000        //节点的服务端口号
	node.Tag = "test"       //节点标签，自定义填写，可以是服务名或者是一些配置信息

	//创建本地节点列表，同时将本地节点信息加入节点列表，设置本地UDP监听服务地址及其他参数
	nodeList := gossip.New(node, "0.0.0.0", 8000, 3, 6, 10, 30, true)
	//参数说明：
	//node 就是本地节点信息
	//"0.0.0.0"和8000 则是本地UDP服务监听端口（一般与上面的node地址信息相同），通过该端口来接收其他节点的同步信息
	//3 表示每次广播最多只给3个节点发送信息
	//6 表示每隔6秒向集群发送同步信息（即心跳信息）
	//10 代表UDP监听缓冲区，节点数量较多时可适当调大该值
	//30 表示一个节点的超时下线时间，如果本地节点列表中的某一个节点信息在30秒内没有收到心跳更新，则删除该节点信息
	//true 表示允许将服务运行信息输出到控制台

	//将本地节点列表加入到Gossip集群中（后台启动Gossip同步协程，与其他服务节点的节点列表进行数据同步）
	nodeList.Join()
	
	//延迟60秒
	time.Sleep(60 * time.Second)

	//60秒后停止本地节点列表与集群其他节点列表的信息同步工作（即节点下线）
	nodeList.Stop()
	
	for {
		time.Sleep(1 * time.Second)
	}
}
```
* Gossip集群的第二个服务节点，节点8001启动，并将8001节点与初始节点8000相连。
```
package main

import (
	"github.com/dpwgc/pekonode/gossip"
	"time"
)

func main() {

	//设置本地节点信息
	var node gossip.Node
	node.Addr = "0.0.0.0"   //公网环境下请填公网IP
	node.Port = 8001
	node.Tag = "test"

	//创建本地节点列表，同时将本地节点信息加入节点列表，设置本地UDP监听服务地址及其他参数
	nodeList := gossip.New(node, "0.0.0.0", 8001, 3, 6, 10, 30, true)

	//将一个新的节点信息加入到本地节点列表（Gossip同步协程启动后，将向这些节点发送同步信息）
	nodeList.Set(gossip.Node{
		Addr: "0.0.0.0",
		Port: 8000, //这里将第一个启动的8000节点信息添加进8001节点的本地节点列表中
		Tag:  "test",
	})
	//可以选择调用Set函数加入多个节点信息
	//......

	//将本地节点列表加入到Gossip集群中（后台启动Gossip同步协程，与其他服务节点（例如8000节点）的节点列表进行数据同步）
	nodeList.Join()
	
	for {
		time.Sleep(1 * time.Second)
	}
}
```

***

### 函数说明
##### New 创建本地节点列表
```
func New(node Node, listenAddr string, listenPort int, amount int, cycle int64, buffer int, timeout int64, isPrint bool) NodeList 
```
* node 节点信息
* listenAddr 本地UDP监听服务地址
* listenPort 本地UDP监听服务端口
* amount 单次广播传染的最大节点数量
* cycle 节点心跳同步周期（每cycle秒向集群的某些节点传染心跳包，越小越精准，不得大于timeout）
* buffer UDP监听缓冲区大小
* timeout 节点超时下线时间（单位：秒，越小越精准，不得小于cycle）
* isPrint 是否在控制台输出信息

##### Join 加入集群
```
func (nodeList *NodeList) Join() 
```

##### Stop 停止广播心跳
```
func (nodeList *NodeList) Stop() 
```

##### Start 重新开始广播心跳
```
func (nodeList *NodeList) Start() {
```

##### Set 向本地节点列表中加入其他节点信息
```
func (nodeList *NodeList) Set(node Node) 
```

##### Get 获取本地节点列表
```
func (nodeList *NodeList) Get() []Node 
```