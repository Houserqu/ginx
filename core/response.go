// 响应封装，统一接口返回格式

package core

import "github.com/gin-gonic/gin"

// Response 统一响应结构体
type Response struct {
	Code int    `json:"code" example:"0"`
	Msg  string `json:"msg" example:"success"`
	Data any    `json:"data"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "success",
		"data": data,
	})
}

func SuccessList(c *gin.Context, list interface{}, total int64) {
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{
			"list":  list,
			"total": total,
		},
	})
}

func SuccessWithMsg(c *gin.Context, msg string, data interface{}) {
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  msg,
		"data": data,
	})
}

func Error(c *gin.Context, msg string) {
	c.JSON(200, gin.H{
		"code": 1,
		"msg":  msg,
		"data": nil,
	})
}

func ErrorWithCode(c *gin.Context, code int, msg string) {
	c.JSON(200, gin.H{
		"code": code,
		"msg":  msg,
		"data": nil,
	})
}
