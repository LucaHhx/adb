package main

import (
	"fmt"

	"github.com/luca/adb/adb"
	"github.com/luca/adb/adb/uixml"
)

func main() {
	dev := adb.NewDevice()
	//err := dev.ClickButton("Login")
	//if err != nil {
	//	return
	//}
	nodes, err := dev.FindNodes(func(n, pn uixml.Node) bool {
		return n.Class == "android.widget.Button"
	})
	if err != nil {
		return
	}
	fmt.Println(nodes)
	//node, err := dev.FindNode(func(n, pn uixml.Node) bool {
	//	if n.NAF == "true" && pn.ContentDesc == "Password" {
	//		return true
	//	}
	//	return false
	//})
	//if err != nil {
	//	return
	//}
	//dev.ClickNode(node)
	//fmt.Println(node)
}
