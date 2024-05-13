// @Author fermin 2024/4/25 17:22:00
package main

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/pkg/errors"
	"os"
	"os/signal"
	"sync"
	"time"
)

func MoreConsumer() {
	// 创建Kafka配置
	config := sarama.NewConfig()
	// config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategySticky

	config.Version = sarama.V2_0_0_0
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Return.Errors = false

	// 定义topic
	topic := []string{TestTopicName}

	// 定义消费者组的名称
	consumerGroup := "your_consumer_group"

	// 创建Kafka消费者组
	group, err := sarama.NewConsumerGroup(KafkaAddr, consumerGroup, config)
	if err != nil {
		panic(err)
	}

	// 定义消费者组成员的数量
	numConsumers := 3

	// 使用WaitGroup等待所有消费者goroutine完成
	var wg sync.WaitGroup
	wg.Add(numConsumers)

	// 启动多个消费者goroutine
	for i := 0; i < numConsumers; i++ {
		go func(i int) {
			fmt.Println("启动多个消费者：", i)
			defer wg.Done()
			for {
				// 每个goroutine都处理分配给它的分区
				err := group.Consume(context.Background(), topic, &consumerGroupHandler{I: i})
				if err != nil {
					panic(errors.Wrap(err, "Error consuming from Kafka"))
				}
			}
		}(i)
	}

	// 捕获中断信号以进行优雅关闭
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	select {
	case <-interrupt:
		fmt.Println("Received interrupt, shutting down...")
		if err := group.Close(); err != nil {
			panic(err)
		}
	}

	// 等待所有消费者goroutine完成
	wg.Wait()
}

// consumerGroupHandler 实现了sarama.ConsumerGroupHandler接口，用于处理从Kafka接收到的消息
type consumerGroupHandler struct {
	I int
}

// Setup 在分配给成员的分区发生变化时调用
func (h *consumerGroupHandler) Setup(sess sarama.ConsumerGroupSession) error {
	fmt.Println(fmt.Sprintf("Consumer group member is set up, i: %v, memberID: %v", h.I, sess.MemberID()))
	return nil
}

// Cleanup 在成员停止消费分区前调用
func (h *consumerGroupHandler) Cleanup(sess sarama.ConsumerGroupSession) error {
	fmt.Println(fmt.Sprintf("Consumer group member is cleaned up, i: %v, memberID: %v", h.I, sess.MemberID()))
	return nil
}

// ConsumeClaim 当分区被分配给该成员时调用
func (h *consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		// fmt.Printf("Message claimed: memberID=%s, topic=%s, partition=%d, offset=%d, key=%s, value=%s\n", sess.MemberID(), msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)

		fmt.Printf("memberID=%s, offset=%d, key=%s, value=%s\n", sess.MemberID(), msg.Offset, msg.Key, msg.Value)
		time.Sleep(2 * time.Second)
		sess.MarkMessage(msg, "")
	}
	return nil
}
