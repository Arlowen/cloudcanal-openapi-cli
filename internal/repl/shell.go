package repl

import (
	"cloudcanal-openapi-cli/internal/app"
	"cloudcanal-openapi-cli/internal/console"
	"cloudcanal-openapi-cli/internal/datajob"
	"cloudcanal-openapi-cli/internal/util"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const prompt = "cloudcanal> "

type Shell struct {
	io      console.IO
	runtime app.RuntimeContext
}

func NewShell(io console.IO, runtime app.RuntimeContext) *Shell {
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
	s.io.Println("Type 'help' to see available commands.")
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
			s.io.Println("Error: " + util.SummarizeError(err))
		}
	}
}

func (s *Shell) handle(commandLine string) error {
	tokens := strings.Fields(commandLine)
	return s.handleTokens(tokens, commandLine)
}

func (s *Shell) handleTokens(tokens []string, commandLine string) error {
	if len(tokens) == 0 {
		return nil
	}

	switch strings.ToLower(tokens[0]) {
	case "help":
		s.printHelp()
		return nil
	case "jobs":
		return s.handleJobs(tokens)
	case "config":
		return s.handleConfig(tokens)
	default:
		s.io.Println("Unknown command: " + commandLine)
		s.io.Println("Use 'help' to see available commands.")
		return nil
	}
}

func (s *Shell) handleJobs(tokens []string) error {
	if len(tokens) < 2 {
		s.io.Println("Usage: jobs list | jobs show <jobId> | jobs start <jobId> | jobs stop <jobId> | jobs delete <jobId> | jobs replay <jobId> [--auto-start] [--reset-to-created]")
		return nil
	}

	switch strings.ToLower(tokens[1]) {
	case "list":
		if len(tokens) != 2 {
			s.io.Println("Usage: jobs list")
			return nil
		}
		return s.printJobs()
	case "show":
		if len(tokens) != 3 {
			s.io.Println("Usage: jobs show <jobId>")
			return nil
		}
		jobID, err := parseJobID(tokens[2])
		if err != nil {
			return err
		}
		return s.printJob(jobID)
	case "start", "stop", "delete":
		if len(tokens) != 3 {
			s.io.Println("Usage: jobs " + strings.ToLower(tokens[1]) + " <jobId>")
			return nil
		}
		jobID, err := parseJobID(tokens[2])
		if err != nil {
			return err
		}
		switch strings.ToLower(tokens[1]) {
		case "start":
			if err := s.runtime.DataJobs().StartJob(jobID); err != nil {
				return err
			}
			s.io.Println(fmt.Sprintf("Job %d started successfully", jobID))
		case "stop":
			if err := s.runtime.DataJobs().StopJob(jobID); err != nil {
				return err
			}
			s.io.Println(fmt.Sprintf("Job %d stopped successfully", jobID))
		default:
			if err := s.runtime.DataJobs().DeleteJob(jobID); err != nil {
				return err
			}
			s.io.Println(fmt.Sprintf("Job %d deleted successfully", jobID))
		}
		return nil
	case "replay":
		if len(tokens) < 3 {
			s.io.Println("Usage: jobs replay <jobId> [--auto-start] [--reset-to-created]")
			return nil
		}
		jobID, err := parseJobID(tokens[2])
		if err != nil {
			return err
		}
		options, err := parseReplayOptions(tokens[3:])
		if err != nil {
			return err
		}
		if err := s.runtime.DataJobs().ReplayJob(jobID, options); err != nil {
			return err
		}
		s.io.Println(fmt.Sprintf("Job %d replay requested successfully", jobID))
		return nil
	default:
		s.io.Println("Usage: jobs list | jobs show <jobId> | jobs start <jobId> | jobs stop <jobId> | jobs delete <jobId> | jobs replay <jobId> [--auto-start] [--reset-to-created]")
		return nil
	}
}

func (s *Shell) handleConfig(tokens []string) error {
	if len(tokens) != 2 {
		s.io.Println("Usage: config show | config init")
		return nil
	}
	switch strings.ToLower(tokens[1]) {
	case "show":
		cfg := s.runtime.Config()
		s.io.Println("apiBaseUrl: " + cfg.APIBaseURL)
		s.io.Println("accessKey: " + util.MaskSecret(cfg.AccessKey))
		return nil
	case "init":
		updated, err := s.runtime.Reinitialize(s.io)
		if err != nil {
			return err
		}
		if updated {
			s.io.Println("Configuration updated.")
		}
		return nil
	default:
		s.io.Println("Usage: config show | config init")
		return nil
	}
}

func (s *Shell) printHelp() {
	s.io.Println("Available commands:")
	s.io.Println("  jobs list")
	s.io.Println("  jobs show <jobId>")
	s.io.Println("  jobs start <jobId>")
	s.io.Println("  jobs stop <jobId>")
	s.io.Println("  jobs delete <jobId>")
	s.io.Println("  jobs replay <jobId> [--auto-start] [--reset-to-created]")
	s.io.Println("  config show")
	s.io.Println("  config init")
	s.io.Println("  help")
	s.io.Println("  exit | quit")
}

