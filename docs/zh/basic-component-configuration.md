部署 Kozex 开源版之后，如需使用图片上传功能，或知识库相关功能，应参考本文档配置功能相关的基础组件。这些组件通常依赖火山引擎等第三方服务，配置组件时需要填写第三方服务的账号密钥等鉴权配置。
## 上传组件
在多模态对话场景，往往需要在对话中上传图片、文件等信息，例如在智能体调试区域中发送图片消息：

![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/6d18160e25a54359b3db9ceee90ec25a~tplv-goo7wpa0wc-image.image)

此功能由上传组件提供。Kozex 上传组件目前支持以下三种，你可以**任选一种**作为上传组件。

* **（默认）minio**：图片和文件自动上传至本地主机的 minio 服务，通过指定端口号可访问。但本地主机必须配置公网域名，否则已上传的图片和文件只能生成内网访问链接，无法被大模型读取和识别。
* **火山引擎对象存储 TOS**：图片和文件自动上传至指定的火山引擎对象存储 TOS 服务，并生成一个公网可访问的 URL。如果选择 TOS，必须先开通 TOS 并在 Kozex 中配置火山密钥。
* **火山引擎 ImageX**：图片和文件自动上传至指定的火山引擎 ImageX，并生成一个公网可访问的 URL。如果选择 ImageX，必须先开通 ImageX 并在 Kozex 中配置火山密钥。

上传组件的配置方式如下：

1. 设置上传组件类型。
   在 docker 目录中打开 `.env` 文件，配置项 "FILE_UPLOAD_COMPONENT_TYPE" 的值表示上传组件类型。
   * **storage**（默认）：表示使用`STORAGE_TYPE` 配置的存储组件，`STORAGE_TYPE` 默认为 minio，也支持配置为 tos。
   * **imagex**：表示使用火山 ImageX 组件。
   ```Bash
   # This Upload component used in Agent / workflow File/Image With LLM  , support the component of imagex / storage
   # default: storage, use the settings of storage component
   # if imagex, you must finish the configuration of <VolcEngine ImageX> 
   export FILE_UPLOAD_COMPONENT_TYPE="storage"
   ```

2. 为上传组件添加秘钥等配置。
   同样在 docker 目录的 `.env` 文件中，根据组件类型填写以下配置。
   | **组件类型** | **配置方式** | **示例** |
   | --- | --- | --- |
   | minio <br> （默认） | 1. 在 Kozex 项目的 `docker/.env` 文件中，FILE_UPLOAD_COMPONENT_TYPE 设置为 storage。 <br> 2. 在 Storage component 区域，设置 STORAGE_TYPE 为 minio。 <br> 3. 在 MiniO 区域，维持默认配置即可。 <br>  <br> 如果你选择在火山引擎等公共云上部署 Kozex，则需要提前为你的云服务器开放 8888、8889 端口的访问权限。 <br>  | ![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/64578ddbcbb94dafa875cf266e49a0fc~tplv-goo7wpa0wc-image.image) <br>  |
   | tos | 1. 开通火山引擎 TOS 产品。 <br> 2. 在 Kozex 项目的 `docker/.env` 文件中，FILE_UPLOAD_COMPONENT_TYPE 设置为 storage。 <br> 3. 在 Storage component 区域，设置 STORAGE_TYPE 为 tos。 <br> 4. 在 TOS 区域，填写以下参数： <br>    * TOS_ACCESS_KEY：火山引擎密钥 AK。获取方式可参考[获取火山引擎 API 密钥](https://www.volcengine.com/docs/6257/64983)。 <br>    * TOS_SECRET_KEY：火山引擎密钥 SK。获取方式可参考[获取火山引擎 API 密钥](https://www.volcengine.com/docs/6257/64983)。 <br>    * TOS_ENDPOINT：TOS 服务的 Endpoint，获取方式可参考[地域和访问域名](https://www.volcengine.com/docs/6349/107356)。 <br>    * TOS_REGION：TOS 服务所在的地域，获取方式可参考[地域和访问域名](https://www.volcengine.com/docs/6349/107356)。 |  <br> ![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/9230b5a67296447ca807ea97ae6fdb62~tplv-goo7wpa0wc-image.image) <br>  |
   | imagex | 1. 开通火山引擎 veImageX 产品，并创建**素材托管**服务。操作步骤可参考[火山 veImageX 官方文档](https://www.volcengine.com/docs/508/8084)。注意创建**素材托管**服务时需要填写域名，建议提前获取一个可用的公开域名。 <br> 2. 在 Kozex 项目的 `docker/.env` 文件中，FILE_UPLOAD_COMPONENT_TYPE 设置为 imagex。 <br> 3. 在 VolcEngine ImageX 区域，填写以下参数： <br>    * VE_IMAGEX_AK：火山引擎密钥 AK。获取方式可参考[获取火山引擎 API 密钥](https://www.volcengine.com/docs/6257/64983)。 <br>    * VE_IMAGEX_SK：火山引擎密钥 SK。获取方式可参考[获取火山引擎 API 密钥](https://www.volcengine.com/docs/6257/64983)。 <br>    * VE_IMAGEX_SERVER_ID：火山引擎 veImageX 产品**服务管理**页面展示的服务 ID。 <br>    * VE_IMAGEX_DOMAIN：创建服务时指定的域名。 <br>    * VE_IMAGEX_TEMPLATE：图片处理配置的模板名称。 |  <br> ![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/30ef93d304d542159141c4a00777c363~tplv-goo7wpa0wc-image.image) <br>  |
3. 执行以下命令重启服务，使以上配置生效。
   ```Shell
   docker compose --profile '*' up -d --force-recreate --no-deps coze-server
   ```

其他组件配置，请访问 http://localhost:8888/admin/ 配置。
