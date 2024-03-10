package api_handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

// ClientErrorHandler 客户端错误处理
func ClientErrorHandler(c *gin.Context, code int) {
	c.JSON(http.StatusOK, Response{Status: code, Message: ErrorCodeToMessage[code].Error()})
	c.Abort()
	return
}

// InternalErrorHandler 服务端错误处理
func InternalErrorHandler(c *gin.Context, err error) {
	log.Printf("服务器异常错误：%s", err)
	c.JSON(http.StatusOK, Response{Status: 500, Message: "服务器异常，请稍后再试！"})
	c.Abort()
	return
}

func ForbiddenHandler(c *gin.Context, errMsg string) {
	c.JSON(http.StatusForbidden, Response{Status: 403, Message: errMsg})
	c.Abort()
	return
}

// GetUintParamFromPath 获取某个路由参数上的值并返回uint类型
func GetUintParamFromPath(c *gin.Context, param string) (uint, error) {
	param = c.Param(param)
	// 转换为uint类型
	paramUint, err := strconv.ParseUint(param, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(paramUint), nil
}
