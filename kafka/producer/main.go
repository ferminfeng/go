package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"log"
)

func main() {
	producer()
}

// 获取生产者
func producer() {
	fmt.Println("kafka:生产者")

	// 获取生产者接口，当为外网主机时修改localhost为主机IP地址
	producer, err := sarama.NewSyncProducer([]string{"localhost:9093"}, nil)
	if err != nil {
		log.Fatal(err)
		return
	}

	defer func() {
		// 关闭生产者
		if err = producer.Close(); err != nil {
			log.Fatal(err)
			return
		}
	}()

	// 定义需要发送的消息
	msg := &sarama.ProducerMessage{
		Topic: "topic",
		Key:   sarama.StringEncoder("key"),
		Value: sarama.StringEncoder("value"),
	}

	// 发送消息，并获取该消息的分片、偏移量
	partition, offset, _ := producer.SendMessage(msg)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Printf("partition:%d offset:%d\n", partition, offset)
}
