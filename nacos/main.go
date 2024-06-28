package main

import (
	"encoding/json"
	"fmt"
	"time"
)

func main() {
	NaCosConfig = NaCosConfigStruct{
		IpAddr:      "127.0.0.1",
		Port:        8848,
		NamespaceId: "9db1dfc0-4ed2-4aa5-a5b5-67af919b971a",
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
