package main

import (
	"flag"
	jsoniter "github.com/json-iterator/go"
	"log"
	"time"
)

var startType string

var KafkaAddr = []string{"127.0.0.1:9092"}
var TestTopicName = "test_topic"

func main() {
	flag.StringVar(&startType, "s", "", "启动生产者还是消费者")

	flag.Parse()

	switch startType {
	case "producer":
		producer()
	case "producerNew":
		producerNew()
	case "consumer":
		consumer()
	case "consumerNew":
		consumerNew()
	case "moreConsumer":
		MoreConsumer()
	default:
		log.Fatal("请指定启动生产者还是消费者，startType:", startType)
	}

	return
}

// Stamp2Str 时间戳 -> 字符串  1660188363 -> 2022-08-11 11:26:03
func Stamp2Str(stamp int32) string {
	timeLayout := "2006-01-02 15:04:05"
	str := time.Unix(int64(stamp), 0).Format(timeLayout)
	return str
}

func Marshal(value interface{}) []byte {
	data, err := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(value)
	if err != nil {
		log.Printf("Marshal json:%v", err.Error())
	}
	return data
}

func Unmarshal(data []byte, value interface{}) error {
	err := jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(data, value)
	if err != nil {
		log.Printf("Unmarshal json:%v", err.Error())
	}
	return err
}
