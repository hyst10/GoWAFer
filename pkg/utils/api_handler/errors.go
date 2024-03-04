package api_handler

import "errors"

var ErrorCodeToMessage = map[int]error{
	40001: errors.New("请求体未通过验证！"),
	40002: errors.New("验证码错误！"),
	40003: errors.New("账号或密码错误！"),
	50001: errors.New("数据库异常！请稍后再试！"),
}
