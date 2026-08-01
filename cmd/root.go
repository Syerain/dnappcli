package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// serverURLFlag 是全局 --server 参数，默认指向本地测试服务。
var serverURLFlag string

// RootCmd 是 dnappcli 的顶层命令。
var RootCmd = &cobra.Command{
	Use:   "dnappcli",
	Short: "daizy-night server 的 CLI 测试客户端",
	Long:  "dnappcli 通过 HTTP 调用 daizy-night-server 提供的 API，用于开发调试与手工测试。",
}

func init() {
	RootCmd.PersistentFlags().StringVar(
		&serverURLFlag,
		"server",
		"http://127.0.0.1:4703",
		"server base URL, e.g. http://127.0.0.1:4703",
	)
}

// Execute 执行根命令。
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// serverURL 返回 server 地址，去除末尾斜杠。
func serverURL() string {
	return strings.TrimRight(serverURLFlag, "/")
}

