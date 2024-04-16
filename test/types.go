package test

import (
	"GoWAFer/internal/repository"
	"GoWAFer/internal/service"
	"GoWAFer/pkg/database"
)

var db, _ = database.InitDB()
var rdb = database.InitRedis()
var (
	adminRepository   = repository.NewAdminRepository(db)
	adminService      = service.NewAdminService(adminRepository)
	ipRepository      = repository.NewIPManageRepository(rdb)
	ipService         = service.NewIPManageService(ipRepository)
	routingRepository = repository.NewRoutingManageRepository(rdb)
	routingService    = service.NewRoutingManageService(routingRepository)
)
