package pekonode

import (
	"fmt"
	"time"
)

//打印信息到控制台
func (nodeList *NodeList) println(a ...interface{}) {
	if nodeList.IsPrint {
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"), a)
	}
}
