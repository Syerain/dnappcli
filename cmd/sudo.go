package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Syerain/dnappcli/internal/client"
	"github.com/Syerain/dnappcli/internal/model"

	"github.com/spf13/cobra"
)

var sudoCmd = &cobra.Command{
	Use:   "sudo",
	Short: "访问管理员接口（需 admin 角色）",
	Long:  "读取 data.yaml 中的 access_token，调用 POST /api/v1/admin/sudo。仅 admin 角色可访问。",
	RunE:  runSudo,
}

func init() {
	RootCmd.AddCommand(sudoCmd)
}

func runSudo(cmd *cobra.Command, args []string) error {
	td, err := LoadTokenData()
	if err != nil {
		return err
	}

	ctx := context.Background()
	c := client.New(ServerURL())
	// 服务端 Sudo 当前仅返回 nil（无响应体），body 传 nil 即可。
	code, respBody, err := c.PostJSON(ctx, "/api/v1/admin/sudo", td.AccessToken, nil)
	if err != nil {
		return fmt.Errorf("sudo request failed: %w", err)
	}

	fmt.Printf("HTTP %d\n", code)
	if len(respBody) == 0 {
		return nil
	}

	var resp model.Response
	if jerr := json.Unmarshal(respBody, &resp); jerr == nil {
		fmt.Printf("message: %s\n", resp.Message)
	} else {
		fmt.Printf("body: %s\n", string(respBody))
	}
	return nil
}
