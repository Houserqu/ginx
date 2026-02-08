// 缓存中间件，基于 Redis 实现
// 通过请求 URL 和查询参数生成唯一的缓存 key，缓存响应数据

package middleware

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"ginx/core"
	"ginx/utils"
	"net/url"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func Cache(seconds int) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开发环境不使用缓存
		if viper.GetBool("dev") {
			c.Next()
			return
		}

		platform := c.Request.Header.Get("Platform")
		if platform == "" {
			platform = "none"
		}

		parsedURL, _ := url.Parse(c.Request.URL.String())

		// 会包含查询参数，导致 key 过长，用 md5 哈希处理
		hash := md5.New()
		hash.Write([]byte(parsedURL.RawQuery))
		hashValue := hash.Sum(nil)

		key := "ApiCache:" + platform + ":" + strconv.FormatUint(c.GetUint64("TeamId"), 10) + ":" + c.Request.URL.Path + ":" + hex.EncodeToString(hashValue)

		cacheData, err := utils.Redis.Get(c, key).Result()
		if err == nil && cacheData != "" {
			var resData map[string]any
			err := json.Unmarshal([]byte(cacheData), &resData)
			if err == nil {
				core.Success(c, resData["data"])
				c.Abort()
				return
			}
		}

		// 创建一个 ResponseRecorder
		recorder := &responseRecorder{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}
		c.Writer = recorder

		c.Next()

		// 读取响应体
		responseBody := recorder.body.String()
		if responseBody != "" {
			// 缓存
			utils.Redis.SetNX(c, key, responseBody, time.Duration(seconds)*time.Second)
		}
	}
}

// 自定义 ResponseRecorder，用于记录响应体
type responseRecorder struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}
