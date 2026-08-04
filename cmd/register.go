package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Syerain/dnappcli/internal/client"
	"github.com/Syerain/dnappcli/internal/model"
	"github.com/Syerain/dnappcli/internal/registercode"

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
	Long:  "调用 server 的 /api/register 接口注册新用户。若未提供注册码则自动用 config 的 key 生成。",
	RunE:  runRegister,
}

func init() {
	f := registerCmd.Flags()
	f.StringVar(&registerFlags.username, "username", "", "用户名（仅字母数字，最长15）")
	f.StringVar(&registerFlags.nickname, "nickname", "", "昵称（最长15）")
	f.StringVar(&registerFlags.password, "password", "", "密码")
	f.StringVar(&registerFlags.registercode, "registercode", "", "注册码（留空则自动生成）")
	f.StringVar(&registerFlags.registerway, "registerway", "legacy", "注册方式：legacy / oauth-github")

	RootCmd.AddCommand(registerCmd)
}

func runRegister(cmd *cobra.Command, args []string) error {
	// REPL 模式：register <username> <password> <nickname> [registercode]
	if len(args) >= 3 {
		rc := ""
		if len(args) >= 4 {
			rc = args[3]
		}
		return DoRegister(ServerURL(), args[0], args[1], args[2], rc)
	}

	// 一次性模式下使用 flags
	if err := validateRegisterFlags(); err != nil {
		return err
	}
	return DoRegister(ServerURL(), registerFlags.username, registerFlags.password, registerFlags.nickname, registerFlags.registercode)
}

// DoRegister 执行注册逻辑。registercode 为空时自动生成。
func DoRegister(serverURL, username, password, nickname, rc string) error {
	if username == "" {
		return fmt.Errorf("用户名不能为空")
	}
	if password == "" {
		return fmt.Errorf("密码不能为空")
	}
	if nickname == "" {
		return fmt.Errorf("昵称不能为空")
	}

	if rc == "" {
		cfg := GetConfig()
		if cfg.Security.RegistercodeEnckey == "" {
			return fmt.Errorf("config 中未配置 registercodeEnckey，无法自动生成注册码")
		}
		signer, err := registercode.NewSigner(cfg.Security.RegistercodeEnckey)
		if err != nil {
			return fmt.Errorf("创建注册码签名器失败: %w", err)
		}
		genCode, err := signer.Sign(registercode.Payload{
			Before: time.Now().Add(1 * time.Hour),
		})
		if err != nil {
			return fmt.Errorf("生成注册码失败: %w", err)
		}
		rc = genCode
		fmt.Printf("自动生成注册码: %s\n", rc)
	}

	body := model.RegisterBody{
		Username:     username,
		Nickname:     nickname,
		Password:     password,
		Registercode: rc,
		Registerway:  "legacy",
	}

	ctx := context.Background()
	c := client.New(serverURL)
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
	// registercode 留空时自动生成，不再校验
	if registerFlags.registerway != "legacy" && registerFlags.registerway != "oauth-github" {
		return fmt.Errorf("--registerway 必须为 legacy 或 oauth-github")
	}
	return nil
}
