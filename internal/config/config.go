package config

import (
	"cloudcanal-openapi-cli/internal/i18n"
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
	Language   string `json:"language,omitempty"`
}

func (c AppConfig) Validate() error {
	language := c.NormalizedLanguage()
	if strings.TrimSpace(c.APIBaseURL) == "" {
		return errors.New(i18n.TFor(language, "config.apiBaseUrlRequired"))
	}
	if strings.TrimSpace(c.AccessKey) == "" {
		return errors.New(i18n.TFor(language, "config.accessKeyRequired"))
	}
	if strings.TrimSpace(c.SecretKey) == "" {
		return errors.New(i18n.TFor(language, "config.secretKeyRequired"))
	}
	if normalized := i18n.NormalizeLanguage(c.Language); normalized == "" && strings.TrimSpace(c.Language) != "" {
		return errors.New(i18n.T("config.languageUnsupported"))
	}

	parsed, err := url.Parse(strings.TrimSpace(c.APIBaseURL))
	if err != nil {
		return errors.New(i18n.TFor(language, "config.apiBaseUrlInvalid"))
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return errors.New(i18n.TFor(language, "config.apiBaseUrlScheme"))
	}
	if strings.TrimSpace(parsed.Host) == "" {
		return errors.New(i18n.TFor(language, "config.apiBaseUrlHost"))
	}
	return nil
}

func (c AppConfig) NormalizedBaseURL() string {
	value := strings.TrimSpace(c.APIBaseURL)
	return strings.TrimRight(value, "/")
}

func (c AppConfig) NormalizedLanguage() string {
	normalized := i18n.NormalizeLanguage(c.Language)
	if normalized == "" {
		return i18n.DefaultLanguage()
	}
	return normalized
}

func (c AppConfig) WithDefaults() AppConfig {
	c.Language = c.NormalizedLanguage()
	return c
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
		return AppConfig{}, errors.New(i18n.T("config.invalidJSON"))
	}
	cfg = cfg.WithDefaults()
	if err := cfg.Validate(); err != nil {
		return AppConfig{}, err
	}
	return cfg, nil
}

func (s *Service) Save(cfg AppConfig) error {
	cfg = cfg.WithDefaults()
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
