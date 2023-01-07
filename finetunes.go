package gogpt

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
)

// FineTuneRequest represents a request to create an OpenAI fine-tune.
type FineTuneRequest struct {
	TrainingFile string `json:"training_file"`
	Model        string `json:"model,omitempty"`
	Suffix       string `json:"suffix,omitempty"`
}

// FineTuneResponse represents a response to a request for creating an OpenAI fine-tune.
type FineTuneResponse struct {
	ID     string          `json:"id"`
	Model  string          `json:"model"`
	Status string          `json:"status"`
	Events []FineTuneEvent `json:"events"`
}

type FineTuneEvent struct {
	Object    string `json:"object"`
	CreatedAt int64  `json:"created_at"`
	Level     string `json:"level"`
	Message   string `json:"message"`
}

// CreateFineTune requests creating an OpenAI fine-tune.
func (c *Client) CreateFineTune(ctx context.Context, request FineTuneRequest) (FineTuneResponse, error) {
	b := bytes.NewBuffer(nil)
	if err := json.NewEncoder(b).Encode(request); err != nil {
		return FineTuneResponse{}, fmt.Errorf("JSON encode request: %w", err)
	}

	u := c.fullURL("/fine-tunes")
	req, err := http.NewRequest(http.MethodPost, u, b)
	if err != nil {
		return FineTuneResponse{}, fmt.Errorf("make %s request to %s: %w", http.MethodPost, u, err)
	}

	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")

	var resp FineTuneResponse
	err = c.sendRequest(req, &resp)
	return resp, err
}

// FineTuneJob represents an OpenAI fine-tune job.
type FineTuneJob struct {
	ID              string          `json:"id"`
	Object          string          `json:"object"`
	Model           string          `json:"model"`
	CreatedAt       int64           `json:"created_at"`
	Events          []FineTuneEvent `json:"events"`
	FineTunedModel  string          `json:"fine_tuned_model"`
	OrganizationID  string          `json:"organization_id"`
	ResultFiles     []File          `json:"result_files"`
	Status          string          `json:"status"`
	ValidationFiles []File          `json:"validation_files"`
	TrainingFiles   []File          `json:"training_files"`
}

// RetrieveFineTune requests the retrieval of an OpenAI fine-tune.
func (c *Client) RetrieveFineTune(ctx context.Context, id string) (FineTuneJob, error) {
	u := c.fullURL(path.Join("/fine-tunes", id))
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return FineTuneJob{}, fmt.Errorf("make %s request to %s: %w", http.MethodGet, u, err)
	}

	req = req.WithContext(ctx)

	var job FineTuneJob
	err = c.sendRequest(req, &job)
	return job, err
}

type fineTuneJobResp struct {
	Object string        `json:"object"`
	Data   []FineTuneJob `json:"data"`
}

// ListFineTunes lists the organization's OpenAI fine-tune jobs.
func (c *Client) ListFineTunes(ctx context.Context) ([]FineTuneJob, error) {
	u := c.fullURL("/fine-tunes")
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("make %s request to %s: %w", http.MethodGet, u, err)
	}

	req = req.WithContext(ctx)

	var resp fineTuneJobResp
	err = c.sendRequest(req, &resp)
	return resp.Data, err
}
