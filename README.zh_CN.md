<div align="center">

# Kozex

<p>
  <a href="#什么是-kozex">什么是 Kozex</a> •
  <a href="#功能清单">功能清单</a> •
  <a href="#快速开始">快速开始</a> •
  <a href="#开发指南">开发指南</a>
</p>
<p>
  <img alt="License" src="https://img.shields.io/badge/license-apache2.0-blue.svg">
  <img alt="Go Version" src="https://img.shields.io/badge/go-%3E%3D%201.23.4-blue">
</p>

[English](README.md) | 中文

</div>

## 什么是 Kozex？

Kozex 是企业级 AI Agent 开发平台，基于 [Coze Studio](https://github.com/coze-dev/coze-studio) fork，在其基础上持续加入生产可用的工程改进。

* **提供 AI Agent 开发所需的全部核心技术**：Prompt、RAG、Plugin、Workflow，使开发者可以聚焦创造 AI 核心价值。
* **开箱即用，低代码**：可视化画布、完整应用模板和编排框架，快速构建 AI Agent 和工作流。

后端采用 Golang 开发，前端使用 React + TypeScript，整体基于微服务架构和领域驱动设计（DDD）。Kozex 在上游基础上新增了异步工作流执行引擎、独立 Python 沙盒服务及企业功能。

**为什么选 Kozex？**
- **全链路治理** — 覆盖设计、执行、发布、数据四个层次，AI Agent 工作流从画布到生产的完整生命周期管理
- **生产就绪架构** — Job、Sandbox 独立部署，各组件可独立扩容、故障隔离，面向 Kubernetes 演进
- **可观测性** — Prometheus 指标、LLM 全链路追踪、结构化日志，将运行时黑盒透明化
- **企业功能** — 吸收社区需求，持续演进
- **开源共建** — 活跃维护，持续跟进上游，欢迎社区共建

## 功能清单
| **功能模块** | **功能点** |
| --- | --- |
| 模型服务 | 管理模型列表，可接入 OpenAI、火山方舟等在线或离线模型服务 |
| 搭建智能体 | * 编排、发布、管理智能体 <br> * 支持配置工作流、知识库等资源 |
| 搭建应用 | * 创建、发布应用 <br> * 通过工作流搭建业务逻辑 |
| 搭建工作流 | 创建、修改、发布、删除工作流 |
| 开发资源 | 支持创建并管理以下资源： <br> * 插件 <br> * 知识库 <br> * 数据库 <br> * 提示词 |
| API 与 SDK | * 创建会话、发起对话等 OpenAPI <br> * 通过 Chat SDK 将智能体或应用集成到自己的应用 |

## 快速开始

环境要求：

* 最低系统要求：2 Core、4 GB RAM
* 提前安装 Docker、Docker Compose，并启动 Docker 服务。

部署步骤：

1. 获取源码。
   ```bash
   git clone https://github.com/kozex-ai/kozex.git
   ```

2. 部署并启动服务。首次部署需要拉取并构建镜像，耗时较久，请耐心等待。
   ```bash
   cd kozex
   # macOS 或 Linux
   make web
   # Windows
   cp ./docker/.env.example ./docker/.env
   docker compose -f ./docker/docker-compose.yml up
   ```

3. 注册账号：访问 `http://localhost:8888/sign`，输入用户名和密码点击注册。
4. 在 `docker/.env` 中配置管理员邮箱，否则无法访问后台：
   ```env
   ADMIN_EMAILS=your@email.com
   ```
5. 配置模型：访问 `http://localhost:8888/admin/#model-management` 新增模型。
6. 访问 Kozex：`http://localhost:8888/`。

> [!WARNING]
> 若要将 Kozex 部署到公网环境，建议提前评估安全风险：账号注册功能暴露、工作流代码节点 Python 执行环境、SSRF 以及部分 API 水平越权风险。

## 开发指南

* **项目配置**：
   * [模型配置](docs/zh/model-configuration.md)：部署前必须配置模型服务。
   * [插件配置](docs/zh/plugin-configuration.md)：添加第三方服务鉴权密钥以使用官方插件。
   * [基础组件配置](docs/zh/basic-component-configuration.md)：配置图片上传等组件。
* [API 参考](docs/zh/api-reference.md)：通过个人访问令牌鉴权，提供对话和工作流相关 API。
* [开发规范](docs/zh/development-standards.md)：项目架构、代码开发与测试、故障排查。
* [常见问题](docs/zh/faq.md)

## License
本项目采用 Apache 2.0 许可证。详情请参阅 [LICENSE](https://github.com/kozex-ai/kozex/blob/main/LICENSE-APACHE) 文件。

## 社区贡献
欢迎社区贡献，贡献指南参见 [CONTRIBUTING](https://github.com/kozex-ai/kozex/blob/main/CONTRIBUTING.md) 和 [Code of Conduct](https://github.com/kozex-ai/kozex/blob/main/CODE_OF_CONDUCT.md)。

## 安全
如果你发现潜在的安全问题，请**不要**创建公开的 GitHub Issue，请通过 [GitHub Security Advisories](https://github.com/kozex-ai/kozex/security/advisories) 报告。

## 加入社区

### 🐛 问题反馈与功能建议
- **GitHub Issues**：[提交 Bug 报告或功能请求](https://github.com/kozex-ai/kozex/issues)
- **Pull Requests**：[贡献代码或文档改进](https://github.com/kozex-ai/kozex/pulls)

### 💬 技术交流
加入飞书群与其他开发者交流，获取项目最新动态。

## 致谢
* [Eino](https://github.com/cloudwego/eino) 框架团队 — 智能体和工作流运行时、模型抽象、知识库索引与检索
* [FlowGram](https://github.com/bytedance/flowgram.ai) 团队 — 工作流画布编辑引擎
* [Hertz](https://github.com/cloudwego/hertz) 团队 — 高性能 Go HTTP 框架
* [Coze Studio](https://github.com/coze-dev/coze-studio) 团队 — 本项目 fork 的上游项目
