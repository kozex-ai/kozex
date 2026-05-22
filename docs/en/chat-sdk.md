This document is written based on the current latest version of kozex. You can upgrade and use the Chat SDK as needed according to the features introduced in version updates. The release history of Kozex can be found on [GitHub Releases](https://github.com/kozex-ai/kozex/releases).

## Before you begin

* Ensure that the agent or chat flow has been published to the Chat SDK channel.
* [Personal Access Tokens](api-reference.md#get-personal-access-token) have been acquired.


## Debug Chat SDK

1. Start the backend service.
2. Modify configuration files.
   Go to the `frontend/packages/studio/open-platform/chat-app-sdk` directory, rename the `.env.default` file to `.env.local`, and enter the following information:
   * Publish the Agent Chat SDK:
      * `CHAT_APP_INDEX_COZE_BOT_ID`: The agent ID. In the URL on the agent build page, the string following the **bot** keyword is the agent ID. For example, `https://localhost:8888/space/341****/bot/73428668*****`. The agent ID is `73428668*****`.
      * `CHAT_APP_COZE_TOKEN`: The [Personal Access Tokens](api-reference.md#get-personal-access-token) obtained before you begin.
   * Publish chat flow Chat SDK:
      * `CHAT_APP_CHATFLOW_COZE_APP_ID`: App ID. When you open a workflow or chat flow in the app, the string following the **project-ide** parameter in the URL is the appId.
      * `CHAT_APP_CHATFLOW_COZE_WORKFLOW_ID`: The workflow or chat flow ID. When you open a workflow or chat flow in the app, the string following the **workflow** parameter in the URL is the workflowId.
      * `CHAT_APP_COZE_TOKEN`: [Personal Access Tokens](api-reference.md#get-personal-access-token) obtained before you begin.
3. Run the following command to start the ChatSDK frontend.
   Run the following commands in the `frontend/packages/studio/open-platform/chat-app-sdk` directory:
   ```TypeScript
   npm run dev
   ```

4. Debug Chat SDK.
   Open `http://localhost:8081/client` in your browser to start debugging.


## Publish Chat SDK

1. Modify the Chat SDK API domain name
   Open the file `frontend/packages/studio/open-platform/open-env-adapter/src/chat/index.ts`, and change `openApiHostByRegion` to your backend API domain name.
2. Build an image.
   Go to the `frontend/packages/studio/open-platform/chat-app-sdk` directory and run the following command:
   ```TypeScript
   npm run build
   ```

3. Upload CDN resources.
   After the image build is complete, the libs directory will be generated automatically. **It is recommended to upload all resources in the libs directory to your own CDN server**.
   You can also upload using npm:
   1. Change the package name and version number in `package.json` to your own package name and version number.
   2. Run `npm publish`, then you can access the resource at `https://cdn.jsdelivr.net/npm/{packageName}@{packageversion}/libs/cn/index.js`.


## Install the ChatSDK

1. Install ChatSDK.
   You can load the Chat SDK JavaScript code directly on the page using a script tag. Caution: Replace the package name and version number.
   ```HTML
   <script src="https://cdn.jsdelivr.net/npm/{packageName}@{packageversion}/libs/cn/index.js"></script>
   ```

2. After the SDK is installed successfully, configure the following parameters in the code:
   The code example is as follows:
   ```TypeScript
   new CozeWebSDK.WebChatClient({
     /**
     * Agent or app settings
     * for agent
     *@param config.bot_id - Agent ID. * for app
     *@param config.type - To integrate a kozex app, you must set the value to app. *@param config.isIframe - Whether to use the iframe method to open the chat box
     *@param config.appInfo.appId - AI app ID. *@param config.appInfo.workflowId - Workflow or chatflow ID. */
     config: {
       type: 'bot',
       bot_id: 'xxxx',
       isIframe: false,
     },
     /**
     * The auth property is used to configure the authentication method. *@param type - Authentication method, default type is 'unauth', which means no authentication is required; it is recommended to set it to 'token', which means authentication is done through PAT (Personal Access Token) or OAuth. *@param token - When the type is set to 'token', you need to configure the PAT (Personal Access Token) or OAuth access key. *@param onRefreshToken - When the access key expires, a new key can be set as needed. */
     auth: {
       type: 'token',
       token: 'xxxx',
       onRefreshToken: async () => 'token'
     },
     /**
     * The userInfo parameter is used to set the display of agent user information in the chat box. *@param userInfo.id - ID of the agent user. *@param userInfo.url - URL address of the user's avatar. *@param userInfo.nickname - Nickname of the agent user. */
     userInfo: {
       id: 'user',
       url: 'https://lf-coze-web-cdn.coze.cn/obj/eden-cn/lm-lgvj/ljhwZthlaukjlkulzlp/coze/coze-logo.png',
       nickname: 'User',
     },
     ui: {
       /**
       * The ui.base parameter is used to add the overall UI effect of the chat window. *@param base.icon - Application icon URL. *@param base.layout - Layout style of the agent chat box window, which can be set as 'pc' or'mobile'. *@param base.lang - System language of the agent, which can be set as 'en' or 'zh-CN'. *@param base.zIndex - The z-index of the chat box. */
       base: {
         icon: 'https://lf-coze-web-cdn.coze.cn/obj/eden-cn/lm-lgvj/ljhwZthlaukjlkulzlp/coze/chatsdk-logo.png',
         layout: 'pc',
         lang: 'en',
         zIndex: 1000
       },
       /**
       * Controls whether to display the top title bar and the close button
       *@param header.isShow - Whether to display the top title bar. *@param header.isNeedClose - Whether to display the close button. */
       header: {
         isShow: true,
         isNeedClose: true,
       },
       /**
       * Controls whether to display the floating ball at the bottom right corner of the page. */
       asstBtn: {
         isNeed: true
       },
       /**
       * The ui.footer parameter is used to add the footer of the chat window. *@param footer.isShow - Whether to display the bottom copy module. *@param footer.expressionText - The text information displayed at the bottom. *@param footer.linkvars - The link copy and link address in the footer. */
       footer: {
         isShow: true,
         expressionText: 'Powered by ...',
       },
       /**
       * Control the UI and basic capabilities of the chat box. *@param chatBot.title - The title of the chatbox
       *@param chatBot.uploadable - Whether file uploading is supported. *@param chatBot.width - The width of the agent window on PC is measured in px, default is 460. *@param chatBot.el - Container for setting the placement of the chat box (Element). */
       chatBot: {
         title: 'kozex Bot',
         uploadable: true,
         width: 390,
       },
     },
   });
   ```
