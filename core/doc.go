package core

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func InitDoc(svr *gin.Engine) {
	if viper.GetBool("dev") {
		// 提供 swagger.json 静态文件访问
		svr.StaticFile("/docs/swagger.json", "./docs/swagger.json")

		// 渲染 Scalar 文档界面
		svr.GET("/docs", func(c *gin.Context) {
			html := fmt.Sprintf(`<!DOCTYPE html>
<html>
  <head>
    <title>matrix-api API 文档</title>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
  </head>
  <body>
    <div id="app"></div>
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
    <script>
      Scalar.createApiReference('#app', {
        url: '%s',
        persistAuth: true
      })
    </script>
  </body>
</html>`, "/docs/swagger.json")
			c.Data(200, "text/html; charset=utf-8", []byte(html))
		})
	}
}
