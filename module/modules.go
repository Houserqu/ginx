// 统一管理所有模块的导入，新增模块时在此文件添加 import
// 在 main.go 中只需 import _ "ginx/module" 即可加载所有模块

package module

import (
	_ "ginx/module/login"
	_ "ginx/module/user"
)
