.PHONY: help run build test clean

APP_NAME := ginx

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