# PekoNode

## 基于Golang整合UDP实现的Gossip协议工具包

### 简单示例
* Gossip集群的第一个服务节点，8000节点启动。
```
package main

import (
	"pekonode/server"
	"time"
)

func main() {

	//设置本地节点信息
	var node server.Node
	node.Addr = "0.0.0.0"   //公网环境下请填公网IP
	node.Port = 8000
	node.Tag = "test"

	//创建本地节点列表，同时将本地节点信息加入节点列表，设置本地UDP监听服务地址及其他参数
	nodeList := server.New(node, "0.0.0.0", 8000, 3, 10, 5, 30, true)

	//将本地节点列表加入到Gossip集群中（后台启动Gossip同步协程，与其他服务节点的节点列表进行数据同步）
	nodeList.Join()
	
	time.Sleep(10 * time.Second)

	//停止本地节点列表与集群其他节点列表的更新工作
	nodeList.Stop()
	
	for {
		time.Sleep(1 * time.Second)
	}
}
```
* Gossip集群的第二个服务节点，节点8001启动，并将8001节点与8000节点相连。
```
package main

import (
	"pekonode/server"
	"time"
)

func main() {

	//设置本地节点信息
	var node server.Node
	node.Addr = "0.0.0.0"   //公网环境下请填公网IP
	node.Port = 8001
	node.Tag = "test"

	//创建本地节点列表，同时将本地节点信息加入节点列表，设置本地UDP监听服务地址及其他参数
	nodeList := server.New(node, "0.0.0.0", 8001, 3, 10, 5, 30, true)

	//将一个新的节点信息加入到本地节点列表（Gossip同步协程启动后，将向这些节点发送同步信息）
	nodeList.Set(server.Node{
		Addr: "0.0.0.0",
		Port: 8000, //这里将第一个启动的8000节点信息添加进8001节点的本地节点列表中
		Tag:  "test",
	})

	//将本地节点列表加入到Gossip集群中（后台启动Gossip同步协程，与其他服务节点（例如8000节点）的节点列表进行数据同步）
	nodeList.Join()
	
	time.Sleep(10 * time.Second)

	//停止本地节点列表与集群其他节点列表的更新工作
	nodeList.Stop()
	
	for {
		time.Sleep(1 * time.Second)
	}
}
```

### 函数说明
##### New 创建本地节点列表
```
func New(node Node, listenAddr string, listenPort int, amount int, cycle int64, buffer int, timeout int64, isPrint bool) NodeList 
```
* node 节点信息
* listenAddr 本地UDP监听服务地址
* listenPort 本地UDP监听服务端口
* amount 单次广播传染的最大节点数量
* cycle 节点心跳同步周期（每cycle秒向集群的某些节点传染心跳包）
* buffer UDP监听缓冲区大小
* timeout 节点超时下线时间（单位：秒）
* isPrint 是否在控制台输出信息

##### Join 加入集群
```
func (nodeList *NodeList) Join() 
```

##### Stop 停止同步
```
func (nodeList *NodeList) Stop() 
```

##### Set 向本地节点列表中加入其他节点信息
```
func (nodeList *NodeList) Set(node Node) 
```

##### Get 获取本地节点列表
```
func (nodeList *NodeList) Get() []Node 
```
