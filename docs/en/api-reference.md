# Before you begin
Before calling the Kozex open source API, ensure you have published the agent as an API service and obtained the access token through different authorization methods.
## Publish agent as an API service
After publishing the agent as an API service, you can use the agent by calling the API, such as viewing the agent's basic settings or initiating a chat with the agent.
To download the example file, follow these steps:

1. Log in to Kozex open source.
2. On the **Development** page, select the target agent.
3. In the upper-right corner of the page, click **publish**.
4. On the **publish** page, select the **API** option, then click **publish**.

   ![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/8b0fd3945757416faab0482a2c927b32~tplv-goo7wpa0wc-image.image)

## Get Personal Access Token

The Kozex API and Chat SDK authenticate via Personal Access Token. Before calling the API, you need to obtain an access token first.

When calling the Coze API, you need to specify the access token through the Authorization parameter in the Header. The Coze server will verify the operation permissions of the caller based on the access token.


The steps to obtain a Personal Access Token are as follows:

1. Log in to Kozex Open Source Edition.
2. Click your personal avatar in the bottom left corner of the page and select **API Authorization**.

   ![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/f5bb30622ad3450088b786cff028a800~tplv-goo7wpa0wc-image.image)
3. On the **Personal Access Token** page, click **Add New Token**.
4. Complete the following configuration on the popped-up page, then click **Confirm**.
   | **Setting** | **Note** |
   | --- | --- |
   | Name | Name of the Personal Access Token. |
   | Expiration time | Duration of validity for the Personal Access Token. Once the token expires, it will become invalid, and you will not be able to continue using it to call Coze API. <br> After generating the token, the expiration time cannot be modified. |
5. Copy and securely save the Personal Access Token.
   The generated token is displayed only once at this time, so copy and save it immediately.

# Create a conversation

Create a conversation.

A conversation is a Q&A interaction for a topic between a user and an agent. A conversation consists of one or more messages. When creating a conversation, Coze simultaneously creates a blank contextual segment within the conversation to store messages related to a specific topic. When initiating subsequent chats, the messages in the contextual segment serve as the message history visible to the model.

You can write messages synchronously when creating a conversation. By default, messages are written into the latest contextual segment and will be used as context during the chat passed to the model.

**Basic information**

| **HTTP method** | POST |
| --- | --- |
| **URI** | &#123;&#123;host&#125;&#125;/v1/conversation/create |
| **API description** | Create a conversation. |

**Request parameters**

**Header**

| **Parameter** | **Value** | **Description** |
| --- | --- | --- |
| Authorization | Bearer *$Access_Token* | A **Personal Access Token** used to verify the identity of the client. You can generate a Personal Access Token on the Kozex platform. For more details, refer to 【Preparations】. |
| Content-Type | application/json | Indicates the format in which the request body is interpreted. |

**Body**

| **Parameter** | **Type** | **Required** | **Example** | **Description** |
| --- | --- | --- | --- | --- |
| bot_id | String | Optional | 730454116184516* | The agent ID corresponding to the conversation. <br> Specify the agent ID when creating the conversation, so you can later view the conversation list corresponding to the specified agent ID. |
| connector_id | String | Optional | 1024 | Which channel the conversation was created on. The following channels are currently supported: <br>  <br> * API: (default) 1024 <br> * ChatSDK：999 |

**Response parameters**

