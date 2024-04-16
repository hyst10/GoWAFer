package test

import (
	"log"
	"testing"
)

func TestIPManage(t *testing.T) {
	ip := "127.0.0.1"
	// 新增黑名单IP
	err := ipService.AddIP(ip, 3, 1)
	if err != nil {
		t.Error(err)
		return
	}
	log.Println("黑名单IP添加成功！")

	// 删除黑名单IP
	err = ipService.DeleteIP(ip, 1)
	if err != nil {
		t.Error(err)
		return
	}
	log.Println("黑名单IP删除成功！")

	// 新增白名单IP
	err = ipService.AddIP(ip, 0, 2)
	if err != nil {
		t.Error(err)
		return
	}
	log.Println("白名单IP添加成功！")

	// 删除白名单IP
	err = ipService.DeleteIP(ip, 1)
	if err != nil {
		t.Error(err)
		return
	}
	log.Println("白名单IP删除成功！")
}
