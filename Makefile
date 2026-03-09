.PHONY: help run build test clean get post module

APP_NAME := ginx

# 捕获命令后面的两个参数作为 module 和 action
# 用法: make get <module> <action> 或 make post <module> <action>
_API_ARGS := $(wordlist 2, 3, $(MAKECMDGOALS))
_API_MODULE := $(word 1, $(_API_ARGS))
_API_ACTION := $(word 2, $(_API_ARGS))

get:
	@chmod +x scripts/new_api.sh
	@./scripts/new_api.sh $(_API_MODULE) $(_API_ACTION) GET

post:
	@chmod +x scripts/new_api.sh
	@./scripts/new_api.sh $(_API_MODULE) $(_API_ACTION) POST

# 捕获 module 后面的一个参数作为模块名
# 用法: make module <module>
_MOD_ARGS := $(wordlist 2, 2, $(MAKECMDGOALS))
_MOD_NAME := $(word 1, $(_MOD_ARGS))

# 防止 Make 将命令行参数当作目标
%:
	@:

build:
	@echo "构建 amd64 版本的 Linux 可执行文件"
	GOOS=linux GOARCH=amd64 go build --tags=production -o app
	@echo "构建完成"
docker:
	@echo "构建 Docker 镜像"
	docker build -t $(APP_NAME):latest .
	@echo "Docker 镜像构建完成"
clean:
	rm -f app
doc: # 生成 API 文档
	@echo "生成 API 文档"
	swag init -g main.go -o docs
	@echo "API 文档生成完成"