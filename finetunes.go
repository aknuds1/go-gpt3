package gogpt

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// FineTuneRequest represents a request to create an OpenAPI fine-tune.
type FineTuneRequest struct {
	TrainingFile string `json:"training_file"`
	Model        string `json:"model,omitempty"`
	Suffix       string `json:"suffix,omitempty"`
}

// FineTuneResponse represents a response to a request for creating an OpenAPI fine-tune.
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
