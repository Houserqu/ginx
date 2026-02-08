// 登录接口实现

package user

import "github.com/gin-gonic/gin"

type UserLoginParams struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func UserLogin(c *gin.Context, params *UserLoginParams) (data any, err error) {
	return
}
