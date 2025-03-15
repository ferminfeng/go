package main

type NaCosConfigStruct struct {
	IpAddr      string
	Port        uint64
	Username    string
	Password    string
	NamespaceId string // 命名空间ID
	TimeoutMs   uint64 // 请求NaCos服务器的超时时间，默认值为10000ms
	LogDir      string // 日志文件目录
	CacheDir    string // 缓存文件目录
	LogLevel    string // 日志级别 debug info warn error panic

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

type Redis struct {
	Address        []string `yaml:"address"`
	Username       string   `yaml:"username"`
	Password       string   `yaml:"password"`
	EnablePipeline bool     `yaml:"enablePipeline"`
	ClusterMode    bool     `yaml:"clusterMode"`
	DB             int      `yaml:"db"`
	MaxRetry       int      `yaml:"MaxRetry"`
}
