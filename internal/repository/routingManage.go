package repository

import (
	"GoWAFer/internal/types"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strings"
)

const (
	blackRoutingKey = "blackRoutingList" // 黑名单路由
	whiteRoutingKey = "whiteRoutingList" // 白名单路由
)

// RoutingManageRepository 路由管理仓库接口
type RoutingManageRepository struct {
	ctx context.Context
	rdb *redis.Client
}

// NewRoutingManageRepository 实例化路由管理仓库接口
func NewRoutingManageRepository(rdb *redis.Client) *RoutingManageRepository {
	return &RoutingManageRepository{
		ctx: context.Background(),
		rdb: rdb,
	}
}

// 生成key
func generateRoutingKey(rType int) string {
	switch rType {
	case 1:
		return blackRoutingKey
	case 2:
		return whiteRoutingKey
	default:
		return blackRoutingKey
	}
}

// Add 新增一条路由管理记录，设置过期时间储存到redis中
func (r *RoutingManageRepository) Add(routing, method string, routingType int) error {
	key := fmt.Sprintf("%s:%s:%s", generateRoutingKey(routingType), routing, method)
	return r.rdb.Set(r.ctx, key, 1, 0).Err()
}

// Del 删除一条路由管理记录
func (r *RoutingManageRepository) Del(routing, method string, routingType int) error {
	key := fmt.Sprintf("%s:%s:%s", generateRoutingKey(routingType), routing, method)
	return r.rdb.Del(r.ctx, key).Err()
}

// GetAllWithPagination 分页查询路由管理记录
func (r *RoutingManageRepository) GetAllWithPagination(page, pageSize, routingType int, keywords string) ([]types.RouteInfo, int) {
	keyHeader := generateRoutingKey(routingType)
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

	routeInfos := make([]types.RouteInfo, 0)
	for _, key := range keys {
		parts := strings.Split(key, ":")
		routeInfo := types.RouteInfo{
			Routing: parts[1],
			Method:  parts[2],
		}
		routeInfos = append(routeInfos, routeInfo)
	}
	return routeInfos, count
}

// IsExist 判断路由记录是否存在
func (r *RoutingManageRepository) IsExist(routing, method string, routingType int) (bool, error) {
	key := fmt.Sprintf("%s:%s:%s", generateRoutingKey(routingType), routing, method)
	result, err := r.rdb.Exists(r.ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result == 1, nil
}
