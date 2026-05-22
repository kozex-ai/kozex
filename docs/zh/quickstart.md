## 环境要求
在参考本文安装 Kozex 之前，确保您的软硬件环境满足以下要求： 
| **项目** | **说明** |
| --- | --- |
| CPU | 2 Core |
| RAM | 4 GiB |
| Docker  | 提前安装 Docker、Docker Compose，并启动 Docker 服务，详细操作请参考 Docker 文档：  <br>  <br> * **macOS**：推荐使用 Docker Desktop 安装，参考 [Docker Desktop For Mac](https://docs.docker.com/desktop/setup/install/mac-install/) 安装指南。  <br> * **Linux**：参考 [Docker 安装指南](https://docs.docker.com/engine/install/) 和 [Docker Compose](https://docs.docker.com/compose/install/) 安装指南。  <br> * **Windows**：推荐使用 Docker Desktop 安装，参考 [Docker Desktop For Windows](https://docs.docker.com/desktop/setup/install/windows-install/) 安装指南。  |

## 安装 Kozex

### 步骤一：获取源码

在本地项目中执行以下命令，获取 Kozex 最新版本的源码。

```Bash
# 克隆代码 
git clone https://github.com/kozex-ai/kozex.git
```

### 步骤二：部署并启动服务

   首次部署并启动 Kozex 需要拉取镜像、构建本地镜像，可能耗时较久，请耐心等待。如果看到提示 "Container coze-server Started"，表示 Kozex 服务已成功启动。 
   
   ```Bash
   cd kozex
   # start service
   # for macOS or Linux
   make web  
   # for windows
   cp .env.example .env
   docker compose -f ./docker/docker-compose.yml up
   ```

### 步骤三：注册账号

访问 `http://localhost:8888/sign` 输入用户名、密码点击注册按钮。


### 步骤三：配置模型

配置模型，`http://localhost:8888/admin/#model-management` 新增模型。


### 步骤五：登录访问

访问 Kozex `http://localhost:8888/`

![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/088ba3638c1b49c1a02bad091e039a22~tplv-goo7wpa0wc-image.image)

## 后续操作

成功部署 Kozex 后，如需使用插件、知识库等功能，你还需要去 http://localhost:8888/admin/ 配置相关组件。

## 公网安全风险
如果要将 Kozex 部署到公网环境下，建议关注以下安全风险：
- 建议关闭注册功能，或开启邮箱白名单，否则任意用户均可通过链接使用。
- 建议为工作流代码节点开启默认沙箱环境，以增强安全性。
- 建议根据业务需求（例如内网访问）为部署环境配置网络，规避 SSRF 风险。
- Kozex 服务端默认只监听 localhost，部署到公网环境下时建议仅在必要时监听 `0.0.0.0`，或者增加额外的安全措施，以免服务直接暴露在公网。
