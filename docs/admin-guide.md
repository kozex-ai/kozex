# Admin Panel Guide

## Accessing the Admin Panel

Navigate to `http://<server-host>/admin/` in a browser. With the default local setup:

```
http://localhost:8888/admin/
```

### Who Can Access

The admin panel requires a logged-in session. Access is restricted to users whose email address appears in the **Admin Emails** list.

Set initial admin emails before first login via the environment variable:

```bash
# docker/.env
export ADMIN_EMAILS="you@example.com,colleague@example.com"
```

Once inside, you can add or remove admin emails through the panel itself (Basic Configuration → Admin Emails). Changes take effect immediately — the list in the database always takes precedence over the environment variable.

---

## Navigation

The left sidebar has one menu group — **Configuration Management** — with three sections:

| Section | What it controls |
|---------|-----------------|
| Basic Configuration | Admin access, registration, code runner, server address, Coze SaaS plugin |
| Model Management | AI model providers available to agents and workflows |
| Knowledge Configuration | Embedding, rerank, OCR, document parsing, built-in model |

---

## Basic Configuration

**Save is live** — changes take effect on the next request; no server restart required (except where noted).

### Admin Emails

Comma-separated list of email addresses allowed to access this panel. Type an address and press **Enter** or **,** to add it. At least one address must remain.

### User Registration

- **Disable User Registration** — when checked, the public sign-up form is hidden.
- **Whitelist Emails** — appears only when registration is disabled. Only these addresses can create accounts; useful for invite-only deployments.

### Python Code Execution Environment

Controls which backend executes Python code in workflow Code nodes. See [code-runner-guide.md](code-runner-guide.md) for full details on each option.

| Option | Overhead | Notes |
|--------|----------|-------|
| Local | ~100–300 ms | No isolation. Default for simple setups. |
| Sandbox | ~1–3 s | Deno + Pyodide WASM. Requires `deno` on PATH. |
| Coze Sandbox | ~10–50 ms | Recommended. Requires the `coze-sandbox` Docker service. |

Selecting **Sandbox** reveals additional permission fields:

| Field | Description |
|-------|-------------|
| Allow Environment Variables | Comma-separated env var names the sandbox may read (e.g. `PATH,HOME`) |
| Allow Read Directories | Filesystem paths readable by user code (e.g. `/tmp`) |
| Allow Write Directories | Filesystem paths writable by user code |
| Allow Run Commands | Executables user code may invoke |
| Allow Network Access | Hostnames/IPs reachable from user code (e.g. `cdn.jsdelivr.net`) |
| Allow FFI Libraries | Shared library paths for FFI use |
| Node Modules Directory | Path to a `node_modules` directory for the Deno runtime |
| Code Execution Timeout | Seconds before a code node is killed (1–3600). Set ≥ 120 for Sandbox cold starts. |
| Memory Limit | MB available to user code (1–10240). |

### Server Host Address

The public-facing address of this server, used when generating callback URLs and shared links. Accepts `hostname:port` or a full URL with scheme.

```
localhost:8888          # local development
https://coze.example.com  # production with TLS
```

### Coze SaaS Plugin

Integrates with the hosted Coze SaaS platform so that Coze SaaS plugins are available inside this instance.

- **Coze API Token** — personal access token from the Coze SaaS developer console.
- **Enable Coze SaaS Plugin** — toggle to activate the integration.
- **Coze SaaS API Base URL** — defaults to `https://api.coze.cn`; change for other regions.

---

## Model Management

Manage the AI models available to agents and workflows. Models are grouped by provider.

### Adding a Model

Click **+ Add** on any provider card to open the add-model form.

**Fields common to all providers:**

| Field | Notes |
|-------|-------|
| Display Name | Human-readable label shown in the model picker |
| Model | Provider model identifier (e.g. `gpt-4o`, `claude-opus-4-7`) |
| API Key | Authentication key; not required for Ollama |
| Base URL | Optional custom endpoint; leave blank to use the provider default |

**Provider-specific fields:**

| Provider | Extra fields |
|----------|-------------|
| Ark (Doubao) | Region; Thinking Type (Disable / Enable / Auto) |
| OpenAI | Use Azure toggle; Azure API Version; Thinking Type |
| Gemini | Backend; Project; Location; Thinking Type |
| Ollama | No API key required; Thinking Type |
| Qwen, Claude, DeepSeek | Thinking Type |

Before saving, the panel runs a quick connectivity test (`1+1=?`) against the model. If the call fails, the model is not saved and the error is shown.

### Deleting a Model

Click the delete icon on any model card and confirm. The model is removed immediately and is no longer selectable in agent/workflow configuration.

---

## Knowledge Configuration

Controls how the knowledge base indexes and retrieves documents.

### Embedding

Converts text into vectors for semantic search. Select a provider and fill in its connection details.

| Provider | Required fields |
|----------|----------------|
| Ark | Base URL, API Key, Model, Dims |
| OpenAI | Base URL, API Key, Model, Dims; optional: Azure toggle + API version |
| Ollama | Base URL, Model, Dims |
| Gemini | Base URL, API Key, Model, Backend, Dims; optional: Project, Location |
| HTTP | Address, Dims |

**Max Batch Size** — how many chunks are embedded in a single API call (default: 100). Reduce if the provider enforces a lower limit.

### Rerank

Reorders candidate chunks by relevance after retrieval.

| Option | Notes |
|--------|-------|
| RRF (default) | Reciprocal Rank Fusion — no external service required |
| VikingDB | Volcengine managed rerank service; requires AK, SK, and optionally Host, Region, Model |

### OCR

Extracts text from images and scanned PDFs uploaded to the knowledge base.

| Option | Notes |
|--------|-------|
| Volcengine | Volcano Engine OCR; requires AK and SK |
| PaddleOCR | Self-hosted PaddleOCR service; requires API URL |

### Document Parser

Parses document structure (tables, headings, layout) during ingestion.

| Option | Notes |
|--------|-------|
| Builtin (default) | No extra service needed |
| PaddleOCR | Uses a PaddleOCR structure API; requires API URL |

### Built-in Model

The model used internally for knowledge base operations: NL-to-SQL, query rewriting, image annotation, and workflow knowledge recall. Select any model already configured in **Model Management**.

Prefix-based overrides are available via environment variables if different operations need different models (see `.env.example` for `NL2SQL_`, `M2Q_`, `IA_`, `WKR_` prefixes).

---

## Priority: Panel vs Environment Variable

Settings saved in the admin panel are stored in the database and **always take precedence** over environment variables. Environment variables act as defaults only when no database value exists.

To reset a setting to its environment-variable default, clear the field in the panel and save.
