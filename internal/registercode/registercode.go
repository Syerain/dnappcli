package registercode

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Syerain/dnappcli/internal/utils"
)

// Payload 与 server internal/model/RegistercodePayload 结构一致。
type Payload struct {
	Magicword string
	Before    time.Time
}

// Signer 使用 ed25519 私钥签发注册码。
type Signer struct {
	enckey ed25519.PrivateKey
}

// NewSigner 从十六进制私钥创建 Signer。
func NewSigner(enckeyHex string) (*Signer, error) {
	key, err := utils.HexToPrivKey(enckeyHex)
	if err != nil {
		return nil, fmt.Errorf("registercodeEnckey 不是合法的 ed25519 私钥 hex: %w", err)
	}
	return &Signer{enckey: key}, nil
}

// Sign 生成与 server 端 SignRegistercode 兼容的注册码。
// 格式：hex(json(payload)).hex(ed25519sig)
func (s *Signer) Sign(payload Payload) (string, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshal payload: %w", err)
	}
	payloadHex := hex.EncodeToString(payloadBytes)

	sig := ed25519.Sign(s.enckey, payloadBytes)
	sigHex := hex.EncodeToString(sig)

	return payloadHex + "." + sigHex, nil
}
