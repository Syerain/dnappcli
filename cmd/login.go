package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/Syerain/dnappcli/internal/client"
	"github.com/Syerain/dnappcli/internal/model"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var loginFlags struct {
	username string
	password string
	loginway string
}

// loginCmd 对应 POST /api/login。
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "用户登录",
	Long:  "调用 server 的 /api/login 接口登录，成功后 token 保存到 data.yaml。",
	RunE:  runLogin,
}

func init() {
	f := loginCmd.Flags()
	f.StringVar(&loginFlags.username, "username", "", "用户名")
	f.StringVar(&loginFlags.password, "password", "", "密码")
	f.StringVar(&loginFlags.loginway, "loginway", "legacy", "登录方式：legacy / oauth-github")

	RootCmd.AddCommand(loginCmd)
}

func runLogin(cmd *cobra.Command, args []string) error {
	if err := validateLoginFlags(); err != nil {
		return err
	}

	body := model.LoginBody{
		Loginway: loginFlags.loginway,
		Username: loginFlags.username,
		Password: loginFlags.password,
	}

	ctx := context.Background()
	c := client.New(serverURL())
	code, respBody, err := c.PostJSON(ctx, "/api/login", body)
	if err != nil {
		return fmt.Errorf("login request failed: %w", err)
	}

	fmt.Printf("HTTP %d\n", code)

	if code == 200 {
		var loginResp model.LoginResponse
		if jerr := json.Unmarshal(respBody, &loginResp); jerr != nil {
			return fmt.Errorf("parse login response: %w", jerr)
		}

		fmt.Printf("access_token:  %s\n", loginResp.AccessToken)
		fmt.Printf("refresh_token: %s\n", loginResp.RefreshToken)

		if err := saveTokenData(loginResp.AccessToken, loginResp.RefreshToken); err != nil {
			return fmt.Errorf("save token data: %w", err)
		}
		fmt.Println("token 已保存到 data.yaml")
	} else {
		var resp model.Response
		if jerr := json.Unmarshal(respBody, &resp); jerr == nil {
			fmt.Printf("message: %s\n", resp.Message)
		} else {
			fmt.Printf("body: %s\n", string(respBody))
		}
	}

	return nil
}

func validateLoginFlags() error {
	if loginFlags.username == "" {
		return fmt.Errorf("--username 必填")
	}
	if loginFlags.password == "" {
		return fmt.Errorf("--password 必填")
	}
	if loginFlags.loginway != "legacy" && loginFlags.loginway != "oauth-github" {
		return fmt.Errorf("--loginway 必须为 legacy 或 oauth-github")
	}
	return nil
}

// saveTokenData 将 access_token 和 refresh_token 写入 data.yaml。
func saveTokenData(accessToken, refreshToken string) error {
	data := model.TokenData{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	yamlBytes, err := yaml.Marshal(&data)
	if err != nil {
		return fmt.Errorf("marshal yaml: %w", err)
	}

	if err := os.WriteFile("data.yaml", yamlBytes, 0600); err != nil {
		return fmt.Errorf("write data.yaml: %w", err)
	}

	return nil
}
