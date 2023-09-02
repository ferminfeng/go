package main

import (
	"fmt"
	"github.com/Shopify/sarama"
)

func main() {
	consumer()
}

// 获取消费者
func consumer() {
	fmt.Println("kafka:消费者")

	// 获取消费者接口，当为外网主机时修改localhost为主机IP地址
	consumer, err := sarama.NewConsumer([]string{"localhost:9093"}, sarama.NewConfig())
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		// 关闭消费者
		if err = consumer.Close(); err != nil {
			fmt.Println(err)
			return
		}
	}()

	// 获取消费者的分片接口，sarama.OffsetNewest 标识获取新的消息
	partitionConsumer, err := consumer.ConsumePartition("topic", 0, sarama.OffsetNewest)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer func() {
		if err = partitionConsumer.Close(); err != nil {
			fmt.Println(err)
			return
		}
	}()

	for msg := range partitionConsumer.Messages() {
		fmt.Printf("分片:%d 偏移:%d Key:%s Value:%s\n", msg.Partition, msg.Offset, msg.Key, msg.Value)
	}

}
