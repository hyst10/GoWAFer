package types

// SqlInjectRule sql注入防护规则
type SqlInjectRule struct {
	Rule string `json:"rule"`
}

// AddSqlInjectRuleRequest 添加sql注入防护规则请求
type AddSqlInjectRuleRequest struct {
	Rule string `json:"rule" binding:"required"`
}

// DeleteSqlInjectRuleRequest 删除sql注入防护规则请求
type DeleteSqlInjectRuleRequest struct {
	Rule string `json:"rule" binding:"required"`
}