func (s *Shell) printJobs() error {
	jobs, err := s.runtime.DataJobs().ListJobs()
	if err != nil {
		return err
	}

	headers := []string{"ID", "Name", "Type", "State", "Source", "Target"}
	rows := make([][]string, 0, len(jobs))
	for _, job := range jobs {
		rows = append(rows, []string{
			strconv.FormatInt(job.DataJobID, 10),
			orDash(job.DataJobName),
			orDash(job.DataJobType),
			orDash(job.DataTaskState),
			instanceDesc(job.SourceDS),
			instanceDesc(job.TargetDS),
		})
	}

	s.io.Println(util.FormatTable(headers, rows))
	s.io.Println(fmt.Sprintf("%d jobs", len(jobs)))
	return nil
}

func (s *Shell) printJob(jobID int64) error {
	job, err := s.runtime.DataJobs().GetJob(jobID)
	if err != nil {
		return err
	}

	s.io.Println("Job details:")
	s.io.Println("  ID: " + strconv.FormatInt(job.DataJobID, 10))
	s.io.Println("  Name: " + orDash(job.DataJobName))
	s.io.Println("  Description: " + orDash(job.DataJobDesc))
	s.io.Println("  Type: " + orDash(job.DataJobType))
	s.io.Println("  State: " + orDash(job.DataTaskState))
	s.io.Println("  Current Task Status: " + orDash(job.CurrTaskStatus))
	s.io.Println("  Lifecycle: " + orDash(job.LifeCycleState))
	s.io.Println("  User: " + orDash(job.UserName))
	s.io.Println("  Console Job ID: " + formatOptionalInt64(job.ConsoleJobID))
	s.io.Println("  Console Task State: " + orDash(job.ConsoleTaskState))
	s.io.Println("  Source: " + sourceSummary(job.SourceDS))
	s.io.Println("  Target: " + sourceSummary(job.TargetDS))
	s.io.Println("  Source Schema: " + orDash(job.SourceSchema))
	s.io.Println("  Target Schema: " + orDash(job.TargetSchema))
	s.io.Println("  Tasks: " + strconv.Itoa(len(job.DataTasks)))
	s.io.Println("  Has Exception: " + formatBool(job.HaveException))

	if len(job.DataTasks) > 0 {
		headers := []string{"Task ID", "Name", "Type", "Status", "Worker IP"}
		rows := make([][]string, 0, len(job.DataTasks))
		for _, task := range job.DataTasks {
			rows = append(rows, []string{
				strconv.FormatInt(task.DataTaskID, 10),
				orDash(task.DataTaskName),
				orDash(task.DataTaskType),
				orDash(task.DataTaskStatus),
				orDash(task.WorkerIP),
			})
		}
		s.io.Println("")
		s.io.Println(util.FormatTable(headers, rows))
	}
	return nil
}

func parseJobID(value string) (int64, error) {
	jobID, err := strconv.ParseInt(value, 10, 64)
	if err != nil || jobID <= 0 {
		return 0, fmt.Errorf("jobId must be a positive integer")
	}
	return jobID, nil
}

func parseReplayOptions(tokens []string) (datajob.ReplayOptions, error) {
	var options datajob.ReplayOptions
	for _, token := range tokens {
		switch strings.ToLower(token) {
		case "--auto-start":
			options.AutoStart = true
		case "--reset-to-created":
			options.ResetToCreated = true
		default:
			return datajob.ReplayOptions{}, fmt.Errorf("unknown replay option: %s", token)
		}
	}
	return options, nil
}

func instanceDesc(source *datajob.Source) string {
	if source == nil {
		return "-"
	}
	return orDash(source.InstanceDesc)
}

func sourceSummary(source *datajob.Source) string {
	if source == nil {
		return "-"
	}

	extras := make([]string, 0, 3)
	if strings.TrimSpace(source.DataSourceType) != "" {
		extras = append(extras, source.DataSourceType)
	}
	if strings.TrimSpace(source.HostType) != "" {
		extras = append(extras, source.HostType)
	}
	if strings.TrimSpace(source.Region) != "" {
		extras = append(extras, source.Region)
	}

	label := orDash(source.InstanceDesc)
	if len(extras) == 0 {
		return label
	}
	return label + " (" + strings.Join(extras, ", ") + ")"
}

func formatOptionalInt64(value int64) string {
	if value == 0 {
		return "-"
	}
	return strconv.FormatInt(value, 10)
}

func formatBool(value bool) string {
	if value {
		return "true"
	}
	return "false"
}

func orDash(value string) string {
	if strings.TrimSpace(value) == "" {
		return "-"
	}
	return value
}
