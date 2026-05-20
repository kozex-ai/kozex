/*
 * Copyright 2025 coze-dev Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cozesandbox

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/kozex-ai/kozex/backend/infra/coderunner"
	"github.com/kozex-ai/kozex/backend/types/consts"
)

type executeRequest struct {
	Code     string         `json:"code"`
	Params   map[string]any `json:"params"`
	Language string         `json:"language"`
}

type executeResponse struct {
	Result map[string]any `json:"result,omitempty"`
	Error  string         `json:"error,omitempty"`
}

type runner struct {
	endpoint string
	client   *http.Client
}

func NewRunner() coderunner.Runner {
	endpoint := os.Getenv(consts.CozeSandboxEndpoint)
	if endpoint == "" {
		endpoint = "http://localhost:8889"
	}

	// HTTP client timeout = sandbox's max exec cap + 5s buffer, so the sandbox
	// always gets a chance to return a clean JSON error before this side cuts
	// the connection. Default matches COZE_SANDBOX_EXEC_TIMEOUT_SECONDS default.
	maxExecTimeout := 300
	if v := os.Getenv(consts.CozeSandboxExecTimeoutSeconds); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			maxExecTimeout = n
		}
	}

	return &runner{
		endpoint: endpoint,
		client: &http.Client{
			Timeout: time.Duration(maxExecTimeout+5) * time.Second,
		},
	}
}

func (r *runner) Run(ctx context.Context, request *coderunner.RunRequest) (*coderunner.RunResponse, error) {
	body, err := json.Marshal(executeRequest{
		Code:     request.Code,
		Params:   request.Params,
		Language: string(request.Language),
	})
	if err != nil {
		return nil, fmt.Errorf("[coze-sandbox] marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, r.endpoint+"/execute", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("[coze-sandbox] create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("[coze-sandbox] http: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusServiceUnavailable {
		return nil, fmt.Errorf("[coze-sandbox] sandbox pool full, try again later")
	}

	var result executeResponse
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("[coze-sandbox] decode response (status=%d): %w", resp.StatusCode, err)
	}
	if result.Error != "" {
		return nil, fmt.Errorf("[coze-sandbox] %s", result.Error)
	}
	return &coderunner.RunResponse{Result: result.Result}, nil
}
