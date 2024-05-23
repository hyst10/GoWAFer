package repository

import (
	"GoWAFer/constants"
	"GoWAFer/internal/types"
	"context"
	"github.com/go-redis/redis/v8"
)

type XssDetectRepository struct {
	ctx context.Context
	rdb *redis.Client
}

func NewXssDetectRepository(rdb *redis.Client) *XssDetectRepository {
	return &XssDetectRepository{
		ctx: context.Background(),
		rdb: rdb}
}

// Add 新增xss攻击防护规则
func (r *XssDetectRepository) Add(rule string) error {
	if err := r.rdb.SAdd(r.ctx, constants.XssDetectKey, rule).Err(); err != nil {
		return err
	}
	return nil
}

// GetAll 查询全部xss攻击防护规则
func (r *XssDetectRepository) GetAll() ([]types.SqlInjectRule, int) {
	// 获取集合中所有成员
	rules, _ := r.rdb.SMembers(r.ctx, constants.XssDetectKey).Result()

	ruleInfos := make([]types.SqlInjectRule, 0)
	for _, rule := range rules {
		ruleInfo := types.SqlInjectRule{
			Rule: rule,
		}
		ruleInfos = append(ruleInfos, ruleInfo)
	}
	return ruleInfos, len(rules)
}

// Delete 删除xss攻击防护规则
func (r *XssDetectRepository) Delete(rule string) error {
	return r.rdb.SRem(r.ctx, constants.XssDetectKey, rule).Err()
}
