package datajob

import (
	"cloudcanal-openapi-cli/internal/config"
	"cloudcanal-openapi-cli/internal/openapi"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServiceListsJobs(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"code":"1","data":[{"dataJobId":11,"dataJobName":"sync-1","dataJobType":"SYNC","dataTaskState":"RUNNING"}]}`))
	}))
	defer server.Close()

	client, err := openapi.NewClient(config.AppConfig{
		APIBaseURL: server.URL,
		AccessKey:  "test-ak",
		SecretKey:  "test-sk",
	})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	service := NewService(client)
	jobs, err := service.ListJobs()
	if err != nil {
		t.Fatalf("ListJobs() error = %v", err)
	}
	if len(jobs) != 1 || jobs[0].DataJobID != 11 || jobs[0].DataJobName != "sync-1" {
		t.Fatalf("jobs = %#v, want single sync-1 job", jobs)
	}
}

func TestServiceGetsJobDetails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"code":"1","data":{"dataJobId":11,"dataJobName":"sync-1","dataJobDesc":"nightly sync","dataJobType":"SYNC","dataTaskState":"RUNNING","currTaskStatus":"FULL_RUNNING","consoleJobId":21,"lifeCycleState":"ACTIVE","sourceDsVO":{"instanceDesc":"src-db","dataSourceType":"MYSQL"},"targetDsVO":{"instanceDesc":"dst-db","dataSourceType":"STARROCKS"},"dataTasks":[{"dataTaskId":101,"dataTaskName":"full-task","dataTaskStatus":"RUNNING"}]}}`))
	}))
	defer server.Close()

	client, err := openapi.NewClient(config.AppConfig{
		APIBaseURL: server.URL,
		AccessKey:  "test-ak",
		SecretKey:  "test-sk",
	})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	service := NewService(client)
	job, err := service.GetJob(11)
	if err != nil {
		t.Fatalf("GetJob() error = %v", err)
	}
	if job.DataJobID != 11 || job.DataJobDesc != "nightly sync" || len(job.DataTasks) != 1 {
		t.Fatalf("job = %#v, want detailed job with one task", job)
	}
}

func TestServiceReplayJobSendsFlags(t *testing.T) {
	var gotBody map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
			t.Fatalf("Decode() error = %v", err)
		}
		_, _ = w.Write([]byte(`{"code":"1","msg":"ok"}`))
	}))
	defer server.Close()

	client, err := openapi.NewClient(config.AppConfig{
		APIBaseURL: server.URL,
		AccessKey:  "test-ak",
		SecretKey:  "test-sk",
	})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	service := NewService(client)
	if err := service.ReplayJob(12, ReplayOptions{AutoStart: true, ResetToCreated: true}); err != nil {
		t.Fatalf("ReplayJob() error = %v", err)
	}
	if gotBody["jobId"] != float64(12) || gotBody["autoStart"] != true || gotBody["resetToCreated"] != true {
		t.Fatalf("request body = %#v, want replay flags", gotBody)
	}
}

func TestServiceRejectsBusinessFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"code":"0","msg":"invalid credentials"}`))
	}))
	defer server.Close()

	client, err := openapi.NewClient(config.AppConfig{
		APIBaseURL: server.URL,
		AccessKey:  "test-ak",
		SecretKey:  "test-sk",
	})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	service := NewService(client)
	if err := service.StartJob(10); err == nil || err.Error() != "invalid credentials" {
		t.Fatalf("StartJob() error = %v, want invalid credentials", err)
	}
}
