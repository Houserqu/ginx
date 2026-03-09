# ginx 项目 Copilot 指令

## 项目概述

ginx 是一个基于 **Go + Gin + GORM** 的 RESTful API 开发脚手架，核心设计理念是：
- 每个 API 对应一个独立文件，通过 `init()` 自动注册路由
- 统一的参数绑定、响应格式与 Swagger 文档
- 使用 `make get/post/put/delete <module> <action>` 一键生成 API 文件

---

## 项目目录结构

```
ginx/
├── main.go              # 入口，负责初始化和优雅停机
├── core/
│   ├── router.go        # Register() / InitRoutes() / Handler() 泛型包装器
│   ├── response.go      # 统一响应：Success / Error / Response 结构体
│   ├── context.go       # GetUserId / SetUserId 等上下文工具
│   ├── config.go        # InitConfig（viper）
│   └── logger.go        # InitLogger
├── middleware/
│   ├── auth.go          # CheckLogin() 登录校验中间件
│   ├── access.go        # AccessMiddleware() 访问日志
│   └── cache.go         # 缓存中间件
├── module/
│   ├── modules.go       # 统一 import 所有子模块（新模块在此添加）
│   └── <module>/        # 每个业务模块一个目录
│       └── api_<action>.go
├── model/               # GORM 数据模型
├── utils/
│   ├── mysql.go         # InitMysql
│   ├── redis.go         # InitRedis
│   └── crud/crud.go     # 通用 CRUD 工具（BuildQueryConditions 等）
└── scripts/
    └── new_api.sh       # API 文件生成脚本
```

---

## 核心约定

### 1. API 文件结构（每个文件必须遵循此模板）

```go
package <module>

import (
    "ginx/core"
    "ginx/middleware"
    "github.com/gin-gonic/gin"
)

func init() {
    core.Register(func(svr *gin.Engine) {
        svr.POST("/api/<module>/<action>", middleware.CheckLogin(), core.Handler(ActionName))
    })
}

type ActionNameParams struct {
    Field string `json:"field" binding:"required" example:"value"` // 字段说明
}

// ActionName godoc
// @Summary     <action> 接口
// @Description <action> 接口描述
// @Tags        <module>
// @Accept      json
// @Produce     json
// @Param       request body ActionNameParams true "请求参数"
// @Success     200 {object} core.Response "成功响应"
// @Security    JWT
// @Router      /api/<module>/<action> [POST]
func ActionName(c *gin.Context, params *ActionNameParams) (data any, err error) {
    return
}
```

### 2. 处理函数签名

所有处理函数签名固定为：
```go
func FuncName(c *gin.Context, params *ParamsStruct) (data any, err error)
```
- 参数由 `core.Handler()` 泛型自动绑定（`ShouldBind`），无需在函数内手动绑定
- 返回 `err != nil` 时自动响应错误；返回 `data` 时自动响应成功
- 不需要在函数内调用 `c.JSON()`

### 3. 统一响应

```go
// 成功（直接 return data 即可，框架自动调用）
core.Success(c, data)

// 成功列表（含分页）
core.SuccessList(c, list, total)

// 错误
core.Error(c, "错误信息")

// 带错误码
core.ErrorWithCode(c, 40001, "自定义错误")
```

响应格式统一为：
```json
{ "code": 0, "msg": "success", "data": {} }
```

### 4. 路由注册

- 路由在各 `api_*.go` 的 `init()` 中通过 `core.Register()` 注册
- **不需要**在 `main.go` 或其他地方手动添加路由
- 新模块创建后，在 `module/modules.go` 中添加 `import _ "ginx/module/<newmodule>"`

### 5. 参数结构体标签规范

```go
type Params struct {
    // binding 校验
    Name  string `json:"name" binding:"required"`
    Email string `json:"email" binding:"required,email"`
    Page  int    `json:"page" binding:"min=1" form:"page"`   // GET 请求用 form 标签
    Size  int    `json:"size" binding:"min=1,max=100" form:"size"`

    // Swagger 示例值
    Phone string `json:"phone" binding:"required" example:"18888888888"`

    // CRUD 查询条件标签（utils/crud 使用）
    Status int `json:"status" query:"eq"`    // WHERE status = ?
    Name   string `json:"name" query:"like"` // WHERE name LIKE ?
}
```

### 6. 上下文工具

```go
userId := core.GetUserId(c)   // 获取当前登录用户 ID
reqId  := core.GetReqId(c)    // 获取请求 ID（用于链路追踪）
```

### 7. 新增模块流程

1. 运行 `make post <module> <action>` 生成 API 文件
2. 如果是新模块目录，脚本会自动更新 `module/modules.go`
3. 运行 `make doc` 重新生成 Swagger 文档

---

## Swagger 注释要点

- GET 请求：`@Param request query ParamsStruct true "请求参数"`
- POST/PUT/DELETE：`@Param request body ParamsStruct true "请求参数"`
- 需要登录的接口加：`@Security JWT`
- 响应引用统一用：`{object} core.Response`
- 分页响应：`{object} core.Response{data=object}`

---

## 文件命名规范

| 类型 | 命名规范 | 示例 |
|------|---------|------|
| API 文件 | `api_<snake_case>.go` | `api_login_by_phone.go` |
| 函数名 | PascalCase | `LoginByPhone` |
| 参数结构体 | `<FuncName>Params` | `LoginByPhoneParams` |
| 路由路径 | `/api/<module>/<snake_case>` | `/api/login/login_by_phone` |
| 模块目录 | 小写单词 | `module/user/` |

---

## 常用命令

```bash
make get <module> <action>     # 生成 GET 请求 API 文件
make post <module> <action>    # 生成 POST 请求 API 文件
make put <module> <action>     # 生成 PUT 请求 API 文件
make delete <module> <action>  # 生成 DELETE 请求 API 文件
make doc                       # 重新生成 Swagger 文档（swag init）
make build                     # 构建 linux/amd64 生产包
make docker                    # 构建 Docker 镜像
```
