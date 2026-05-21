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
package impl

import (
	"os"
	"strings"

	"github.com/kozex-ai/kozex/backend/api/model/admin/config"
	"github.com/kozex-ai/kozex/backend/infra/coderunner"
	"github.com/kozex-ai/kozex/backend/infra/coderunner/impl/cozesandbox"
	"github.com/kozex-ai/kozex/backend/infra/coderunner/impl/direct"
	"github.com/kozex-ai/kozex/backend/infra/coderunner/impl/sandbox"
	"github.com/kozex-ai/kozex/backend/pkg/logs"
)

type Runner = coderunner.Runner

// New initialises the code runner selected by the admin configuration.
//
// Runner comparison (framework overhead per request, excluding user code):
//
//	Type 0 — Local/Direct  ~100–300 ms   fork a Python subprocess per request
//	Type 1 — Sandbox       ~1000–3000 ms new Deno process + Pyodide WASM init per request
//	Type 2 — Coze Sandbox  ~10–50 ms     pre-warmed worker pool, HTTP + stdin/stdout IPC
func New(conf *config.BasicConfiguration) Runner {
	logs.Infof("init code runner, type=%v", conf.CodeRunnerType)
	switch conf.CodeRunnerType {
	case config.CodeRunnerType_Sandbox:
		// Sandbox: Deno + Pyodide WASM isolation. Strong security boundary but
		// ~1-3s cold-start per request and no process reuse. Best for platforms
		// open to untrusted users; overkill for internal/trusted deployments.
		getAndSplit := func(key string) []string {
			v := os.Getenv(key)
			if v == "" {
				return nil
			}
			return strings.Split(v, ",")
		}
		config := &sandbox.Config{
			AllowEnv:       getAndSplit(conf.SandboxConfig.AllowEnv),
			AllowRead:      getAndSplit(conf.SandboxConfig.AllowRead),
			AllowWrite:     getAndSplit(conf.SandboxConfig.AllowWrite),
			AllowNet:       getAndSplit(conf.SandboxConfig.AllowNet),
			AllowRun:       getAndSplit(conf.SandboxConfig.AllowRun),
			AllowFFI:       getAndSplit(conf.SandboxConfig.AllowFfi),
			NodeModulesDir: conf.SandboxConfig.NodeModulesDir,
			TimeoutSeconds: conf.SandboxConfig.TimeoutSeconds,
			MemoryLimitMB:  conf.SandboxConfig.MemoryLimitMb,
		}

		return sandbox.NewRunner(config)
	case config.CodeRunnerType_CozeSandbox:
		// CozeSandbox: dedicated service (backend/cmd/coze-sandbox) with pre-warmed worker pool.
		// ~10-50ms per request. Configure endpoint via COZE_SANDBOX_ENDPOINT.
		return cozesandbox.NewRunner()
	default:
		// Local/Direct: fork a Python subprocess per request. ~100-300ms overhead,
		// no isolation. Suitable for trusted internal deployments.
		return direct.NewRunner()
	}
}
