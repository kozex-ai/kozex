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
	_ "embed"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

//go:embed worker.py
var workerScript []byte

func main() {
	pyPath := envStr("COZE_SANDBOX_PYTHON_PATH", "python3")
	poolSize := envInt("COZE_SANDBOX_POOL_SIZE", 8)
	port := envStr("COZE_SANDBOX_PORT", "8889")

	// Write embedded worker script to a temp file so the Python process can be
	// started with a file path (required by exec.Command).
	tmpFile, err := os.CreateTemp("", "coze-sandbox-worker-*.py")
	if err != nil {
		fmt.Fprintf(os.Stderr, "coze-sandbox: create temp file: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(tmpFile.Name())
	if _, err = tmpFile.Write(workerScript); err != nil {
		fmt.Fprintf(os.Stderr, "coze-sandbox: write worker script: %v\n", err)
		os.Exit(1)
	}
	tmpFile.Close()

	pool := NewPool(poolSize, pyPath, tmpFile.Name())
	defer pool.Close()

	mux := http.NewServeMux()
	mux.Handle("/execute", handleExecute(pool))
	mux.Handle("/health", handleHealth())

	srv := &http.Server{Addr: ":" + port, Handler: mux}

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		fmt.Println("coze-sandbox: shutting down")
		srv.Shutdown(context.Background())
	}()

	fmt.Printf("coze-sandbox listening on :%s (pool_size=%d)\n", port, poolSize)
	if err = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Fprintf(os.Stderr, "coze-sandbox: %v\n", err)
		os.Exit(1)
	}
}

func envStr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func envInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return n
		}
	}
	return def
}
