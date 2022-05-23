package test

import (
	"fmt"
	"github.com/dpwgc/pekonode"
	"testing"
	"time"
)

//测试用例-启动四个节点，构成一个Gossip集群
func TestFourNode(t *testing.T) {

	//先启动节点A（初始节点）
	nodeA()
	//启动节点B
	nodeB()
	//启动节点C
	nodeC()
	//启动节点D
	nodeD()

	for {
		time.Sleep(10 * time.Second)
	}
}

//运行节点A（初始节点）
func nodeA() {
	//配置节点A的本地节点列表nodeList参数
	nodeList := pekonode.NodeList{
		SecretKey: "test_key",
		IsPrint:   true, //是否在控制台输出日志信息
	}

	//创建节点A及其本地节点列表
	nodeList.New(pekonode.Node{
		Addr: "0.0.0.0",  //本地节点IP地址，公网环境下请填写公网IP
		Port: 8000,       //本地节点端口号
		Name: "A-server", //节点名称，自定义填写
		Tag:  "A",        //节点标签，自定义填写，可以填一些节点基本信息
	})

	//因为是第一个启动的节点，所以不需要用Set函数添加其他节点

	//本地节点加入Gossip集群，本地节点列表与集群中的各个节点所存储的节点列表进行数据同步
	nodeList.Join()

	//延迟3秒
	time.Sleep(3 * time.Second)
}

//运行节点B
func nodeB() {
	//配置节点B的本地节点列表nodeList参数
	nodeList := pekonode.NodeList{
		SecretKey: "test_key",
		IsPrint:   true,
	}

	//创建节点B及其本地节点列表
	nodeList.New(pekonode.Node{
		Addr: "0.0.0.0",
		Port: 8001,
		Name: "B-server",
		Tag:  "B",
	})

	//将初始节点A的信息加入到B节点的本地节点列表当中
	nodeList.Set(pekonode.Node{
		Addr: "0.0.0.0",
		Port: 8000,
		Name: "A-server",
		Tag:  "A", //将节点A的信息添加进节点B的本地节点列表
	})

	//调用Join后，节点B会自动与节点A进行数据同步
	nodeList.Join()

	//延迟10秒
	time.Sleep(10 * time.Second)
}

//运行节点C
func nodeC() {
	nodeList := pekonode.NodeList{
		SecretKey: "test_key",
		IsPrint:   true,
	}

	//创建节点C及其本地节点列表
	nodeList.New(pekonode.Node{
		Addr: "0.0.0.0",
		Port: 8002,
		Name: "C-server",
		Tag:  "C",
	})

	//也可以在加入集群之前，在本地节点列表中添加多个节点信息
	nodeList.Set(pekonode.Node{
		Addr: "0.0.0.0",
		Port: 8000,
		Name: "A-server",
		Tag:  "A", //将节点A的信息添加进节点C的本地节点列表
	})
	nodeList.Set(pekonode.Node{
		Addr: "0.0.0.0",
		Port: 8001,
		Name: "B-server",
		Tag:  "B", //将节点B的信息添加进节点C的本地节点列表
	})

	//在加入集群后，节点C将会与上面的节点A及节点B进行数据同步
	nodeList.Join()

	//延迟10秒
	time.Sleep(10 * time.Second)

	//获取本地节点列表
	list := nodeList.Get()
	fmt.Println("Node list::", list) //打印节点列表

	//在集群中发布新的元数据信息
	nodeList.Publish([]byte("test metadata"))

	//读取本地元数据信息
	metadata := nodeList.Read()
	fmt.Println("Metadata:", string(metadata)) //打印元数据信息

	//停止节点C的心跳广播服务（节点C暂时下线）
	nodeList.Stop()

	//延迟30秒
	time.Sleep(10 * time.Second)

	//因为之前节点C下线，C的本地节点列表无法接收到各节点的心跳数据包，列表被清空
	//所以要先往C的本地节点列表中添加一些集群节点，再调用Start()重启节点D的同步工作
	nodeList.Set(pekonode.Node{
		Addr: "0.0.0.0",
		Port: 8001,
		Name: "B-server",
		Tag:  "B", //这里添加节点B的信息
	})

	//重启节点C的心跳广播服务（节点C重新上线）
	nodeList.Start()
}

//运行节点D
func nodeD() {
	//配置节点D的本地节点列表nodeList参数
	nodeList := pekonode.NodeList{
		SecretKey: "test_key",
		IsPrint:   true,
	}

	//创建节点D及其本地节点列表
	nodeList.New(pekonode.Node{
		Addr: "0.0.0.0",
		Port: 8003,
		Name: "D-server",
		Tag:  "D",
	})

	//将初始节点A的信息加入到D节点的本地节点列表当中
	nodeList.Set(pekonode.Node{
		Addr: "0.0.0.0",
		Port: 8000,
		Name: "A-server",
		Tag:  "A", //将节点A的信息添加进节点D的本地节点列表
	})

	//调用Join后，节点D会自动与节点A进行数据同步
	nodeList.Join()
}
