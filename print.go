package pekonode

import (
	"fmt"
	"log"
	"time"
)

//打印信息到控制台
func (nodeList *NodeList) println(a ...interface{}) {
	//输出错误信息
	if a[0] == "[Error]:" {
		log.Println(a[1])
	}
	if nodeList.IsPrint {
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"), a)
	}
}
