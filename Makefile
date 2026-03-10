.PHONY: help run build test clean api get post put delete

APP_NAME := ginx

# 捕获命令后面的路径参数
# 用法: make api <module/action>     (默认 POST)
#       make get <module/action>
#       make post <module/action>
#       make put <module/action>
#       make delete <module/action>
# 示例: make api product/book/create
#       make get product/book/list
_API_PATH := $(word 2, $(MAKECMDGOALS))

api:
	@chmod +x scripts/new_api.sh
	@./scripts/new_api.sh "$(_API_PATH)" POST

get:
	@chmod +x scripts/new_api.sh
	@./scripts/new_api.sh "$(_API_PATH)" GET

post:
	@chmod +x scripts/new_api.sh
	@./scripts/new_api.sh "$(_API_PATH)" POST

put:
	@chmod +x scripts/new_api.sh
	@./scripts/new_api.sh "$(_API_PATH)" PUT

delete:
	@chmod +x scripts/new_api.sh
	@./scripts/new_api.sh "$(_API_PATH)" DELETE

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