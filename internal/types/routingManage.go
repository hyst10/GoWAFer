package types

// RouteInfo 路由信息
type RouteInfo struct {
	Routing string `json:"routing"` // 路由名
	Method  string `json:"method"`  // HTTP 方法
}

// AddRoutingRequest 添加路由请求
type AddRoutingRequest struct {
	Routing string `json:"routing" binding:"required"`
	Type    int    `json:"type" binding:"required"`
	Method  string `json:"method" binding:"required"`
}

// DeleteRoutingRequest 删除路由请求
type DeleteRoutingRequest struct {
	Routing string `json:"routing" binding:"required"`
	Type    int    `json:"type" binding:"required"`
	Method  string `json:"method" binding:"required"`
}
