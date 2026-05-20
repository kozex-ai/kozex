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
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"sync"
	"time"
)

type workerReq struct {
	Code   string         `json:"code"`
	Params map[string]any `json:"params"`
}

type workerResp struct {
	Result map[string]any `json:"result,omitempty"`
	Error  string         `json:"error,omitempty"`
}

type worker struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout *bufio.Reader
	mu     sync.Mutex
}

func newWorker(pyPath, scriptPath string) (*worker, error) {
	cmd := exec.Command(pyPath, scriptPath) // ignore_security_alert RCE
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("stdin pipe: %w", err)
	}
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("stdout pipe: %w", err)
	}
	if err = cmd.Start(); err != nil {
		return nil, fmt.Errorf("start: %w", err)
	}
	return &worker{
		cmd:    cmd,
		stdin:  stdin,
		stdout: bufio.NewReader(stdoutPipe),
	}, nil
}

// execute runs one code request and returns the result.
// dead=true means the worker process has crashed and must be replaced.
// A user-code error returns dead=false (the worker is still healthy).
func (w *worker) execute(code string, params map[string]any) (result map[string]any, dead bool, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	reqBytes, err := json.Marshal(workerReq{Code: code, Params: params})
	if err != nil {
		return nil, false, fmt.Errorf("marshal: %w", err)
	}

	if _, err = w.stdin.Write(append(reqBytes, '\n')); err != nil {
		return nil, true, fmt.Errorf("write to worker: %w", err)
	}

	line, err := w.stdout.ReadString('\n')
	if err != nil {
		return nil, true, fmt.Errorf("read from worker: %w", err)
	}

	var resp workerResp
	if err = json.Unmarshal([]byte(line), &resp); err != nil {
		return nil, true, fmt.Errorf("unmarshal response: %w", err)
	}
	if resp.Error != "" {
		return nil, false, fmt.Errorf("%s", resp.Error)
	}
	return resp.Result, false, nil
}

func (w *worker) kill() {
	_ = w.stdin.Close()
	_ = w.cmd.Process.Kill()
	_ = w.cmd.Wait()
}

// Pool maintains a fixed set of pre-warmed Python worker processes.
// Workers are returned to the pool after each execution; a dead worker
// is replaced asynchronously so pool size stays stable over time.
type Pool struct {
	ch         chan *worker
	pyPath     string
	scriptPath string
}

func NewPool(size int, pyPath, scriptPath string) *Pool {
	p := &Pool{
		ch:         make(chan *worker, size),
		pyPath:     pyPath,
		scriptPath: scriptPath,
	}
	for i := range size {
		w, err := newWorker(pyPath, scriptPath)
		if err != nil {
			fmt.Printf("coze-sandbox: failed to start worker %d: %v\n", i, err)
			continue
		}
		p.ch <- w
	}
	return p
}

// Run acquires a worker, executes the code, and returns the worker to the pool.
// Blocks until a worker is available or ctx is cancelled.
func (p *Pool) Run(ctx context.Context, code string, params map[string]any) (map[string]any, error) {
	select {
	case w := <-p.ch:
		result, dead, err := w.execute(code, params)
		if dead {
			w.kill()
			go p.replenish()
		} else {
			p.ch <- w
		}
		return result, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (p *Pool) replenish() {
	time.Sleep(200 * time.Millisecond)
	w, err := newWorker(p.pyPath, p.scriptPath)
	if err != nil {
		fmt.Printf("coze-sandbox: replenish failed: %v\n", err)
		return
	}
	p.ch <- w
}

// Close kills all idle workers. Workers currently executing a request
// are not tracked here and will be orphaned — call Close only on shutdown.
func (p *Pool) Close() {
	for {
		select {
		case w := <-p.ch:
			w.kill()
		default:
			return
		}
	}
}
