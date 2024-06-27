package main

import (
	"encoding/json"
)

func main() {
	InitNaCosClient()

	kafkaConfig := &Kafka{}
	GetConfig("kafka.yml", "DEFAULT_GROUP", kafkaConfig)

	redisConfig := &Redis{}
	GetConfig("redis.yml", "DEFAULT_GROUP", redisConfig)

	ch := make(chan int)
	ch <- 1
}

func StructToJsonString(data any) string {
	jsonData, _ := json.Marshal(data)
	return string(jsonData)
}
