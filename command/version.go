package command

import (
	"github.com/astaxie/beego/logs"
	"os"
)

// 获取当前版本
func Version()  {
	if len(os.Args) > 1 && os.Args[1] == "version" {
		logs.Info("%s: ", "v0.1")
		os.Exit(2)
	}
}
