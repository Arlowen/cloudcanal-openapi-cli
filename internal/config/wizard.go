package config

import (
	"cloudcanal-openapi-cli/internal/console"
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
	w.io.Println("CloudCanal CLI initialization")
	w.io.Println("Type exit at any prompt to cancel.")
	w.io.Println("apiHost must be a full URL, for example: https://cc.example.com")
	if w.initial.APIBaseURL != "" || w.initial.AccessKey != "" || w.initial.SecretKey != "" {
		w.io.Println("Press Enter to keep the current value.")
	}

	current := w.initial

	for {
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
		}

		if err := current.Validate(); err != nil {
			w.io.Println("Invalid configuration: " + err.Error())
			continue
		}
		w.io.Println("Checking OpenAPI connection...")
		if err := w.validator(current); err != nil {
			w.io.Println("Configuration validation failed: " + err.Error())
			w.io.Println("Press Enter to reuse the current values, or type new ones.")
			continue
		}
		if err := w.service.Save(current); err != nil {
			return nil, err
		}
		w.io.Println("Configuration saved to " + w.service.Path())
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
			w.io.Println("Invalid " + label + ": " + err.Error())
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
				w.io.Println("Invalid " + label + ": secretKey is required")
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

func validateAPIBaseURL(value string) error {
	return AppConfig{APIBaseURL: value, AccessKey: "ak", SecretKey: "sk"}.Validate()
}

func validateAccessKey(value string) error {
	return AppConfig{APIBaseURL: "https://cc.example.com", AccessKey: value, SecretKey: "sk"}.Validate()
}
