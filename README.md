<div align="center">

# Kozex

<p>
<a href="#what-is-kozex">What is Kozex</a> •
<a href="#feature-list">Feature list</a> •
<a href="#quickstart">Quickstart</a> •
<a href="#developer-guide">Developer Guide</a>
</p>
<p>
  <img alt="License" src="https://img.shields.io/badge/license-apache2.0-blue.svg">
  <img alt="Go Version" src="https://img.shields.io/badge/go-%3E%3D%201.23.4-blue">
</p>

English | [中文](README.zh_CN.md)

</div>

## What is Kozex?

Kozex is an enterprise-grade AI agent development platform, forked from [Coze Studio](https://github.com/coze-dev/coze-studio) with production-readiness improvements.

* **All core technologies for AI agent development**: prompt, RAG, plugin, workflow — developers can focus on creating AI value.
* **Low-code, ready to use**: visual canvas, complete app templates and build frameworks to quickly construct AI agents and workflows.

The backend is developed in Golang, the frontend uses React + TypeScript, and the overall architecture is based on microservices and domain-driven design (DDD). Kozex adds an async workflow execution engine, isolated Python sandbox, and enterprise features on top of the upstream foundation.

**Why Kozex?**
- **Full lifecycle governance** — Design, execution, release, and data: complete AI Agent workflow lifecycle management
- **Production-ready architecture** — Job and Sandbox deployed independently, each component can scale and fail in isolation, evolving toward Kubernetes
- **Observability** — Prometheus metrics, LLM full-chain tracing, structured logging to demystify the runtime black box
- **Enterprise features** — Driven by community needs, continuously evolving
- **Open source co-building** — Actively maintained, tracking upstream, welcoming community contributions

## Feature list
| **Module** | **Feature** |
| --- | --- |
| Model service | Manage the model list, integrate services such as OpenAI and Volcengine |
| Build agent | * Build, publish, and manage agent <br> * Support configuring workflows, knowledge bases, and other resources |
| Build apps | * Create and publish apps <br> * Build business logic through workflows |
| Build a workflow | Create, modify, publish, and delete workflows |
| Develop resources | Support creating and managing the following resources: <br> * Plugins <br> * Knowledge bases <br> * Databases <br> * Prompts |
| API and SDK | * Create conversations, initiate chats, and other OpenAPI <br> * Integrate agents or apps into your own app through Chat SDK |

## Quickstart

Environment requirements:

* Minimum system requirements: 2 Core, 4 GB RAM
* Docker and Docker Compose installed and running

Deployment steps:

1. Clone the repository.
   ```bash
   git clone https://github.com/kozex-ai/kozex.git
   ```
2. Deploy and start the service. The first run may take a while to pull and build images.

   ```bash
   cd kozex
   # macOS or Linux
   make web
   # Windows
   cp ./docker/.env.example ./docker/.env
   docker compose -f ./docker/docker-compose.yml up
   ```

3. Register an account at `http://localhost:8888/sign`.
4. Set admin credentials in `docker/.env` before accessing the admin panel:
   ```env
   ADMIN_EMAILS=your@email.com
   ```
5. Configure a model at `http://localhost:8888/admin/#model-management`.
6. Access Kozex at `http://localhost:8888/`.

> [!WARNING]
> If deploying to a public network, assess security risks beforehand: account registration exposure, Python execution in workflow code nodes, SSRF, and API privilege escalation risks.

## Developer Guide

* **Project Configuration**:
   * [Model Configuration](docs/en/model-configuration.md): Configure the model service before deploying.
   * [Plugin Configuration](docs/en/plugin-configuration.md): Add authentication keys for third-party services to use official plugins.
   * [Basic Component Configuration](docs/en/basic-component-configuration.md): Configure components such as image uploaders.
* [API Reference](docs/en/api-reference.md): API and Chat SDK authenticated via Personal Access Token.
* [Development Guidelines](docs/en/development-standards.md):
   * [Project Architecture](docs/en/development-standards.md#project-architecture)
   * [Code Development and Testing](docs/en/development-standards.md#code-development-and-testing)
   * [Troubleshooting](docs/en/development-standards.md#troubleshooting)
* [FAQ](docs/en/faq.md)

## License
This project uses the Apache 2.0 license. See the [LICENSE](https://github.com/kozex-ai/kozex/blob/main/LICENSE-APACHE) file for details.

## Community contributions
Contributions are welcome. See [CONTRIBUTING](https://github.com/kozex-ai/kozex/blob/main/CONTRIBUTING.md) and [Code of Conduct](https://github.com/kozex-ai/kozex/blob/main/CODE_OF_CONDUCT.md).

## Security
If you discover a potential security issue, please **do not** create a public GitHub Issue. Report it via [GitHub Security Advisories](https://github.com/kozex-ai/kozex/security/advisories) instead.

## Join Community

### 🐛 Issue Reports & Feature Requests
- **GitHub Issues**: [Submit bug reports or feature requests](https://github.com/kozex-ai/kozex/issues)
- **Pull Requests**: [Contribute code or documentation](https://github.com/kozex-ai/kozex/pulls)

### 💬 Discussion
Join the Feishu group for technical discussion and project updates.

## Acknowledgments
* The [Eino](https://github.com/cloudwego/eino) framework team — agent and workflow runtime, model abstractions, knowledge base indexing and retrieval
* The [FlowGram](https://github.com/bytedance/flowgram.ai) team — workflow canvas editor engine
* The [Hertz](https://github.com/cloudwego/hertz) team — high-performance Go HTTP framework
* The [Coze Studio](https://github.com/coze-dev/coze-studio) team — the upstream project this fork is based on
