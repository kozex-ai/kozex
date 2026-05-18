# Local Development Guide

## Overview

Local development follows three steps: start middleware → start backend → start frontend.

---

## Step 1: Start Middleware

```bash
make middleware
```

This starts all required services (MySQL, Redis, NSQ, Milvus, Elasticsearch, MinIO) via `docker-compose-debug.yml`, with all ports mapped to `127.0.0.1`.

On first run, `docker/.env.debug` is automatically created from `docker/.env.debug.example`.

Wait until all containers show `healthy` before proceeding:

```bash
docker ps --format "table {{.Names}}\t{{.Status}}"
```

---

## Step 2: Start the Backend

**Pick one of the two options** — they cannot run at the same time as both bind to `:8888`.

### Option A: Normal Start (no breakpoints needed)

```bash
make server
```

`scripts/setup/server.sh` does the following:
1. Compiles the Go binary to `bin/opencoze`
2. Copies `docker/.env.debug` into `bin/`
3. Copies `backend/conf/` and `backend/static/` into `bin/resources/`
4. Starts the server from `bin/`, listening on `:8888`

### Option B: VS Code Debugger (mutually exclusive with Option A)

**Prerequisite**: Run `make server` at least once to populate `bin/` with `.env.debug` and `resources/`. You can Ctrl+C immediately after it starts — the files are what matter.

Once `bin/` is set up, press `F5` in VS Code:
1. The `build-opencoze-debug` task compiles the binary to `bin/opencoze-debug` with `-gcflags=all=-N -l` (disables inlining, enables breakpoints)
2. The server starts with `APP_ENV=debug`, CWD set to `bin/`, reading `bin/.env.debug`, listening on `:8888`

Set breakpoints in the code, then trigger a request from the frontend to hit them.

> **Note**: If you change anything under `backend/conf/`, stop the debugger, run `make server` once (Ctrl+C after it starts), then press `F5` again to sync the config files into `bin/resources/`.

---

## Step 3: Start the Frontend

```bash
cd frontend/apps/coze-studio
npm run dev
```

The dev server starts at `http://127.0.0.1:8080` and proxies `/api` and `/v1` requests to the backend at `http://localhost:8888`.

---

## Debugging Tips

### Inspect NSQ Message State

When debugging async workflow execution, open the NSQ Admin UI to check whether messages are being published and consumed correctly:

```
http://127.0.0.1:4171
```

Check the `kozex_workflow_executor` topic: `depth` shows backlog, `in_flight` shows messages currently being processed.

### Workflow Execution Stuck in Queued

Confirm the backend startup log contains the following line, which means the NSQ consumer registered successfully:

```
INF [kozex_workflow_executor/cg_workflow_executor] (127.0.0.1:4150) connecting to nsqd
```

---

## Troubleshooting

**`make server` fails to build**

Ensure Go ≥ 1.21 is installed. Run `go build ./...` inside `backend/` to see the full error.

**Breakpoints not hit in VS Code**

Make sure the Go extension (`golang.go`) is installed and `dlv` (Delve debugger) is available:
```bash
which dlv || go install github.com/go-delve/delve/cmd/dlv@latest
```

**Frontend `/api` requests return 502**

Confirm the backend is running and listening on `:8888`. If you changed the port, pass it when starting the frontend:
```bash
WEB_SERVER_PORT=8889 npm run dev
```
