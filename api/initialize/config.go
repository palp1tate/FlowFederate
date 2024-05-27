package initialize

import (
	"gopkg.in/yaml.v3"

	"github.com/palp1tate/FlowFederate/api/global"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"
)

func InitConfig() {
	configFileName := "api/config.yaml"
	v := viper.New()
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := v.Unmarshal(&global.NacosConfig); err != nil {
		panic(err)
	}

	sc := []constant.ServerConfig{
		{
			IpAddr: global.NacosConfig.Host,
			Port:   uint64(global.NacosConfig.Port),
		},
	}

	cc := constant.ClientConfig{
		NamespaceId:         global.NacosConfig.Namespace,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "api/tmp/nacos/log",
		CacheDir:            "api/tmp/nacos/cache",
		LogLevel:            "warn",
	}

	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})
	if err != nil {
		panic(err)
	}
	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: global.NacosConfig.DataId,
		Group:  global.NacosConfig.Group,
	})
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal([]byte(content), &global.ServerConfig)
	if err != nil {
		panic(err)
	}
}
