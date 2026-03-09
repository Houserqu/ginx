#!/bin/bash
# 用法: ./scripts/new_api.sh <module> <action> [method]
# 示例: make get order list 或 make post order create

set -e

MODULE=$1
ACTION=$2
METHOD=${3:-POST}
METHOD=$(echo "$METHOD" | tr '[:lower:]' '[:upper:]')

if [ -z "$MODULE" ] || [ -z "$ACTION" ]; then
  echo "Usage: make get <module> <action> 或 make post <module> <action>"
  echo "Example: make get order list"
  echo "Example: make post order create"
  exit 1
fi

# 将 snake_case 转为 PascalCase
to_pascal() {
  echo "$1" | awk -F'_' '{res=""; for(i=1;i<=NF;i++) res=res toupper(substr($i,1,1)) substr($i,2); print res}'
}

FUNC_NAME=$(to_pascal "$ACTION")
PARAMS_STRUCT="${FUNC_NAME}Params"
DIR="module/$MODULE"
FILE="$DIR/api_${ACTION}.go"
ROUTE="/api/${MODULE}/${ACTION}"

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
MODULE_IMPORT="\t_ \"ginx/module/$MODULE\""
if ! grep -qF "ginx/module/$MODULE" "$MODULES_FILE"; then
  # 检查是否已有 import 块
  if grep -q '^import' "$MODULES_FILE"; then
    # 在 import 块的最后一个 _ 导入后插入
    sed -i '' "/^import/,/^)/{/^)/i\\
$MODULE_IMPORT
}" "$MODULES_FILE"
  else
    # 追加新的 import 块
    printf "\nimport (\n\t_ \"ginx/module/%s\"\n)\n" "$MODULE" >> "$MODULES_FILE"
  fi
  echo "已更新: $MODULES_FILE (添加 ginx/module/$MODULE)"
fi

cat > "$FILE" << EOF
package $MODULE

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
// @Tags        ${MODULE}
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
