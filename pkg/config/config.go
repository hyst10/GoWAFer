package config

import (
	"context"
	"encoding/csv"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"os"
)

// Config 配置文件结构
type Config struct {
	Server struct {
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

// ImportRulesFromCSV 从 CSV 文件中读取规则并导入到 Redis 集合中
func ImportRulesFromCSV(rdb *redis.Client, filename, setName string) error {
	// 打开 CSV 文件
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// 创建 CSV 读取器
	reader := csv.NewReader(file)

	// 读取 CSV 头部
	headers, err := reader.Read()
	if err != nil {
		return err
	}

	// 确保 CSV 文件包含 "rule" 列头
	if len(headers) != 1 || headers[0] != "rule" {
		return fmt.Errorf("CSV file must contain a single 'rule' header")
	}

	ctx := context.Background()

	// 逐行读取 CSV 文件内容
	for {
		// 读取一行 CSV 记录
		record, err := reader.Read()
		if err != nil {
			// 判断是否到达文件末尾
			if err.Error() == "EOF" {
				break
			}
			return err
		}

		// 获取规则并添加到 Redis 集合中
		rule := record[0]
		if err := rdb.SAdd(ctx, setName, rule).Err(); err != nil {
			return err
		}
	}

	return nil
}
