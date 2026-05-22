Kozex is an AI app development platform based on LLM. Before running the Kozex open-source version for the first time, you need to clone the project to your local machine and configure the required models. During normal project operations, you can also add new model services as needed or delete unnecessary model services at any time.
## Model list
The model services supported by Kozex are as follows:

* Volcengine Ark | BytePlus ModelArk
* OpenAI
* DeepSeek
* Claude
* Ollama
* Qwen
* Gemini

## Model configuration instructions
Before filling out the model configuration file, ensure you have understood the following Important notes:

* Before deleting the model, ensure that it is no longer receiving online traffic.
* Agents or workflows call models based on model IDs. For models that have already been launched, do not modify their IDs; otherwise, it may result in model call failures.

## Configure the model for kozex.
Kozex is an AI app development platform based on LLMs. Before deploying and launching the open-source version of Kozex for the first time, you need to configure the model service in the Kozex project, otherwise, you won't be able to properly select a model during the creation of agents or workflows.

Configure the model at `http://localhost:8888/admin/#model-management` by adding a new model.


### Third-Party Model Service

| **Platform** | **Protocol** | **Base_url** | **Special Instructions** |
| --- | --- | --- | --- |
| Volcengine Ark |  ark | Volcengine Engine: https://ark.cn-beijing.volces.com/api/v3/ <br> Overseas BytePlus: https://ark.ap-southeast.bytepluses.com/api/v3/ | None |
| Alibaba Bai Lian |  openai or <br> qwen | https://dashscope.aliyuncs.com/compatible-mode/v1 | The qwen3 series does not support thinking in non-streaming calls. If used, you need to set enable_thinking: false in conn_config. Kozex will adapt this capability in future versions. |
| Silicon-based Flow |  openai | https://api.siliconflow.cn/v1 | None |
| Other Third-Party API Relay |  openai | The address provided in the API document <br> Note that the path usually has a /v1 suffix and does not have a /chat/completions suffix | If the platform only relays or proxies model services and the model is not an openai model, please configure the protocol according to the documentation in the [Official Model Service] section. |

### Open-Source Framework
| **Framework** |  **Protocol** | **Base_url** | **Special Instructions** |
| --- | --- | --- | --- |
| Ollama |  ollama | http://${ip}:11434 | 1. When the mirror network mode is bridge, localhost in the coze-server mirror is not the localhost of the host. It needs to be modified to the ip of the Ollama deployment machine or `http://host.docker.internal:11434`. <br> 2. Check the api_key: When the api_key is not set, this parameter is left blank. <br> 3. Confirm whether the firewall of the Ollama-deployed host has opened port 11434. <br> 4. Confirm that the Ollama network has enabled **External Exposure**. <br>  <br>![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/c16a26be91c442b1ad7f74ce3c088fb5~tplv-goo7wpa0wc-image.image) |
| vllm |  openai | http://${ip}:8000/v1 (specified when starting the port) | None |
| xinference |  openai | http://${ip}:9997/v1 (specified when starting the port) | None |
| sglang |  openai | http://${ip}:35140/v1 (specified when starting the port) | None |
| LMStudio |  openai | http://${ip}:${port}/v1 | None |

### Official Model Service
| **Model** |  **Protocol** | **Base_url** | **Special Instructions** |
| --- | --- | --- | --- |
| Doubao |  ark | https://ark.cn-beijing.volces.com/api/v3/ | None |
| OpenAI |  openai | https://api.openai.com/v1 | Check the by_azure field configuration. If the model is a model service provided by Microsoft Azure, this parameter should be set to true. |
| Deepseek |  deepseek | https://api.deepseek.com/ | None |
| Qwen |  qwen | https://dashscope.aliyuncs.com/compatible-mode/v1 | The qwen3 series does not support thinking in non-streaming calls. If used, you need to set enable_thinking: false in conn_config. Kozex will adapt this capability in future versions. |
| Gemini |  gemini | https://generativelanguage.googleapis.com/ | None |
| Claude |  claude | https://api.anthropic.com/v1/ | None |
