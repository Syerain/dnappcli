package utils

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
)

// HexToPrivKey 将十六进制字符串解码为 ed25519 私钥。
// 支持两种情况：
//   - 64 字节：标准 ed25519 私钥（seed + 公钥）
//   - 32 字节：仅 seed，自动通过 ed25519.NewKeyFromSeed 展开
func HexToPrivKey(s string) (ed25519.PrivateKey, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}
	switch len(b) {
	case ed25519.SeedSize:
		return ed25519.NewKeyFromSeed(b), nil
	case ed25519.PrivateKeySize:
		return ed25519.PrivateKey(b), nil
	default:
		return nil, fmt.Errorf("非法 ed25519 私钥长度: %d 字节", len(b))
	}
}

