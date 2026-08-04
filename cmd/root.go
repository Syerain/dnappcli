package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/Syerain/dnappcli/internal/config"
	"github.com/Syerain/dnappcli/internal/utils"

	"github.com/spf13/cobra"
)

// serverURLFlag 是全局 --server 参数，默认指向本地测试服务。
var serverURLFlag string

// appCfg 是全局加载的配置。
var appCfg *config.Config

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

// initViper loads app-level config & logger once.
func appInit() {
	if appCfg != nil {
		return
	}
	appCfg = config.MustLoadConfig(".")
	utils.InitModuleLogger(appCfg.Main.IsDebugMode, "appcli")
}

// GetConfig 返回全局配置。
func GetConfig() *config.Config {
	if appCfg == nil {
		appInit()
	}
	return appCfg
}

// Execute 执行根命令（一次性模式）。
func Execute() {
	appInit()
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// ExecuteArgs 用给定的 args 执行根命令（供 REPL 复用）。
func ExecuteArgs(args []string) {
	RootCmd.SetArgs(args)
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

// ServerURL 返回 server 地址，去除末尾斜杠。
func ServerURL() string {
	return strings.TrimRight(serverURLFlag, "/")
}


