package pekonode

import (
	"log"
)

//打印信息到控制台
func (nodeList *NodeList) println(a ...interface{}) {
	//输出错误信息
	if a[0] == "[Error]:" && !nodeList.IsPrint {
		log.Println(a)
	}
	if nodeList.IsPrint {
		log.Println(a)
	}
}
