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

	// 限速器配置
	RateLimiter struct {
		MaxCounter  int `mapstructure:"maxCounter" json:"maxCounter"`   // 计数器最大值，超过该值拉入永久IP黑名单
		BanCounter  int `mapstructure:"banCounter" json:"banCounter"`   // 计数器常规值，超过触发常规封禁时间
		BanDuration int `mapstructure:"banDuration" json:"banDuration"` // 常规封禁时间，单位秒
		Mode        int `mapstructure:"mode" json:"mode"`               // 限速算法，1:令牌桶2:漏桶3:固定窗口4:滑动窗口
		// 令牌桶模式配置
		TokenBucket struct {
			MaxToken       int `mapstructure:"maxToken" json:"maxToken"`             // 最大令牌数
			TokenPerSecond int `mapstructure:"tokenPerSecond" json:"tokenPerSecond"` //令牌生成速度，单位秒
		} `mapstructure:"tokenBucket" json:"tokenBucket"`
		// 漏桶模式配置
		LeakyBucket struct {
			Capacity       int `mapstructure:"capacity" json:"capacity"`              // 桶容量
			LeakyPerSecond int ` mapstructure:"leakyPerSecond" json:"leakyPerSecond"` // 漏水速度，单位秒
		} `mapstructure:"leakyBucket" json:"leakyBucket"`
		// 固定窗口模式配置
		FixedWindow struct {
			WindowSize int `mapstructure:"windowSize" json:"windowSize"` // 窗口大小,单位秒
			MaxRequest int `mapstructure:"maxRequest" json:"maxRequest"` // 窗口内最大请求数
		} `mapstructure:"fixedWindow" json:"fixedWindow"`
		// 滑动窗口模式配置
		SlideWindow struct {
			WindowSize int `mapstructure:"windowSize" json:"windowSize"` // 窗口大小，单位秒
			MaxRequest int `mapstructure:"maxRequest" json:"maxRequest"` // 窗口内最大队列数
		} `mapstructure:"slideWindow" json:"slideWindow"`
	} `mapstructure:"rateLimiter" json:"rateLimiter"`
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
