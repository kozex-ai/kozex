插件工具可以扩展 LLM 的能力，比如联网搜索、科学计算或绘制图片，赋予并增强了 LLM 连接外部世界的能力。部署 Kozex 开源版之后，你可以参考本文档配置插件工具，否则这些插件工具可能无法正常运行。
## 插件接入方式
Kozex 提供了两种插件工具类型，即官方内置插件、自定义插件。
| **插件类型** | **说明** | **接入方式** |
| --- | --- | --- |
| 官方内置插件 | 由后台统一配置，支持 Kozex 开源版中所有用户使用。 | 参考[接入官方插件](#接入官方插件)。 |
| 自定义插件 | 由开发者自行创建与配置，使用范围为当前工作空间。 | 在 Kozex 平台中创建插件。 |

## 接入官方插件
### 官方插件工具列表
登录 Kozex 开源版之后，你可以在左侧导航栏单击探索，并在插件页面查看所有官方插件工具。这些插件可通过 API 的方式调用第三方服务。
插件页面中，如果插件带有**未授权**的标签，表示此插件需要授权，且尚未配置授权信息。配置方式可参考[配置官方插件工具](#配置官方插件工具)。没有此标签的插件可以直接使用。

![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/816629dfc61e4264a140afcb829ee5c5~tplv-goo7wpa0wc-image.image)

Kozex 开源版提供的官方插件工具列表如下：
| **插件名称** | **插件描述** |
| --- | --- |
| Wolfram Alpha | 强大的数学工具，可计算各种数学表达式。 |
| 搜狐热闻 | 帮助用户获取搜狐网上的每日热闻。 |
| 什么值得买 | 帮助用户查询商品优惠信息，提供价格、购买渠道和性价比推荐。 |
| 板栗看板 | 任务拆解工具，根据用户提问拆解任务，生成任务看板。 |
| 文库搜索 | 根据文库文档标题关键字，搜索豆柴文库内容和网页URL。 |
| 高德地图 | 路线规划、位置查询、地理信息转换等功能 |
| 飞书电子表格 | 对飞书电子表格进行创建、编辑和查询等操作。 |
| 博查搜索 | 从全网搜索任何网页信息和网页链接，结果准确、摘要完整，更适合AI使用。 |
| 图片压缩 | 上传图片链接，返回压缩后的base64格式图片。 |
| 飞书认证及授权 | 提供飞书应用的认证及授权功能。 |
| 飞书多维表格 | 支持创建和管理飞书多维表格，进行数据操作和搜索。 |
| 飞书日历 | 在飞书日历上创建、更新、删除、查询日程信息。 |
| 飞书云文档 | 支持创建文档、添加内容、获取文档信息和搜索文档。 |
| 飞书消息 | 支持发送飞书消息和获取消息历史记录。 |
| 飞书任务 | 调用飞书开放平台任务相关API，创建、更新、删除和查询任务。 |
| 飞书知识库 | 搜索飞书知识库Wiki、获取Wiki全部子节点列表。 |
| 天眼查 | 利用天眼查API检索企业信息，查询企业基本信息、新闻舆情等。 |
### 配置官方插件工具
如果你需要直接使用 Kozex 提供的官方内置插件工具，在使用前，应参考以下步骤配置相应插件的授权凭证。
此处以"飞书云文档"插件为例，其他官方插件请参考对应第三方服务的 API 文档获取授权凭证。


1. 进入目录 `backend/conf/plugin/pluginproduct`。
2. 打开 `plugin_meta.yaml` 文件。
3. 找到你想要使用的内置工具，例如搜索插件名称"飞书云文档"。
   1. 该插件的授权方式为 OAuth2.0 授权码模式，需要配置公网 host。修改 `kozex/docker/.env` 中的`SERVER_HOST` 为公网 host。
   2. 找到对应插件的 `payload` 字段。
   3. 进入飞书开放平台[开发者后台](https://open.feishu.cn/app) > 应用详情页 > 凭证与基础信息 > **应用凭证**，获取应用凭证 App ID 和 App Secret ，修改 `payload` 中的 `client_id`、`client_secret` 字段。
   4. 进入飞书开放平台[开发者后台](https://open.feishu.cn/app) > 应用详情页 > **安全设置**，配置重定向 URL 为 `https://{SERVER_HOST}/api/plugin_oauth/:plugin_id/authorization_code`，`plugin_id` 从插件 url 上获取。
   5. 进入飞书开放平台[开发者后台](https://open.feishu.cn/app) > 应用详情页 > **权限管理**，申请 `payload` 中要求的 `scope` 权限。
      关于插件的其他参数配置说明，可参考[插件配置 schema](#插件配置-schema)；payload 的其他字段说明，可参考[授权配置 Payload](#授权配置-Payload)。
      ![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/89bb7471de0e48f89e1ff3d8fd63129d~tplv-goo7wpa0wc-image.image)
   6. 保存文件。
4. 执行以下命令重新启动服务。
   ```Bash
   #  重启服务
   docker compose --profile '*' up -d --force-recreate --no-deps coze-server
   ```


服务启动之后，在智能体编辑页面，选择飞书云文档插件，就可以正常使用了。
![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/77fa8788688048ffba34a50e37c0518e~tplv-goo7wpa0wc-image.image)
### 插件配置 schema
以飞书云文档插件为例，其配置 schema 详细字段说明如下。你可以根据业务需求修改各个可选参数配置。关于 payload 部分的详细字段说明，可参考[授权配置 Payload](#授权配置-Payload)。
```YAML
- plugin_id: 19 #必填，类型int。插件唯一标识ID。
  product_id: 7395041536909574154 #必填，预留字段，保持当前值。
  deprecated: false #可选，类型bool。插件是否已废弃，当官方插件需要从商店下架时，设置为true。
  version: v1.0.0 #必填，类型string。插件版本，格式必须为vX.X.X，每一次插件升级都应该升级版本号。
  openapi_doc_file: lark_docx.yaml #必填，类型string。描述插件工具API的OpenAPI3文档。
  plugin_type: 1 #必填，类型int。目前固定为1，表示插件类型为HTTP。
  manifest: #必填，类型object。插件配置元信息。
    schema_version: v1 #必填，类型string。插件配置元信息的版本号，目前固定为v1。
    name_for_model: 飞书云文档 #必填，类型string，暂时无用。插件名称，提供给大模型理解。
    name_for_human: 飞书云文档 #必填，类型string。插件公开展示的名称。
    description_for_model: 这是一个飞书云文档 #必填，类型string，暂时无用。插件描述，提供给大模型理解。
    description_for_human: 这是一个飞书云文档 #必填，类型string。插件公开展示的描述。
    auth: #必填，类型object。插件鉴权配置。
      type: oauth #必填，类型string。插件鉴权类型，可选none(无鉴权)；oauth(标准oauth2.0)；service_http(简单鉴权)。
      sub_type: authorization_code #可选，类型string。插件鉴权子类型。鉴权类型为none时无效；鉴权类型为oauth时，仅支持authorization_code(授权码模式)；授权类型为service_http时，仅支持token/api_key(APIToken或者APIKey等单字段鉴权信息)；
      payload: '{"client_id":"","client_secret":"","client_url":"https://accounts.feishu.cn/open-apis/authen/v1/authorize","scope":"drive:drive","authorization_url":"https://open.larkoffice.com/open-apis/authen/v2/oauth/token","authorization_content_type":"application/json"}' #可选，类型string。鉴权的具体配置，json格式，具体字段根据鉴权类型和子类型而定。
    logo_url: official_plugin_icon/plugin_lark_docx.png #必填，类型string。插件图标URI，需要提前在./docker/volumes/minio/official_plugin_icon配置。
    api: #必填，类型object。插件额外描述信息。
      type: openapi #必填，类型string。插件类型，目前固定为openapi，表示插件类型为HTTP，与plugin_type含义相同。
    common_params: #必填，类型object。工具HTTP请求的通用参数。
      header: #可选，类型array。工具HTTP请求的通用头参数。User-Agent 固定为 kozex/1.0。
        - name: User-Agent
          value: kozex/1.0
      path: [ ] #可选，类型array。工具HTTP请求的通用路径参数。
      query: [ ] #可选，类型array。工具HTTP请求的通用查询参数。
      body: [ ] #可选，类型array。工具HTTP请求的通用body参数。
  tools: #必填，类型array。插件工具列表。
    - tool_id: 190001 #必填，类型int。工具唯一标识ID。
      deprecated: false #可选，类型bool。工具是否已废弃，当需要从商店下架工具时，设置为true。
      method: post #必填，类型string。工具HTTP请求方法。
      sub_url: /document/create_document #必填，类型string。工具HTTP请求子路径。
    - tool_id: 190002
      deprecated: false
      method: get
      sub_url: /document/get_document_content
```

### 授权配置 Payload
以飞书云文档插件为例，其授权配置 payload 详细字段说明如下。你可以根据业务需求修改各个可选参数配置。
```Go
type AuthOfAPIToken struct {
    // Location is the location of the parameter.
    // It can be "header" or "query".
    Location HTTPParamLocation `json:"location"`
    // Key is the name of the parameter.
    Key string `json:"key"`
    // ServiceToken is the simple authorization information for the service.
    ServiceToken string `json:"service_token"`
}

type OAuthAuthorizationCodeConfig struct {
    ClientID     string `json:"client_id"`
    ClientSecret string `json:"client_secret"`
    // ClientURL is the URL of authorization endpoint.
    ClientURL string `json:"client_url"`
    // Scope is the scope of the authorization request.
    // If multiple scopes are requested, they must be separated by a space.
    Scope string `json:"scope,omitempty"`
    // AuthorizationURL is the URL of token exchange endpoint.
    AuthorizationURL string `json:"authorization_url"`
    // AuthorizationContentType is the content type of the authorization request, and it must be "application/json".
    AuthorizationContentType string `json:"authorization_content_type"`
}
```


