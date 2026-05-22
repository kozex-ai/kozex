Kozex 是基于大语言模型的 AI 应用开发平台，首次部署运行 Kozex 开源版之前，你需要先克隆到本地的项目中，配置所需要的模型。项目正常运行过程中，也可以随时按需添加新的模型服务、删除不需要的模型服务。

## 模型列表 
Kozex 支持的模型服务如下： 

* 火山方舟 | Byteplus ModelArk
* OpenAI 
* DeepSeek 
* Claude 
* Ollama 
* Qwen
* Gemini

## 模型配置注意事项 
在填写模型配置文件之前，请确保你已经了解了以下注意事项： 

* 在删除模型之前，请确保此模型已无线上流量。
* 智能体或工作流根据模型 ID 来调用模型。对于已上线的模型，请勿修改模型 ID，否则可能导致模型调用失败。 

## 为 Kozex 配置模型

Kozex 是基于大语言模型的 AI 应用开发平台，首次部署并启动 Kozex 开源版之前，你需要先在 Kozex 项目里配置模型服务，否则创建智能体或者工作流时，无法正常选择模型。

可以通过配置管理后台 `http://localhost:8888/admin/#model-management` 来添加或删除模型。


### 第三方模型服务

| **平台** | **protocol** | **base_url** | **特别说明** |
| --- | --- | --- | --- |
| 火山方舟 |  ark | 国内火山引擎：https://ark.cn-beijing.volces.com/api/v3/ <br> 海外 BytePlus：https://ark.ap-southeast.bytepluses.com/api/v3/ | 无 |
| 阿里百炼 | openai 或 <br> qwen | https://dashscope.aliyuncs.com/compatible-mode/v1 | qwen3 系列在非流式调用时不支持 thinking，如果使用需要设置 conn_config 中 enable_thinking: false，Kozex 后续版本会适配此能力。 |
| 硅基流动 | openai | https://api.siliconflow.cn/v1 | 无 |
| 其他第三方 api 中转 | openai | api 文档中提供的地址 <br> 注意路径通常带 /v1 后缀，不带 /chat/completions 后缀 | 如果平台仅中转或代理模型服务，非 openai 模型请按照【官方模型服务】部分文档配置 protocol |

### 开源框架
| **框架** |  **protocol** | **base_url** | **特别说明** |
| --- | --- | --- | --- |
| ollama | ollama | http://${ip}:11434 | 1. 镜像网络模式是 bridge，coze-server 镜像内 localhost 不是主机的 localhost。需要修改为 ollama 部署机器的 ip，或 `http://host.docker.internal:11434` <br> 2. 检查 api_key：未设置 api_key 时，此参数置空。 <br> 3. 确认部署 Ollama 主机的防火墙是否已开放 11434 端口。 <br> 4. 确认 ollama 网络已开启**对外暴露**。 <br>  <br> ![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/c16a26be91c442b1ad7f74ce3c088fb5~tplv-goo7wpa0wc-image.image) |
| vllm | openai | http://${ip}:8000/v1（port 启动时指定） | 无 |
| xinference | openai | http://${ip}:9997/v1 （port 启动时指定） | 无 |
| sglang | openai | http://${ip}:35140/v1（port 启动时指定） | 无 |
| LMStudio | openai | http://${ip}:${port}/v1 | 无 |

### 官方模型服务
| **模型** |  **protocol** | **base_url** | **特别说明** |
| --- | --- | --- | --- |
| Doubao | ark | https://ark.cn-beijing.volces.com/api/v3/ | 无 |
| OpenAI | openai | https://api.openai.com/v1 | 检查 by_azure 字段配置，如果模型是微软 azure 提供的模型服务，此参数应设置为 true。 |
| Deepseek | deepseek | https://api.deepseek.com/ | 无 |
| Qwen | qwen | https://dashscope.aliyuncs.com/compatible-mode/v1 | qwen3 系列在非流式调用时不支持 thinking，如果使用需要设置 conn_config 中 enable_thinking: false，Kozex 后续版本会适配此能力。 |
| Gemini | gemini | https://generativelanguage.googleapis.com/ | 无 |
| Claude | claude | https://api.anthropic.com/v1/ | 无 |
