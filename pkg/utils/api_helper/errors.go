package api_helper

import "errors"

var ErrorCodeToMessage = map[int]error{
	40001: errors.New("请求体未通过验证！"),
	40002: errors.New("验证码错误！"),
	40003: errors.New("账号或密码错误！"),
	40004: errors.New("路由参数不合法！"),
	40005: errors.New("不存在此数据！"),
	40006: errors.New("请输入正确的路由格式！示例：/admin/index"),
	40007: errors.New("请输入正确的ipv4格式！示例：192.168.1.1"),
	50001: errors.New("数据库异常！请稍后再试！"),
}
