// 常用的上下文处理函数

package core

import "github.com/gin-gonic/gin"

func GetUserId(c *gin.Context) int64 {
	return c.GetInt64("UserId")
}

func SetUserId(c *gin.Context, userId int64) {
	c.Set("UserId", userId)
}

func GetReqId(c *gin.Context) string {
	return c.GetString("ReqId")
}

func SetReqId(c *gin.Context, reqId string) {
	c.Set("ReqId", reqId)
}
