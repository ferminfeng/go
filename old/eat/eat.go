package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	fmt.Println("今天吃啥呢？\n\n")
	time.Sleep(time.Second * 1)
	fmt.Println("你有如下几种选择：\n\n")
	foodList := []string{"黄焖鸡", "麻辣烫", "饺子+炸串", "咖喱鸡饭", "称菜", "炒菜"}
	for _, food := range foodList {
		fmt.Println(food)
		time.Sleep(time.Second * 1)
	}

	fmt.Println("\n\n思考一下\n\n")
	time.Sleep(time.Second * 1)
	fmt.Println("。。。")
	time.Sleep(time.Second * 2)
	fmt.Println("。。。。。。")
	time.Sleep(time.Second * 3)
	fmt.Println("。。。。。。。。。")

	rander := rand.New(rand.NewSource(time.Now().UnixNano()))
	random := rander.Intn(len(foodList) - 1)
	fmt.Println("\n\n结果是今天吃：", foodList[random])
}
