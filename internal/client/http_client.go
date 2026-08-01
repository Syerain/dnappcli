package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client 是 dnappcli 发起 HTTP 请求的封装。
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// New 构造一个 HTTP 客户端，baseURL 形如 "http://127.0.0.1:4703"。
func New(baseURL string) *Client {
	return &Client{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// PostJSON 向 path 发送 JSON POST 请求，body 会序列化为 JSON。
// 返回 HTTP 状态码与原始响应体字节。
func (c *Client) PostJSON(ctx context.Context, path string, body any) (statusCode int, respBody []byte, err error) {
	payload, err := json.Marshal(body)
	if err != nil {
		return 0, nil, fmt.Errorf("marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, bytes.NewReader(payload))
	if err != nil {
		return 0, nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, nil, fmt.Errorf("read response body: %w", err)
	}

	return resp.StatusCode, data, nil
}
