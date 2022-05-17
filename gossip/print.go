package gossip

import "fmt"

//打印信息到控制台
func (nodeList *NodeList) println(a ...interface{}) {
	if nodeList.isPrint {
		fmt.Println(a)
	}
}
