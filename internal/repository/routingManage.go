package repository

import (
	"GoWAFer/constants"
	"GoWAFer/internal/types"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strings"
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

func choosePathKeyHeader(isBlack bool) string {
	if isBlack {
		return constants.BlackPathKey
	}
	return constants.WhitePathKey
}

// Add 新增一条路由管理记录，设置过期时间储存到redis中
func (r *RoutingManageRepository) Add(path string, isBlack bool) error {
	key := fmt.Sprintf("%s:%s", choosePathKeyHeader(isBlack), path)
	return r.rdb.Set(r.ctx, key, 1, 0).Err()
}

// Del 删除一条路由管理记录
func (r *RoutingManageRepository) Del(path string, isBlack bool) error {
	key := fmt.Sprintf("%s:%s", choosePathKeyHeader(isBlack), path)
	return r.rdb.Del(r.ctx, key).Err()
}

// GetAllWithPagination 分页查询路由管理记录
func (r *RoutingManageRepository) GetAllWithPagination(page, pageSize int, isBlack bool, query string) ([]types.RouteInfo, int) {
	startIndex := (page - 1) * pageSize
	endIndex := startIndex + pageSize - 1
	var paths []string
	var count int
	if query != "" {
		paths, _ = r.rdb.Keys(r.ctx, fmt.Sprintf("%s:*%s*", choosePathKeyHeader(isBlack), query)).Result()
	} else {
		paths, _ = r.rdb.Keys(r.ctx, choosePathKeyHeader(isBlack)+":*").Result()
	}
	count = len(paths)

	// 设置索引
	if startIndex >= count {
		paths = []string{}
	} else {
		endIndex = startIndex + pageSize
		if endIndex > count {
			endIndex = count
		}
		paths = paths[startIndex:endIndex]
	}

	routeInfos := make([]types.RouteInfo, 0)
	for _, key := range paths {
		parts := strings.Split(key, ":")
		routeInfo := types.RouteInfo{
			Path: parts[1],
		}
		routeInfos = append(routeInfos, routeInfo)
	}
	return routeInfos, count
}

// IsExist 判断路由记录是否存在
func (r *RoutingManageRepository) IsExist(path string) string {
	key := fmt.Sprintf("%s:%s", constants.WhitePathKey, path)
	result, err := r.rdb.Exists(r.ctx, key).Result()
	if err != nil {
		return ""
	}
	if result == 1 {
		return constants.WhitePathKey
	}
	key = fmt.Sprintf("%s:%s", constants.BlackPathKey, path)
	result, err = r.rdb.Exists(r.ctx, key).Result()
	if err != nil {
		return ""
	}
	if result == 1 {
		return constants.BlackPathKey
	}
	return ""
}
