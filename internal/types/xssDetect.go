package types

// XssDetectRule xss攻击防护规则
type XssDetectRule struct {
	Rule string `json:"rule"`
}

// AddXssDetectRuleRequest 添加xss攻击防护规则请求
type AddXssDetectRuleRequest struct {
	Rule string `json:"rule" binding:"required"`
}

// DeleteXssDetectRuleRequest 删除xss攻击防护规则请求
type DeleteXssDetectRuleRequest struct {
	Rule string `json:"rule" binding:"required"`
}
