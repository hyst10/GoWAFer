package repository

import (
	"GoWAFer/internal/types"
	"context"
	"github.com/go-redis/redis/v8"
)

const (
	SqlInjectRules = "sqlInjectRules" // sql防护规则集合
)

type SqlInjectRepository struct {
	ctx context.Context
	rdb *redis.Client
}

func NewSqlInjectRepository(rdb *redis.Client) *SqlInjectRepository {
	return &SqlInjectRepository{
		ctx: context.Background(),
		rdb: rdb,
	}
}

// Add 添加sql注入防护规则到规则集合中
func (r *SqlInjectRepository) Add(rule string) error {
	if err := r.rdb.SAdd(r.ctx, SqlInjectRules, rule).Err(); err != nil {
		return err
	}
	return nil
}

// GetAll 查询全部sql注入防护规则
func (r *SqlInjectRepository) GetAll() ([]types.SqlInjectRule, int) {
	// 获取集合中所有成员
	rules, _ := r.rdb.SMembers(r.ctx, SqlInjectRules).Result()

	ruleInfos := make([]types.SqlInjectRule, 0)
	for _, rule := range rules {
		ruleInfo := types.SqlInjectRule{
			Rule: rule,
		}
		ruleInfos = append(ruleInfos, ruleInfo)
	}
	return ruleInfos, len(rules)
}

// Delete 删除sql注入规则
func (r *SqlInjectRepository) Delete(rule string) error {
	return r.rdb.SRem(r.ctx, SqlInjectRules, rule).Err()
}
