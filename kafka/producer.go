package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"log"
	"time"
)

// 获取生产者
func producer() {
	log.Printf("kafka: 生产者")

	// 获取生产者接口，当为外网主机时修改localhost为主机IP地址
	producer, err := sarama.NewSyncProducer(KafkaAddr, nil)
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

	for i := 0; i < 100; i++ {

		// 定义需要发送的消息
		msg := &sarama.ProducerMessage{
			Topic:     TestTopicName,
			Key:       sarama.StringEncoder("test_kafka_key"),
			Value:     sarama.StringEncoder(fmt.Sprintf("这里是kafka消息-%v", i)),
			Timestamp: time.Now(),
		}

		log.Printf("生产者-待发送的消息：%s", string(Marshal(msg)))

		// 发送消息，并获取该消息的分片、偏移量
		partition, offset, _ := producer.SendMessage(msg)
		if err != nil {
			log.Fatal(err)
			return
		}

		log.Printf("partition:%d offset:%d\n", partition, offset)
	}
}

func producerNew() {
	log.Printf("kafka: 生产者 new")

	config := sarama.NewConfig()
	// 设置
	// ack应答机制
	config.Producer.RequiredAcks = sarama.WaitForAll

	// 发送分区
	config.Producer.Partitioner = sarama.NewRandomPartitioner

	// 回复确认
	config.Producer.Return.Successes = true

	// 构造一个消息
	msg := &sarama.ProducerMessage{}
	msg.Topic = TestTopicName
	msg.Value = sarama.StringEncoder("test:weatherStation device")

	// 连接kafka
	client, err := sarama.NewSyncProducer(KafkaAddr, config)
	if err != nil {
		log.Printf("producer closed,err:%v", err)
	}
	defer client.Close()

	// 发送消息
	pid, offset, err := client.SendMessage(msg)
	if err != nil {
		log.Printf("send msg failed,err:%v", err)
		return
	}

	log.Printf("pid:%v offset: %v", pid, offset)

}
