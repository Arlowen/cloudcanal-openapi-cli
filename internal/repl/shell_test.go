package repl

import (
	"cloudcanal-openapi-cli/internal/app"
	"cloudcanal-openapi-cli/internal/config"
	"cloudcanal-openapi-cli/internal/console"
	"cloudcanal-openapi-cli/internal/datajob"
	"cloudcanal-openapi-cli/internal/testsupport"
	"strings"
	"testing"
)

func TestShellHandlesHappyPathCommands(t *testing.T) {
	service := &fakeDataJobs{
		jobs: []datajob.Job{
			{
				DataJobID:     11,
				DataJobName:   "sync-job",
				DataJobType:   "SYNC",
				DataTaskState: "RUNNING",
				SourceDS:      &datajob.Source{InstanceDesc: "src-db"},
				TargetDS:      &datajob.Source{InstanceDesc: "dst-db"},
			},
		},
		job: datajob.Job{
			DataJobID:      11,
			DataJobName:    "sync-job",
			DataJobDesc:    "nightly sync",
			DataJobType:    "SYNC",
			DataTaskState:  "RUNNING",
			CurrTaskStatus: "FULL_RUNNING",
			LifeCycleState: "ACTIVE",
			UserName:       "admin",
			ConsoleJobID:   21,
			SourceDS: &datajob.Source{
				InstanceDesc:   "src-db",
				DataSourceType: "MYSQL",
			},
			TargetDS: &datajob.Source{
				InstanceDesc:   "dst-db",
				DataSourceType: "STARROCKS",
			},
			DataTasks: []datajob.Task{
				{DataTaskID: 101, DataTaskName: "full-task", DataTaskType: "FULL", DataTaskStatus: "RUNNING"},
			},
		},
	}
	runtime := &fakeRuntime{
		cfg:      config.AppConfig{APIBaseURL: "https://cc.example.com", AccessKey: "abcdefghijkl", SecretKey: "qrstuvwxyz1234"},
		dataJobs: service,
	}
	io := testsupport.NewTestConsole("help", "jobs list", "jobs show 11", "jobs replay 11 --auto-start --reset-to-created", "jobs start 11", "jobs stop 11", "jobs delete 11", "config show", "exit")

	shell := NewShell(io, runtime)
	if err := shell.Run(); err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	out := io.Output()
	if service.lastStartedJobID != 11 || service.lastStoppedJobID != 11 || service.lastDeletedJobID != 11 || service.lastReplayedJobID != 11 {
		t.Fatalf("unexpected job actions: start=%d stop=%d delete=%d replay=%d", service.lastStartedJobID, service.lastStoppedJobID, service.lastDeletedJobID, service.lastReplayedJobID)
	}
	if !service.lastReplayOptions.AutoStart || !service.lastReplayOptions.ResetToCreated {
		t.Fatalf("unexpected replay options: %+v", service.lastReplayOptions)
	}
	for _, want := range []string{
		"Available commands:",
		"sync-job",
		"Job details:",
		"nightly sync",
		"Job 11 replay requested successfully",
		"Job 11 started successfully",
		"Job 11 stopped successfully",
		"Job 11 deleted successfully",
		"apiBaseUrl: https://cc.example.com",
		"accessKey: abcd****ijkl",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("output missing %q in %q", want, out)
		}
	}
	if strings.Contains(out, "secretKey:") {
		t.Fatalf("output unexpectedly contains secret key line: %q", out)
	}
}

func TestShellReportsInvalidCommandsWithoutExiting(t *testing.T) {
	service := &fakeDataJobs{}
	runtime := &fakeRuntime{
		cfg:               config.AppConfig{APIBaseURL: "https://cc.example.com", AccessKey: "abcdefghijkl", SecretKey: "qrstuvwxyz1234"},
		dataJobs:          service,
		reinitializeValue: true,
	}
	io := testsupport.NewTestConsole("jobs start abc", "jobs replay 11 --bad", "unknown", "config init", "exit")

	shell := NewShell(io, runtime)
	if err := shell.Run(); err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	out := io.Output()
	if runtime.reinitializeCalls != 1 {
		t.Fatalf("reinitializeCalls = %d, want 1", runtime.reinitializeCalls)
	}
	for _, want := range []string{
		"jobId must be a positive integer",
		"unknown replay option: --bad",
		"Unknown command: unknown",
		"Configuration updated.",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("output missing %q in %q", want, out)
		}
	}
}

func TestShellExecutesArgsWithoutInteractiveLoop(t *testing.T) {
	service := &fakeDataJobs{
		jobs: []datajob.Job{
			{
				DataJobID:     22,
				DataJobName:   "batch-job",
				DataJobType:   "CHECK",
				DataTaskState: "CREATED",
				SourceDS:      &datajob.Source{InstanceDesc: "src-check"},
				TargetDS:      &datajob.Source{InstanceDesc: "dst-check"},
			},
		},
	}
	runtime := &fakeRuntime{
		cfg:      config.AppConfig{APIBaseURL: "https://cc.example.com", AccessKey: "abcdefghijkl", SecretKey: "qrstuvwxyz1234"},
		dataJobs: service,
	}
	io := testsupport.NewTestConsole()

	shell := NewShell(io, runtime)
	if err := shell.ExecuteArgs([]string{"jobs", "list"}); err != nil {
		t.Fatalf("ExecuteArgs() error = %v", err)
	}

	out := io.Output()
	for _, want := range []string{
		"batch-job",
		"1 jobs",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("output missing %q in %q", want, out)
		}
	}
	if strings.Contains(out, "cloudcanal>") || strings.Contains(out, "Type 'help'") {
		t.Fatalf("output unexpectedly contains interactive text: %q", out)
	}
}

type fakeRuntime struct {
	cfg               config.AppConfig
	dataJobs          datajob.Operations
	reinitializeCalls int
	reinitializeValue bool
}

func (f *fakeRuntime) Config() config.AppConfig {
	return f.cfg
}

func (f *fakeRuntime) DataJobs() datajob.Operations {
	return f.dataJobs
}

func (f *fakeRuntime) Reinitialize(io console.IO) (bool, error) {
	f.reinitializeCalls++
	return f.reinitializeValue, nil
}

var _ app.RuntimeContext = (*fakeRuntime)(nil)

type fakeDataJobs struct {
	jobs              []datajob.Job
	job               datajob.Job
	lastStartedJobID  int64
	lastStoppedJobID  int64
	lastDeletedJobID  int64
	lastReplayedJobID int64
	lastReplayOptions datajob.ReplayOptions
}

func (f *fakeDataJobs) ListJobs() ([]datajob.Job, error) {
	return f.jobs, nil
}

func (f *fakeDataJobs) GetJob(jobID int64) (datajob.Job, error) {
	return f.job, nil
}

func (f *fakeDataJobs) StartJob(jobID int64) error {
	f.lastStartedJobID = jobID
	return nil
}

func (f *fakeDataJobs) StopJob(jobID int64) error {
	f.lastStoppedJobID = jobID
	return nil
}

func (f *fakeDataJobs) DeleteJob(jobID int64) error {
	f.lastDeletedJobID = jobID
	return nil
}

func (f *fakeDataJobs) ReplayJob(jobID int64, options datajob.ReplayOptions) error {
	f.lastReplayedJobID = jobID
	f.lastReplayOptions = options
	return nil
}
