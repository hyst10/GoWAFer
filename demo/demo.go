package main

import (
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

type Response struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
	Data   data   `json:"data"`
}
type data struct {
	ClientIP string `json:"clientIP"`
	Method   string `json:"method"`
	Body     string `json:"body"`
}

// 一个用于测试的demo，有基础的HTTP方法响应
func main() {
	r := gin.Default()

	r.GET("", func(c *gin.Context) {
		c.JSON(http.StatusOK, Response{
			Status: 0,
			Msg:    "success",
			Data: data{
				ClientIP: c.ClientIP(),
				Method:   "GET",
			},
		})
	})
	r.POST("", func(c *gin.Context) {
		// 读取请求体
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, Response{Status: 400, Msg: "请求体出错"})
			return
		}
		// 响应请求者，回显请求体内容
		c.JSON(http.StatusOK, Response{
			Status: 0,
			Msg:    "success",
			Data: data{
				ClientIP: c.ClientIP(),
				Method:   "POST",
				Body:     string(bodyBytes),
			},
		})
	})
	r.PUT("", func(c *gin.Context) {
		// 读取请求体
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, Response{Status: 400, Msg: "请求体出错"})
			return
		}
		// 响应请求者，回显请求体内容
		c.JSON(http.StatusOK, Response{
			Status: 0,
			Msg:    "success",
			Data: data{
				ClientIP: c.ClientIP(),
				Method:   "PUT",
				Body:     string(bodyBytes),
			},
		})
	})
	r.DELETE("", func(c *gin.Context) {
		// 读取请求体
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, Response{Status: 400, Msg: "请求体出错"})
			return
		}
		// 响应请求者，回显请求体内容
		c.JSON(http.StatusOK, Response{
			Status: 0,
			Msg:    "success",
			Data: data{
				ClientIP: c.ClientIP(),
				Method:   "DELETE",
				Body:     string(bodyBytes),
			},
		})
	})
	r.PATCH("", func(c *gin.Context) {
		// 读取请求体
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, Response{Status: 400, Msg: "请求体出错"})
			return
		}
		// 响应请求者，回显请求体内容
		c.JSON(http.StatusOK, Response{
			Status: 0,
			Msg:    "success",
			Data: data{
				ClientIP: c.ClientIP(),
				Method:   "PATCH",
				Body:     string(bodyBytes),
			},
		})
	})
	r.Run(":5353")
}
