package initialize

import (
	"github.com/palp1tate/FlowFederate/api/global"

	"github.com/spf13/viper"
)

func InitConfig() {
	configFileName := "api/config.yaml"
	v := viper.New()
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := v.Unmarshal(&global.ServerConfig); err != nil {
		panic(err)
	}
}
