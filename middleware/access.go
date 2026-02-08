// 请求处理中间件
// 记录请求日志

package middleware

import (
	"ginx/core"
	"log/slog"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func AccessMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录请求 ID
		reqId := c.GetHeader("x-request-id")
		if reqId == "" {
			reqId = uuid.NewV4().String()
		}
		core.SetReqId(c, reqId)

		// 处理用户ID
		userId, err := strconv.ParseInt(c.Request.Header.Get("x-user-id"), 10, 64)
		if err == nil {
			core.SetUserId(c, userId)
		}

		logAttrs := []any{
			slog.String("req_id", reqId),
			slog.String("client_ip", c.ClientIP()),
			slog.String("req_method", c.Request.Method),
			slog.String("req_uri", c.Request.RequestURI),
			slog.String("ua", c.Request.UserAgent()),
			slog.String("platform", c.GetHeader("platform")),
			slog.String("device_id", c.GetHeader("device-id")),
			slog.Int64("user_id", userId),
		}

		// 记录开始时间
		startTime := time.Now()

		defer func() {
			// 请求耗时
			endTime := time.Now()
			logAttrs = append(logAttrs, slog.Duration("latency_time", endTime.Sub(startTime)))

			core.Log.Info(
				"ACCESS",
				logAttrs...,
			)
		}()

		c.Next() // 处理请求

		logAttrs = append(logAttrs,
			slog.Int("status_code", c.Writer.Status()),
		)
	}
}
