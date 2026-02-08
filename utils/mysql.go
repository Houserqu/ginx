package utils

import (
	"context"
	"fmt"
	"ginx/core"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Mysql *gorm.DB

type GormLogger struct {
	logger.Interface // 嵌入官方 Logger
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	// 动态从 context 取出 reqid
	reqID, _ := ctx.Value("ReqId").(string)

	// 执行原本的打印逻辑
	sql, rows := fc()
	if err != nil {
		core.Log.With(
			slog.String("req_id", reqID),
		).Error("MYSQL", "error", err, "sql", sql, "rows", rows)
	} else {
		core.Log.With(
			slog.String("req_id", reqID),
		).Error("MYSQL", "sql", sql, "rows", rows)
	}
}

func InitMysql() {
	if viper.GetString("mysql.host") == "" {
		return
	}

	host := viper.GetString("mysql.host")
	port := viper.GetString("mysql.port")
	user := viper.GetString("mysql.user")
	password := viper.GetString("mysql.password")
	database := viper.GetString("mysql.database")

	dsn := fmt.Sprint(user, ":", password, "@tcp(", host, ":", port, ")/", database, "?charset=utf8mb4&parseTime=True&loc=Local")

	LogLevel := logger.Warn   // 默认只输出慢SQL和错误
	if viper.GetBool("dev") { // 开发环境输出所有SQL
		LogLevel = logger.Info
	}

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,           // Slow SQL threshold
			LogLevel:                  LogLevel,              // Log level
			IgnoreRecordNotFoundError: true,                  // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      !viper.GetBool("dev"), // Don't include params in the SQL log
			Colorful:                  viper.GetBool("dev"),  // Disable color
		},
	)

	db, err := gorm.Open(mysql.New(mysql.Config{DSN: dsn}), &gorm.Config{
		Logger: &GormLogger{newLogger},
	})

	if err != nil {
		log.Fatal(err.Error())
	}

	Mysql = db
}

func MysqlC(c *gin.Context) *gorm.DB {
	return Mysql.WithContext(c)
}
