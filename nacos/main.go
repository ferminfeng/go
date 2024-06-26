package main

import (
	"encoding/json"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"gopkg.in/yaml.v2"
)

func main() {
	ch := make(chan int)
	ch <- 1
}

type Kafka struct {
	Username       string    `yaml:"username"`
	Password       string    `yaml:"password"`
	ProducerAck    string    `yaml:"producerAck"`
	CompressType   string    `yaml:"compressType"`
	Address        []string  `yaml:"address"`
	ToRedisTopic   string    `yaml:"toRedisTopic"`
	ToMongoTopic   string    `yaml:"toMongoTopic"`
	ToPushTopic    string    `yaml:"toPushTopic"`
	ToRedisGroupID string    `yaml:"toRedisGroupID"`
	ToMongoGroupID string    `yaml:"toMongoGroupID"`
	ToPushGroupID  string    `yaml:"toPushGroupID"`
	Tls            TLSConfig `yaml:"tls"`
	Test           string    `yaml:"test"`
}

type TLSConfig struct {
	EnableTLS          bool   `yaml:"enableTLS"`
	CACrt              string `yaml:"caCrt"`
	ClientCrt          string `yaml:"clientCrt"`
	ClientKey          string `yaml:"clientKey"`
	ClientKeyPwd       string `yaml:"clientKeyPwd"`
	InsecureSkipVerify bool   `yaml:"insecureSkipVerify"`
}

func init() {
	sc := []constant.ServerConfig{{
		IpAddr: "127.0.0.1",
		Port:   8848,
	}}

	cc := constant.ClientConfig{
		// 如果需要支持多namespace，我们可以创建多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		NamespaceId:         "9db1dfc0-4ed2-4aa5-a5b5-67af919b971a",
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "./tmp/log",
		CacheDir:            "./tmp/cache",
		LogLevel:            "debug",
	}

	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})
	if err != nil {
		fmt.Println(err.Error())
	}

	kafkaConfig := &Kafka{}

	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: "kafka.yml",
		Group:  "DEFAULT_GROUP",
	})
	UnmarshalConfig(content, kafkaConfig)
	fmt.Printf("\n\n\n\n config:\n%+v\n", StructToJsonString(kafkaConfig))

	if err != nil {
		fmt.Println(err.Error())
	}
	err = configClient.ListenConfig(vo.ConfigParam{
		DataId: "kafka.yml",
		Group:  "DEFAULT_GROUP",
		OnChange: func(namespace, group, dataId, data string) {
			fmt.Println("配置文件发生了变化...")
			fmt.Println("group:" + group + ", dataId:" + dataId + ", data:" + data)

			UnmarshalConfig(data, kafkaConfig)

			fmt.Printf("\n\n\n\n config:\n%+v\n", StructToJsonString(kafkaConfig))
		},
	})
}

func UnmarshalConfig(content string, c any) {

	fmt.Println("\n\ncontent:")
	fmt.Println(content)
	err := yaml.Unmarshal([]byte(content), c)
	if err != nil {
		fmt.Println("\n\n\n\nerr:", err)
	}
}

func StructToJsonString(data any) string {
	jsonData, _ := json.Marshal(data)
	return string(jsonData)
}
