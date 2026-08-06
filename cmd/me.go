package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Syerain/dnappcli/internal/client"
	"github.com/Syerain/dnappcli/internal/model"

	"github.com/spf13/cobra"
)

var meCmd = &cobra.Command{
	Use:   "me",
	Short: "获取当前登录用户信息",
	Long:  "读取 data.yaml 中的 access_token，调用 GET /api/v1/user/me 获取当前用户信息。",
	RunE:  runMe,
}

func init() {
	RootCmd.AddCommand(meCmd)
}

func runMe(cmd *cobra.Command, args []string) error {
	td, err := LoadTokenData()
	if err != nil {
		return err
	}

	ctx := context.Background()
	c := client.New(ServerURL())
	code, respBody, err := c.GetJSON(ctx, "/api/v1/user/me", td.AccessToken)
	if err != nil {
		return fmt.Errorf("me request failed: %w", err)
	}

	fmt.Printf("HTTP %d\n", code)
	if code == http.StatusOK {
		var me model.MeResponse
		if jerr := json.Unmarshal(respBody, &me); jerr != nil {
			return fmt.Errorf("parse me response: %w", jerr)
		}
		printInfoMe(me.InfoMe)
		return nil
	}

	// 非 200：按统一错误外壳打印
	var resp model.Response
	if jerr := json.Unmarshal(respBody, &resp); jerr == nil {
		fmt.Printf("message: %s\n", resp.Message)
	} else {
		fmt.Printf("body: %s\n", string(respBody))
	}
	return nil
}

func printInfoMe(m model.InfoMe) {
	fmt.Printf("uid:           %d\n", m.Uid)
	fmt.Printf("username:      %s\n", m.Username)
	fmt.Printf("nickname:      %s\n", m.Nickname)
	fmt.Printf("email:         %s\n", m.Email)
	fmt.Printf("register_time: %s\n", m.RegisterTime)
	fmt.Printf("role:          %s\n", m.Role)
	if m.GitHubID != nil {
		fmt.Printf("github_id:     %d\n", *m.GitHubID)
	} else {
		fmt.Println("github_id:     -")
	}
	fmt.Printf("github_login:  %s\n", m.GitHubLogin)
}
