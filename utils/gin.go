package utils

import (
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

// 在同一进程内，模拟发起一次 HTTP 请求
func InternalRequest(router *gin.Engine, method, path string, body io.Reader) *httptest.ResponseRecorder {
	// 构造一个模拟的响应记录器
	w := httptest.NewRecorder()
	// 构造一个模拟的 HTTP 请求
	req, _ := http.NewRequest(method, path, body)

	// 直接让 Gin 引擎处理这个请求，不走网卡
	router.ServeHTTP(w, req)

	return w
}
