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

// DoJSON 是通用请求入口：method 为 http.MethodGet/Post 等；
// path 如 "/api/v1/user/me"；token 非空时自动带 Authorization: Bearer <token>；
// body 为 nil 时不带请求体，否则序列化为 JSON。
// 返回 HTTP 状态码与原始响应体字节。
func (c *Client) DoJSON(ctx context.Context, method, path, token string, body any) (statusCode int, respBody []byte, err error) {
	var reader io.Reader
	if body != nil {
		payload, err := json.Marshal(body)
		if err != nil {
			return 0, nil, fmt.Errorf("marshal request body: %w", err)
		}
		reader = bytes.NewReader(payload)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reader)
	if err != nil {
		return 0, nil, fmt.Errorf("build request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

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

// PostJSON 向 path 发送带可选 token 的 JSON POST 请求。
func (c *Client) PostJSON(ctx context.Context, path, token string, body any) (statusCode int, respBody []byte, err error) {
	return c.DoJSON(ctx, http.MethodPost, path, token, body)
}

// GetJSON 向 path 发送带可选 token 的 GET 请求。
func (c *Client) GetJSON(ctx context.Context, path, token string) (statusCode int, respBody []byte, err error) {
	return c.DoJSON(ctx, http.MethodGet, path, token, nil)
}
