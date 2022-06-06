# PekoNode

## 一个基于Golang整合UDP/TCP实现的Gossip协议工具包

![MIT](https://img.shields.io/static/v1?label=LICENSE&message=MIT&color=red)
![Go](https://img.shields.io/static/v1?label=Go&message=v1.17&color=blue)
[![github](https://img.shields.io/static/v1?label=Github&message=pekonode&color=blue)](https://github.com/dpwgc/pekonode)

#### [English document: README.md](./README.md)
***
### 实现功能
##### 集群节点列表共享
* 通过谣言传播方式同步集群节点列表 `NodeList`（每个节点最终都会存储完整的节点列表，可用于服务注册发现场景）
##### 集群元数据信息共享
* 通过谣言传播方式发布集群元数据信息 `Metadata`（集群公共数据，每个节点的本地元数据信息最终一致，可以自定义存放内容，例如存放一些公共配置信息，充当配置中心），通过反熵传播方式实现集群各个节点的元数据校验同步功能。
##### 自定义配置
* 节点列表 `NodeList` 提供一系列参数供用户自定义配置，用户可以使用默认缺省参数，也可根据需求自行填写参数。
***

### 实现原理
##### `NodeList` 节点列表信息同步
* 每个节点都有一个本地节点列表NodeList。
* 每个节点的后台同步协程定时将节点信息封装成心跳数据包，并广播给集群中部分未被传染的节点。
* 其他节点接收到心跳数据包后，更新自己的本地节点列表NodeList，再将该心跳数据包广播给集群中部分未被传染的节点。
* 重复前一个广播步骤（谣言传播方式），直至所有节点都被传染，本次心跳传染结束。
* 如果本地节点列表NodeList中，存在超时未发送心跳更新的节点，则删除该超时节点数据。

##### `Metadata` 元数据信息同步
* 在某一节点调用Publish()函数发布新的元数据后，新数据会扩散到各个节点，然后覆盖他们的本地元数据信息。
* 每个节点都会定期选取一个随机的节点进行元数据交换检查操作，如发现某个节点上的元数据是旧的，则将其覆盖（反熵传播方式）。
* 当有新节点加入到集群时，该节点会通过数据交换功能获取到最新的集群元数据信息。

***

### 导入包
```
//Goland终端输入
go get github.com/dpwgc/pekonode
```
```
//程序中引入
import "github.com/dpwgc/pekonode"
```

***
### 使用方法
* 配置本地节点列表`nodeList`参数
```
nodeList := pekonode.NodeList{}
```

* 初始化本地节点列表，传入本地节点信息
```
nodeList.New(pekonode.Node{
	Addr: "0.0.0.0",  //本地节点IP地址，公网环境下请填写公网IP
	Port: 8000,       //本地节点端口号
})
```
* 往本地节点列表中添加新的节点信息
```
nodeList.Set(pekonode.Node{
	Addr: "0.0.0.0",
	Port: 9999,
})
```
* 将该节点加入Gossip集群（在后台启动心跳广播与监听协程）
```
nodeList.Join()
```
* 获取本地节点列表
```
list := nodeList.Get()
fmt.Println(list)
```
* 节点停止发布心跳
```
nodeList.Stop()
```
* 节点重新开始发布心跳
```
nodeList.Start()
```
* 在集群中发布新的元数据信息
```
nodeList.Publish([]byte("test metadata"))
```
* 获取本地元数据信息
```
metadata := nodeList.Read()
fmt.Println(string(metadata))
```
***

### 简单使用示例
#### 启动一个节点
```
package main

import (
	"github.com/dpwgc/pekonode"
	"time"
)

//简单示例，启动一个节点
func main()  {

	//配置该节点的本地节点列表nodeList参数
	nodeList := pekonode.NodeList{
		IsPrint:    true,           //是否在控制台输出日志信息，不填默认为false
	}

	//创建本地节点列表，传入本地节点信息
	nodeList.New(pekonode.Node{
		Addr: "0.0.0.0",		//本地节点IP地址，公网环境下请填写公网IP
		Port: 8000,				//本地节点端口号
		Name: "Test",           //节点名称，自定义填写
		Tag: "Test",			//节点标签，自定义填写，可以填一些节点基本信息
	})

	//往本地节点列表中添加新的节点信息，可添加多个节点，本地节点将会与这些新节点同步信息
	//如果启动的是集群中的第一个节点，可不进行Set添加操作
	nodeList.Set(pekonode.Node{
		Addr: "0.0.0.0",
		Port: 9999,
		Name: "Hello",
		Tag: "Hello",
	})
	nodeList.Set(pekonode.Node{
		Addr: "0.0.0.0",
		Port: 7777,
		Name: "Hi",
		Tag: "Hi",
	})

	//将该节点加入Gossip集群（在后台启动心跳广播与监听协程）
	nodeList.Join()
	
	//获取本地节点列表
	list := nodeList.Get()
	//打印节点列表
	fmt.Println(list)
	
	//在集群中发布新的元数据信息
	nodeList.Publish([]byte("test metadata"))
	
	//读取本地元数据信息
	metadata := nodeList.Read()
	//打印元数据信息
	fmt.Println(string(metadata))
	
	//因为心跳广播这些工作都是在后台协程进行的，所以在调用Join函数后不能让主协程关闭，否则程序将直接退出
	//无限循环
	for {
		time.Sleep(10*time.Second)
	}
}
```

***

### 完整使用示例
示例文件：`/test` 目录下

***
### 模板说明
```
// NodeList 节点列表
type NodeList struct {
	nodes   sync.Map        //节点集合（key为Node结构体，value为节点最近更新的秒级时间戳）
	Amount  int             //每次给多少个节点发送同步信息（扇出度）
	Cycle   int64           //同步时间周期（每隔多少秒向其他节点发送一次列表同步信息）
	Buffer  int             //UDP/TCP接收缓冲区大小（决定UDP/TCP监听服务可以异步处理多少个请求）
	Size    int             //单个UDP/TCP心跳数据包的最大容量，默认16k，如果需要同步较大的Metadata，请自行调大（单位：字节）
	Timeout int64           //单个节点的过期删除界限（多少秒后删除）
	SecretKey string        //集群密钥，同一集群内的各个节点密钥应该保持一致
	localNode Node          //本地节点信息
	Protocol string         //集群连接使用的网络协议，UDP或TCP，默认UDP
	ListenAddr string       //本地UDP/TCP监听地址，用这个监听地址接收其他节点发来的心跳包（一般填0.0.0.0即可）
	status atomic.Value     //本地节点列表更新状态（true：正常运行，false：停止发布心跳）
	IsPrint bool            //是否打印列表同步信息到控制台
	metadata atomic.Value   //元数据，集群中各个节点的元数据内容一致，相当于集群的公共数据（可存储一些公共配置信息），可以通过广播更新各个节点的元数据内容
}

// Node 节点
type Node struct {
	Addr string     //节点IP地址（公网环境下填公网IP）
	Port int        //端口号
	Name string     //节点名称（自定义）
	Tag  string     //节点标签（自定义，可以写一些基本信息）
}
```

*** 

### 控制台打印信息
#### 当NodeList的IsPrint参数被设为true时，程序会在控制台打印出运行信息
##### 当节点加入集群时，打印：
```
2022-05-19 14:51:23 [[Join]: {0.0.0.0 8000 A-server A}]
```
* 表示节点0.0.0.0:8000加入集群

##### 当节点发布心跳时，打印：
```
2022-05-19 14:52:23 [[Listen]: 0.0.0.0:8000 / [Node list]: [{0.0.0.0 8000 A-server A} {0.0.0.0 8001 B-server B}]]
```
* Listen表示本地UDP监听服务地址与端口，Node list表示当前本地节点列表。

##### 当暂停节点心跳广播时，打印：
```
2022-05-19 14:52:06 [[Stop]: {0.0.0.0 8002 C-server C}]
```
* 表示节点0.0.0.0:8002停止广播心跳数据包

##### 当重新开始节点心跳广播时，打印：
```
2022-05-19 14:52:36 [[Start]: {0.0.0.0 8002 C-server C}]
```
* 表示节点0.0.0.0:8002重新开始广播心跳数据包

##### 当两个节点之间进行元数据交换时，打印：
```
2022-05-20 13:12:26 [[Swap Request]: 0.0.0.0:8002 -> 0.0.0.0:8000]
2022-05-20 13:12:26 [[Swap Response]: 0.0.0.0:8002 <- 0.0.0.0:8000]
```
* 8002节点向8000节点发起数据交换请求
* 8000节点回应8002节点的交换请求，数据交换工作完成

***

### 项目结构
* pekonode
  * test `测试文件`
  * model.go `结构体模板`
  * opt.go `提供给外部的系列操作函数`
  * print.go `控制台输出`
  * sync.go `集群同步服务`
  * net.go `网络服务`
  * udp.go `UDP收发服务`
  * tcp.go `TCP收发服务`
  * md5.go `MD5密钥生成`

