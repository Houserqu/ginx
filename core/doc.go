package core

import (
	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func InitDoc(svr *gin.Engine) {
	if viper.GetBool("dev") {
		// 渲染 Scalar 文档界面
		svr.GET("/docs", func(c *gin.Context) {
			html, _ := scalar.ApiReferenceHTML(&scalar.Options{
				SpecURL: "./docs/swagger.json", // 指向你生成的 spec 文件
				CustomOptions: scalar.CustomOptions{
					PageTitle: "ginx API 文档",
				},
			})
			c.Data(200, "text/html; charset=utf-8", []byte(html))
		})
	}
}
