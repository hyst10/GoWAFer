package types

import "errors"

var (
	ErrInvalidBody     = errors.New("请求体未通过验证")
	ErrRoutingNotMatch = errors.New("请输入正确的路由格式！示例：/admin/index")
	ErrIPNotMatch      = errors.New("请输入正确的ipv4格式！示例：192.168.1.1")
)
