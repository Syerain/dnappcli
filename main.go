package main

import (
	"os"

	"github.com/Syerain/dnappcli/cmd"
)

func main() {
	// 有命令行参数时：一次性模式（兼容原有用法）
	if len(os.Args) > 1 {
		cmd.Execute()
		return
	}

	// 无参数时：启动交互式 REPL
	cmd.StartREPL()
}
