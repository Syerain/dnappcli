package cmd

import (
	"fmt"
	"time"

	"github.com/Syerain/dnappcli/internal/registercode"

	"github.com/spf13/cobra"
)

var registercodeFlags struct {
	magicword string
	expire    string // 只有数字会被当作分钟；也可写 units 如 30m
}

// registercodeCmd 生成注册码。
var registercodeCmd = &cobra.Command{
	Use:   "registercode",
	Short: "生成注册码",
	Long:  "使用 config 中的 registercodeEnckey 签发一个注册码（与 server 兼容）。",
	RunE:  runRegistercode,
}

func init() {
	f := registercodeCmd.Flags()
	f.StringVar(&registercodeFlags.magicword, "magicword", "", "magicword（可选）")
	f.StringVar(&registercodeFlags.expire, "expire", "60m", "有效期（如 30m / 2h）")

	RootCmd.AddCommand(registercodeCmd)
}

func runRegistercode(cmd *cobra.Command, args []string) error {
	cfg := GetConfig()

	d, err := time.ParseDuration(registercodeFlags.expire)
	if err != nil {
		return fmt.Errorf("--expire 格式非法（应为如 30m / 2h）: %w", err)
	}

	signer, err := registercode.NewSigner(cfg.Security.RegistercodeEnckey)
	if err != nil {
		return err
	}

	code, err := signer.Sign(registercode.Payload{
		Magicword: registercodeFlags.magicword,
		Before:    time.Now().Add(d),
	})
	if err != nil {
		return err
	}

	fmt.Printf("registercode: %s\n", code)
	return nil
}
