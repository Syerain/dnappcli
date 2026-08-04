package config

// Config 是 dnappcli 的配置文件结构。
// 注意：server 端配置与客户端配置即便字段同名也各自独立。
type Config struct {
	Main struct {
		IsDebugMode bool `mapstructure:"isDebugMode" default:"false"`
	} `mapstructure:"main"`

	Log struct {
		IsColored bool `mapstructure:"isColored" default:"false"`
	} `mapstructure:"log"`

	Security struct {
		// RegistercodeEnckey 是签发注册码所用的 ed25519 私钥（十六进制）。
		RegistercodeEnckey string `mapstructure:"registercodeEnckey"`
	} `mapstructure:"security"`
}
