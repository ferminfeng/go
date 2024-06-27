package main

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

type Redis struct {
	Address        []string `yaml:"address"`
	Username       string   `yaml:"username"`
	Password       string   `yaml:"password"`
	EnablePipeline bool     `yaml:"enablePipeline"`
	ClusterMode    bool     `yaml:"clusterMode"`
	DB             int      `yaml:"db"`
	MaxRetry       int      `yaml:"MaxRetry"`
}
