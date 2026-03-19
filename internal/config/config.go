package config

import (
	"encoding/json"
	"errors"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type AppConfig struct {
	APIBaseURL string `json:"apiBaseUrl"`
	AccessKey  string `json:"accessKey"`
	SecretKey  string `json:"secretKey"`
}

func (c AppConfig) Validate() error {
	if strings.TrimSpace(c.APIBaseURL) == "" {
		return errors.New("apiBaseUrl is required")
	}
	if strings.TrimSpace(c.AccessKey) == "" {
		return errors.New("accessKey is required")
	}
	if strings.TrimSpace(c.SecretKey) == "" {
		return errors.New("secretKey is required")
	}

	parsed, err := url.Parse(strings.TrimSpace(c.APIBaseURL))
	if err != nil {
		return errors.New("apiBaseUrl is not a valid URL")
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return errors.New("apiBaseUrl must start with http:// or https://")
	}
	if strings.TrimSpace(parsed.Host) == "" {
		return errors.New("apiBaseUrl must contain a host")
	}
	return nil
}

func (c AppConfig) NormalizedBaseURL() string {
	value := strings.TrimSpace(c.APIBaseURL)
	return strings.TrimRight(value, "/")
}

type Service struct {
	path string
}

func NewService(path string) *Service {
	if strings.TrimSpace(path) == "" {
		path = DefaultPath()
	}
	return &Service{path: path}
}

func DefaultPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".cloudcanal/config.json"
	}
	return filepath.Join(home, ".cloudcanal", "config.json")
}

func (s *Service) Path() string {
	return s.path
}

func (s *Service) Exists() bool {
	_, err := os.Stat(s.path)
	return err == nil
}

func (s *Service) Load() (AppConfig, error) {
	data, err := os.ReadFile(s.path)
	if err != nil {
		return AppConfig{}, err
	}

	var cfg AppConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return AppConfig{}, errors.New("configuration file is not valid JSON")
	}
	if err := cfg.Validate(); err != nil {
		return AppConfig{}, err
	}
	return cfg, nil
}

func (s *Service) Save(cfg AppConfig) error {
	if err := cfg.Validate(); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}
	content, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, content, 0o600)
}
