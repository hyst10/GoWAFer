package repository

import (
	"GoWAFer/constants"
	"GoWAFer/internal/types"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strings"
	"time"
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

func chooseKeyHeader(isBlack bool) string {
	if isBlack {
		return constants.BlackIPKey
	}
	return constants.WhiteIPKey
}

// Add 新增一条IP管理记录，设置过期时间储存到redis中
func (r *IPManageRepository) Add(ip string, expiration int, isBlack bool) error {
	key := fmt.Sprintf("%s:%s", chooseKeyHeader(isBlack), ip)
	return r.rdb.Set(r.ctx, key, 1, time.Duration(expiration)*time.Second).Err()
}

// Del 删除一条IP管理记录
func (r *IPManageRepository) Del(ip string, isBlack bool) error {
	key := fmt.Sprintf("%s:%s", chooseKeyHeader(isBlack), ip)
	return r.rdb.Del(r.ctx, key).Err()
}

// GetAllWithPagination 分页查询IP管理记录，并查询生存时间
func (r *IPManageRepository) GetAllWithPagination(page, pageSize int, isBlack bool, query string) ([]types.IPInfo, int) {
	startIndex := (page - 1) * pageSize
	endIndex := startIndex + pageSize - 1
	var ips []string
	var count int
	if query != "" {
		ips, _ = r.rdb.Keys(r.ctx, fmt.Sprintf("%s:*%s*", chooseKeyHeader(isBlack), query)).Result()
	} else {
		ips, _ = r.rdb.Keys(r.ctx, chooseKeyHeader(isBlack)+":*").Result()
	}
	count = len(ips)

	// 设置索引
	if startIndex >= count {
		ips = []string{}
	} else {
		endIndex = startIndex + pageSize
		if endIndex > count {
			endIndex = count
		}
		ips = ips[startIndex:endIndex]
	}

	infos := make([]types.IPInfo, 0)
	for _, item := range ips {
		ttlResult := r.rdb.TTL(r.ctx, item).Val()
		expirationStr := ""
		if ttlResult == -1 {
			expirationStr = "永久"
		} else {
			expiration := time.Now().Add(ttlResult)
			expirationStr = expiration.Format("2006-01-02 15:04:05")
		}
		parts := strings.Split(item, ":")
		info := types.IPInfo{IP: parts[1], Expiration: expirationStr}
		infos = append(infos, info)
	}
	return infos, count
}

// IsExist 判断IP记录是否存在
func (r *IPManageRepository) IsExist(ip string) string {
	key := fmt.Sprintf("%s:%s", constants.WhiteIPKey, ip)
	result, err := r.rdb.Exists(r.ctx, key).Result()
	if err != nil {
		return ""
	}
	if result == 1 {
		return constants.WhiteIPKey
	}
	key = fmt.Sprintf("%s:%s", constants.BlackIPKey, ip)
	result, err = r.rdb.Exists(r.ctx, key).Result()
	if err != nil {
		return ""
	}
	if result == 1 {
		return constants.BlackIPKey
	}
	return ""
}
