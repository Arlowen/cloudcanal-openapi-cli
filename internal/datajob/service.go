package datajob

import (
	"cloudcanal-openapi-cli/internal/openapi"
	"errors"
)

const (
	listPath   = "/cloudcanal/console/api/v1/openapi/datajob/list"
	queryPath  = "/cloudcanal/console/api/v1/openapi/datajob/queryjob"
	startPath  = "/cloudcanal/console/api/v1/openapi/datajob/start"
	stopPath   = "/cloudcanal/console/api/v1/openapi/datajob/stop"
	deletePath = "/cloudcanal/console/api/v1/openapi/datajob/delete"
	replayPath = "/cloudcanal/console/api/v1/openapi/datajob/replay"
)

type Operations interface {
	ListJobs() ([]Job, error)
	GetJob(jobID int64) (Job, error)
	StartJob(jobID int64) error
	StopJob(jobID int64) error
	DeleteJob(jobID int64) error
	ReplayJob(jobID int64, options ReplayOptions) error
}

type Service struct {
	client *openapi.Client
}

func NewService(client *openapi.Client) *Service {
	return &Service{client: client}
}

type response struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}

type listJobsRequest struct {
	DataJobName      string `json:"dataJobName,omitempty"`
	DataJobType      string `json:"dataJobType,omitempty"`
	Desc             string `json:"desc,omitempty"`
	SourceInstanceID *int64 `json:"sourceInstanceId,omitempty"`
	TargetInstanceID *int64 `json:"targetInstanceId,omitempty"`
}

type jobActionRequest struct {
	JobID int64 `json:"jobId"`
}

type replayJobRequest struct {
	JobID          int64 `json:"jobId"`
	AutoStart      *bool `json:"autoStart,omitempty"`
	ResetToCreated *bool `json:"resetToCreated,omitempty"`
}

type ReplayOptions struct {
	AutoStart      bool
	ResetToCreated bool
}

type Source struct {
	InstanceDesc   string `json:"instanceDesc"`
	InstanceID     string `json:"instanceId"`
	DataSourceType string `json:"dataSourceType"`
	HostType       string `json:"hostType"`
	DeployType     string `json:"deployType"`
	Region         string `json:"region"`
	LifeCycleState string `json:"lifeCycleState"`
}

type Task struct {
	DataTaskID     int64  `json:"dataTaskId"`
	DataTaskType   string `json:"dataTaskType"`
	DataTaskName   string `json:"dataTaskName"`
	DataTaskStatus string `json:"dataTaskStatus"`
	WorkerIP       string `json:"workerIp"`
}

type Job struct {
	DataJobID        int64   `json:"dataJobId"`
	DataJobName      string  `json:"dataJobName"`
	DataJobDesc      string  `json:"dataJobDesc"`
	UserName         string  `json:"userName"`
	DataJobType      string  `json:"dataJobType"`
	DataTaskState    string  `json:"dataTaskState"`
	CurrTaskStatus   string  `json:"currTaskStatus"`
	SourceDS         *Source `json:"sourceDsVO"`
	TargetDS         *Source `json:"targetDsVO"`
	SourceSchema     string  `json:"sourceSchema"`
	TargetSchema     string  `json:"targetSchema"`
	ConsoleJobID     int64   `json:"consoleJobId"`
	ConsoleTaskState string  `json:"consoleTaskState"`
	LifeCycleState   string  `json:"lifeCycleState"`
	HaveException    bool    `json:"haveException"`
	DataTasks        []Task  `json:"dataTasks"`
}

type listJobsResponse struct {
	response
	Data []Job `json:"data"`
}

type queryJobResponse struct {
	response
	Data Job `json:"data"`
}

func (s *Service) ListJobs() ([]Job, error) {
	var out listJobsResponse
	if err := s.client.PostJSON(listPath, listJobsRequest{}, &out); err != nil {
		return nil, err
	}
	if err := ensureSuccess(out.response, "failed to list jobs"); err != nil {
		return nil, err
	}
	if out.Data == nil {
		return []Job{}, nil
	}
	return out.Data, nil
}

func (s *Service) GetJob(jobID int64) (Job, error) {
	var out queryJobResponse
	if err := s.client.PostJSON(queryPath, jobActionRequest{JobID: jobID}, &out); err != nil {
		return Job{}, err
	}
	if err := ensureSuccess(out.response, "failed to query job"); err != nil {
		return Job{}, err
	}
	return out.Data, nil
}

func (s *Service) StartJob(jobID int64) error {
	var out response
	if err := s.client.PostJSON(startPath, jobActionRequest{JobID: jobID}, &out); err != nil {
		return err
	}
	return ensureSuccess(out, "failed to start job")
}

func (s *Service) StopJob(jobID int64) error {
	var out response
	if err := s.client.PostJSON(stopPath, jobActionRequest{JobID: jobID}, &out); err != nil {
		return err
	}
	return ensureSuccess(out, "failed to stop job")
}

func (s *Service) DeleteJob(jobID int64) error {
	var out response
	if err := s.client.PostJSON(deletePath, jobActionRequest{JobID: jobID}, &out); err != nil {
		return err
	}
	return ensureSuccess(out, "failed to delete job")
}

func (s *Service) ReplayJob(jobID int64, options ReplayOptions) error {
	var out response
	if err := s.client.PostJSON(replayPath, newReplayJobRequest(jobID, options), &out); err != nil {
		return err
	}
	return ensureSuccess(out, "failed to replay job")
}

func newReplayJobRequest(jobID int64, options ReplayOptions) replayJobRequest {
	req := replayJobRequest{JobID: jobID}
	if options.AutoStart {
		value := true
		req.AutoStart = &value
	}
	if options.ResetToCreated {
		value := true
		req.ResetToCreated = &value
	}
	return req
}

func ensureSuccess(resp response, fallback string) error {
	if resp.Code == "1" {
		return nil
	}
	if resp.Msg != "" {
		return errors.New(resp.Msg)
	}
	return errors.New(fallback)
}
