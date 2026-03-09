package main

import (
	"ginx/core"
	"ginx/middleware"
	"ginx/utils"

	_ "ginx/module/user"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// @title           ginx 项目 API 文档
// @version         1.0
// @description     基于 Gin + Gorm 框架的 API 开发脚手架

// @securityDefinitions.apikey JWT
// @in header
// @name Authorization
// @description Bearer token 格式: "Bearer {token}"

// @schemes http https
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

	// 注册路由
	core.InitRoutes(svr)

	// 注册 API 文档路由
	core.InitDoc(svr)

	// 启动服务器
	svr.Run(viper.GetString("server.addr"))
}
