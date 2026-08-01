package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/atomreforge/dnappcli/internal/client"
	"github.com/atomreforge/dnappcli/internal/model"

	"github.com/spf13/cobra"
)

var registerFlags struct {
	username     string
	nickname     string
	password     string
	registercode string
	registerway  string
}

// registerCmd 对应 POST /api/register。
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "注册新用户",
	Long:  "调用 server 的 /api/register 接口注册新用户。",
	RunE:  runRegister,
}

func init() {
	f := registerCmd.Flags()
	f.StringVar(&registerFlags.username, "username", "", "用户名（仅字母数字，最长15）")
	f.StringVar(&registerFlags.nickname, "nickname", "", "昵称（最长15）")
	f.StringVar(&registerFlags.password, "password", "", "密码")
	f.StringVar(&registerFlags.registercode, "registercode", "", "注册码（由 server 签发）")
	f.StringVar(&registerFlags.registerway, "registerway", "legacy", "注册方式：legacy / oauth-github")

	RootCmd.AddCommand(registerCmd)
}

func runRegister(cmd *cobra.Command, args []string) error {
	// 参数检查（与 server 校验逻辑对应，便于提前暴露错误）
	if err := validateRegisterFlags(); err != nil {
		return err
	}

	body := model.RegisterBody{
		Username:     registerFlags.username,
		Nickname:     registerFlags.nickname,
		Password:     registerFlags.password,
		Registercode: registerFlags.registercode,
		Registerway:  registerFlags.registerway,
	}

	ctx := context.Background()
	c := client.New(serverURL())
	code, respBody, err := c.PostJSON(ctx, "/api/register", body)
	if err != nil {
		return fmt.Errorf("register request failed: %w", err)
	}

	fmt.Printf("HTTP %d\n", code)
	if len(respBody) > 0 {
		var resp model.Response
		if jerr := json.Unmarshal(respBody, &resp); jerr == nil {
			fmt.Printf("message: %s\n", resp.Message)
		} else {
			fmt.Printf("body: %s\n", string(respBody))
		}
	}
	return nil
}

func validateRegisterFlags() error {
	if registerFlags.username == "" {
		return fmt.Errorf("--username 必填")
	}
	if registerFlags.nickname == "" {
		return fmt.Errorf("--nickname 必填")
	}
	if registerFlags.password == "" {
		return fmt.Errorf("--password 必填")
	}
	if registerFlags.registercode == "" {
		return fmt.Errorf("--registercode 必填")
	}
	if registerFlags.registerway != "legacy" && registerFlags.registerway != "oauth-github" {
		return fmt.Errorf("--registerway 必须为 legacy 或 oauth-github")
	}
	return nil
}
