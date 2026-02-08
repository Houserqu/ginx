package main

import (
	"ginx/core"
	"ginx/middleware"
	"ginx/module/user"
	"ginx/utils"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	// 初始化
	core.InitConfig()
	core.InitLogger()
	utils.InitMysql()
	utils.InitRedis()

	// 创建 gin 引擎
	svr := gin.New()

	// 注册中间件
	svr.Use(gin.Recovery())
	svr.Use(middleware.AccessMiddleware())

	// 注册模块路由
	user.InitRouter(svr)

	// 启动服务器
	svr.Run(viper.GetString("server.addr"))
}
