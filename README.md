# Ginx2 - Golang + Gin API 项目模板

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go" alt="Go Version" />
  <img src="https://img.shields.io/badge/Gin-1.10+-00ADD8?style=flat&logo=go" alt="Gin Version" />
  <img src="https://img.shields.io/badge/License-MIT-green.svg" alt="License" />
</p>

一个开箱即用的 Golang + Gin Web API 项目模板，封装了 API 开发中的核心能力，帮助开发者快速启动新项目。

## ✨ 核心特性

- 🚀 **开箱即用** - 完整的项目结构，克隆即可开发
- 🔧 **核心封装** - 统一的配置、日志、响应、错误处理
-  **数据库集成** - GORM + MySQL，支持 CRUD 封装
- ⚡ **缓存支持** - Redis 集成，缓存中间件
- 🆔 **JS安全ID生成器** - 雪花算法变体，前端安全整数范围

## 📁 项目结构

```
ginx/
├── main.go                 # 程序入口
├── config.yaml             # 配置文件
├── Dockerfile              # Docker 构建文件
├── Makefile                # 构建脚本
├── k8s.sh                  # Kubernetes 部署脚本
│
├── core/                  # 核心功能
│   ├── config.go          # 配置管理
│   ├── context.go         # 上下文封装
│   ├── handler.go         # 处理器基础
│   ├── logger.go          # 日志管理
│   └── response.go        # 统一响应
│
├── middleware/            # 中间件
│   ├── access.go          # 访问日志
│   ├── auth.go            # 认证中间件
│   └── cache.go           # 缓存中间件
│
├── model/                 # 数据模型
│   └── user.go            # 用户模型
│
├── module/                # 业务模块
│   ├── login/             # 登录模块
│   │   ├── api_login_by_password.go
│   │   ├── router.go
│   │   ├── service.go
│   │   └── types.go
│   └── user/                # 用户管理模块
│       ├── api_user_list.go # 用户列表接口（一个接口一个文件）
│       ├── service.go       # 复用逻辑
│       │── types.go         # 常量和类型声明
│       └── router.go        # 路由绑定
│
└── utils/                 # 工具库
    ├── crud            # CRUD 工具函数封装
    ├── gin.go             # Gin 辅助
    ├── gorm.go            # GORM 辅助
    ├── mysql.go           # MySQL 连接
    ├── redis.go           # Redis 连接
    └── id.go              # JS安全雪花ID生成器
```

## 🚀 快速开始

### 环境要求

- Go 1.23+

### 安装步骤

1. **克隆项目**

```bash
git clone https://github.com/yourusername/ginx2.git
cd ginx2
```

2. **安装依赖**

```bash
go mod download
```

3. **配置文件**

创建 `config.yaml` 文件：

```yaml
server:
  port: 8080
  mode: debug  # debug/release/test

mysql:
  host: localhost
  port: 3306
  database: ginx
  username: root
  password: your_password
  charset: utf8mb4
  max_idle_conns: 10
  max_open_conns: 100

redis:
  host: localhost
  port: 6379
  password: ""
  db: 0
  pool_size: 10
```

## 📝 开发指南

### 添加新模块

1. 在 `module/` 下创建模块目录
2. 创建必要文件：
   - `router.go` - 路由定义
   - `api_*.go` - API 处理器
   - `service.go` - 业务逻辑
   - `types.go` - 数据结构定义

示例：

```go
// module/example/router.go
package example

import "github.com/gin-gonic/gin"

func InitRouter(r *gin.Engine) {
    g := r.Group("/api/example")
    {
        g.GET("/list", GetList)
        g.POST("/create", CreateItem)
    }
}
```

3. 在 `main.go` 中注册模块路由：

```go
example.InitRouter(svr)
```

### 添加中间件

在 `middleware/` 目录下创建中间件：

```go
func YourMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 前置处理
        c.Next()
        // 后置处理
    }
}
```

### 添加工具函数

在 `utils/` 目录下添加工具函数，按功能分类到不同文件。

## 🐳 部署

### Docker 部署

1. **构建镜像**

```bash
make docker
```

或手动构建：

```bash
docker build -t ginx2:latest .
```

2. **运行容器**

```bash
docker run -d \
  -p 8080:8080 \
  -v $(pwd)/config.yaml:/app/config.yaml \
  --name ginx2 \
  ginx2:latest
```

## 🛠️ 依赖库

| 库 | 用途 |
|---|---|
| [gin-gonic/gin](https://github.com/gin-gonic/gin) | Web 框架 |
| [gorm.io/gorm](https://gorm.io) | ORM 框架 |
| [go-redis/redis](https://github.com/go-redis/redis) | Redis 客户端 |
| [spf13/viper](https://github.com/spf13/viper) | 配置管理 |
| [lmittmann/tint](https://github.com/lmittmann/tint) | 彩色日志 |

## 📄 License

[MIT License](LICENSE)

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📧 联系方式

如有问题，请提交 Issue。

---

⭐ 如果这个项目对你有帮助，请给个 Star！
