package main

import (
	"fmt"
	naCosClient "github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	naCosConstant "github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"gopkg.in/yaml.v2"
)

var naCosConfigClient config_client.IConfigClient

func InitNaCosClient() {
	var err error
	naCosConfigClient, err = naCosClient.CreateConfigClient(map[string]interface{}{
		"serverConfigs": []naCosConstant.ServerConfig{{
			IpAddr: "127.0.0.1",
			Port:   8848,
		}},
		"clientConfig": naCosConstant.ClientConfig{
			NamespaceId:         "9db1dfc0-4ed2-4aa5-a5b5-67af919b971a",
			TimeoutMs:           5000,
			NotLoadCacheAtStart: true,
			LogDir:              "./tmp/log",
			CacheDir:            "./tmp/cache",
			LogLevel:            "debug",
		},
	})

	if err != nil {
		fmt.Println(err.Error())
	}
}

func GetConfig(dataId string, group string, config any) {
	configContent, err := naCosConfigClient.GetConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
	})

	UnmarshalConfig(dataId, configContent, config)

	if err != nil {
		fmt.Println(err.Error())
	}

	err = naCosConfigClient.ListenConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
		OnChange: func(namespace, group, dataId, configContent string) {
			UnmarshalConfig(dataId, configContent, config)

		},
	})
}

func UnmarshalConfig(dataId, in string, out any) {
	err := yaml.Unmarshal([]byte(in), out)
	if err != nil {
		fmt.Println("\nerr:", err)
	}

	fmt.Printf("\nload %s config:\n%+v\n", dataId, StructToJsonString(out))
}
