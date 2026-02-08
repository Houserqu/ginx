// 日志封装
// 请使用 LogC(c) 来打印带请求信息的日志，便于追踪问题

package core

import (
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/lmittmann/tint"
	"github.com/spf13/viper"
)

var Log *slog.Logger

func InitLogger() {
	var handler slog.Handler

	if viper.GetBool("dev") {
		// 开发环境使用 tint 输出彩色日志
		handler = tint.NewHandler(os.Stdout, &tint.Options{
			Level: slog.LevelDebug,
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				// 将时间格式化为更易读的格式
				if a.Key == slog.TimeKey {
					return slog.String(a.Key, a.Value.Time().Format("2006-01-02 15:04:05"))
				}
				return a
			},
		})
	} else {
		// 生产环境使用 JSON 格式日志
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,            // 显示文件名和行号
			Level:     slog.LevelDebug, // 设置最低日志级别
		})
	}

	Log = slog.New(handler)
}

// LogC 返回一个带有请求上下文信息的 Logger，在日志中添加请求 ID 和用户 ID 以便追踪
func LogC(c *gin.Context) *slog.Logger {
	return Log.With(
		slog.String("req_id", GetReqId(c)),
		slog.Int64("user_id", GetUserId(c)),
	)
}
