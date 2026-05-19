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
	"github.com/kozex-ai/kozex/backend/infra/coderunner/impl/direct"
	"github.com/kozex-ai/kozex/backend/infra/coderunner/impl/sandbox"
)

type Runner = coderunner.Runner

func New(conf *config.BasicConfiguration) Runner {
	switch conf.CodeRunnerType {
	case config.CodeRunnerType_Sandbox:
		// Sandbox mode runs user code inside Deno + Pyodide (Python in WASM).
		// Its only advantage over direct mode is security isolation: user code
		// cannot access the host filesystem, environment variables, or network
		// beyond what is explicitly allowed via the AllowXxx fields below.
		//
		// Trade-offs to be aware of:
		//   - Each invocation cold-starts a new Deno process and initializes the
		//     Pyodide WASM runtime, adding roughly 1-3 seconds of overhead before
		//     user code begins executing. There is no worker pool or process reuse.
		//   - No concurrency limit: N concurrent code nodes spawn N Deno processes
		//     simultaneously, each consuming 50MB+ of memory for the WASM runtime
		//     alone. Under heavy load this can exhaust host memory.
		//
		// Use sandbox mode when the platform is open to untrusted users who can
		// write arbitrary code. For internal/trusted users, direct mode is
		// significantly more efficient.
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
	default:
		return direct.NewRunner()
	}
}
