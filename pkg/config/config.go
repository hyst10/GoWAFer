package config

import (
	"fmt"
	"github.com/spf13/viper"
)

// Config 配置文件结构
type Config struct {
	Server struct {
		WafPort       int    `mapstructure:"wafPort" json:"wafPort"`
		TargetAddress string `mapstructure:"targetAddress" json:"targetAddress"`
	} `mapstructure:"server" json:"server"`

	Secret struct {
		JwtSecretKey     string `mapstructure:"jwtSecretKey" json:"jwtSecretKey"`
		SessionSecretKey string `mapstructure:"sessionSecretKey" json:"sessionSecretKey"`
	} `mapstructure:"secret" json:"secret"`
}

// ReadConfig 读取配置文件
func ReadConfig() *Config {
	viper.SetConfigName("config") // 配置文件名
	viper.SetConfigType("yaml")   // 配置文件类型
	viper.AddConfigPath(".")      // 搜索配置文件的路径
	// 读取配置文件
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Sprintf("无法读取配置文件：%v", err))
	}
	// 将配置信息解析到结构体中
	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		panic(fmt.Sprintf("无法解析配置文件：%v", err))
	}
	return &config
}
