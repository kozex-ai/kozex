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

package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"
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

func handleExecute(pool *Pool, maxExecTimeout time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req executeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Respect the caller's deadline (set by the workflow node timeout) but
		// cap at maxExecTimeout to prevent a single request from holding a worker
		// indefinitely when no deadline is set or the deadline is unreasonably long.
		ctx := r.Context()
		if deadline, ok := ctx.Deadline(); !ok || time.Until(deadline) > maxExecTimeout {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(r.Context(), maxExecTimeout)
			defer cancel()
		}

		result, err := pool.Run(ctx, req.Code, req.Params)
		if err != nil {
			if errors.Is(err, ErrPoolFull) {
				http.Error(w, "sandbox pool full, try again later", http.StatusServiceUnavailable)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(executeResponse{Error: err.Error()})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(executeResponse{Result: result})
	}
}

func handleHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}
