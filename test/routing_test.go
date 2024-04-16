package test

import (
	"log"
	"testing"
)

func TestRoutingManage(t *testing.T) {
	routing := "/admin"
	// 新增黑名单路由
	err := routingService.AddRouting(routing, "GET", 1)
	if err != nil {
		t.Error(err)
		return
	}
	log.Println("黑名单路由添加成功！")

	// 删除黑名单路由
	err = routingService.DeleteRouting(routing, "GET", 1)
	if err != nil {
		t.Error(err)
		return
	}
	log.Println("黑名单路由删除成功！")

	// 新增白名单路由
	err = routingService.AddRouting(routing, "GET", 2)
	if err != nil {
		t.Error(err)
		return
	}
	log.Println("白名单路由添加成功！")

	// 删除白名单路由
	err = routingService.DeleteRouting(routing, "GET", 2)
	if err != nil {
		t.Error(err)
		return
	}
	log.Println("白名单路由删除成功！")
}
