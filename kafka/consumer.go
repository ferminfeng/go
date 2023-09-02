package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"log"
	"sync"
)

// 获取消费者
func consumer() {
	log.Printf("kafka: 消费者")

	// 获取消费者接口，当为外网主机时修改localhost为主机IP地址
	consumer, err := sarama.NewConsumer(KafkaAddr, sarama.NewConfig())
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
	partitionConsumer, err := consumer.ConsumePartition(TestTopicName, 0, sarama.OffsetNewest)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer func() {
		if err = partitionConsumer.Close(); err != nil {
			log.Fatal(err)
			return
		}
	}()

	for msg := range partitionConsumer.Messages() {
		log.Printf("\n分片:%d \n偏移:%d \nKey:%s \nValue:%s \nTime:%s \nMQData:%s \n", msg.Partition, msg.Offset, msg.Key, msg.Value, msg.Timestamp, string(Marshal(msg)))
	}

}

var wg sync.WaitGroup

func consumerNew() {
	log.Printf("kafka: 消费者 new")

	// 创建新的消费者
	consumer, err := sarama.NewConsumer(KafkaAddr, nil)
	if err != nil {
		log.Printf("fail to start consumer,%v", err)
		return
	}

	// 根据topic获取所有的分区列表
	partitionList, err := consumer.Partitions(TestTopicName)
	if err != nil {
		log.Printf("fail to get list of partition,err:%v", err)
		return
	}

	log.Printf("分区列表: %v", partitionList)

	// 遍历所有的分区
	for p := range partitionList {
		// 针对每一个分区创建一个对应分区的消费者
		pc, err := consumer.ConsumePartition(TestTopicName, int32(p), sarama.OffsetNewest)
		if err != nil {
			log.Printf("failed to start consumer for partition %d,err:%v\n", p, err)
			return
		}

		defer func() {
			pc.AsyncClose()
		}()

		wg.Add(1)

		// 异步从每个分区消费信息
		go func(sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				log.Printf("\npartition: %d \nOffse: %d \nKey: %v \nValue: %s  \nMQData: %s \n\n",
					msg.Partition, msg.Offset, msg.Key, msg.Value, string(Marshal(msg)))
			}
		}(pc)
	}

	wg.Wait()
}
