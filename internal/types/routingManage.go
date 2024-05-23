package types

// RouteInfo 路由信息
type RouteInfo struct {
	Path string `json:"path"`
}

// AddRoutingRequest 添加路由请求
type AddRoutingRequest struct {
	Path    string `json:"path" binding:"required"`
	IsBlack bool   `json:"isBlack"`
}

// DeleteRoutingRequest 删除路由请求
type DeleteRoutingRequest struct {
	Path    string `json:"path" binding:"required"`
	IsBlack bool   `json:"isBlack"`
}
