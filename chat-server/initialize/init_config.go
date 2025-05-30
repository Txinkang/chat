package initialize

import (
	"chat-server/config"
	"fmt"
	"github.com/spf13/viper"
)

var AppConfig config.AppConfig

func InitConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("viper 读取配置失败, err:%w", err)
	}
	if err := viper.Unmarshal(&AppConfig); err != nil {
		return fmt.Errorf("viper 配置反序列化失败, err:%w", err)
	}

	return nil
}
