// Package main @Author fermin 2024/5/13 18:32:00
package main

import (
	"fmt"
	// "github.com/samuel/go-zookeeper/zk"
	"github.com/go-zookeeper/zk"
	"time"
)

func main() {
	conn, _, err := zk.Connect([]string{"127.0.0.1:2181"}, 5*time.Second)
	if err != nil {
		panic(err)
	}

	// 1.验证根是否存在
	if hasRoot, _, _ := conn.Exists("/root1"); !hasRoot {
		// 2.新增根
		_, err = conn.Create("/root1", []byte("root_content"), 0, zk.WorldACL(zk.PermAll))
		if err != nil {
			fmt.Println("failed add root, info: ", err.Error())
		}
		fmt.Println("add /root1 success")
	}

	// 3.查询根
	data, stat, err := conn.Get("/root1")
	if err != nil {
		fmt.Println("failed get root, info: ", err.Error())
	}
	fmt.Println("text: ", string(data), stat.Version)

	// 4.修改根
	if _, err = conn.Set("/root1", []byte("update text"), stat.Version); err != nil {
		fmt.Println("failed update root")
	}

	// 5.设置子节点(必须要有根/父节点)
	if _, err = conn.Create("/root1/subnode", []byte("node_text"), 0, zk.WorldACL(zk.PermAll)); err != nil {
		fmt.Println("failed add subnode, info: ", err.Error())
	}
	// 6.获取子节点列表
	childNodes, _, err := conn.Children("/root1")
	if err != nil {
		fmt.Println("failed get node list, info: ", err.Error())
	} else {
		fmt.Println("node list: ", childNodes)
	}

	// 6.删除根(必须先查后删, 删完子才能删父节点)
	_, stat, _ = conn.Get("/root1/subnode")
	if err := conn.Delete("/root1/subnode", stat.Version); err != nil {
		fmt.Println("falied delete node, info: ", stat.Version, err.Error())
	}
	_, stat, _ = conn.Get("/root1")
	if err := conn.Delete("/root1", stat.Version); err != nil {
		fmt.Println("falied delete root, info: ", stat.Version, err.Error())
	}
}
