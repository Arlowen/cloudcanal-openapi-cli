package repl

import (
	"cloudcanal-openapi-cli/internal/app"
	"cloudcanal-openapi-cli/internal/console"
	"cloudcanal-openapi-cli/internal/i18n"
	"cloudcanal-openapi-cli/internal/util"
	"io"
	"strings"
)

const prompt = "cloudcanal> "

type Shell struct {
	io      console.IO
	runtime app.RuntimeContext
}

func NewShell(io console.IO, runtime app.RuntimeContext) *Shell {
	_ = i18n.SetLanguage(runtime.Config().NormalizedLanguage())
	return &Shell{io: io, runtime: runtime}
}

func (s *Shell) ExecuteArgs(args []string) error {
	if len(args) == 0 {
		return nil
	}
	commandLine := strings.Join(args, " ")
	return s.handleTokens(args, commandLine)
}

func (s *Shell) Run() error {
	s.io.Println(i18n.T("common.typeHelp"))
	for {
		line, err := s.io.ReadLine(prompt)
		if err != nil {
			if err == io.EOF {
				s.io.Println("")
				return nil
			}
			return err
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.EqualFold(line, "exit") || strings.EqualFold(line, "quit") {
			return nil
		}

		if err := s.handle(line); err != nil {
			s.io.Println(i18n.T("common.errorPrefix", util.SummarizeError(err)))
		}
	}
}

func (s *Shell) handle(commandLine string) error {
	tokens, err := splitCommandLine(commandLine)
	if err != nil {
		return err
	}
	return s.handleTokens(tokens, commandLine)
}

func (s *Shell) handleTokens(tokens []string, commandLine string) error {
	if len(tokens) == 0 {
		return nil
	}

	switch strings.ToLower(tokens[0]) {
	case "help":
		s.printHelp(tokens[1:])
		return nil
	case "clear", "cls":
		s.io.ClearScreen()
		return nil
	case "jobs":
		return s.handleJobs(tokens)
	case "datasources":
		return s.handleDataSources(tokens)
	case "clusters":
		return s.handleClusters(tokens)
	case "workers":
		return s.handleWorkers(tokens)
	case "consolejobs":
		return s.handleConsoleJobs(tokens)
	case "job-config", "jobconfig":
		return s.handleJobConfig(tokens)
	case "config":
		return s.handleConfig(tokens)
	case "lang", "language":
		return s.handleLang(tokens)
	default:
		s.io.Println(i18n.T("common.unknownCommand", commandLine))
		s.io.Println(i18n.T("common.useHelp"))
		return nil
	}
}

func (s *Shell) handleConfig(tokens []string) error {
	if len(tokens) != 2 {
		s.io.Println(s.usageConfig())
		return nil
	}
	switch strings.ToLower(tokens[1]) {
	case "show":
		cfg := s.runtime.Config()
		s.io.Println(i18n.T("config.apiBaseUrlLabel") + ": " + cfg.APIBaseURL)
		s.io.Println(i18n.T("config.accessKeyLabel") + ": " + util.MaskSecret(cfg.AccessKey))
		s.io.Println(i18n.T("config.languageLabel") + ": " + cfg.NormalizedLanguage())
		return nil
	case "init":
		updated, err := s.runtime.Reinitialize(s.io)
		if err != nil {
			return err
		}
		if updated {
			s.io.Println(i18n.T("runtime.configUpdated"))
		}
		return nil
	default:
		s.io.Println(s.usageConfig())
		return nil
	}
}

func (s *Shell) handleLang(tokens []string) error {
	if len(tokens) == 1 || strings.EqualFold(tokens[1], "show") {
		if len(tokens) > 2 {
			s.io.Println(i18n.T("lang.usage"))
			return nil
		}
		s.io.Println(i18n.T("lang.current", s.runtime.Config().NormalizedLanguage()))
		s.io.Println(i18n.T("common.supportedLanguages"))
		return nil
	}
	if strings.EqualFold(tokens[1], "set") {
		if len(tokens) != 3 {
			s.io.Println(i18n.T("lang.usage"))
			return nil
		}
		if err := s.runtime.SetLanguage(tokens[2]); err != nil {
			return err
		}
		s.io.Println(i18n.T("lang.updated", i18n.DisplayName(s.runtime.Config().NormalizedLanguage())))
		return nil
	}
	s.io.Println(i18n.T("lang.usage"))
	return nil
}
