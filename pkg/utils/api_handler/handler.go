package api_handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// Response 通用响应结构
type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
}

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