| **Parameter** | **Type** | **Example** | **Description** |
| --- | --- | --- | --- |
| data | Object of [ConversationData](#conversationdata) | See the return example for details | Basic information of the conversation. Detailed instructions can refer to [ConversationData](#conversationdata). |
| code | Long | 0 | Status code. 0 indicates a successful call, other values indicate a failed call. You can determine the detailed error reason through the msg field. |
| msg | String | "Success" | This field provides detailed information about the result of the API call. If the API call fails, msg contains a description of the error or the reason for the failure. |

**ConversationData**

| **Parameter** | **Type** | **Example** | **Description** |
| --- | --- | --- | --- |
| id | String | 737999610479815**** | Conversation ID, which is the unique identifier of the conversation. |
| created_at | Long | 1718289297 | The creation time of the conversation. The format is a 10-digit Unixtime timestamp in seconds. |
| last_section_id | String | 7495664347616952360 | The latest context fragment ID in the conversation. |

**Example**

**Request example**

Create an empty conversation:

```shell
curl --location '{{host}}/v1/conversation/create' \
--header 'Authorization: Bearer pat_xxxxx' \
--header 'Content-Type: application/json' \
--data '{
    "bot_id":"7531406778865549312"
}'
```



**Response example**

```json
{
    "code": 0,
    "data": {
        "created_at": 1742820175,
        "id": "748535563872637****",
        "last_section_id": "748535563872637****",
    },
    "msg": ""
}
```



# Get conversation list
View the conversation list of a specified agent
* Only supports viewing the conversations generated by the agent in the API or Chat SDK channel through this API.
* Only supports querying conversations created by oneself.


**Basic information**

| **HTTP method** | GET |
| --- | --- |
| **URI** | &#123;&#123;host&#125;&#125;/v1/conversations |
| **API description** | Get the conversation list of the specified agent |

**Request parameters**

**Header**

| **Parameter** | **Value** | **Description** |
| --- | --- | --- |
| Authorization | Bearer *$Access_Token* | Personal Access Token used to verify the identity of the client You can generate a Personal Access Token on the Kozex platform. For details, refer to [Preparation] |
| Content-Type | application/json | Indicates the format in which the request body is interpreted. |

**Query**

| **Parameter** | **Type** | **Required** | **Example** | **Description** |
| --- | --- | --- | --- | --- |
| bot_id | String | Required | 73428668***** | Agent ID, the acquisition method is as follows: <br> Navigate to the agent's build page, the number after the `bot` parameter in the build page URL is the agent ID. For example, `https://www.coze.cn/space/341****/bot/73428668*****`, the agent ID is `73428668*****`. |
| page_num | Integer | Optional | 1 | Page number, starting from 1, default is 1. |
| page_size | Integer | Optional | 40 | Number of data entries per page, default is 50, maximum is 50. |
| sort_order | String | Optional | ASC | Sorting method for the conversation list: <br>  <br> * **ASC**: Sort by creation time in ascending order, with the earliest created conversations at the top. <br> * **DESC**: (Default) Sort by creation time in descending order, with the most recently created conversations at the top. |
| connector_id | String | Optional | 999 | Publish channel ID, used to filter specific channels' conversations. Only supports viewing conversations for the following channels: <br>  <br> * (Default) API channel, channel ID is 1024. <br> * Chat SDK channel, channel ID is 999. |

**Response parameters**

| **Parameter** | **Type** | **Example** | **Description** |
| --- | --- | --- | --- |
| code | Long | 0 | Status code. 0 means the call is successful, other values indicate the call failed. You can determine the detailed reason for the error through the msg field. |
| msg | String | Success | This field provides detailed information about the result of the API call. If the API call fails, msg contains a description of the error or the reason for the failure. |
| data | Object of [ListConversationData](#listconversationdata) | [ { "created_at": 1731575569, "id": "123456789123456789", "meta_data": {}, } ] | Detailed conversation list |

**ListConversationData**

| **Parameter** | **Type** | **Example** | **Description** |
| --- | --- | --- | --- |
| has_more | Boolean | false | Are there more conversations that were not returned in this request? <br>  <br> * true: There are more conversations not returned. <br> * false: All conversations matching the filter criteria have been returned. |
| conversations | Array of [ConversationData](#conversationdata) | { "created_at": 1731575569, "id": "12345456789*****", "meta_data": {}, } | Details of the conversation. |

**ConversationData**

| **Parameter** | **Type** | **Example** | **Description** |
| --- | --- | --- | --- |
| id | String | 737999610479815**** | Conversation ID, which is the unique identifier of the conversation. |
| created_at | Long | 1718289297 | The time when the conversation was created. The format is a 10-digit Unixtime timestamp in seconds. |
| last_section_id | String | 7495664347616952360 | The latest context fragment ID in the conversation. |

**Example**

**Request example**

```shell
curl --location --request curl -X GET {{host}}/v1/conversations?bot_id=73428668*****&page_num=1&page_size=20 \
  -H "Authorization: Bearer pat_OYDacMzM3WyOWV3Dtj2bHRMymzxP****"
```


**Response example**

```json
{
  "code": 0,
  "msg": "",
  "data": {
    "has_more": false,
    "conversations": [
      {
        "created_at": 1731575569,
        "id": "12345***3456789",
      }
    ]
  }
}
```

# ChatV3

Call this API to initiate a chat, supporting the addition of context and streaming responses.

**API description**

The Initiate a chat API is used to initiate a chat with a specified agent, supporting the addition of contextual messages in the conversation so that the agent can provide reasonable replies based on historical messages.
Kozex Community Edition's chat API only supports streaming response; while generating replies, the agent sends reply messages to the client one by one in the form of a data stream. After processing ends, the server returns a complete reply from the agent. 

The distinction between the **Create a conversation** API and the **Initiate a chat** API is as follows:

* Create a conversation: Primarily used to initialize a new conversation environment.
* Initiate a chat:
   * Used to initiate a chat within an already existing conversation.
   * Supports adding context and streaming responses
   * Can perform context association based on historical messages, providing more contextually appropriate replies


**Basic information**

| **HTTP method** | POST |
| --- | --- |
| **URI** | &#123;&#123;host&#125;&#125;/v3/chat |
| **API description** | Initiate a chat, allowing for adding context and streaming responses. |

**Header**

| **Parameter** | **Value** | **Description** |
| --- | --- | --- |
| Authorization | Bearer <span style="color: #D83931"><em>$Access_Token</em></span> | These credentials are used to authenticate the **Personal Access Token**. You can generate a Personal Access Token on the Kozex platform. For more information, see Before you begin. |
| Content-Type | application/json | Indicates the format in which the request body is interpreted. |

**Query**

| **Parameter** | **Type** | **Required** | **Description** |
| --- | --- | --- | --- |
| conversation_id | String | Optional | Identifies which conversation the chat occurred in. <br> A conversation is a Q&A interaction for a topic between a user and an agent. A conversation contains one or multiple messages. A chat is an invocation of the agent within a conversation, and the agent adds messages generated during the chat to the conversation. <br>  <br> * An existing conversation can be used, and messages already in the conversation will be passed as context to the model. <br> * For scenarios like single Q&A where distinguishing conversations is unnecessary, this parameter can be omitted, and the system will automatically generate a conversation. |

**Body**

| **Parameter** | **Type** | **Required** | **Description** |
| --- | --- | --- | --- |
| bot_id | String | Required | The agent ID for the conversation. <br> Enter the development page of the agent. The number after the parameter `bot` in the development page URL is the agent ID. For example, `https://www.coze.cn/space/341****/bot/73428668*****`, the agent ID is `73428668*****`. <br>  |
| user_id | String | Required | Identifies the user currently chatting with the agent, defined, generated, and maintained by the user The user_id is used to identify different users in a chat. The context messages, database, and other chat memory data of different user_ids are isolated from each other If user data isolation is not required, this parameter can be fixed as any arbitrary string, such as `123`, `abc`, etc. <br>  |
| additional_messages <br>  <br>  | Array of object <br>  | Optional <br>  | Chat input information. You can use this field to pass the user's question for this chat <br> * additional_messages 仅支持 role=user 的记录，以免影响模型效果。 <br>  |
| stream <br>  | Boolean <br>  | Optional <br>  | Whether to enable streaming return. Currently, **only streaming responses** are supported. <br>  <br> * **true**: "Streaming response" provides the model's real-time response to the client, similar to a typewriter effect. You can retrieve chat or message events returned by the server in real time, synchronize them in the client for real-time display, or directly retrieve the agent's final reply in the completed event. |
| shortcut_command | Object | Optional | Shortcut command information. You can specify the shortcut command to be executed during this conversation through this parameter. The shortcut command must be one already linked with the agent. <br> The message structure can refer to **ShortcutCommandDetail Object**. <br> 调用快捷指令，会自动根据快捷指令配置信息生成本次对话中的用户问题，并放入 additional_messages 最后一条消息作为本次对话的用户输入。 <br>  |

**EnterMessage Object**

| **Parameter** | **Type** | **Required** | **Description** |
| --- | --- | --- | --- |
| role | String | Required | The entity that sends this message. Valid values: <br>  <br> * **user**: indicates that the content of the message is sent by the user. <br> * **assistant**: Indicates that the message content is sent by the agent. |
| type <br>  | String | Optional <br>  | Message type. The default value is **question.** <br>  <br> * **question**: User input content. <br> * **answer**: Message content returned by the agent to the user, supports incremental returns. If the workflow is bound to a message node, there may be multiple answers. The end flag of the stream response can be used to determine when all answers are completed. <br> * **function_call**: Intermediate result of a function call during the chat with the agent. <br> * **tool_response**: The result returned after calling a tool (function call). <br> * **follow_up**: If the user question suggestion switch is enabled for the agent, replies related to recommended questions will be returned. Not supported as an input parameter in the request. <br> * **verbose**: In multiple answer scenarios, the server will return a verbose package, with the corresponding content in JSON format. `content.msg_type =generate_answer_finish` indicates that all answer replies are completed. Not supported as an input parameter in the request. <br>  <br> 仅发起会话（v3）接口支持将此参数作为入参，且： <br>  <br> * 如果 autoSaveHistory=true，type 支持设置为 question 或 answer。 <br> * If autoSaveHistory=false, type can be set to question, answer, function_call, or tool_output/tool_response. <br>  <br> Among them, type=question can only correspond to role=user, meaning only the user role can initiate messages of the question type. <br>  |
| content | String | Optional | The content of the message, supporting multiple content types such as plain text, multi-modal (text, image, file mixed input) content, and card content. <br>  <br> * When content_type is object_string, content is the JSON String serialized from the **object_string object**  array. For detailed descriptions, refer to **object_string object.** <br> * When content_type **=** text **** , content is plain text, for example, `"content" :"Hello!"`. |
| content_type | String | Optional | The type of message content, which can be set as: <br>  <br> * text: text format. <br> * object_string: Multi-modal content, a combination of text and files or text and images. <br> * Card: card. This enumeration value only appears in API responses and is not supported as an input parameter. <br>  <br> When content is not empty, this parameter is required. <br>  |

**object_string object**

| **Parameter** | **Type** | **Required** | **Description** |
| --- | --- | --- | --- |
| type | String | Required | Multimodal message content type, supported settings include: <br>  <br> * text: text type. <br> * file: file type. <br> * image: image type. <br> * audio: audio type. |
| text | String | Optional | text content. |
| file_url | String | Optional | Online address of file, image, or audio file. Must be a publicly accessible valid address. <br> When the type is file, image, or audio, at least one of file_id and file_url should be specified. |
* An object_string array can contain at most one `text` type message, but can include multiple `file` and `image` type messages.
* When there is a `text` type message in the object_string array, there must also be at least one `file` or `image` message. Pure text messages (only containing `text` type) need to be directly specified using the `content_type: text` field and cannot use the `object_string` array.
* Pure image or pure file messages are supported, but the previous or next message for each pure image or pure file message must contain a pure text message with `content_type: text`, as context for the user's query. For example, `"content": "[{\"type\":\"image\",\"file_id\":\"&#123;&#123;file_id_1&#125;&#125;\"}]"` is a pure image message. The previous or next message of this pure image message must be a pure text message; otherwise, the API will return a 4000 parameter error.

For example, the following array is a complete multimodal content:

<div style="display: flex;">
<div style="flex-shrink: 0;width: calc((100% - 16px) * 0.5000);">

Before serialization:
```JSON
[
    {
        "type": "text",
        "text": "你好我有一个帽衫，我想问问它好看么，你帮我看看"
    }, 
        {
        "type": "file",
        "file_url": "{{file_url_1}}"
    }
]
```



</div>
<div style="flex-shrink: 0;width: calc((100% - 16px) * 0.5000);margin-left: 16px;">

After serialization:
```JSON
"[{\"type\":\"text\",\"text\":\"你好我有一个帽衫，我想问问它好看么，你帮我看看\"},{\"type\":\"file\",\"file_url\":\"{{file_url_1}}\"}]"
```




</div>
</div>


**ShortcutCommandDetail Object**

| **Parameter** | **Type** | **Required** | **Description** |
| --- | --- | --- | --- |
| command_id | String | Required | The quick command ID to be executed by the chat, which must be a quick command already bound to the agent. |
| parameters | Map&lt;String, String&gt; | Optional | Quick command component Parameter information entered by the user. <br> Custom key-value pairs, where the key is the name of the quick command component and the value is the user input corresponding to the component, serialized as a JSON String after the **object_string object** array. For detailed Description, refer to **object_string object.** |

**Streaming response**

In a streaming response, the server sends data incrementally to the client, formatted as a data stream. This stream includes various events triggered during the chat until the process is either completed or interrupted. Upon completion, the server returns the concatenated and complete reply from the model through the conversation.message.completed event. Descriptions of each event can be found in **streaming response events**.
Streaming responses allow the client to start processing data before receiving the complete data stream, for example, displaying the agent's response content in a chat interface in real time, reducing the client's wait time for a complete response from the model.
The overall process of a streaming response is as follows:

* Streaming response process:
   ```Plain Text
   ######### 整体概览 （chat, MESSAGE 两级）
   # chat - 开始
   # chat - 处理中
   #   MESSAGE - 知识库召回
   #   MESSAGE - function_call
   #   MESSAGE - tool_output
   #   MESSAGE - answer is card
   #   MESSAGE - answer is normal text
   #   MESSAGE - 多 answer 的情况下，会继续有 message.delta
   #   MESSAGE - verbose （多 answer、Multi-agent 跳转等场景）
   #   MESSAGE - suggestion
   # chat - 完成
   # 流结束 event: done
   #########
   ```

* Streaming response example:
   ```Plain Text
   # chat - 开始
   event: conversation.chat.created
   // In the chat event, the id in the data field represents the Chat ID, i.e., the session ID.
   data: {"id": "123", "conversation_id":"123", "bot_id":"222", "created_at":1710348675,compleated_at:null, "last_error": null, "meta_data": {}, "status": "created","usage":null}
   
   # chat - 处理中
   event: conversation.chat.in_progress
   data: {"id": "123", "conversation_id":"123", "bot_id":"222", "created_at":1710348675, compleated_at: null, "last_error": null,"meta_data": {}, "status": "in_progress","usage":null}
   
   # MESSAGE - 知识库召回
   event: conversation.message.completed
   data: {"id": "msg_001", "role":"assistant","type":"knowledge","content":"---\nrecall slice 1:xxxxxxx\n","content_type":"text","chat_id": "123", "conversation_id":"123", "bot_id":"222"}
   
   # MESSAGE - function_call
   event: conversation.message.completed
   data: {"id": "msg_002", "role":"assistant","type":"function_call","content":"{\"name\":\"toutiaosousuo-search\",\"arguments\":{\"cursor\":0,\"input_query\":\"今天的体育新闻\",\"plugin_id\":7281192623887548473,\"api_id\":7288907006982012986,\"plugin_type\":1","content_type":"text","chat_id": "123", "conversation_id":"123", "bot_id":"222"}
   
   # MESSAGE - toolOutput
   event: conversation.message.completed
   data: {"id": "msg_003", "role":"assistant","type":"tool_output","content":"........","content_type":"card","chat_id": "123", "conversation_id":"123", "bot_id":"222"}
   
   # MESSAGE - answer is card
   event: conversation.message.completed
   data: {"id": "msg_004", "role":"assistant","type":"answer","content":"{{card_json}}","content_type":"card","chat_id": "123", "conversation_id":"123", "bot_id":"222"}
   
   # MESSAGE - answer is normal text
   event: conversation.message.delta
   data:{"id": "msg_005", "role":"assistant","type":"answer","content":"以下","content_type":"text","chat_id": "123", "conversation_id":"123", "bot_id":"222"}
   
   event: conversation.message.delta
   data:{"id": "msg_005", "role":"assistant","type":"answer","content":"是","content_type":"text","chat_id": "123", "conversation_id":"123", "bot_id":"222"}
   
   ...... {{ N 个 delta 消息包}} ......
   
   event: conversation.message.completed
   data:{"id": "msg_005", "role":"assistant","type":"answer","content":"{{msg_005 完整的结果。即之前所有 msg_005 delta 内容拼接的结果}}","content_type":"text","chat_id": "123", "conversation_id":"123", "bot_id":"222"}
   
   
   # MESSAGE - 多 answer 的情况,会继续有 message.delta
   event: conversation.message.delta
   data:{"id": "msg_006", "role":"assistant","type":"answer","content":"你好你好","content_type":"text","chat_id": "123", "conversation_id":"123", "bot_id":"222"}
   
   ...... {{ N 个 delta 消息包}} ......
   
   event: conversation.message.completed
   data:{"id": "msg_006", "role":"assistant","type":"answer","content":"{{msg_006 完整的结果。即之前所有 msg_006 delta 内容拼接的结果}}","content_type":"text","chat_id": "123", "conversation_id":"123", "bot_id":"222"}
   
   # MESSAGE - Verbose （流式 plugin, 多 answer 结束，Multi-agent 跳转等场景）
   event: conversation.message.completed
   data:{"id": "msg_007", "role":"assistant","type":"verbose","content":"{\"msg_type\":\"generate_answer_finish\",\"data\":\"\"}","content_type":"text","chat_id": "123", "conversation_id":"123", "bot_id":"222"}
   
   # MESSAGE - suggestion
   event: conversation.message.completed
   data: {"id": "msg_008", "role":"assistant","type":"follow_up","content":"朗尼克的报价是否会成功？","content_type":"text","chat_id": "123", "conversation_id":"123", "bot_id":"222"}
   event: conversation.message.completed
   data: {"id": "msg_009", "role":"assistant","type":"follow_up","content":"中国足球能否出现？","content_type":"text","chat_id": "123", "conversation_id":"123", "bot_id":"222"}
   event: conversation.message.completed
   data: {"id": "msg_010", "role":"assistant","type":"follow_up","content":"羽毛球种子选手都有谁？","content_type":"text","chat_id": "123", "conversation_id":"123", "bot_id":"222"}
   
   # chat - 完成
   event: conversation.chat.completed （chat完成）
   data: {"id": "123", "chat_id": "123", "conversation_id":"123", "bot_id":"222", "created_at":1710348675, compleated_at:1710348675, "last_error":null, "meta_data": {}, "status": "compleated", "usage":{"token_count":3397,"output_tokens":1173,"input_tokens":2224}}
   
   event: done （stream流结束）
   data: [DONE]
   
   # chat - 失败
   event: conversation.chat.failed
   data: {
       "code":701231,
       "msg":"error"
   }
   ```


The structure of the returned event message body is as follows:
| **Parameter** | **Type** | **Description** |
| --- | --- | --- |
| event | String | The current data packet event of the streaming response. For detailed information, see **streaming response events**. |
| data | Object | Message content. The formats of the chat event and the message event differ from each other. <br>  <br> * In the chat event, the data is represented by a **Chat Object**. <br> * In the message and audio events, the data is represented by **Message Object**. |

**Streaming response events**

| **Event** | **Description** |
| --- | --- |
| conversation.chat.created | Event indicating the creation and start of a new chat. |
| conversation.chat.in_progress | Event indicating that the server is currently processing the chat. |
| conversation.message.delta | Event for incremental message updates, typically occurring when type=answer. |
| conversation.audio.delta | Incremental audio message updates, typically occurring when type=answer. audio.delta is returned only when the input contains audio messages. |
| conversation.message.completed | Event indicating that the message has been fully replied. The streaming packet at this point contains the concatenated results of all message.deltas, and each message is marked as completed. |
| conversation.chat.completed | Event indicating that the chat has been completed. Tool type |
| conversation.chat.failed | Event indicating that the chat has failed. |
| conversation.chat.requires_action | Event indicating that the chat is interrupted and requires user action to report the execution result of a tool. |
| error | Event indicating an error during the streaming response process. |
| done | Event indicating that the streaming response of this conversation has ended normally. |
**API response parameters:**

| **Parameter** | **Type** | **Description** |
| --- | --- | --- |
| data | Object | The basic information of this chat. Detailed description can be referenced at **Chat Object**. |
| code | Integer | Status code. <br> `0` represents a successful call. |
| msg | String | This field provides detailed information about the result of the API call. If the API call fails, msg contains a description of the error or the reason for the failure. |

**Message Object**

| **Parameter** | **Type** | **Description** |
| --- | --- | --- |
| id | String | Message ID, the unique identifier of the message. |
| conversation_id | String | The ID of the conversation that this message belongs to. |
| bot_id | String | The ID of the agent that returns this message. This parameter is only returned in messages generated from the chat. |
| chat_id | String | Chat ID。 This parameter is only returned in messages generated from the chat. |
| meta_data | Map | Additional message data when creating a message, which will also be returned when retrieving the message. |
| role | String | The entity that sends this message. Valid values: <br>  <br> * **user**: indicates that the content of the message is sent by the user. <br> * **assistant**: Indicates that the message content is sent by the agent. |
| content | String <br>  | The content of the message, supporting various types such as plain text, multimodal (text, images, mixed file input), and cards. |
| content_type | String | The type of message content includes: <br>  <br> * text: text format. <br> * object_string: Multi-modal content, a combination of text and files or text and images. <br> * card: Card. This enumeration value only appears in the API response and is not supported as an input parameter. <br> * audio: audio type. This enumeration value only appears in the API response and is not supported as an input parameter. This type will only be returned if the input contains an audio file. When the content_type is audio, the content is audio data encoded in base64. The encoding of the audio varies depending on the input audio file: <br>    * If the input is a wav format audio file, the content is **a base64 string of pcm audio fragments with a sampling rate of 24kHz, raw 16 bit, 1 channel, little-endian**. <br>    * If the input is an ogg_opus format audio file, the content is **a base64 string of opus format audio fragments with a sampling rate of 48kHz, 1 channel, and a frame length of 10ms**. |
| created_at | Integer | The time when the message was created, represented as a 10-digit Unix timestamp in seconds (s). |
| updated_at | Integer | The update time of the message, formatted as a 10-digit Unix timestamp in seconds (s). |
| type | String | Message type. <br>  <br> * **question**: User input content. <br> * **answer**: Message content returned by the agent to the user, supports incremental returns. If the workflow is bound to a message node, there may be multiple answers. The end flag of the stream response can be used to determine when all answers are completed. <br> * **function_call**: Intermediate result of a function call during the chat with the agent. <br> * **tool_response**: The result returned after calling a tool (function call). <br> * **follow_up**: If the recommended question toggle is enabled on the agent, it will return relevant reply content for the suggested questions. <br> * **verbose**: In a multiple-answer scenario, the server will return a verbose package, with the corresponding content in JSON format. `content.msg_type =generate_answer_finish` represents all answers have been completed. <br>  <br> Only the conversation (v3) API supports using this parameter as an input, and: <br>  <br> * If autoSaveHistory=true, type can be set to question or answer. <br> * If autoSaveHistory=false, type can be set to question, answer, function_call, or tool_response. <br>  <br> Among them, type=question can only correspond to role=user, meaning only the user role can initiate messages of the question type. <br>  |
| section_id | String | Context fragment ID.  |
| reasoning_content | String |   <br> 该参数仅在使用 DeepSeek-R1 模型时才会返回。 |

**Example**

**Streaming response**

* **Request**
   ```Shell
   curl --location --request POST '{{host}}/v3/chat?conversation_id=7374752000116113452' \
   --header 'Authorization: Bearer pat_OYDacMzM3WyOWV3Dtj2bHRMymzxP****' \
   --header 'Content-Type: application/json' \
   --data-raw '{
       "bot_id": "734829333445931****",
       "user_id": "123456789",
       "stream": true,
       "auto_save_history":true,
       "additional_messages":[
           {
               "role":"user",
               "content":"2024年10月1日是星期几",
               "content_type":"text"
           }
       ]
   }'
   ```

* **Response**
   ```Plain Text
   event:conversation.chat.created
   // 在 chat 事件里，data 字段中的 id 为 Chat ID，即会话 ID。
   data:{"id":"7382159487131697202","conversation_id":"7381473525342978089","bot_id":"7379462189365198898","completed_at":1718792949,"last_error":{"code":0,"msg":""},"status":"created","usage":{"token_count":0,"output_count":0,"input_count":0}}
   
   event:conversation.chat.in_progress
   data:{"id":"7382159487131697202","conversation_id":"7381473525342978089","bot_id":"7379462189365198898","completed_at":1718792949,"last_error":{"code":0,"msg":""},"status":"in_progress","usage":{"token_count":0,"output_count":0,"input_count":0}}
   
   event:conversation.message.delta
   data:{"id":"7382159494123470858","conversation_id":"7381473525342978089","bot_id":"7379462189365198898","role":"assistant","type":"answer","content":"2","content_type":"text","chat_id":"7382159487131697202"}
   
   event:conversation.message.delta
   data:{"id":"7382159494123470858","conversation_id":"7381473525342978089","bot_id":"7379462189365198898","role":"assistant","type":"answer","content":"0","content_type":"text","chat_id":"7382159487131697202"}
   
   //省略模型回复的部分中间事件event:conversation.message.delta
   ......
   
   event:conversation.message.delta
   data:{"id":"7382159494123470858","conversation_id":"7381473525342978089","bot_id":"7379462189365198898","role":"assistant","type":"answer","content":"星期三","content_type":"text","chat_id":"7382159487131697202"}
   
   event:conversation.message.delta
   data:{"id":"7382159494123470858","conversation_id":"7381473525342978089","bot_id":"7379462189365198898","role":"assistant","type":"answer","content":"。","content_type":"text","chat_id":"7382159487131697202"}
   
   event:conversation.message.completed
   data:{"id":"7382159494123470858","conversation_id":"7381473525342978089","bot_id":"7379462189365198898","role":"assistant","type":"answer","content":"2024 年 10 月 1 日是星期三。","content_type":"text","chat_id":"7382159487131697202"}
   
   event:conversation.message.completed
   data:{"id":"7382159494123552778","conversation_id":"7381473525342978089","bot_id":"7379462189365198898","role":"assistant","type":"verbose","content":"{\"msg_type\":\"generate_answer_finish\",\"data\":\"\",\"from_module\":null,\"from_unit\":null}","content_type":"text","chat_id":"7382159487131697202"}
   
   event:conversation.chat.completed
   data:{"id":"7382159487131697202","conversation_id":"7381473525342978089","bot_id":"7379462189365198898","completed_at":1718792949,"last_error":{"code":0,"msg":""},"status":"completed","usage":{"token_count":633,"output_count":19,"input_count":614}}
   
   event:done
   data:"[DONE]"
   ```

# Get the message list

Ciew the list of messages for the specified conversation.

**View the message list** API and **View chat message details** API differ in the following ways:

* **View the message list** API is used to query the message records in the specified conversation, including every message manually inserted by the developer in the conversation and the queries sent by the user, as well as the agent replies of type=answer obtained from calling the **Initiate chat** API, but excluding intermediate messages of type=function_call, tool_response, and **follow-up**.
* **View chat message details** API is typically used in non-streaming chat scenarios to view the agent replies of type=answer and intermediate messages of type=function_call, tool_response, and **follow-up** in a specific chat. Excludes user-sent queries.

**Basic information**
| **HTTP method** | POST |
| --- | --- |
| **URI** | &#123;&#123;host&#125;&#125;/v1/conversation/message/list |
| **API description** | Ciew the list of messages for the specified conversation. |

**Request parameters**

**Header**

| **Parameter** | **Value** | **Description** |
| --- | --- | --- |
| Authorization | Bearer $AccessToken | These credentials are used to authenticate the **Personal Access Token**. You can generate a personal access token on the Kozex platform. For more information, see 【Before you begin】. |
| Content-Type | application/json | Indicates the format in which the request body is interpreted. |

**Query**

| **Parameter** | **Type** | **Required** | **Example** | **Description** |
| --- | --- | --- | --- | --- |
| conversation_id | Integer | Required | 737363834493434**** | Conversation ID, the unique identifier of the conversation. You can view the conversation_id field in the Response of the initiate chat API. |

**Body**

| **Parameter** | **Type** | **Required** | **Example** | **Description** |
| --- | --- | --- | --- | --- |
| order | String | Optional | desc | The sorting method of the message list. <br>  <br> * desc: (default) Sort in descending order of creation time, with the latest messages appearing first. <br> * asc: Sort in ascending order of creation time, with the earliest messages appearing first. |
| chat_id | String | Optional | 737999610479815**** | Filter the message list in the specified Chat ID. The chat API Response field data.id of the Chat event is the Chat ID. |
| before_id | String | Optional | 737363834493437**** | View messages before the specified position. <br> Defaults to 0, indicating no specific position. To page backward, specify the first_id from the returned results. |
| after_id | String | Optional | 737363834493437**** | View messages after the specified position. <br> The default is 0, indicating no specified position. To paginate backward, specify the last_id from the returned results. |
| limit | Long | Optional | 30 | The data volume returned with each query. The default value is 50, with an adjustable range from 1 to 50. |

**Response parameters**

| **Parameter** | **Type** | **Example** | **Description** |
| --- | --- | --- | --- |
| data | Array of [OpenMessageApi](#openmessageapi) | [ { "bot_id": "", "chat_id": "", "content": "What is your name", "content_type": "text", "conversation_id": "737363834493434****", "created_at": 1716809829, "id": "737363834493437****", "meta_data": {}, "role": "user", "type": "", "updated_at": "1716809829" }] | Message details. |
| has_more | Boolean | true | Has all messages been returned. <br>  <br> * true: Not all messages have been returned; you can call this API again to view other pages. <br> * false: All messages have been returned. |
| first_id | String | 737363834493437**** | The Message ID of the first message in the returned message list. |
| last_id | String | 737363834493440**** | The Message ID of the last message in the returned message list. |
| code | Long | 0 | Status code. 0 indicates the call was successful, other values indicate a failure; you can determine the detailed error reason through the msg field. |
| msg | String | Success | This field provides detailed information about the result of the API call. If the API call fails, msg contains a description of the error or the reason for the failure. |

**OpenMessageApi**

| **Parameter** | **Type** | **Example** | **Description** |
| --- | --- | --- | --- |
| id | String | 738130009748252**** | Message ID, the unique identifier of the message. |
| conversation_id | String | 737999610479815**** | The ID of the conversation that this message belongs to. |
| bot_id | String | 747363834493437**** | The ID of the agent that returns this message. This parameter is only returned in messages generated from the chat. |
| chat_id | String | 757363834493437**** | Chat ID。 This parameter is only returned in messages generated from the chat. <br> In different chats, the system will generate a new `chat_id`. For the same user, the `chat_id` will be different between the first chat and the second chat. |
| meta_data | JSON Map | {} | Additional message data when creating a message, which will also be returned when retrieving the message. |
| role | String | user | The entity that sends this message. Valid values: <br>  <br> * **user**: indicates that the content of the message is sent by the user. <br> * **assistant**: Indicates that the content of this message is sent by the Bot. |
| content | String | Good morning, what day is it today? | The content of the message, supporting multiple content types such as plain text, multi-modal (text, image, file mixed input) content, and cards. |
| content_type | String | text | The type of message content includes: <br>  <br> * text: text format. <br> * object_string: Multi-modal content, a combination of text and files or text and images. <br> * card: card. This enumeration value only appears in API responses and does not support being used as an input parameter. |
| created_at | Long | 1718592898 | The time when the message was created, represented as a 10-digit Unix timestamp in seconds (s). |
| updated_at | Long | 1718592898 | The update time of the message, formatted as a 10-digit Unix timestamp in seconds (s). |
| type | String | question | Message type. <br>  <br> * **question**: User input content. <br> * **answer**: Message content returned by the agent to the user, supports incremental returns. If the workflow is bound to a messge node, there may be multiple answers. The end flag of the stream response can be used to determine when all answers are completed. <br> * **function_call**: Intermediate result of a function call during the chat with the agent. <br> * **tool_response**: The result returned after calling a tool (function call). <br> * **follow_up**: If the user suggestion switch is enabled on the bot, the reply content related to recommended questions will be returned. <br> * **verbose**: In the multi-answer scenario, the server will return a verbose package, and the corresponding content is in JSON format. `content.msg_type =generate_answer_finish` indicates that all answer replies are completed. <br>  <br> Only the conversation (v3) API supports using this parameter as an input argument, and: <br>  <br> * If autoSaveHistory=true, type can be set to question or answer. <br> * If autoSaveHistory=false, type can be set to question, answer, function_call, tool_output/tool_response. <br>    Among them, type=question can only correspond to role=user, meaning only the user role can and must initiate question-type messages. <br>  |
| section_id | String | 757363834493437**** | Context fragment ID. A new section_id is generated each time the context is cleared. |
| reasoning_content | String | Okay, I now need to provide study advice for a 13-year-old college student. First, I need to consider the user's situation... | DeepSeek-R1 model's chain of thought (CoT). The model will decompose complex problems step by step into multiple simple steps and derive the final answer based on these steps. <br> This parameter will only be returned when using the DeepSeek-R1 model. |

**Example**

**Request example**

```shell
curl --location --request POST {{host}}/v1/conversation/message/list?&conversation_id=737363834493434**** \
--data-raw '{
    "limit": null,
    "order": "asc",
    "chat_id": "737363834493437****",
    "before_id": "737363834493437****"
}'
```


**Response example**

```json
{
    "code": 0,
    "data": [
        {
            "bot_id": "737363834493434****",
            "chat_id": "747363834493434****",
            "content": "你的名字叫什么",
            "content_type": "text",
            "conversation_id": "737363834493434****",
            "created_at": 1716809829,
            "id": "737363834493437****",
            "meta_data": {},
            "role": "user",
            "type": "",
            "updated_at": "1716809829"
        },
        {
            "bot_id": "737363834493434****",
            "chat_id": "747363834493434****",
            "content": "我的名字叫bot",
            "content_type": "text",
            "conversation_id": "737363834493434****",
            "created_at": "1716809829",
            "id": "737363834493440****",
            "meta_data": {},
            "role": "assistant",
            "type": "",
            "updated_at": "1716936779"
        }
    ],
    "first_id": "737363834493437****",
    "has_more": true,
    "last_id": "737363834493440****",
    "msg": ""

}
```


# Clear context

Clear the context in a specified conversation.

**API description**

In the agent chat scenario, context refers to the historical messages and chat records that the model can reference while processing the current input. It determines the model's depth of understanding and coherence in the current task, directly impacting the accuracy and relevance of the output. In multi-turn chats, if you need to start a new topic, you can call this API to clear the context. After clearing the context, historical messages in the chat will not be input as context to the model, and subsequent chats will no longer be influenced by previous historical messages. This avoids information interference and ensures the accuracy and coherence of the chat.
Messages in the conversation are stored in context sections. Each time this API is called to clear the context, the system will automatically delete the old context sections and create new ones to store new messages. When starting a new chat, messages in the new section will be passed to the model as the new context for reference.
* The Clear Context API only removes the message history visible to the model but does not actually delete messages in the conversation. Messages can still be viewed through APIs like the View Message List.
* Only the context of conversations created by oneself can be cleared.
* For the explanation of terms such as conversation, chat, message, and context paragraphs, please refer to [Basic Concepts](https://github.com/kozex-ai/kozex/wiki#fed4a54c).


**Basic information**

| **HTTP method** | POST |
| --- | --- |
| **URI** | &#123;&#123;host&#125;&#125;/v1/conversations/:conversation_id/clear |
| **API description** | Clear the context in the specified conversation. |

**Request parameters**

**Header**

| **Parameter** | **Value** | **Description** |
| --- | --- | --- |
| Authorization | Bearer *$Access_Token* | The **Personal Access Token** used to authenticate the identity of the client. You can generate a Personal Access Token on the Kozex platform. For more details, refer to [Preparation]. |
| Content-Type | application/json | Indicates the format in which the request body is interpreted. |

**Path**

| **Parameter** | **Type** | **Required** | **Example** | **Description** |
| --- | --- | --- | --- | --- |
| conversation_id | String | Required | 737989918257**** | The conversation ID whose context needs to be cleared. You can obtain the conversation ID through the conversation_id field in the Response of the [chat initiation](https://www.coze.cn/docs/developer_guides/chat_v3) API. |

**Response parameters**

| **Parameter** | **Type** | **Example** | **Description** |
| --- | --- | --- | --- |
| code | Long | 0 | Status code. 0 means the call was successful, other values indicate failure. You can use the `msg` field to determine the detailed error reason. |
| msg | String | "Success" | This field provides detailed information about the result of the API call. If the API call fails, msg contains a description of the error or the reason for the failure. |
| data | Object of [Section](#section) | { "id": "1234567890123456789", "conversation_id": "1234567890123456789" } | Detailed information about the contextual paragraph (section). <br> A Section is an independent contextual paragraph used to separate different chat stages or topics. The Section includes contextual messages; when the user clears the context, the system will create a new Section to ensure the new chat is not affected by historical messages. |

**Section**

| **Parameter** | **Type** | **Example** | **Description** |
| --- | --- | --- | --- |
| id | String | 1234567890123456789 | Session ID, which is the unique identifier for a newly created context section after clearing the context. <br> Each context section corresponds to a batch of independent context messages. Every time the context is cleared, the system creates a new context section to store new context messages. |
| conversation_id | String | 737999610479815**** | Conversation ID, which is the unique identifier of the conversation. |

**Example**

**Request example**

```shell
curl --location --request POST {{host}}/v1/conversations/:conversation_id/clear \
--header 'Authorization: Bearer pat_OYDacMzM3WyOWV3Dtj2bHRMymzxP****' \
--header 'Content-Type: application/json' \
--data-raw '{}'
```


**Response example**

```json
{
    "code": 0,
    "msg": "",
    "data": {
        "id": "12345678****56789",
        "conversation_id": "1234567****456789"
    }
}
```


# Execute workflow

Execute the published workflow.

**API description**

This API operates in a non-streaming response mode. For nodes supporting streaming output, use the **Execute workflow (streaming response)** API to obtain streaming responses. After calling the API, you can obtain the debug_url from the response. By accessing the link, you can view the trial run process of the workflow through a visual interface, including detailed information on the input and output of each execution node, helping you debug or troubleshoot online.

**Restrictions**

| **Restrictions** | **Description** |
| --- | --- |
| Workflow publish status | Must be published. Executing an unpublished workflow will return error code 4200. The operations to create and publish a workflow can refer to [using workflows](https://www.coze.cn/docs/guides/use_workflow). |
| Node limitations | The workflow cannot include message nodes, end nodes with streaming output enabled, or question nodes. |
| Associate agent | Before calling this API, you should first test-run this workflow on the Kozex platform. If the test run requires an associated agent, you also need to specify the agent ID when calling this API to execute the workflow. In general, executing workflows with database nodes, variable nodes, and similar nodes requires an associated agent. |
| Request size limit | 20 MB, including input parameters, message history generated during operation, and all related data. |
| Timeout duration | * When asynchronous operation of the workflow is not enabled, the overall workflow timeout is 10 minutes. It is recommended to control the execution time within 5 minutes to ensure result accuracy. For detailed explanation, please refer to [workflow usage limitations](https://www.coze.cn/docs/guides/workflow_limits). <br> * After enabling workflow asynchronous operation, the overall workflow timeout is 24 hours. |

**Basic information**

| **HTTP method** | POST |
| --- | --- |
| **URI** | &#123;&#123;host&#125;&#125;/v1/workflow/run |
| **API description** | Execute the published workflow. |

**Request parameters**

**Header**

| **Parameter** | **Value** | **Description** |
| --- | --- | --- |
| Authorization | Bearer *$Access_Token* | Used to verify the identity of the **Personal Access Token**. You can generate a Personal Access Token on the Kozex platform. For detailed information, refer to 【Preparation Work】. |
| Content-Type | application/json | Indicates the format in which the request body is interpreted. |

**Body**

| **Parameter** | **Type** | **Required** | **Example** | **Description** |
| --- | --- | --- | --- | --- |
| workflow_id | String | Required | 73664689170551***** | The Workflow ID to be executed, this workflow should already be published. <br> Go to the Workflow build page, and the number after the `workflow` parameter in the page URL is the Workflow ID. For example, `https://www.coze.com/work_flow?space_id=42463***&workflow_id=73505836754923***`, the Workflow ID is `73505836754923***`. |
| parameters | Map[string]any <br>  | Optional | { <br> "user_id":"12345", <br> "user_name":"George" <br> } | Input parameters and values for the workflow's start node can be viewed on the build page of the specified workflow. <br> If the input parameter of the workflow is a file type such as Image, you can call the [file upload](https://www.coze.cn/open/docs/developer_guides/upload_files) API to obtain the file_id. When calling this API, pass the file_id in the parameters in serialized JSON format. For example, `"parameters" : { "input": "{\\\"file_id\\\": \\\"xxxxx\\\"}" }`. |
| bot_id | String | Optional | 73428668***** | The agent ID to be associated. Some workflows require specifying the associated agent, such as workflows containing database nodes, variable nodes, and other nodes. <br> ![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/55746fa5540b488ea83a79064a223500~tplv-goo7wpa0wc-image.image) <br> Enter the development page of the agent, the number following the bot parameter in the URL of the development page is the agent ID. For example `https://www.coze.com/space/341****/bot/73428668*****`, the agent ID is `73428668*****`. <br> Ensure that the token used to call this API has the permission of the workspace where the agent is located. <br> Ensure that the agent has been published as an API service. <br>  |
| app_id | String | Optional | 749081945898306**** | The ID of the associated app in the workflow |

**Response parameters**

| **Parameter** | **Type** | **Example** | **Description** |
| --- | --- | --- | --- |
| code | Long | 0 | Status code. <br>  <br> * A value of 0 indicates a successful call. <br> * Any other value indicates a failed call. You can determine the detailed error reason through the msg field. |
| msg | String | Success | This field provides detailed information about the result of the API call. If the API call fails, msg contains a description of the error or the reason for the failure. |
| data | String |  | Workflow execution results, typically serialized as JSON strings. In some scenarios, it may return strings with non-JSON structures. |
| execute_id | String | 741364789030728**** | Event ID for asynchronous execution. |
| token | Long |  | Reserved field, no need to pay attention. |
| cost | String | 0 | Reserved field, no need to pay attention. |
| debug_url | String |  | Workflow trial run debugging page. Visit this page to view the execution results, inputs and outputs of each workflow node. |

**Example**

**Request example**

```shell
curl --location --request POST 'https://{{host}}/v1/workflow/run' \
--header 'Authorization: Bearer pat_hfwkehfncaf****' \
--header 'Content-Type: application/json' \
--data-raw '{
    "workflow_id": "73664689170551*****",
    "parameters": "{\"user_id\":\"12345\",\"user_name\":\"George\"}"
}'
```


**Response example**

```json
{
    "code": 0,
    "cost": "0",
    "data": "{\"output\":\"北京的经度为116.4074°E，纬度为39.9042°N。\"}",
    "debug_url": "https://www.coze.cn/work_flow?execute_id=741364789030728****&space_id=736142423532160****&workflow_id=738958910358870****",
    "msg": "Success",
    "token": 98
}
```

# Execute workflow (streaming response)

Execute the published workflow, response mode is streaming response.

**API description**

When invoking the API to execute a workflow that supports streaming output, it is often necessary to use a streaming response mode to receive response data, such as real-time display of workflow output information or presenting a typing effect.
In a streaming response mode, the server does not send all data at once but sends data piece by piece in the form of a data stream to the client. The data stream includes various events triggered during workflow execution until processing is complete or interrupted. After processing is completed, the server will notify the workflow execution completion through the `event: Done` event.
The workflow nodes currently supporting streaming responses include **output nodes**, **question nodes**, and **end nodes with streaming output enabled**. For workflows without these nodes, you can use the **Execute Workflow** API to receive response data at once.


**Restrictions**


* Before executing a workflow through the API method, please confirm that the workflow has been published. If attempting to execute an unpublished workflow, an error code 4200 will be returned.
* Before calling this API, you should first run this workflow on the Kozex platform.
   * If an agent needs to be associated during trial run, bot_id must also be specified when calling this API to execute the workflow. Typically, executing workflows with nodes like database nodes or variable nodes requires associating an agent.
   * When executing a workflow in the app, app_id must be specified.
   * Do not specify both bot_id and app_id simultaneously, otherwise the API will throw error 4000, indicating a request parameter error.
* This API is a synchronous API. If the workflow as a whole or certain nodes timeout during execution, the agent may fail to provide the expected response. It is recommended to limit the execution time of the workflow to within 5 minutes.
* The maximum request size supported by the workflow is 20MB, including input parameters and all related data such as message history generated during execution.


**Basic information**

| **HTTP method** | POST |
| --- | --- |
| **URI** | &#123;&#123;host&#125;&#125;/v1/workflow/stream_run |
| **API description** | Execute the published workflow with a streaming response. |

**Header**

| **Parameter** | **Value** | **Description** |
| --- | --- | --- |
| Authorization | Bearer <span style="color: #D83931"><em>$Access_Token</em></span> | Personal Access Token for verifying the identity of the **client**. You can generate a Personal Access Token on the Kozex platform. For more details, refer to 【Preparations】. |
| Content-Type | application/json | Indicates the format in which the request body is interpreted. |

**Body**

| **Parameter** | **Type** | **Required** | **Description** |
| --- | --- | --- | --- |
| workflow_id | String | Required | The Workflow ID to be executed, and this workflow should have been published. <br> Navigate to the `workflow` build page, and in the page URL, the number after the workflow parameter is the Workflow ID. For example, `https://www.coze.com/work_flow?space_id=42463***&workflow_id=73505836754923***`, the workflow ID is `73505836754923***`. |
| parameters | Map[string]any | Optional | The input parameters and values of the workflow's start node can be viewed in the parameter list on the build page of the specified workflow. |
| bot_id <br>  | String <br>  | Optional <br>  | The agent ID to be associated. Some workflows require specifying the associated agent during execution, such as workflows containing database nodes, variable nodes, etc. <br> ![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/55746fa5540b488ea83a79064a223500~tplv-goo7wpa0wc-image.image) <br> Enter the build page of the agent, the number following the bot parameter in the URL of the build page is the agent ID. For example, in the URL `https://www.coze.com/space/341****/bot/73428668*****`, the Bot ID is `73428668*****`. <br> * Ensure that the token used to call this API has the permission for the workspace where this agent is located. <br> * Ensure that the agent has been published as an API service. <br>  |
| app_id | String | Optional | The application ID where the workflow resides. <br> You can get the app ID from the URL of the app's build page. The app ID is the string of characters following the project-ide parameter in the URL. For example, in the URL https://www.coze.cncom/space/739174157340921****/project-ide/743996105122521****/workflow/744102227704147****, the app ID is 743996105122521 **** . <br> Setting app_id is only required when running workflows in the Coze app. Workflows bound to the agent or workflows in the workspace repository do not require app_id to be set. <br>  |

**Response**


In streaming output responses, developers need to pay attention to whether there is any packet loss.

* The event ID (id) starts counting from 0 by default, and ends with the event that includes `event: Done`. Developers should confirm that the response message has no packet loss based on the id.
* The message ID of the Message event starts counting from 0 by default, and ends with the event that includes `node_is_finish : true`. Developers should confirm that each message in the Message event has been completely returned with no packet loss, based on node_seq_id.

| **Parameter name** | **Parameter type** | **Parameter description** |
| --- | --- | --- |
| id | Integer | The event ID of this message in the API response. Start with 0 |
| event | String | The current data packet event of the streaming response. Including the following types: <br>  <br> * Message: Indicates output messages from workflow nodes, such as message nodes or end nodes. You can view the specific message content in the data field. <br> * Error: Indicates an error occurrence. You can view error_code and error_message in the data field to troubleshoot the issue. <br> * Done: Indicates completion. This event signifies that the workflow has finished execution, and the data field contains a debug URL. <br> * Interrupt: Indicates an interruption. When this event is observed, it indicates that the workflow is interrupted, and the data field contains specific interrupt information. <br> * PING: Heartbeat signal. This event signifies that the workflow is processing, the message content is empty, and it is used to maintain the connection. |
| data | Object | Event content. The event content format varies for different event types. |

**Message event**

The structure of the data in the Message event is as follows:
| **Parameter name** | **Parameter type** | **Parameter description** |
| --- | --- | --- |
| content | String | Message content for streaming output. |
| node_title | String | Node name of the output message, such as message node or end node. |
| node_seq_id | String | The message ID in this node starts counting from 0, for example, the 5th message in the message node. |
| node_is_finish | Boolean | Is the current message the last packet for this node? |
| ext | Map[String]String | Additional fields. |
| cost | String | Reserved field, no need to focus on. |

**Interrupt event**

The structure of the data in the Interrupt event is as follows:
| **Parameter name** | **Parameter type** | **Parameter description** |
| --- | --- | --- |
| interrupt_data | Object | Interruption control content. |
| interrupt_data.event_id | String | Workflow interruption event ID, this field should be passed back when resuming operation. |
| interrupt_data.type | Integer | Interrupt type in the workflow, this field should be passed back when resuming operation. |
| node_title | String | The name of the output message node, such as "question". |

**Error event**

The structure of the data in the Error event is as follows:
| **Parameter name** | **Parameter type** | **Parameter description** |
| --- | --- | --- |
| error_code | Integer | Status code. <br>  <br> * A value of 0 indicates a successful call. <br> * Any other value indicates a failed call. You can determine the detailed error reason through the error_message field. |
| error_message | String | This field provides detailed information about the result of the API call. If the API call fails, msg contains a description of the error or the reason for the failure. |

**Example**

**Request example**

```Shell
curl --location --request POST 'https://{{host}}/v1/workflow/stream_run' \
--header 'Authorization: Bearer pat_fhwefweuk****' \
--header 'Content-Type: application/json' \
--data-raw '{
    "workflow_id": "73664689170551*****",
    "{\"user_id\":\"12345\",\"user_name\":\"George\"}"
}'
```


**Response example**


* Message event
   ```Plain Text
   id: 0
   event: Message
   data: {"content":"msg","node_is_finish":false,"node_seq_id":"0","node_title":"Message"}
   
   id: 1
   event: Message
   data: {"content":"为","node_is_finish":false,"node_seq_id":"1","node_title":"Message"}
   
   id: 2
   event: Message
   data: {"content":"什么小明要带一把尺子去看电影？\n因","node_is_finish":false,"node_seq_id":"2","node_title":"Message"}
   
   id: 3
   event: Message
   data: {"content":"为他听说电影很长，怕","node_is_finish":false,"node_seq_id":"3","node_title":"Message"}
   
   id: 4
   event: Message
   data: {"content":"坐不下！","node_is_finish":true,"node_seq_id":"4","node_title":"Message"}
   
   id: 5
   event: Message
   data: {"content":"{\"output\":\"为什么小明要带一把尺子去看电影？\\n因为他听说电影很长，怕坐不下！\"}","cost":"0.00","node_is_finish":true,"node_seq_id":"0","node_title":"","token":0}
   
   id: 6
   event: Done
   data: {}
   ```

* Error event
   ```Plain Text
   id: 0
   event: Error
   data: {"error_code":4000,"error_message":"Request parameter error"}
   ```

* Interrupt event
   ```Plain Text
   // Streaming execution workflow, triggers the question node, and the Bot asks a question
   id: 0
   event: Message
   data: {"content":"请问你想查看哪个城市、哪一天的天气呢","content_type":"text","node_is_finish":true,"node_seq_id":"0","node_title":"问答"}
   
   id: 1
   event: Interrupt
   data: {"interrupt_data":{"data":"","event_id":"7404830425073352713/2769808280134765896","type":2},"node_title":"问答"}
   ```


# Resume running the workflow

Resume running the interrupted workflow.

**API description**

When executing a workflow that contains a question node, the agent will ask the user designated questions and wait for their response. At this time, the workflow is in an interrupted state, and the developer needs to call this API to answer the question and resume running the workflow. If the user's response does not match the information expected to be extracted by the agent, such as missing required fields or inconsistent data types, the workflow will interrupt again and inquire further. If the expected response is not received after three attempts, the workflow execution is deemed to have failed.

**Restrictions**


* The API can be called to resume running up to three times. If, during the third attempt, the agent still does not receive the expected response, the workflow execution is deemed to have failed.
* After resuming operation, both the index and the node index will be reset.
* Resuming operation of the workflow also incurs token consumption, and the number of tokens consumed is identical to that during **workflow execution (streaming response)**.


**Basic information**

| **HTTP method** | POST |
| --- | --- |
| **URI** | &#123;&#123;host&#125;&#125;/v1/workflow/stream_resume |
| **API description** | Resume running the interrupted workflow. |

**Header**

| **Parameter** | **Value** | **Description** |
| --- | --- | --- |
| Authorization | Bearer <span style="color: #D83931"><em>$Access_Token</em></span> | Used to verify the identity of the **Personal Access Token**. You can generate a Personal Access Token on the Kozex platform. For detailed information, refer to 【Preparation】. |
| Content-Type | application/json | Indicates the format in which the request body is interpreted. |

**Body**

| **Parameter** | **Type** | **Required** | **Description** |
| --- | --- | --- | --- |
| workflow_id | String | Required | The Workflow ID to be executed, this workflow must have been published. <br> Go to the Workflow build page, and the number after the `workflow` parameter in the page URL is the Workflow ID. For example, `https://www.coze.com/work_flow?space_id=42463***&workflow_id=73505836754923***`, the Workflow ID is `73505836754923***`. |
| event_id | String | Required | Workflow execution interrupt event ID. You can obtain the interrupt event ID from the response information of **executing workflows (streamed responses)**. |
| resume_data | String | Required | The user's reply to the agent's specified question when resuming execution. <br> If the interruption is caused by a question node, the reply should contain required parameters from the question node; otherwise, the workflow will interrupt again and prompt for the correct information. |
| interrupt_type | Integer | Required | Interrupt type, you can obtain the interrupt type of the interrupt time from the response information of **executing the workflow (streaming response)**. |

**Response**


In a streaming response, developers need to pay attention to whether there is any packet loss.

* Event ID (id) starts counting from 0 by default and ends with the event containing `event: Done`. Developers should confirm based on the id that there is no packet loss in the overall response message.
* The message ID of the Message event starts counting from 0 by default and ends with the event containing `node_is_finish : true`. Developers should confirm based on the node_seq_id that every message in the Message event is completely returned without packet loss.

| **Parameter name** | **Parameter type** | **Description of the parameter.** |
| --- | --- | --- |
| id | Integer | The event ID in the API response. Start from 0. |
| event | String | The current data packet event of the streaming response. Includes the following types: <br>  <br> * Message: Indicates output messages from workflow nodes, such as message nodes or end nodes. You can view the specific message content in the data field. <br> * Error: Indicates an error occurrence. You can check the error_code and error_message in the data field to troubleshoot the issue <br> * Done: Indicates completion. This event signifies that the workflow has finished execution, and the data field will return a debug_url <br> * Interrupt: Indicates an interruption. This event signifies that the workflow has been interrupted, and the data field contains specific interruption information <br> * PING: Heartbeat signal This event signifies that the workflow is in execution. The message content is empty and is used to maintain the connection |
| data | Object | Event content. The event content format of each event type is different. |

**Message event**

The structure of the data in the Message event is as follows:
| **Parameter name** | **Parameter type** | **Parameter description** |
| --- | --- | --- |
| content | String | Message content for streaming output. |
| node_title | String | The name of the node that outputs the message, such as a message node or an end node. |
| node_seq_id | String | The message ID within the node, starting from 0. For example, the 5th message of a message node. |
| node_is_finish | Boolean | Whether the current message is the last data packet of this node. |
| ext | Map[String]String | Extra fields. |
| cost | String | Reserved field, no need to pay attention. |

**Interrupt event**

The structure of the data in the Interrupt event is as follows:

| **Parameter name** | **Parameter type** | **Parameter description** |
| --- | --- | --- |
| interrupt_data | Object | Interrupt control content. |
| interrupt_data.event_id | String | The workflow interruption event ID, this field should be returned when resuming operation. |
| interrupt_data.type | Integer | The workflow interruption type, this field should be returned when resuming operation. |
| node_title | String | The name of the node that outputs the message, such as "Q&A". |

**Error event**

The structure of data in the Error event is as follows:
| **Parameter name** | **Parameter type** | **Description of the parameter** |
| --- | --- | --- |
| error_code | Integer | Status code. <br>  <br> * A value of 0 indicates a successful call. <br> * Any other value indicates a failed call. You can determine the specific error cause through the error_message field. |
| error_message | String | This field provides detailed information about the result of the API call. If the API call fails, msg contains a description of the error or the reason for the failure. |

**Example**

In the case of workflow execution interruption, taking the weather-checking workflow as an example, a complete API call example is as follows.

1. Call API **to execute the workflow (streamed response)**, requesting to check the weather.
   An request example is as follows:



   ```Plain Text
   curl --location 'https://{{host}}/v1/workflow/stream_run' \
   --header 'Authorization: Bearer pat_vTG1****' \
   --header 'Content-Type: application/json' \
   --data '{
       "workflow_id": "739739507914235****",
       "parameters": "{\"BOT_USER_INPUT\":\"查看天气\"}"
   }'
   ```

2. When a Question node is triggered and the workflow is interrupted, the response will return the question raised by the agent, requesting the user to provide the city and date.



   A response example is as follows:
   ```Plain Text
   id: 0
   event: Message
   data: {"content":"请问你想查看哪个城市、哪一天的天气呢","content_type":"text","node_is_finish":true,"node_seq_id":"0","node_title":"问答"}
   
   id: 1
   event: Interrupt
   data: {"interrupt_data":{"data":"","event_id":"7404831988202520614/6302059919516746633","type":2},"node_title":"问答"}
   ```

3. Call API to resume the workflow and reply to the agent with the city and date.
   An request example is as follows:



   ```Plain Text
   curl --location 'https://{{host}}/v1/workflow/stream_resume' \
   --header 'Authorization: Bearer pat_vTG1****' \
   --header 'Content-Type: application/json' \
   --data '{
       "event_id":"740483727529459****/433802199567434****",
       "interrupt_type":2,
       "resume_data":"杭州，2024-08-20",
       "workflow_id":"739739507914235****"
   }'
   ```

4. Once the workflow is executed and the weather query is completed, the workflow will return the output message.
   A response example is as follows:
   ```Plain Text
   id: 0
   event: Message
   data: {"content":"{\"output\":[{\"condition\":\"中到大雨\",\"humidity\":72,\"predict_date\":\"2024-08-20\",\"temp_high\":35,\"temp_low\":26,\"weather_day\":\"中到大雨\",\"wind_dir_day\":\"西风\",\"wind_dir_night\":\"西风\",\"wind_level_day\":\"3\",\"wind_level_night\":\"3\"}]}","content_type":"text","cost":"0","node_is_finish":true,"node_seq_id":"0","node_title":"End","token":386}
   
   id: 1
   event: Done
   data: {}
   ```

# Execute chat flow (streaming response)
Execute the published chat flow with the response method set to streaming response.
## **API description**

* A chatflow is a specialized workflow designed for chat scenarios, efficiently managing chat requests. Chatflows interact with users through chats and manage complex business logic. By adding a chatflow to your app, you can break down user instructions into step-by-step nodes and create an intuitive user interface. This allows you to build a conversational AI app for mobile or web, enabling automated and intelligent chat interactions. For more information about chatflows, see [Workflow and chatflow](/guides/workflow_and_chatflow).
* This API operates in streaming response mode, enabling the client to process data as it is received, rather than waiting for the entire data stream. For example, it can display replies in the chat interface in real-time, minimizing the client's waiting time for a complete response.
* The API also supports nodes that can interrupt the chat, such as question nodes and input nodes. If the chat is interrupted, you can resume the conversation by calling the chatflow again and specifying the input content in additional_messages to continue the chat.

**If your chatflow input includes multimodal content like files and images, first upload the content to a third-party storage tool to obtain a publicly accessible URL. Use this URL as the input for the chatflow.**
When calling the API, you'll receive a response with a debug_url in the Done event. Open this link in your browser to access a visual interface that displays the test run of the chatflow. This interface provides detailed insights into the input and output of each execution node, making online debugging and troubleshooting straightforward.


This API can be used to call chat flows in the workspace resource library or in the Coze app. When invoking these two types of chat flows, the input parameters are different:
| **Input parameters** | **Resource library chat flows** |  | **Chat flow in the Coze app** |
| --- | --- | --- | --- |
|  | **Execute in an agent** | **Execute in the Coze app** |  |
| workflow_id | Required | Required | Required |
| app_id | Do not send | Required | Required |
| bot_id | Required | Leave blank | Do not send |
| conversation_id | Optional | Optional | Optional |

## **Basic information**
| **HTTP method** | POST |
| --- | --- |
| **URI** | &#123;&#123;host&#125;&#125;/v1/workflows/chat |
| **API description** | Execute the published chat flow in streaming response mode. |
### Header
| **Parameter** | **Value** | **Note** |
| --- | --- | --- |
| Authorization | Bearer *$Access_Token* | A **Personal Access Token** used to verify client identity. You can generate a Personal Access Token on the Kozex platform. For more details, see before you begin. |
| Content-Type | application/json | Indicates the format in which the request body is interpreted. |
### **Body**
| **Parameter** | **Type** | **Required** | **Note** |
| --- | --- | --- | --- |
| workflow_id | String | Required | Workflow ID to be executed. This workflow should already be published. <br> Go to the Workflow build page. In the page URL, the number after the `workflow` parameter is the Workflow ID. For example, `https://www.coze.com/work_flow?space_id=42463***&workflow_id=73505836754923***`, where the Workflow ID is `73505836754923***`. |
| additional_messages | array&lt;Object&gt; <br>  <br> * Object | Required | User questions and historical messages in chat. The array length is limited to 50, which means a maximum of 50 messages can be passed. |
| parameters | Object | Required | Set custom parameters within the chat flow's input parameters. You can pass custom parameters and their corresponding values through the `parameters` argument based on your business requirements. <br>  <br> * The USER_INPUT parameter for the chat flow should be passed in additional_messages. USER_INPUT in parameters does not take effect. <br> * If CONVERSATION_NAME or custom input parameters are not specified in parameters, the chat flow runs using the default parameter values; if these parameters are specified, the specified values are used. |
| app_id | String | Optional | The Coze app ID to be associated. When invoking a chat flow, you must specify either app_id or bot_id so that the model can access the agent's or app's database, variables, and other data processing tasks. <br> Go to the app development interface. The number following the project-ide parameter in the development page URL is the AppID. For example, `https://www.coze.cn/space/74421656*****/project-ide/744208683**`, the Coze app ID is `744208683**`**.** |
| bot_id <br>  | String <br>  | Optional <br>  | Agent ID to be associated. Some workflows require specifying the associated agent during execution, such as workflows that include database nodes, variable nodes, and other nodes. <br> ![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/55746fa5540b488ea83a79064a223500~tplv-goo7wpa0wc-image.image) <br> Enter the build page of the agent, the number following the bot parameter in the URL of the build page is the agent ID. For example, `https://www.coze.com/space/341****/bot/73428668*****`, the Bot ID is `73428668*****`. <br> * Ensure that the token used to call this API has the required permissions for the workspace where this agent resides. <br> * Ensure the agent has been published as an API service. <br>  |
| conversation_id | String | Optional | The conversation ID associated with this chatflow. All generated messages are stored within this specific conversation. To specify a conversation, you can either use the default CONVERSATION_NAME set at the start node or define a conversation using the conversation_id parameter. <br>  <br> * If you specify conversation_id, the CONVERSATION_NAME set in parameters does not take effect. <br> * The creator of the conversation must be the same as the user executing the chatflow, specifically the creator of the API access token, otherwise, the chat flow cannot be executed. <br> * Conversations are matched with app_id and channels, and conversations from different channels are isolated. <br> * When bot_id is specified, if conversation_id is not provided, Coze will create a new conversation. Specifying both bot_id and app_id at the same time is not supported. |
| workflow_version | String | Optional | The version number of the workflow is only valid when the running workflow is a repository workflow. If no version number is specified, the latest version of the workflow is executed by default. <br>  |
additional_messages object
| Parameter | Type | Is it required? | Note |
| --- | --- | --- | --- |
| content | String | Required | The content of a message. Only plain text is supported. <br> Currently, multi-modal input (text, image, file mixed input), cards, and other types of content are not supported. <br>  |
| content_type | String | Required | Type of message content. <br> content_type is fixed as text, indicating plain text. <br> When specifying content, you should also set content_type. |
| role | String | Required | * **user**: indicates that the content of the message is sent by the user. <br> * **assistant**: indicates that the message content sent by the model. |
| type |  |  | Message type. The default value is **question.** <br>  <br> * **question**: User input content. <br> * **answer**: Message content returned by the model to the user, supports incremental returns. If the chatflow is bound to a message node, there may be multiple answers. The end flag of the stream response can be used to determine when all answers are completed. |
## Response parameters

* Streaming responses allow the client to start processing data before receiving the complete data stream, for example, displaying response content in a chat interface in real time, reducing the client's wait time for a complete response from the model.

The overall process of a streaming response is as follows:
### **Streaming response process**

   ```Python
   Overview (two levels: chat and MESSAGE)
   Chat - start
   # chat - In progress
   MESSAGE - Knowledge base retrieval
   #   MESSAGE - function_call
   #   MESSAGE - tool_output
   #   MESSAGE - answer is normal text
   #   MESSAGE - When there are multiple answers, message.delta will continue to be present
   Completed
   Flow end event: done
   #########
   ```


#### **List of streaming response events**

   | **Event** | **Note** |
   | --- | --- |
   | conversation.chat.created | Event indicating the creation and start of a new chat. |
   | conversation.chat.in_progress | Event indicating that the server is currently processing the chat. |
   | conversation.message.delta | Event for incremental message updates, typically occurring when type=answer. |
   | conversation.message.completed | Event indicating that the message has been fully replied. The streaming packet at this point contains the concatenated results of all message.deltas, and each message is marked as completed. |
   | conversation.chat.completed | Event indicating that the chat has been completed. |
   | conversation.chat.failed | Event indicating that the chat has failed. |
   | conversation.chat.requires_action | Event indicating that the chat is interrupted and requires user action to report the execution result of a tool. This usually occurs when a question node or input node is triggered, requiring the API to be called again to continue the chat. |
   | error | Event indicating an error during the streaming response process. |
   | done | Event indicating that the streaming response of this conversation has ended normally. |

#### **Streaming response example**

   ```Python
   Chat - Start
   event: conversation.chat.created
   data: {"id":"120","conversation_id":"456","created_at":1733407180,"last_error":{"code":0,"msg":""},"status":"created","usage":{"token_count":0,"output_count":0,"input_count":0},"section_id":"789"}
   Chat - In progress
   event: conversation.chat.in_progress
   data: {"id":"121","conversation_id":"456","created_at":1733407180,"last_error":{"code":0,"msg":""},"status":"in_progress","usage":{"token_count":0,"output_count":0,"input_count":0},"section_id":"789"}
   # MESSAGE - answer is normal text
   event: conversation.message.delta
   data: {"id":"122","conversation_id":"456","role":"assistant","type":"answer","content":"中午吃啥了","content_type":"text","chat_id":"567","section_id":"789","created_at":1733407182}
   
   End of message
   event: conversation.message.completed
   data: {"id":"124","conversation_id":"456","role":"assistant","type":"answer","content":"中午吃啥了","content_type":"text","chat_id":"567","section_id":"789","created_at":1733407182}
   
   event: conversation.message.completed
   data: {"id":"125","conversation_id":"456","role":"assistant","type":"verbose","content":"{\"msg_type\":\"interrupt\",\"data\":\"\",\"from_module\":null,\"from_unit\":null}","content_type":"text","chat_id":"567","section_id":"789","created_at":1733407182,"updated_at":1733407182}
   
   event: conversation.message.completed
   data: {"id":"130","conversation_id":"456","role":"assistant","type":"verbose","content":"{\"msg_type\":\"generate_answer_finish\",\"data\":\"{\\\"finish_reason\\\":1,\\\"FinData\\\":\\\"\\\"}\",\"from_module\":null,\"from_unit\":null}","content_type":"text","chat_id":"567","section_id":"789","created_at":1733407182,"updated_at":1733407182}
   # Chat - Action required (interruption, usually triggered by a question node or input node)
   event: conversation.chat.requires_action
   data: {"id":"131","conversation_id":"456","created_at":1733407180,"completed_at":1733407182,"last_error":{"code":0,"msg":""},"status":"requires_action","usage":{"token_count":0,"output_count":0,"input_count":0},"required_action":{"type":"submit_tool_outputs","submit_tool_outputs":{"tool_calls":[{"id":"","type":"reply_message","function":null,"require_info":null}]}},"section_id":"789"}
   
   event: done
   data: {"debug_url":"http://{{HOST}}/work_flow?execute_id=74449256856****\u0026space_id=7442165654356*****\u0026workflow_id=744224337778*****"}
   ```




* #### **Event message body structure**
   | **Parameter** | **Type** | **Note** |
   | --- | --- | --- |
   | event | String | The current data packet event of the streaming response. In a streaming response, the server sends data incrementally to the client, formatted as a data stream. This stream includes various events triggered during the chat until the process is either completed or interrupted. Upon completion, the server returns the concatenated and complete reply from the model through the conversation.message.completed event. You can refer to the table below for these events. |
   | data | Object | Message content. The formats of the chat event and the message event differ from each other. <br>  <br> * In the chat event, the data is represented by a **Chat Object**. <br> * In the message event, the data is a **Message Object**. |
* ##### **Chat Object**
   | Parameter | Type | Required | Note |
   | --- | --- | --- | --- |
   | id | String | Required | Chat ID, which is the unique identifier of the chat. |
   | conversation_id | String | Required | Conversation ID, which is the unique identifier of the conversation. |
   | bot_id | String | Required | The agent ID for the conversation. |
   | status <br>  | String | Required | The status of a running chat. Valid values: <br>  <br> * created: Chat is created. <br> * in_progress: The agent is processing the chat. <br> * completed: The agent has completed processing the chat. This chat has ended. <br> * failed: the chat failed. <br> * requires_action: The chat was interrupted and requires further action. <br> * canceled: The chat is canceled. |
   | required_action | Object | Optional | Details of the action that needs to be run. |
   | usage | Object <br> ```JSON <br> { <br> "token_count":123, Total number of tokens <br> "output_count":100, Output token consumption <br> "input_count":23 // Enter token <br> } <br> ``` <br>  <br>  | Optional | Detailed information about token consumption. The actual token consumption will be calculated and provided once the chat ends. |

##### Message Object
| Parameter | Type | Note |
| --- | --- | --- |
| id | String | Message ID, the unique identifier of the message. |
| conversation_id | String | The ID of the conversation that this message belongs to. |
| bot_id | String | The ID of the agent that returns this message. This parameter is only returned in messages generated from the chat. |
| chat_id | String | Chat ID。 This parameter is only returned in messages generated from the chat. |
| role | String | The entity that sends this message. Valid values: <br>  <br> * **user**: indicates that the content of the message is sent by the user. <br> * **assistant**: Indicates that the message content is sent by the agent. |
| content | String | The content of the message, supporting multiple content types such as plain text, multi-modal (text, image, file mixed input) content. |
| content_type | String | The type of message content includes: <br>  <br> * text: text format. <br> * object_string: Multi-modal content, a combination of text and files or text and images. |
| type | String | Message type. <br>  <br> * **question**: User input content. <br> * **answer**: Message content returned by the agent to the user, supports incremental returns. If the chatflow is bound to a message node, there may be multiple answers. The end flag of the stream response can be used to determine when all answers are completed. |

















