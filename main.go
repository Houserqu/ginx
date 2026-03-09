package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ginx/core"
	"ginx/middleware"
	"ginx/utils"

	_ "ginx/module"

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

	// 健康检查
	svr.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// 注册路由
	core.InitRoutes(svr)

	// 注册 API 文档路由
	core.InitDoc(svr)

	// 优雅停机
	srv := &http.Server{
		Addr:    viper.GetString("server.addr"),
		Handler: svr,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("正在关闭服务器...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("服务器强制关闭:", err)
	}
	log.Println("服务器已退出")
}
