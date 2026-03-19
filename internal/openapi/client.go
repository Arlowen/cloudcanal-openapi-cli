package openapi

import (
	"bytes"
	"cloudcanal-openapi-cli/internal/config"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type ServerError struct {
	StatusCode   int
	ResponseBody string
}

func (e *ServerError) Error() string {
	return fmt.Sprintf("HTTP %d: %s", e.StatusCode, strings.TrimSpace(e.ResponseBody))
}

type Client struct {
	config     config.AppConfig
	httpClient *http.Client
}

func NewClient(cfg config.AppConfig) (*Client, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return &Client{
		config: cfg,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}, nil
}

func NewClientWithHTTP(cfg config.AppConfig, httpClient *http.Client) (*Client, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	if httpClient == nil {
		return nil, errors.New("http client is required")
	}
	return &Client{config: cfg, httpClient: httpClient}, nil
}

func (c *Client) PostJSON(path string, payload any, out any) error {
	bodyValue := payload
	if bodyValue == nil {
		bodyValue = map[string]any{}
	}

	bodyBytes, err := json.Marshal(bodyValue)
	if err != nil {
		return fmt.Errorf("failed to encode request body: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, c.fullURL(path), bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("failed to build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call OpenAPI: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read OpenAPI response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &ServerError{StatusCode: resp.StatusCode, ResponseBody: string(responseBody)}
	}
	if out == nil {
		return nil
	}
	if err := json.Unmarshal(responseBody, out); err != nil {
		return fmt.Errorf("failed to parse OpenAPI response: %w", err)
	}
	return nil
}

func (c *Client) fullURL(path string) string {
	params := c.commonParams()
	return c.config.NormalizedBaseURL() + path + "?" + GenSortedParamsStr(params)
}

func (c *Client) commonParams() map[string]string {
	params := map[string]string{
		"SignatureMethod": "HmacSHA1",
		"SignatureNonce":  randomNonce(),
		"AccessKeyId":     c.config.AccessKey,
	}
	signature := SignString(ComposeStringToSign(params), c.config.SecretKey)
	params["Signature"] = signature
	return params
}

func randomNonce() string {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(buf)
}
