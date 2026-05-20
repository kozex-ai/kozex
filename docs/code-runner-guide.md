# Code Runner Guide

## Overview

The Code Runner is the backend responsible for executing user-written Python code in workflow Code nodes. Three implementations are available, selectable via the Admin panel or the `CODE_RUNNER_TYPE` environment variable.

| Type | Name | Overhead per request | Isolation |
|------|------|----------------------|-----------|
| 0 | Local / Direct | ~100–300 ms | None — subprocess shares host env |
| 1 | Sandbox | ~1000–3000 ms | Strong — Deno + Pyodide WASM, fine-grained permissions |
| 2 | Coze Sandbox | ~10–50 ms | Container-level — separate Docker service with resource limits |

---

## Type 0 — Local / Direct

**How it works:** Forks a new Python subprocess for each request, executes the code, and returns the result. No process reuse.

**Pros:**
- Zero additional infrastructure — works out of the box
- No cold-start cost beyond Python process creation

**Cons:**
- No isolation: user code can read environment variables, access the filesystem, and make network requests freely
- ~100–300 ms overhead per request due to process fork

**When to use:** Internal or trusted deployments where users are not writing arbitrary untrusted code.

**Configuration:**
```bash
export CODE_RUNNER_TYPE=0
```

---

## Type 1 — Sandbox (Deno + Pyodide)

**How it works:** Each request launches a new `deno run` subprocess. Deno loads `jsr:@langchain/pyodide-sandbox`, initializes the Pyodide WASM runtime, and executes the Python code inside it. Permission flags (`--allow-read`, `--allow-net`, etc.) restrict what user code can access.

**Pros:**
- Strongest isolation of the three: filesystem, network, and environment access are controlled at the Deno/WASM boundary
- Fine-grained `allow_*` permission fields configurable per deployment

**Cons:**
- Requires `deno` installed and on `PATH` for the server process
- First ever run downloads `jsr:@langchain/pyodide-sandbox` (~100 MB) — requires internet access
- ~1–3 s overhead per request due to Deno process startup + Pyodide WASM initialization; no process reuse
- Deno version compatibility: the sandbox was developed against Deno 1.x. Deno 2.x may require adjustment

**When to use:** Public-facing deployments where untrusted users write arbitrary code and strong isolation is a hard requirement.

**Prerequisites:**
```bash
# Install Deno (adds to ~/.deno/bin)
curl -fsSL https://deno.land/install.sh | sh

# Make deno available to the server process (pick one):
sudo ln -s ~/.deno/bin/deno /usr/local/bin/deno   # system-wide
# OR add to .env:
export PATH="$HOME/.deno/bin:$PATH"
```

**Configuration:**
```bash
export CODE_RUNNER_TYPE=1

# Permission controls (comma-separated; leave empty to deny)
export CODE_RUNNER_ALLOW_ENV=""          # e.g. PATH,HOME
export CODE_RUNNER_ALLOW_READ=""         # e.g. /tmp
export CODE_RUNNER_ALLOW_WRITE=""        # e.g. /tmp
export CODE_RUNNER_ALLOW_NET="cdn.jsdelivr.net"
export CODE_RUNNER_ALLOW_RUN=""          # e.g. python
export CODE_RUNNER_ALLOW_FFI=""
export CODE_RUNNER_NODE_MODULES_DIR=""
export CODE_RUNNER_TIMEOUT_SECONDS=60   # set high enough for Pyodide cold start
export CODE_RUNNER_MEMORY_LIMIT_MB=100
```

Or configure through the Admin panel at `http://localhost:8888/admin/admin/` after selecting **Sandbox**.

---

## Type 2 — Coze Sandbox (recommended)

**How it works:** A standalone Go service (`backend/cmd/coze-sandbox/`) maintains a pool of pre-warmed Python worker processes. Each worker communicates with the service over stdin/stdout JSON. The main server sends code execution requests to this service over HTTP (`POST /execute`). Dead workers are replaced asynchronously; pool size stays stable.

```
coze-server  ──HTTP──▶  coze-sandbox  ──stdin/stdout──▶  python worker × N
                         (Go, :8889)                       (pre-warmed pool)
```

**Pros:**
- Lowest overhead: ~10–50 ms per request (no process creation per request)
- Worker pool size and container resources (CPU, memory) are tunable independently of the main server
- Self-contained Docker image; no host dependencies beyond Python

**Cons:**
- Requires deploying a separate container (`coze-sandbox`)
- Container-level isolation only — worker processes share the container filesystem; no fine-grained permission controls yet

**When to use:** Default recommendation for most deployments. Best balance of performance and operational simplicity.

**Docker Compose:** The `coze-sandbox` service is included in `docker/docker-compose.yml` and `coze-server` depends on it. No manual wiring needed.

**Configuration:**
```bash
export CODE_RUNNER_TYPE=2
export COZE_SANDBOX_ENDPOINT="http://coze-sandbox:8889"   # internal service URL
export COZE_SANDBOX_POOL_SIZE=8        # pre-warmed Python worker processes
export COZE_SANDBOX_MAX_QUEUE=32       # max requests queued while workers are busy; excess → 503
export COZE_SANDBOX_EXEC_TIMEOUT_SECONDS=300 # max execution time; node timeout takes priority if shorter
export COZE_SANDBOX_CPU_LIMIT=2.0      # CPU cores for the coze-sandbox container
export COZE_SANDBOX_MEM_LIMIT=1G       # memory limit for the coze-sandbox container
```

**Resource sizing:** Each idle Python worker uses ~30–50 MB. A pool of 8 workers needs ~400 MB at idle. Set `COZE_SANDBOX_MEM_LIMIT` accordingly.

