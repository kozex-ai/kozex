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
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kozex-ai/kozex/backend/infra/coderunner"
)

func newTestRunner(srv *httptest.Server) *runner {
	return &runner{endpoint: srv.URL, client: srv.Client()}
}

func TestRunner_Run_success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/execute", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)

		var req executeRequest
		require.NoError(t, json.NewDecoder(r.Body).Decode(&req))
		assert.Equal(t, "Python", req.Language)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(executeResponse{
			Result: map[string]any{"output": "hello"},
		})
	}))
	defer srv.Close()

	resp, err := newTestRunner(srv).Run(context.Background(), &coderunner.RunRequest{
		Code:     `async def main(args): return {"output": "hello"}`,
		Params:   map[string]any{},
		Language: coderunner.Python,
	})
	require.NoError(t, err)
	assert.Equal(t, "hello", resp.Result["output"])
}

func TestRunner_Run_executionError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(executeResponse{
			Error: "NameError: name 'x' is not defined",
		})
	}))
	defer srv.Close()

	_, err := newTestRunner(srv).Run(context.Background(), &coderunner.RunRequest{
		Code:     `async def main(args): return x`,
		Language: coderunner.Python,
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "NameError")
}

func TestRunner_Run_contextCancelled(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// never respond
		<-r.Context().Done()
	}))
	defer srv.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := newTestRunner(srv).Run(ctx, &coderunner.RunRequest{
		Code:     `async def main(args): return {}`,
		Language: coderunner.Python,
	})
	require.Error(t, err)
}

func TestRunner_Run_malformedParams(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(executeResponse{
			Result: map[string]any{"value": 42},
		})
	}))
	defer srv.Close()

	resp, err := newTestRunner(srv).Run(context.Background(), &coderunner.RunRequest{
		Code:     `async def main(args): return {"value": 42}`,
		Params:   map[string]any{"key": "val"},
		Language: coderunner.Python,
	})
	require.NoError(t, err)
	assert.EqualValues(t, 42, resp.Result["value"])
}
