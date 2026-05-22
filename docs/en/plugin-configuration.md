Plugin tools can extend the capabilities of LLMs, for example, by searching the internet, performing scientific calculations, or drawing images. They empower and enhance the ability of LLMs to connect to the external world. After you deploy the open-source version of Kozex, you can refer to this document to configure plugin tools. Otherwise, these plugin tools may not function properly.
## Plugin integration methods
Kozex provides three types of plugin tools: official built-in plugins, custom plugins, and third-party plugins.
| **Plugin Type** | **Description** | **Integration Method** |
| --- | --- | --- |
| Official Built-in Plugins | Configured uniformly by the backend, and available to all users in Kozex Open Source Edition. | Refer to [the guide for integrating official plugins](#configure-the-official-plugin-tool). |
| Custom Plugins | Created and configured independently by developers, with usage scope limited to the current workspace. | Create plugins following the standard plugin development guide. |

## Integrating Official plugin
### List of official plugin tools
After logging into the open-source version of Kozex, you can click "Explore" in the left navigation bar and view all official plugin tools on the plugin page. These plugins invoke third-party services via APIs.
On the plugin page, if the plugin has a label of **unauthorized**, it means this plugin requires authentication, and the authentication information has not been configured yet. The configuration method can be referenced at **Configure the official plugin tool**. Plugins without this label can be used directly.

![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/0074d2b8b6a94d9281f9df4c5f4067a3~tplv-goo7wpa0wc-image.image)

The official plugin tool list provided by Kozex Open Source Edition is as follows:
| **Plugin name** | **Plugin description** |
| --- | --- |
| Wolfram Alpha | A powerful mathematical tool that can calculate various mathematical expressions. |
| Sohu Hot News | Helps users access daily hot news on Sohu's website. |
| Sohu Hot News | Help users access daily hot news on Sohu.com. |
| What is worth buying | Help users look up product discount information, providing prices, purchase channels, and cost-performance recommendations. |
| Chestnut Kanban | A task breakdown tool that dissects tasks based on user queries and generates task kanbans. |
| Document library search | Search Douchai document library contents and webpage URLs based on document title keywords. |
| Amap | Route planning, location inquiry, geographic information conversion, and other functions. |
| Lark spreadsheets | Create, edit, and query Lark spreadsheets. |
| Bocha search | Search for any webpage information and links from the entire web, with accurate results and complete summaries, making it more suitable for AI use. |
| Image compression | Upload image links and return compressed images in base64 format. |
| Lark authentication and authorization | Provide authentication and authorization features for Lark apps. |
| Lark multidimensional spreadsheet | Support the creation and management of Lark multidimensional spreadsheets, as well as data operations and searches. |
| Lark calendar | Create, update, delete, and query schedule information on the Lark calendar |
| Lark cloud documents | Support creating documents, adding content, retrieving document information, and searching documents |
| Lark messages | Support sending Lark messages and retrieving message history |
| Lark tasks | Call task-related APIs on the Lark Open Platform to create, update, delete, and query tasks |
| Lark knowledge base | Search the Lark knowledge base Wiki and retrieve the full list of subnode Wiki entries |
| Tianyancha | Use the Tianyancha API to retrieve enterprise information, including basic information and news sentiment. |
### Configure the official plugin tool
If you need to use Kozex's built-in official plugin tools directly, before use, you should follow these steps to configure the corresponding plugin authentication credentials.
Here, the example of the "Lark Docs" plugin is provided. For other official plugins, please refer to the API documentation of the respective third-party services to obtain authorization credentials.


1. Go to the `backend/conf/plugin/pluginproduct` directory.
2. Open the `plugin_meta.yaml` file.
3. Find the built-in tool you want to use, for example, search for the plugin name "Lark Docs".
   1. This plugin uses the OAuth 2.0 Authorization Code grant type, which requires a public host to be configured. Modify `SERVER_HOST` in `kozex/docker/.env` to your public host.
   2. Find the `payload` field for the corresponding plugin.
   3. Go to the [Lark Open Platform Developer Console](https://open.feishu.cn/app) > App Details > Credentials & Basic Info > **App Credentials**, get the App ID and App Secret, and modify the `client_id` and `client_secret` fields in the `payload`.
   4. Go to the [Lark Open Platform Developer Console](https://open.feishu.cn/app) > App Details > **Security Settings**, and configure the Redirect URL as `https://{SERVER_HOST}/api/plugin_oauth/:plugin_id/authorization_code`, where `plugin_id` is obtained from the plugin URL.
   5. Go to the [Lark Open Platform Developer Console](https://open.feishu.cn/app) > App Details > **Permissions & Scopes**, and apply for the `scope` permissions required in the `payload`.
      For other parameter descriptions of the plugin, refer to [Plugin Configuration Schema](#plugin-configuration-schema); for other field descriptions of the payload, refer to [Authorization Configuration Payload](#authorization-configuration-payload).
      ![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/89bb7471de0e48f89e1ff3d8fd63129d~tplv-goo7wpa0wc-image.image)
   6. Save the file.
4. Execute the following command to restart the service.
   ```Bash
   # Restart the service
   docker compose --profile '*' up -d --force-recreate --no-deps coze-server
   ```


After the service starts, go to the agent editing page, select the Lark Docs plugin, and you can use it normally.

![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/63fd8db3bd9446f48250f612e8f3a6b1~tplv-goo7wpa0wc-image.image)

### Plugin configuration schema
Using the Lark Docs plugin as an example, the detailed field description of its configuration schema is as follows. You can modify each optional parameter configuration according to business needs. For detailed field descriptions of the payload section, refer to [Authorization Configuration Payload](#authorization-configuration-payload).
```YAML
- plugin_id: 19 # Mandatory, type int. Unique identifier for the plugin.
  product_id: 7395041536909574154 # Required, reserved field, keep the current value.
  deprecated: false # Optional, type bool. Whether the plugin has been deprecated. When the official plugin needs to be removed from the store, set it to true.
  version: v1.0.0 # Mandatory, type string. Plugin version, the format must be vX.X.X, and the version number should be updated with each plugin upgrade.
  openapi_doc_file: lark_docx.yaml #Required, type string. Describes the OpenAPI3 document for the plugin tool API.
  plugin_type: 1 #Required, type int. Currently fixed at 1, indicating the plugin type is HTTP.
  manifest: #Required, type object. Plugin configuration metadata.
    schema_version: v1 #Required, type string. The version number of the plugin configuration metadata, currently fixed at v1.
    name_for_model: 飞书云文档 #Required, type string, currently unused. Plugin name, provided for understanding by the large model.
    name_for_human: 飞书云文档 #Required, type string. The name of the plugin displayed publicly.
    description_for_model: 这是一个飞书云文档 #Required, type string, temporarily unused. Plugin description, provided for understanding by the large model.
    description_for_human: 这是一个飞书云文档 #Required, type string. The description of the plugin displayed publicly.
    auth: #Required, type object. Plugin authentication configuration.
      type: oauth #Required, type string. Plugin authentication types, optional: none (no authentication); oauth (standard OAuth 2.0); service_http (simple authentication).
      sub_type: authorization_code #Optional, type: string. Plugin authentication subtypes. Authentication type is invalid when set to none; when using OAuth authentication type, only authorization_code (OAuth code mode) is supported; when using service_http authentication type, only token/api_key (API Token or API Key as single field authentication information) is supported.
      payload: '{"client_id":"","client_secret":"","client_url":"https://accounts.feishu.cn/open-apis/authen/v1/authorize","scope":"drive:drive","authorization_url":"https://open.larkoffice.com/open-apis/authen/v2/oauth/token","authorization_content_type":"application/json"}' #Optional, type: string. Specific configuration for authentication, in JSON format, with specific fields determined by the authentication type and subtype.
    logo_url: official_plugin_icon/plugin_lark_docx.png #Required, type: string. Plugin icon URI, which must be pre-configured in ./docker/volumes/minio/official_plugin_icon.
    api: #Required, type: object. Additional description information for the plugin.
      type: openapi #Required, type string. Plugin type, currently fixed as openapi, indicating the plugin type is HTTP, which is the same as the meaning of plugin_type.
    common_params: #Required, type object. General parameters for the tool's HTTP request.
      header: #Optional, type array. General header parameters for the tool's HTTP request. User-Agent is fixed as kozex/1.0.
        - name: User-Agent
          value: kozex/1.0
      path: [ ] #Optional, type array. General path parameters for the tool's HTTP request.
      query: [ ] #Optional, type array. General query parameters for the tool's HTTP requests.
      body: [ ] #Optional, type: array. General body parameters for the tool's HTTP requests.
  tools: #Required, type: array. Tools of the plugin.
    - tool_id: 190001 #Required, type: int. The tool's unique identifier ID.
      deprecated: false #Optional, type: bool. Whether the tool is deprecated. When it needs to be removed from the store, set to true.
      method: post #Required, type: string. Tool HTTP request method.
      sub_url: /document/create_document # Required, type string. Tool HTTP request subpath.
    - tool_id: 190002
      deprecated: false
      method: get
      sub_url: /document/get_document_content
```

### Authorization Configuration Payload
Using the Lark Docs plugin as an example, the detailed field description of its authorization configuration payload is as follows. You can modify each optional parameter configuration according to business needs.
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
