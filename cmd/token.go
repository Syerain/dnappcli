package cmd

import (
	"fmt"
	"os"

	"github.com/Syerain/dnappcli/internal/model"

	"gopkg.in/yaml.v3"
)

// tokenFilePath 是 access/refresh token 保存的文件名（与 login 命令写入一致）。
const tokenFilePath = "data.yaml"

// LoadTokenData 从 data.yaml 读取 token。文件不存在或为空时返回 error。
func LoadTokenData() (*model.TokenData, error) {
	data, err := os.ReadFile(tokenFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("未找到 %s，请先执行 login", tokenFilePath)
		}
		return nil, fmt.Errorf("读取 %s 失败: %w", tokenFilePath, err)
	}

	var td model.TokenData
	if err := yaml.Unmarshal(data, &td); err != nil {
		return nil, fmt.Errorf("解析 %s 失败: %w", tokenFilePath, err)
	}
	if td.AccessToken == "" {
		return nil, fmt.Errorf("%s 中缺少 access_token，请先执行 login", tokenFilePath)
	}
	return &td, nil
}

// SaveTokenData 将 access/refresh token 写入 data.yaml。
func SaveTokenData(accessToken, refreshToken string) error {
	td := model.TokenData{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	raw, err := yaml.Marshal(&td)
	if err != nil {
		return fmt.Errorf("marshal yaml: %w", err)
	}
	if err := os.WriteFile(tokenFilePath, raw, 0600); err != nil {
		return fmt.Errorf("write %s: %w", tokenFilePath, err)
	}
	return nil
}
