package cache

import (
	"context"
	"encoding/csv"
	"fmt"
	"github.com/go-redis/redis/v8"
	"os"
)

func ImportRulesFromCSV(rdb *redis.Client, filename, setName string) error {
	// 打开 CSV 文件
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// 读取 CSV 头部
	headers, err := reader.Read()
	if err != nil {
		return err
	}

	// 如果第一个字符是 BOM，则将其去除

	headers[0] = headers[0][3:]

	// 确保 CSV 文件包含 "rule" 列头
	if len(headers) != 1 || headers[0] != "rule" {
		return fmt.Errorf("csv配置文件格式错误！")
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
