package types

// IPInfo IP信息
type IPInfo struct {
	IP         string `json:"ip"`         // IP
	Expiration string `json:"expiration"` // 过期时间
}

// AddIPRequest 添加IP记录请求
type AddIPRequest struct {
	IP         string `json:"ip" binding:"required"`
	Type       int    `json:"type" binding:"required"`
	Expiration int    `json:"expiration"`
}

// DeleteIPRequest 删除IP记录请求
type DeleteIPRequest struct {
	IP   string `json:"ip" binding:"required"`
	Type int    `json:"type" binding:"required"`
}