**Build:**
```bash
# Build the coze-sandbox image from source
docker build -f backend/Dockerfile.coze-sandbox -t coze-sandbox .
```

**Python packages:** Packages available to user code are installed in `/app/.venv` inside the image. To add packages, edit `Dockerfile.coze-sandbox` and rebuild:
```dockerfile
RUN pip install --no-cache-dir \
    numpy==2.3.1 \
    httpx==0.28.1 \
    requests==2.32.3 \
    your-package==x.y.z   # add here
```

**Execution timeout**

Two independent layers control how long a Code node is allowed to run:

1. **Node timeout (primary)** — Set in the workflow UI under each Code node's **Exception Handling → Timeout**. Stored as `ExceptionConfigs.TimeoutMS` and applied by the Eino framework as a `context.WithTimeout` before the node's `Invoke` call. This deadline propagates through the entire call chain:
   ```
   node_runner → code.Invoke → runner.Run(ctx) → HTTP request ctx → pool.Run(ctx)
   ```
   When it fires mid-execution, the sandbox worker is killed and a replacement is replenished asynchronously.

2. **`COZE_SANDBOX_EXEC_TIMEOUT_SECONDS` (safety cap, default 300 s)** — Applied in the sandbox server only when the request ctx has no deadline, or its deadline exceeds the cap. It is a backstop, not an override:
   ```go
   if deadline, ok := ctx.Deadline(); !ok || time.Until(deadline) > maxExecTimeout {
       ctx, cancel = context.WithTimeout(r.Context(), maxExecTimeout)
   }
   ```

| Node timeout (UI) | `COZE_SANDBOX_EXEC_TIMEOUT_SECONDS` | Effective limit |
|-------------------|--------------------------------------|-----------------|
| Not set | 300 (default) | 300 s — cap is the fallback |
| Not set | 600 | 600 s |
| 60 s | 300 (default) | 60 s — node timeout wins |
| 60 s | 600 | 60 s — node timeout wins |
| 400 s | 300 (default) | 300 s — cap truncates node |
| 400 s | 600 | 400 s — node within cap |

**`time.sleep()` and long-running code** — `time.sleep(N)` works as long as the effective limit exceeds `N`. For `time.sleep(120)`: set the Code node's exception timeout to at least 130 s in the workflow UI. The default cap of 300 s already covers this; no env var change needed unless the node timeout itself exceeds 300 s.

**Pool capacity** — Before execution begins, a request must acquire a worker from the pool. If all workers are busy and the queue (`COZE_SANDBOX_MAX_QUEUE`) is full, the sandbox returns HTTP 503. The Code node surfaces this as an error rather than retrying silently.

---

## Selecting a Runner

### Via Admin Panel

Navigate to `http://localhost:8888/admin/admin/` → **Python Environment** section → select **Local**, **Sandbox**, or **Coze Sandbox** → Save.

Settings saved through the Admin panel are persisted in the database and take effect immediately on the next code node execution (no restart required).

### Via Environment Variable

Set `CODE_RUNNER_TYPE` in `.env` before starting the server. This value is used only when no Admin panel configuration exists in the database.

```bash
# .env
export CODE_RUNNER_TYPE=2   # 0=Local, 1=Sandbox, 2=Coze Sandbox
```

### Priority

**Admin panel (database)** takes precedence over the environment variable.

---

## Switching Between Runners

Switching is live — the runner is re-initialized on the next request after the Admin panel save. For Coze Sandbox, the service must be running before switching to type 2.

---

## Troubleshooting

### `No such file or directory: 'deno'` (Sandbox)
Deno is not on the PATH seen by the server process. Install deno and ensure it is accessible:
```bash
sudo ln -s ~/.deno/bin/deno /usr/local/bin/deno
```

### `V8 did not recognize flag '--experimental-wasm-stack-switching'` (Sandbox)
Deno 2.x removed this flag. The Sandbox implementation targets Deno 1.x; use an older Deno version or switch to Coze Sandbox (type 2).

### Execution timed out (Sandbox)
Increase `CODE_RUNNER_TIMEOUT_SECONDS` in the Admin panel. The Pyodide WASM runtime requires 30–120 s on first run (package download). After the package is cached, subsequent runs need only 1–3 s.

### `coze-sandbox` health check failing
```bash
curl http://localhost:8889/health
docker logs coze-sandbox
```
Verify the container started and the Python pool initialized. Check `COZE_SANDBOX_POOL_SIZE` and memory limits.

### Code node reports timeout error (Coze Sandbox)
1. Check the node's Exception Handling config in the workflow UI — the timeout may be shorter than the code needs.
2. Check `COZE_SANDBOX_EXEC_TIMEOUT_SECONDS` — the cap may be truncating the node's configured timeout (see behavior matrix above).
3. Inspect the code for `time.sleep`, blocking I/O, or infinite loops.

If `context deadline exceeded` appears much earlier than expected, a workflow-level timeout (set in `workflow_run.go` via `ForegroundRunTimeout` / `BackgroundRunTimeout`) may have already expired before the Code node was reached.

### Coze Sandbox returns HTTP 503
All workers are busy and the queue is full. Options:
- Increase `COZE_SANDBOX_POOL_SIZE` (and raise `COZE_SANDBOX_MEM_LIMIT` — each idle worker uses ~30–50 MB)
- Reduce concurrent workflow executions reaching Code nodes
- Check for stuck workers caused by code that ignores signals; they will be reaped when `COZE_SANDBOX_EXEC_TIMEOUT_SECONDS` fires
