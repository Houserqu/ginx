#!/bin/bash
# 用法: ./scripts/new_api.sh <path> [method]
# 示例: make api product/book/create
#       make get product/book/list
#       make post order/create
#
# <path> 格式: <模块路径>/<action>
#   - 最后一段为 action（接口名）
#   - 前面部分为模块路径（支持多级，如 product/book）
#   - package 名取模块路径最后一段

set -e

FULL_PATH=$1
METHOD=${2:-POST}
METHOD=$(echo "$METHOD" | tr '[:lower:]' '[:upper:]')

if [ -z "$FULL_PATH" ]; then
  echo "Usage: make api <path> [method]"
  echo "       make get <path>"
  echo "       make post <path>"
  echo ""
  echo "Examples:"
  echo "  make api product/book/create"
  echo "  make get product/book/list"
  echo "  make post order/create"
  exit 1
fi

# 提取 action（路径最后一段）和模块路径（其余部分）
ACTION=$(basename "$FULL_PATH")
MODULE_PATH=$(dirname "$FULL_PATH")

# 兼容单段路径（如 "login"，此时 dirname 返回 "."）
if [ "$MODULE_PATH" = "." ]; then
  echo "Error: 路径至少需要两段，格式为 <module>/<action>"
  echo "示例: make api login/login_by_phone"
  exit 1
fi

# package 名取模块路径最后一段
PACKAGE=$(basename "$MODULE_PATH")

# 将 snake_case 转为 PascalCase
to_pascal() {
  echo "$1" | awk -F'_' '{res=""; for(i=1;i<=NF;i++) res=res toupper(substr($i,1,1)) substr($i,2); print res}'
}

FUNC_NAME=$(to_pascal "$ACTION")
PARAMS_STRUCT="${FUNC_NAME}Params"
DIR="module/${MODULE_PATH}"
FILE="${DIR}/api_${ACTION}.go"
ROUTE="/api/${FULL_PATH}"
# Swagger Tags 使用完整模块路径（斜杠替换为/，便于分组）
TAGS="${MODULE_PATH}"

# GET 请求使用 query 参数，其他使用 body
if [ "$METHOD" = "GET" ]; then
  SWAGGER_PARAM="// @Param       request query  $PARAMS_STRUCT true \"请求参数\""
else
  SWAGGER_PARAM="// @Param       request body   $PARAMS_STRUCT true \"请求参数\""
fi

if [ ! -d "$DIR" ]; then
  mkdir -p "$DIR"
  echo "已创建目录: $DIR"
fi

if [ -f "$FILE" ]; then
  echo "Error: 文件 $FILE 已存在"
  exit 1
fi

# 如果 module/modules.go 不存在则创建
MODULES_FILE="module/modules.go"
if [ ! -f "$MODULES_FILE" ]; then
  cat > "$MODULES_FILE" << MODEOF
// 统一管理所有模块的导入，新增模块时在此文件添加 import
package module
MODEOF
fi

# 如果模块未被导入则自动添加
MODULE_IMPORT=$'\t'"_ \"ginx/module/${MODULE_PATH}\""
if ! grep -qF "ginx/module/${MODULE_PATH}" "$MODULES_FILE"; then
  if grep -q '^import' "$MODULES_FILE"; then
    sed -i '' "/^import/,/^)/{/^)/i\\
$MODULE_IMPORT
}" "$MODULES_FILE"
  else
    printf "\nimport (\n\t_ \"ginx/module/%s\"\n)\n" "$MODULE_PATH" >> "$MODULES_FILE"
  fi
  echo "已更新: $MODULES_FILE (添加 ginx/module/${MODULE_PATH})"
fi

cat > "$FILE" << EOF
package $PACKAGE

import (
	"ginx/core"
	"ginx/middleware"

	"github.com/gin-gonic/gin"
)

func init() {
	core.Register(func(svr *gin.Engine) {
		svr.${METHOD}("${ROUTE}", middleware.CheckLogin(), core.Handler($FUNC_NAME))
	})
}

type $PARAMS_STRUCT struct {
}

// $FUNC_NAME godoc
// @Summary     ${ACTION} 接口
// @Description ${ACTION} 接口描述
// @Tags        ${TAGS}
// @Accept      json
// @Produce     json
${SWAGGER_PARAM}
// @Success     200 {object} core.Response "成功响应"
// @Security    JWT
// @Router      ${ROUTE} [${METHOD}]
func $FUNC_NAME(c *gin.Context, params *$PARAMS_STRUCT) (data any, err error) {

	return
}
EOF

echo "已创建: $FILE"
