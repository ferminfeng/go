package main

import (
	"encoding/json"
	"fmt"
	"time"
)

func main() {
	// NaCosConfig = NaCosConfigStruct{
	// 	IpAddr:      "127.0.0.1",
	// 	Port:        8848,
	// 	NamespaceId: "9db1dfc0-4ed2-4aa5-a5b5-67af919b971a",
	// 	TimeoutMs:   5000,
	// 	LogDir:      "./tmp/log",
	// 	CacheDir:    "./tmp/cache",
	// 	LogLevel:    "debug",
	// }
	NaCosConfig = NaCosConfigStruct{
		IpAddr:      "http://hw-test-nacos.aidyd.com",
		Port:        8848,
		NamespaceId: "a3f10455-9f7a-4e6f-bd14-95bc1aa23587",
		TimeoutMs:   5000,
		LogDir:      "./tmp/log",
		CacheDir:    "./tmp/cache",
		LogLevel:    "debug",
	}
	InitNaCosClient()

	kafkaConfig := &Kafka{}
	GetConfig("kafka.yml", "DEFAULT_GROUP", kafkaConfig)

	redisConfig := &Redis{}
	GetConfig("redis.yml", "DEFAULT_GROUP", redisConfig)

	for {
		fmt.Println("kafkaConfig Test的值为：", kafkaConfig.Test)
		time.Sleep(5 * time.Second)
	}

	ch := make(chan int)
	ch <- 1
}

func StructToJsonString(data any) string {
	jsonData, _ := json.Marshal(data)
	return string(jsonData)
}
