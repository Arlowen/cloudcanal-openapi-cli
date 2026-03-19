package config

import (
	"cloudcanal-openapi-cli/internal/console"
	"cloudcanal-openapi-cli/internal/i18n"
	"errors"
	"io"
	"strings"
)

type Validator func(AppConfig) error

type Wizard struct {
	io        console.IO
	service   *Service
	validator Validator
	initial   AppConfig
}

func NewWizard(io console.IO, service *Service, validator Validator, initial AppConfig) *Wizard {
	return &Wizard{
		io:        io,
		service:   service,
		validator: validator,
		initial:   initial,
	}
}

func (w *Wizard) Run() (*AppConfig, error) {
	current := w.initial.WithDefaults()
	_ = i18n.SetLanguage(current.Language)
	w.io.Println(i18n.T("wizard.title"))
	w.io.Println(i18n.T("wizard.cancelHint"))
	w.io.Println(i18n.T("wizard.languageHint"))
	w.io.Println(i18n.T("wizard.apiHostHint"))
	if w.hasInitialValue(w.initial) {
		w.io.Println(i18n.T("wizard.keepCurrent"))
	}

	for {
		language, cancelled, err := w.promptLanguage(current.Language)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil, nil
			}
			return nil, err
		}
		if cancelled {
			return nil, nil
		}
		current.Language = language
		_ = i18n.SetLanguage(language)

		apiBaseURL, cancelled, err := w.promptRequired("apiHost", current.APIBaseURL, validateAPIBaseURL)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil, nil
			}
			return nil, err
		}
		if cancelled {
			return nil, nil
		}

		accessKey, cancelled, err := w.promptRequired("ak", current.AccessKey, validateAccessKey)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil, nil
			}
			return nil, err
		}
		if cancelled {
			return nil, nil
		}

		secretKey, cancelled, err := w.promptSecret("sk", current.SecretKey)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil, nil
			}
			return nil, err
		}
		if cancelled {
			return nil, nil
		}

		current = AppConfig{
			APIBaseURL: apiBaseURL,
			AccessKey:  accessKey,
			SecretKey:  secretKey,
			Language:   language,
		}

		if err := current.Validate(); err != nil {
			w.io.Println(i18n.T("wizard.invalidConfig", err.Error()))
			continue
		}
		w.io.Println(i18n.T("wizard.checkingConnection"))
		if err := w.validator(current); err != nil {
			w.io.Println(i18n.T("wizard.validationFailed", err.Error()))
			w.io.Println(i18n.T("wizard.reuseValues"))
			continue
		}
		if err := w.service.Save(current); err != nil {
			return nil, err
		}
		w.io.Println(i18n.T("wizard.savedTo", w.service.Path()))
		return &current, nil
	}
}

func (w *Wizard) promptRequired(label, current string, validate func(string) error) (string, bool, error) {
	for {
		value, cancelled, err := w.promptWithDefault(label, current)
		if err != nil || cancelled {
			return "", cancelled, err
		}
		if err := validate(value); err != nil {
			w.io.Println(i18n.T("wizard.invalidField", label, err.Error()))
			continue
		}
		return value, false, nil
	}
}

func (w *Wizard) promptSecret(label, current string) (string, bool, error) {
	for {
		prompt := label + ": "
		if strings.TrimSpace(current) != "" {
			prompt = label + " [hidden]: "
		}

		value, err := w.io.ReadSecret(prompt)
		if err != nil {
			return "", false, err
		}
		trimmed := strings.TrimSpace(value)
		if strings.EqualFold(trimmed, "exit") {
			return "", true, nil
		}
		if trimmed == "" {
			if strings.TrimSpace(current) == "" {
				w.io.Println(i18n.T("wizard.invalidField", label, i18n.T("config.secretKeyRequired")))
				continue
			}
			return current, false, nil
		}
		return trimmed, false, nil
	}
}

func (w *Wizard) promptWithDefault(label, current string) (string, bool, error) {
	prompt := label + ": "
	if strings.TrimSpace(current) != "" {
		prompt = label + " [" + current + "]: "
	}

	value, err := w.io.ReadLine(prompt)
	if err != nil {
		return "", false, err
	}
	trimmed := strings.TrimSpace(value)
	if strings.EqualFold(trimmed, "exit") {
		return "", true, nil
	}
	if trimmed == "" {
		return strings.TrimSpace(current), false, nil
	}
	return trimmed, false, nil
}

func (w *Wizard) promptLanguage(current string) (string, bool, error) {
	for {
		value, cancelled, err := w.promptWithDefault("language", current)
		if err != nil || cancelled {
			return "", cancelled, err
		}
		normalized := i18n.NormalizeLanguage(value)
		if normalized == "" {
			w.io.Println(i18n.T("wizard.invalidField", "language", i18n.T("config.languageUnsupported")))
			continue
		}
		return normalized, false, nil
	}
}

func (w *Wizard) hasInitialValue(cfg AppConfig) bool {
	return strings.TrimSpace(cfg.APIBaseURL) != "" ||
		strings.TrimSpace(cfg.AccessKey) != "" ||
		strings.TrimSpace(cfg.SecretKey) != "" ||
		strings.TrimSpace(cfg.Language) != ""
}

func validateAPIBaseURL(value string) error {
	return AppConfig{APIBaseURL: value, AccessKey: "ak", SecretKey: "sk"}.Validate()
}

func validateAccessKey(value string) error {
	return AppConfig{APIBaseURL: "https://cc.example.com", AccessKey: value, SecretKey: "sk"}.Validate()
}
