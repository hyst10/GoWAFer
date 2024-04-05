package repository

import (
	"GoWAFer/internal/types"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strings"
	"time"
)

const (
	blackIPKey = "blackIPList" // 黑名单IP
	whiteIPKey = "whiteIPList" // 白名单IP
)

// IPManageRepository IP管理仓库接口
type IPManageRepository struct {
	ctx context.Context
	rdb *redis.Client
}

// NewIPManageRepository 实例化IP管理仓库接口
func NewIPManageRepository(rdb *redis.Client) *IPManageRepository {
	return &IPManageRepository{
		ctx: context.Background(),
		rdb: rdb,
	}
}

// 生成key
func generateIPKey(rType int) string {
	switch rType {
	case 1:
		return blackIPKey
	case 2:
		return whiteIPKey
	default:
		return blackIPKey
	}
}

// Add 新增一条IP管理记录，设置过期时间储存到redis中
func (r *IPManageRepository) Add(ip string, expiration, ipType int) error {
	key := fmt.Sprintf("%s:%s", generateIPKey(ipType), ip)
	duration := time.Duration(expiration) * time.Second
	return r.rdb.Set(r.ctx, key, 1, duration).Err()
}

// Del 删除一条IP管理记录
func (r *IPManageRepository) Del(ip string, ipType int) error {
	key := fmt.Sprintf("%s:%s", generateIPKey(ipType), ip)
	return r.rdb.Del(r.ctx, key).Err()
}

// GetAllWithPagination 分页查询IP管理记录，并查询生存时间
func (r *IPManageRepository) GetAllWithPagination(page, pageSize, ipType int, keywords string) ([]types.IPInfo, int) {
	keyHeader := generateIPKey(ipType)
	startIndex := (page - 1) * pageSize
	endIndex := startIndex + pageSize - 1
	var keys []string
	var count int
	// 获取集合中所有成员
	if keywords != "" {
		keys, _ = r.rdb.Keys(r.ctx, fmt.Sprintf("%s:*%s*", keyHeader, keywords)).Result()
	} else {
		keys, _ = r.rdb.Keys(r.ctx, fmt.Sprintf("%s:*", keyHeader)).Result()
	}
	count = len(keys)

	// 检查索引
	if startIndex >= count {
		keys = []string{}
	} else {
		endIndex = startIndex + pageSize
		if endIndex > count {
			endIndex = count
		}
		keys = keys[startIndex:endIndex]
	}

	ipInfos := make([]types.IPInfo, 0)
	for _, key := range keys {
		ttlResult := r.rdb.TTL(r.ctx, key).Val()
		expirationStr := ""
		if ttlResult == -1 {
			expirationStr = "永久"
		} else {
			expiration := time.Now().Add(ttlResult)
			expirationStr = expiration.Format("2006-01-02 15:04:05")
		}
		parts := strings.Split(key, ":")
		ipInfo := types.IPInfo{
			IP:         parts[1],
			Expiration: expirationStr,
		}
		ipInfos = append(ipInfos, ipInfo)
	}
	return ipInfos, count
}

// IsExist 判断IP记录是否存在
func (r *IPManageRepository) IsExist(ip string, ipType int) (bool, error) {
	key := fmt.Sprintf("%s:%s", generateIPKey(ipType), ip)
	result, err := r.rdb.Exists(r.ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result == 1, nil
}
