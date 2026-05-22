# 准备工作
在调用 Kozex 社区版 API  前，确保你已经将智能体发布为了 API 服务，并通过不同的授权方式获取了访问令牌。
## 发布智能体为 API 服务
智能体发布为 API 服务之后，才能通过调用 API 的方式使用这个智能体，例如查看智能体的基本设置、发起一个智能体对话等。
操作步骤如下：

1. 登录 Kozex 社区版。
2. 在**项目开发**页面，选择智能体。
3. 在页面右上角，单击**发布**。
4. 在**发布**页面，选择 **API** 选项，然后单击**发布**。

   ![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/8b0fd3945757416faab0482a2c927b32~tplv-goo7wpa0wc-image.image)

## 获取访问令牌
Kozex 社区版 API 和 Chat SDK 通过个人访问令牌鉴权。调用 API 之前，你需要先获得访问令牌。

调用扣子 API 时，你需要在 Header 中通过 Authorization 参数指定访问令牌（Access token），扣子服务端会根据访问令牌验证调用方的操作权限。


获取访问令牌的操作步骤如下：

1. 登录 Kozex 社区版。
2. 在页面左下角单击个人头像，并选择 **API 授权**。

   ![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/f5bb30622ad3450088b786cff028a800~tplv-goo7wpa0wc-image.image)
3. 在**个人访问令牌**页面中，单击**添加新令牌**。
4. 在弹出的页面完成以下配置，然后单击**确定**。
   | **配置项** | **说明** |
   | --- | --- |
   | 名称 | 个人访问令牌的名称。 |
   | 过期时间 | 个人访问令牌的有效期时长。令牌过期后将会失效，无法继续用它来调用扣子 API。 <br> 生成令牌后，无法修改过期时间。 |
5. 复制并妥善保存个人访问令牌。
   生成的令牌仅在此时展示一次，请即刻复制并保存。

# 创建会话
创建一个会话。

会话是智能体和用户之间的基于一个或多个主题的问答交互，一个会话包含一条或多条消息。创建会话时，扣子会同时在会话中创建一个空白的上下文片段，用于存储某个主题的消息。后续发起对话时，上下文片段中的消息是模型可见的消息历史。

你可以在创建会话时同步写入消息，消息默认写入到最新的一个上下文片段中，对话时将作为上下文传递给模型。

**基础信息**

| **请求方式** | POST |
| --- | --- |
| **请求地址** | &#123;&#123;host&#125;&#125;/v1/conversation/create |
| **接口说明** | 创建一个会话。 |

**请求参数**

**Header**

| **参数** | **取值** | **说明** |
| --- | --- | --- |
| Authorization | Bearer *$Access_Token* | 用于验证客户端身份的**个人访问令牌**。你可以在扣子平台中生成个人访问令牌，详细信息，参考【准备工作】。 |
| Content-Type | application/json | 解释请求正文的方式。 |

**Body**

| **参数** | **类型** | **是否必选** | **示例** | **说明** |
| --- | --- | --- | --- | --- |
| bot_id | String | 可选 | 730454116184516* | 会话对应的智能体 ID。 <br> 创建会话时指定智能体 ID，后续才能通过 查看会话列表 查看指定智能体 ID 对应的会话列表。 |
| connector_id | String | 可选 | 1024 | 该会话在哪个渠道创建。目前支持如下渠道： <br>  <br> * API：（默认）1024 <br> * ChatSDK：999 |

**返回参数**

| **参数** | **类型** | **示例** | **说明** |
| --- | --- | --- | --- |
| data | Object of [ConversationData](#conversationdata) | 详见返回示例 | 会话的基础信息。详细说明可参考[ConversationData](#conversationdata)。 |
| code | Long | 0 | 调用状态码。0 表示调用成功，其他值表示调用失败，你可以通过 msg 字段判断详细的错误原因。 |
| msg | String | "Success" | 状态信息。API 调用失败时可通过此字段查看详细错误信息。 |

**ConversationData**

| **参数** | **类型** | **示例** | **说明** |
| --- | --- | --- | --- |
| id | String | 737999610479815**** | Conversation ID，即会话的唯一标识。 |
| created_at | Long | 1718289297 | 会话创建的时间。格式为 10 位的 Unixtime 时间戳，单位为秒。 |
| last_section_id | String | 7495664347616952360 | 会话中最新的一个上下文片段 ID。 |

**示例**

**请求示例**

创建空会话：

```shell
curl --location '{{host}}/v1/conversation/create' \
--header 'Authorization: Bearer pat_xxxxx' \
--header 'Content-Type: application/json' \
--data '{
    "bot_id":"7531406778865549312"
}'
```



**返回示例**

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


# 查看会话列表
查看指定智能体的会话列表。
* 仅支持通过此 API 查看智能体在 API 或 Chat SDK 渠道产生的会话。
* 仅支持查询本人创建的会话。


**基础信息**

| **请求方式** | GET |
| --- | --- |
| **请求地址** | &#123;&#123;host&#125;&#125;/v1/conversations  |
| **接口说明** | 获取指定智能体的会话列表。 |

**请求参数**

**Header**

| **参数** | **取值** | **说明** |
| --- | --- | --- |
| Authorization | Bearer *$Access_Token* | 用于验证客户端身份的**个人访问令牌**。你可以在扣子平台中生成个人访问令牌，详细信息，参考【准备工作】。 |
| Content-Type | application/json | 解释请求正文的方式。 |

**Query**

| **参数** | **类型** | **是否必选** | **示例** | **说明** |
| --- | --- | --- | --- | --- |
| bot_id | String | 必选 | 73428668***** | 智能体 ID，获取方法如下： <br> 进入智能体的 开发页面，开发页面 URL 中 `bot` 参数后的数字就是智能体 ID。例如`https://www.coze.cn/space/341****/bot/73428668*****`，智能体 ID 为`73428668*****`。 |
| page_num | Integer | 可选 | 1 | 页码，从 1 开始，默认为 1。 |
| page_size | Integer | 可选 | 40 | 每一页的数据条目个数，默认为 50，最大为 50。 |
| sort_order | String | 可选 | ASC | 会话列表的排序方式： <br>  <br> * **ASC**：按创建时间升序排序，最早创建的会话排序最前。 <br> * **DESC**：（默认）按创建时间降序排序，最近创建的会话排序最前。 |
| connector_id | String | 可选 | 999 | 发布渠道 ID，用于筛选指定渠道的会话。仅支持查看以下渠道的会话： <br>  <br> * （默认）API 渠道，渠道 ID 为 1024。 <br> * Chat SDK 渠道，渠道 ID 为 999。 |

**返回参数**

| **参数** | **类型** | **示例** | **说明** |
| --- | --- | --- | --- |
| code | Long | 0 | 调用状态码。0 表示调用成功，其他值表示调用失败，你可以通过 msg 字段判断详细的错误原因。 |
| msg | String | Success | 状态信息。API 调用失败时可通过此字段查看详细错误信息。 |
| data | Object of [ListConversationData](#listconversationdata) | [ { "created_at": 1731575569, "id": "123456789123456789", "meta_data": {}, } ] | 会话列表的详细。 |

**ListConversationData**

| **参数** | **类型** | **示例** | **说明** |
| --- | --- | --- | --- |
| has_more | Boolean | false | 是否还有更多会话未在本次请求中返回。 <br>  <br> * true：还有更多未返回的会话。 <br> * false：已返回符合筛选条件的全部会话。 |
| conversations | Array of [ConversationData](#conversationdata) | { "created_at": 1731575569, "id": "12345456789*****", "meta_data": {}, } | 会话的详细信息。 |

**ConversationData**

| **参数** | **类型** | **示例** | **说明** |
| --- | --- | --- | --- |
| id | String | 737999610479815**** | Conversation ID，即会话的唯一标识。 |
| created_at | Long | 1718289297 | 会话创建的时间。格式为 10 位的 Unixtime 时间戳，单位为秒。 |
| last_section_id | String | 7495664347616952360 | 会话中最新的一个上下文片段 ID。 |

**示例**

**请求示例**

```shell
curl --location --request curl -X GET {{host}}/v1/conversations?bot_id=73428668*****&page_num=1&page_size=20 \
  -H "Authorization: Bearer pat_OYDacMzM3WyOWV3Dtj2bHRMymzxP****"
```


**返回示例**

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

# 发起对话
调用此接口发起一次对话，支持添加上下文和流式响应。

**接口说明**

发起对话接口用于向指定智能体发起一次对话，支持在对话时添加对话的上下文消息，以便智能体基于历史消息做出合理的回复。
Kozex 社区版发起对话 API 仅支持流式响应，智能体在生成回复的同时，将回复消息以数据流的形式逐条发送给客户端。处理结束后，服务端会返回一条完整的智能体回复。

**创建会话** API 和**发起对话** API 的区别如下：

* 创建会话：主要用于初始化一个新的会话环境。
* 发起对话：
   * 用于在已经存在的会话中发起一次对话。
   * 支持添加上下文和流式响应。
   * 可以基于历史消息进行上下文关联，提供更符合语境的回复。


**基础信息**

| **请求方式** | POST |
| --- | --- |
| **请求地址** | &#123;&#123;host&#125;&#125;/v3/chat |
| **接口说明** | 调用此接口发起一次对话，支持添加上下文和流式响应。 |

**Header**

| **参数** | **取值** | **说明** |
| --- | --- | --- |
| Authorization | Bearer <span style="color: #D83931"><em>$Access_Token</em></span> | 用于验证客户端身份的**个人访问令牌**。你可以在扣子平台中生成个人访问令牌，详细信息，参考【准备工作】。 |
| Content-Type | application/json | 解释请求正文的方式。 |

**Query**

| **参数** | **类型** | **是否必选** | **说明** |
| --- | --- | --- | --- |
| conversation_id | String | 可选 | 标识对话发生在哪一次会话中。 <br> 会话是智能体和用户之间的一段问答交互。一个会话包含一条或多条消息。对话是会话中对智能体的一次调用，智能体会将对话中产生的消息添加到会话中。 <br>  <br> * 可以使用已创建的会话，会话中已存在的消息将作为上下文传递给模型。 <br> * 对于一问一答等不需要区分 conversation 的场合可不传该参数，系统会自动生成一个会话。  |

**Body**

| **参数** | **类型** | **是否必选** | **说明** |
| --- | --- | --- | --- |
| bot_id | String | 必选 | 要进行会话聊天的智能体ID。 <br> 进入智能体的 开发页面，开发页面 URL 中 `bot` 参数后的数字就是智能体ID。例如`https://www.coze.cn/space/341****/bot/73428668*****`，智能体ID 为`73428668*****`。 <br>  |
| user_id | String | 必选 | 标识当前与智能体对话的用户，由使用方自行定义、生成与维护。user_id 用于标识对话中的不同用户，不同的 user_id，其对话的上下文消息、数据库等对话记忆数据互相隔离。如果不需要用户数据隔离，可将此参数固定为一个任意字符串，例如 `123`，`abc` 等。 <br>  |
| additional_messages <br>  <br>  | Array of object <br>  | 可选 <br>  | 对话输入信息，你可以通过此字段传入本次对话中用户的问题。 <br> * additional_messages 仅支持 role=user 的记录，以免影响模型效果。 <br>  |
| stream <br>  | Boolean <br>  | 可选 <br>  | 是否启用流式返回，目前**仅支持流式响应**。 <br>  <br> * **true**： “流式响应”将模型的实时响应提供给客户端，类似打字机效果。你可以实时获取服务端返回的对话、消息事件，并在客户端中同步处理、实时展示，也可以直接在 completed 事件中获取智能体最终的回复。 |
| shortcut_command | Object | 可选 | 快捷指令信息。你可以通过此参数指定此次对话执行的快捷指令，必须是智能体已绑定的快捷指令。 <br> 消息结构可参考 **ShortcutCommandDetail Object**。 <br> 调用快捷指令，会自动根据快捷指令配置信息生成本次对话中的用户问题，并放入 additional_messages 最后一条消息作为本次对话的用户输入。 <br>  |

**EnterMessage Object**

| **参数** | **类型** | **是否必选** | **说明** |
| --- | --- | --- | --- |
| role | String | 必选 | 发送这条消息的实体。取值： <br>  <br> * **user**：代表该条消息内容是用户发送的。 <br> * **assistant**：代表该条消息内容是智能体发送的。 |
| type <br>  | String | 可选 <br>  | 消息类型。默认为 **question。** <br>  <br> * **question**：用户输入内容。 <br> * **answer**：智能体返回给用户的消息内容，支持增量返回。如果工作流绑定了消息节点，可能会存在多 answer 场景，此时可以用流式返回的结束标志来判断所有 answer 完成。 <br> * **function_call**：智能体对话过程中调用函数（function call）的中间结果。  <br> * **tool_response**：调用工具 （function call）后返回的结果。 <br> * **follow_up**：如果在 智能体上配置打开了用户问题建议开关，则会返回推荐问题相关的回复内容。不支持在请求中作为入参。 <br> * **verbose**：多 answer 场景下，服务端会返回一个 verbose 包，对应的 content 为 JSON 格式，`content.msg_type =generate_answer_finish` 代表全部 answer 回复完成。不支持在请求中作为入参。 <br>  <br> 仅发起会话（v3）接口支持将此参数作为入参，且： <br>  <br> * 如果 autoSaveHistory=true，type 支持设置为 question 或 answer。 <br> * 如果 autoSaveHistory=false，type 支持设置为 question、answer、function_call、tool_output/tool_response。 <br>  <br> 其中，type=question 只能和 role=user 对应，即仅用户角色可以且只能发起 question 类型的消息。 <br>  |
| content | String | 可选 | 消息的内容，支持纯文本、多模态（文本、图片、文件混合输入）、卡片等多种类型的内容。 <br>  <br> * content_type 为 object_string 时，content 为 **object_string object** 数组序列化之后的 JSON String，详细说明可参考 **object_string object。** <br> * 当 content_type **=** text **** 时，content 为普通文本，例如 `"content" :"Hello!"`。 |
| content_type | String | 可选 | 消息内容的类型，支持设置为： <br>  <br> * text：文本。 <br> * object_string：多模态内容，即文本和文件的组合、文本和图片的组合。 <br> * card：卡片。此枚举值仅在接口响应中出现，不支持作为入参。 <br>  <br> content 不为空时，此参数为必选。 <br>  |

**object_string object**

| **参数** | **类型** | **是否必选** | **说明** |
| --- | --- | --- | --- |
| type | String | 必选 | 多模态消息内容类型，支持设置为： <br>  <br> * text：文本类型。 <br> * file：文件类型。 <br> * image：图片类型。 <br> * audio：音频类型。 |
| text | String | 可选 | 文本内容。 |
| file_url | String | 可选 | 文件、图片或语音文件的在线地址。必须是可公共访问的有效地址。 <br> 在 type 为 file、image 或 audio 时，file_id 和 file_url 应至少指定一个。 |
* 一个 object_string 数组中最多包含一条 `text` 类型消息，但可以包含多个 `file`、`image` 类型的消息。
* 当 object_string 数组中存在 `text` 类型消息时，必须同时存在至少 1 条 `file` 或 `image` 消息，纯文本消息（仅包含 `text` 类型）需要使用 `content_type: text` 字段直接指定，不能使用 `object_string` 数组。
* 支持发送纯图片或纯文件消息，但每条纯图片或纯文件消息的前一条或后一条消息中，必须包含一条 `content_type: text` 的纯文本消息，作为用户查询的上下文。例如， `"content": "[{\"type\":\"image\",\"file_id\":\"&#123;&#123;file_id_1&#125;&#125;\"}]"` 为一条纯图片消息，该纯图片消息的前一条或后一条消息必须是一条纯文本消息，否则接口会报 4000 参数错误。

例如，以下数组是一个完整的多模态内容：

<div style="display: flex;">
<div style="flex-shrink: 0;width: calc((100% - 16px) * 0.5000);">

序列化前：
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

序列化后：
```JSON
"[{\"type\":\"text\",\"text\":\"你好我有一个帽衫，我想问问它好看么，你帮我看看\"},{\"type\":\"file\",\"file_url\":\"{{file_url_1}}\"}]"
```




</div>
</div>


**ShortcutCommandDetail Object**

| **参数** | **类型** | **是否必选** | **说明** |
| --- | --- | --- | --- |
| command_id | String | 必选 | 对话要执行的快捷指令 ID，必须是智能体已绑定的快捷指令。 |
| parameters | Map&lt;String, String&gt; | 可选 | 用户输入的快捷指令组件参数信息。 <br> 自定义键值对，其中键（key）为快捷指令组件的名称，值（value）为组件对应的用户输入，为 **object_string object** 数组序列化之后的 JSON String，详细说明可参考 **object_string object。** |

**流式响应**

在流式响应中，服务端不会一次性发送所有数据，而是以数据流的形式逐条发送数据给客户端，数据流中包含对话过程中触发的各种事件（event），直至处理完毕或处理中断。处理结束后，服务端会通过 conversation.message.completed 事件返回拼接后完整的模型回复信息。各个事件的说明可参考**流式响应事件**。
流式响应允许客户端在接收到完整的数据流之前就开始处理数据，例如在对话界面实时展示智能体的回复内容，减少客户端等待模型完整回复的时间。
流式响应的整体流程如下：

* 流式响应流程：
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

* 流式响应示例：
   ```Plain Text
   # chat - 开始
   event: conversation.chat.created
   // 在 chat 事件里，data 字段中的 id 为 Chat ID，即会话 ID。
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


返回的事件消息体结构如下：
| **参数** | **类型** | **说明** |
| --- | --- | --- |
| event | String | 当前流式返回的数据包事件。详细说明可参考 **流式响应事件**。 |
| data | Object | 消息内容。其中，chat 事件和 message 事件的格式不同。 <br>  <br> * chat 事件中，data 为 **Chat** **Object**。 <br> * message、audio 事件中，data 为 **Message** **Object**。 |

**流式响应事件**

| **事件（event）名称** | **说明** |
| --- | --- |
| conversation.chat.created | 创建对话的事件，表示对话开始。 |
| conversation.chat.in_progress | 服务端正在处理对话。 |
| conversation.message.delta | 增量消息，通常是 type=answer 时的增量消息。 |
| conversation.audio.delta | 增量语音消息，通常是 type=answer 时的增量消息。只有输入中包含语音消息时，才会返回 audio.delta。 |
| conversation.message.completed | message 已回复完成。此时流式包中带有所有 message.delta 的拼接结果，且每个消息均为 completed 状态。 |
| conversation.chat.completed | 对话完成。工具类型 |
| conversation.chat.failed | 此事件用于标识对话失败。 |
| conversation.chat.requires_action | 对话中断，需要使用方上报工具的执行结果。 |
| error | 流式响应过程中的错误事件。 |
| done | 本次会话的流式返回正常结束。 |
**接口响应参数：**

| **参数** | **类型** | **说明** |
| --- | --- | --- |
| data | Object | 本次对话的基本信息。详细说明可参考 **Chat** **Object**。 |
| code | Integer | 状态码。 <br> `0` 代表调用成功。 |
| msg | String | 状态信息。API 调用失败时可通过此字段查看详细错误信息。 |

**Message Object**

| **参数** | **类型** | **说明** |
| --- | --- | --- |
| id | String | Message ID，即消息的唯一标识。 |
| conversation_id | String | 此消息所在的会话 ID。 |
| bot_id | String | 编写此消息的智能体ID。此参数仅在对话产生的消息中返回。 |
| chat_id | String | Chat ID。此参数仅在对话产生的消息中返回。 |
| meta_data | Map | 创建消息时的附加消息，获取消息时也会返回此附加消息。 |
| role | String | 发送这条消息的实体。取值： <br>  <br> * **user**：代表该条消息内容是用户发送的。 <br> * **assistant**：代表该条消息内容是智能体发送的。 |
| content | String <br>  | 消息的内容，支持纯文本、多模态（文本、图片、文件混合输入）、卡片等多种类型的内容。 |
| content_type | String | 消息内容的类型，取值包括： <br>  <br> * text：文本。 <br> * object_string：多模态内容，即文本和文件的组合、文本和图片的组合。 <br> * card：卡片。此枚举值仅在接口响应中出现，不支持作为入参。 <br> * audio：音频。此枚举值仅在接口响应中出现，不支持作为入参。仅当输入有 audio 文件时，才会返回此类型。当 content_type 为 audio 时，content 为 base64 后的音频数据。音频的编码根据输入的 audio 文件的不同而不同： <br>    * 输入为 wav 格式音频时，content 为**采样率 24kHz，raw 16 bit, 1 channel, little-endian 的 pcm 音频片段 base64 后的字符串** <br>    * 输入为 ogg_opus 格式音频时，content 为**采样率 48kHz，1 channel，10ms 帧长的 opus 格式音频片段base64 后的字符串** |
| created_at | Integer | 消息的创建时间，格式为 10 位的 Unixtime 时间戳，单位为秒（s）。 |
| updated_at | Integer | 消息的更新时间，格式为 10 位的 Unixtime 时间戳，单位为秒（s）。 |
| type | String | 消息类型。 <br>  <br> * **question**：用户输入内容。 <br> * **answer**：智能体返回给用户的消息内容，支持增量返回。如果工作流绑定了 messge 节点，可能会存在多 answer 场景，此时可以用流式返回的结束标志来判断所有 answer 完成。 <br> * **function_call**：智能体对话过程中调用函数（function call）的中间结果。 <br> * **tool_response**：调用工具 （function call）后返回的结果。 <br> * **follow_up**：如果在智能体上配置打开了用户问题建议开关，则会返回推荐问题相关的回复内容。 <br> * **verbose**：多 answer 场景下，服务端会返回一个 verbose 包，对应的 content 为 JSON 格式，`content.msg_type =generate_answer_finish` 代表全部 answer 回复完成。 <br>  <br> 仅发起会话（v3）接口支持将此参数作为入参，且： <br>  <br> * 如果 autoSaveHistory=true，type 支持设置为 question 或 answer。 <br> * 如果 autoSaveHistory=false，type 支持设置为 question、answer、function_call、tool_response。 <br>  <br> 其中，type=question 只能和 role=user 对应，即仅用户角色可以且只能发起 question 类型的消息。 <br>  |
| section_id | String | 上下文片段 ID。每次清除上下文都会生成一个新的 section_id。 |
| reasoning_content | String | DeepSeek-R1 模型的思维链（CoT）。模型会将复杂问题逐步分解为多个简单步骤，并按照这些步骤逐一推导出最终答案。 <br> 该参数仅在使用 DeepSeek-R1 模型时才会返回。 |

**示例**

**流式响应**

* **Request**
   ```Shell
   curl --location --request POST 'https://api.kozex.ai/v3/chat?conversation_id=7374752000116113452' \
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


# 查看消息列表

查看指定会话的消息列表。

**查看消息列表** API 与**查看对话消息详情** API 的区别在于：

* **查看消息列表** API 用于查询指定会话（conversation）中的消息记录，不仅包括开发者在会话中手动插入的每一条消息和用户发送的 Query，也包括调用**发起对话** API 得到的 type=answer 的智能体回复，但不包括 type=function_call、tool_response 和 **** follow-up 类型的对话中间态消息。
* **查看对话消息详情** API 通常用于非流式对话场景中，查看某次对话（chat）中 type=answer 的智能体回复及 type=function_call、tool_response 和 **** follow-up 类型类型的对话中间态消息。不包括用户发送的 Query。

**基础信息**
| **请求方式** | POST |
| --- | --- |
| **请求地址** | &#123;&#123;host&#125;&#125;/v1/conversation/message/list |
| **接口说明** | 调用接口查看指定会话的消息列表。 |

**请求参数**

**Header**

| **参数** | **取值** | **说明** |
| --- | --- | --- |
| Authorization | Bearer $AccessToken | 用于验证客户端身份的**个人访问令牌**。你可以在扣子平台中生成个人访问令牌，详细信息，参考【准备工作】。 |
| Content-Type | application/json | 解释请求正文的方式。 |

**Query**

| **参数** | **类型** | **是否必选** | **示例** | **说明** |
| --- | --- | --- | --- | --- |
| conversation_id | Integer | 必选 | 737363834493434**** | Conversation ID，即会话的唯一标识。可以在发起对话接口 Response 中查看 conversation_id 字段。 |

**Body**

| **参数** | **类型** | **是否必选** | **示例** | **说明** |
| --- | --- | --- | --- | --- |
| order | String | 可选 | desc | 消息列表的排序方式。 <br>  <br> * desc：（默认）按创建时间降序排序，最新的消息排序最前。 <br> * asc：按创建时间升序排序，最早的消息排序最前。 |
| chat_id | String | 可选 | 737999610479815**** | 筛选指定 Chat ID 中的消息列表。发起对话接口 Response 中 Chat 事件的 data.id 字段即为 Chat ID。 |
| before_id | String | 可选 | 737363834493437**** | 查看指定位置之前的消息。 <br> 默认为 0，表示不指定位置。如需向前翻页，则指定为返回结果中的 first_id。 |
| after_id | String | 可选 | 737363834493437**** | 查看指定位置之后的消息。 <br> 默认为 0，表示不指定位置。如需向后翻页，则指定为返回结果中的 last_id。 |
| limit | Long | 可选 | 30 | 每次查询返回的数据量。默认为 50，取值范围为 1~50。 |

**返回参数**

| **参数** | **类型** | **示例** | **说明** |
| --- | --- | --- | --- |
| data | Array of [OpenMessageApi](#openmessageapi) | [ { "bot_id": "", "chat_id": "", "content": "你的名字叫什么", "content_type": "text", "conversation_id": "737363834493434****", "created_at": 1716809829, "id": "737363834493437****", "meta_data": {}, "role": "user", "type": "", "updated_at": "1716809829" }] | 消息详情。 |
| has_more | Boolean | true | 是否已返回全部消息。 <br>  <br> * true：未返回全部消息，可再次调用此接口查看其他分页。 <br> * false：已返回全部消息。 |
| first_id | String | 737363834493437**** | 返回的消息列表中，第一条消息的 Message ID。 |
| last_id | String | 737363834493440**** | 返回的消息列表中，最后一条消息的 Message ID。 |
| code | Long | 0 | 调用状态码。0 表示调用成功，其他值表示调用失败，你可以通过 msg 字段判断详细的错误原因。 |
| msg | String | Success | 状态信息。API 调用失败时可通过此字段查看详细错误信息。 |

**OpenMessageApi**

| **参数** | **类型** | **示例** | **说明** |
| --- | --- | --- | --- |
| id | String | 738130009748252**** | Message ID，即消息的唯一标识。 |
| conversation_id | String | 737999610479815**** | 此消息所在的会话 ID。 |
| bot_id | String | 747363834493437**** | 编写此消息的智能体 ID。此参数仅在对话产生的消息中返回。 |
| chat_id | String | 757363834493437**** | Chat ID。此参数仅在对话产生的消息中返回。 <br> 不同的对话中，系统会生成新的`chat_id`。同一个用户在第一次对话和第二次对话时，`chat_id`不一样。 |
| meta_data | JSON Map | {} | 创建消息时的附加消息，获取消息时也会返回此附加消息。 |
| role | String | user | 发送这条消息的实体。取值： <br>  <br> * **user**：代表该条消息内容是用户发送的。 <br> * **assistant**：代表该条消息内容是 Bot 发送的。 |
| content | String | 早上好，今天星期几？ | 消息的内容，支持纯文本、多模态（文本、图片、文件混合输入）、卡片等多种类型的内容。 |
| content_type | String | text | 消息内容的类型，取值包括： <br>  <br> * text：文本。 <br> * object_string：多模态内容，即文本和文件的组合、文本和图片的组合。 <br> * card：卡片。此枚举值仅在接口响应中出现，不支持作为入参。 |
| created_at | Long | 1718592898 | 消息的创建时间，格式为 10 位的 Unixtime 时间戳，单位为秒（s）。 |
| updated_at | Long | 1718592898 | 消息的更新时间，格式为 10 位的 Unixtime 时间戳，单位为秒（s）。 |
| type | String | question | 消息类型。 <br>  <br> * **question**：用户输入内容。 <br> * **answer**：智能体返回给用户的消息内容，支持增量返回。如果工作流绑定了 messge 节点，可能会存在多 answer 场景，此时可以用流式返回的结束标志来判断所有 answer 完成。 <br> * **function_call**：智能体对话过程中调用函数（function call）的中间结果。 <br> * **tool_response**：调用工具 （function call）后返回的结果。 <br> * **follow_up**：如果在 Bot 上配置打开了用户问题建议开关，则会返回推荐问题相关的回复内容。 <br> * **verbose**：多 answer 场景下，服务端会返回一个 verbose 包，对应的 content 为 JSON 格式，`content.msg_type =generate_answer_finish` 代表全部 answer 回复完成。 <br>  <br> 仅发起会话（v3）接口支持将此参数作为入参，且： <br>  <br> * 如果 autoSaveHistory=true，type 支持设置为 question 或 answer。 <br> * 如果 autoSaveHistory=false，type 支持设置为 question、answer、function_call、tool_output/tool_response。 <br>    其中，type=question 只能和 role=user 对应，即仅用户角色可以且只能发起 question 类型的消息。 <br>  |
| section_id | String | 757363834493437**** | 上下文片段 ID。每次清除上下文都会生成一个新的 section_id。 |
| reasoning_content | String | 好的，我现在需要给一个13岁的大学生提供学习建议。首先，我得考虑用户的情况…… | DeepSeek-R1 模型的思维链（CoT）。模型会将复杂问题逐步分解为多个简单步骤，并按照这些步骤逐一推导出最终答案。 <br> 该参数仅在使用 DeepSeek-R1 模型时才会返回。 |

**示例**

**请求示例**

```shell
curl --location --request POST {{host}}/v1/conversation/message/list?&conversation_id=737363834493434**** \
--data-raw '{
    "limit": null,
    "order": "asc",
    "chat_id": "737363834493437****",
    "before_id": "737363834493437****"
}'
```


**返回示例**

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

# 清除上下文

清除指定会话中的上下文。

**接口说明**

在智能体对话场景中，上下文指模型在处理当前输入时所能参考的历史消息、对话记录，它决定了模型对当前任务的理解深度和连贯性，直接影响输出结果的准确性和相关性。多轮对话中，如果需要开启新的话题，可以调用此 API 清除上下文。清除上下文后，对话中的历史消息不会作为上下文被输入给模型，后续对话不再受之前历史消息的影响，避免信息干扰，确保对话的准确性和连贯性。
会话中的消息存储在上下文段落（section）中，每次调用此 API 清除上下文时，系统会自动删除旧的上下文段落，并创建新的上下文段落用于存储新的消息。再次发起对话时，新分区中的消息会作为新的上下文传递给模型参考。
* 清除上下文 API 只是清空模型可见的消息历史，不会实际删除会话中的消息，通过查看消息列表等 API 仍能查看到消息内容。
* 仅支持清除本人创建的会话的上下文。
* 会话、对话、消息和上下文段落的术语解释请参见[基础概念](https://www.coze.cn/open/docs/developer_guides/coze_api_overview#fed4a54c)。

![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/e4b55816254c4446ae59bbca33ca8e1d~tplv-goo7wpa0wc-image.image)


**基础信息**

| **请求方式** | POST |
| --- | --- |
| **请求地址** | &#123;&#123;host&#125;&#125;/v1/conversations/:conversation_id/clear |
| **接口说明** | 清除指定会话中的上下文。 |

**请求参数**

**Header**

| **参数** | **取值** | **说明** |
| --- | --- | --- |
| Authorization | Bearer *$Access_Token* | 用于验证客户端身份的**个人访问令牌**。你可以在扣子平台中生成个人访问令牌，详细信息，参考【准备工作】。 |
| Content-Type | application/json | 解释请求正文的方式。 |

**Path**

| **参数** | **类型** | **是否必选** | **示例** | **说明** |
| --- | --- | --- | --- | --- |
| conversation_id | String | 必选 | 737989918257**** | 待清除上下文的会话 ID。你可以在[发起对话](https://www.coze.cn/docs/developer_guides/chat_v3)接口的 Response 中通过 conversation_id 字段获取会话 ID。 |

**返回参数**

| **参数** | **类型** | **示例** | **说明** |
| --- | --- | --- | --- |
| code | Long | 0 | 调用状态码。0 表示调用成功，其他值表示调用失败，你可以通过 `msg` 字段判断详细的错误原因。 |
| msg | String | "Success" | 状态信息。API 调用失败时可通过此字段查看详细错误信息。 |
| data | Object of [Section](#section) | { "id": "1234567890123456789", "conversation_id": "1234567890123456789" } | 上下文段落（section ）的详细信息。 <br> Section 是一个独立的上下文段落，用于分隔不同的对话阶段或主题。Section 中包括上下文消息，当用户清除上下文时，系统会创建一个新的 Section，从而确保新的对话不受历史消息的影响。 |

**Section**

| **参数** | **类型** | **示例** | **说明** |
| --- | --- | --- | --- |
| id | String | 1234567890123456789 | Session ID，即清除上下文后新创建的上下文段落（section）的唯一标识符。 <br> 每个上下文段落对应一批独立的上下文消息。每次清除上下文时，系统会新建一个上下文段落用于存储新的上下文消息。 |
| conversation_id | String | 737999610479815**** | Conversation ID，即会话的唯一标识。 |

**示例**

**请求示例**

```shell
curl --location --request POST {{host}}/v1/conversations/:conversation_id/clear \
--header 'Authorization: Bearer pat_OYDacMzM3WyOWV3Dtj2bHRMymzxP****' \
--header 'Content-Type: application/json' \
--data-raw '{}'
```


**返回示例**

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

# 执行工作流

执行已发布的工作流。

**接口说明**

此接口为非流式响应模式，对于支持流式输出的节点，应使用接口**执行工作流（流式响应）**获取流式响应。调用接口后，你可以从响应中获得 debug_url，访问链接即可通过可视化界面查看工作流的试运行过程，其中包含每个执行节点的输入输出等详细信息，帮助你在线调试或排障。

**限制说明**

|  **限制项** |  **说明**  |
| --- | --- |
| 工作流发布状态 |  必须为已发布。执行未发布的工作流会返回错误码 4200。 创建并发布工作流的操作可参考[使用工作流](https://www.coze.cn/docs/guides/use_workflow)。 |
| 节点限制 | 工作流中不能包含消息节点、开启了流式输出的结束节点、问答节点。 |
| 关联智能体 | 调用此 API 之前，应先在扣子平台中试运行此工作流，如果试运行时需要关联智能体，则调用此 API 执行工作流时，也需要指定智能体ID。通常情况下，执行存在数据库节点、变量节点等节点的工作流需要关联智能体。 |
| 请求大小上限 |  20 MB，包括输入参数及运行期间产生的消息历史等所有相关数据。  |
| 超时时间  | * 未开启工作流异步运行时，工作流整体超时时间为 10 分钟，建议执行时间控制在 5 分钟以内，否则不保障执行结果的准确性。 详细说明可参考[工作流使用限制](https://www.coze.cn/docs/guides/workflow_limits)。 <br> * 开启工作流异步运行后，工作流整体超时时间为 24 小时。 |

**基础信息**

| **请求方式** | POST |
| --- | --- |
| **请求地址** |  https://&#123;&#123;host&#125;&#125;/v1/workflow/run  |
| **接口说明** | 执行已发布的工作流。 |

**请求参数**

**Header**

| **参数** | **取值** | **说明** |
| --- | --- | --- |
| Authorization | Bearer *$Access_Token* | 用于验证客户端身份的**个人访问令牌**。你可以在扣子平台中生成个人访问令牌，详细信息，参考【准备工作】。 |
| Content-Type | application/json | 解释请求正文的方式。 |

**Body**

| **参数** | **类型** | **是否必选** | **示例** | **说明** |
| --- | --- | --- | --- | --- |
| workflow_id | String | 必选 | 73664689170551***** | 待执行的 Workflow ID，此工作流应已发布。 <br> 进入 Workflow 编排页面，在页面 URL 中，`workflow` 参数后的数字就是 Workflow ID。例如 `https://www.coze.com/work_flow?space_id=42463***&workflow_id=73505836754923***`，Workflow ID 为 `73505836754923***`。 |
| parameters | Map[string]any <br>  | 可选 | { <br> "user_id":"12345", <br> "user_name":"George" <br> } | 工作流开始节点的输入参数及取值，你可以在指定工作流的编排页面查看参数列表。 <br> 如果工作流输入参数为 Image 等类型的文件，可以调用[上传文件](https://www.coze.cn/open/docs/developer_guides/upload_files) API 获取 file_id，在调用此 API 时，在 parameters 中以序列化之后的 JSON 格式传入 file_id。例如  `“parameters” : { "input": "{\"file_id\": \"xxxxx\"}" }`。 |
| bot_id | String | 可选 | 73428668***** | 需要关联的智能体 ID。 部分工作流执行时需要指定关联的智能体，例如存在数据库节点、变量节点等节点的工作流。 <br> ![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/55746fa5540b488ea83a79064a223500~tplv-goo7wpa0wc-image.image) <br> 进入智能体的开发页面，开发页面 URL 中 bot 参数后的数字就是智能体t ID。例如 `https://www.coze.com/space/341****/bot/73428668*****`，智能体 ID 为 `73428668*****`。 <br> 确保调用该接口使用的令牌开通了此智能体所在空间的权限。 <br> 确保该智能体已发布为 API 服务。 <br>  |
| app_id | String | 可选 | 749081945898306**** | 该工作流关联的应用的 ID |

**返回参数**

| **参数** | **类型** | **示例** | **说明** |
| --- | --- | --- | --- |
| code | Long | 0 | 调用状态码。 <br>  <br> * 0 表示调用成功。 <br> * 其他值表示调用失败。你可以通过 msg 字段判断详细的错误原因。 |
| msg | String | Success | 状态信息。API 调用失败时可通过此字段查看详细错误信息。 |
| data | String |  | 工作流执行结果，通常为 JSON 序列化字符串，部分场景下可能返回非 JSON 结构的字符串。 |
| execute_id | String | 741364789030728**** | 异步执行的事件 ID。 |
| token | Long |  | 预留字段，无需关注。 |
| cost | String | 0 | 预留字段，无需关注。 |
| debug_url | String |  | 工作流试运行调试页面。访问此页面可查看每个工作流节点的运行结果、输入输出等信息。 |

**示例**

**请求示例**

```shell
curl --location --request POST 'https://{{host}}/v1/workflow/run' \
--header 'Authorization: Bearer pat_hfwkehfncaf****' \
--header 'Content-Type: application/json' \
--data-raw '{
    "workflow_id": "73664689170551*****",
    "parameters": "{\"user_id\":\"12345\",\"user_name\":\"George\"}"
}'
```


**返回示例**

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

# 执行工作流（流式响应）

执行已发布的工作流，响应方式为流式响应。

**接口说明**

调用 API 执行工作流时，对于支持流式输出的工作流，往往需要使用流式响应方式接收响应数据，例如实时展示工作流的输出信息、呈现打字机效果等。
在流式响应中，服务端不会一次性发送所有数据，而是以数据流的形式逐条发送数据给客户端，数据流中包含工作流执行过程中触发的各种事件（event），直至处理完毕或处理中断。处理结束后，服务端会通过 `event: Done` 事件提示工作流执行完毕。
目前支持流式响应的工作流节点包括**输出节点**、**问答节点**和**开启了流式输出的结束节点**。对于不包含这些节点的工作流，可以使用**执行工作流**接口一次性接收响应数据。


**限制说明**


* 通过 API 方式执行工作流前，应确认此工作流已发布，执行从未发布过的工作流时会返回错误码 4200。
* 调用此 API 之前，应先在扣子平台中试运行此工作流。
   * 如果试运行时需要关联智能体，则调用此 API 执行工作流时，也需要指定 bot_id。通常情况下，执行存在数据库节点、变量节点等节点的工作流需要关联智能体。
   * 执行应用中的工作流时，需要指定 app_id。
   * 请勿同时指定 bot_id 和 app_id，否则 API 会报错 4000，表示请求参数错误。
* 此接口为同步接口，如果工作流整体或某些节点运行超时，智能体可能无法提供符合预期的回复，建议将工作流的执行时间控制在 5 分钟以内。
* 工作流支持的请求大小上限为 20MB，包括输入参数以及运行期间产生的消息历史等所有相关数据。


**基础信息**

| **请求方式** | POST |
| --- | --- |
| **请求地址** |  https://&#123;&#123;host&#125;&#125;/v1/workflow/stream_run |
| **接口说明** | 执行已发布的工作流，响应方式为流式响应。 |

**Header**

| **参数** | **取值** | **说明** |
| --- | --- | --- |
| Authorization | Bearer <span style="color: #D83931"><em>$Access_Token</em></span> | 用于验证客户端身份的**个人访问令牌**。你可以在扣子平台中生成个人访问令牌，详细信息，参考【准备工作】。 |
| Content-Type | application/json | 解释请求正文的方式。 |

**Body**

| **参数** | **类型** | **是否必选** | **说明** |
| --- | --- | --- | --- |
| workflow_id | String  | 必选 | 待执行的 Workflow ID，此工作流应已发布。 <br> 进入 Workflow 编排页面，在页面 URL 中，`workflow` 参数后的数字就是 Workflow ID。例如 `https://www.coze.com/work_flow?space_id=42463***&workflow_id=73505836754923***`，Workflow ID 为 `73505836754923***`。 |
| parameters | Map[string]any | 可选 | 工作流开始节点的输入参数及取值，你可以在指定工作流的编排页面查看参数列表。 |
| bot_id <br>  | String  <br>  | 可选 <br>  | 需要关联的智能体ID。 部分工作流执行时需要指定关联的智能体，例如存在数据库节点、变量节点等节点的工作流。 <br> ![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/55746fa5540b488ea83a79064a223500~tplv-goo7wpa0wc-image.image) <br> 进入智能体的开发页面，开发页面 URL 中 `bot` 参数后的数字就是智能体ID。例如 `https://www.coze.com/space/341****/bot/73428668*****`，Bot ID 为 `73428668*****`。  <br> * 确保调用该接口使用的令牌开通了此智能体所在空间的权限。 <br> * 确保该智能体已发布为 API 服务。 <br>  |
| app_id | String | 可选 | 工作流所在的应用 ID。 <br> 你可以通过应用的业务编排页面 URL 中获取应用 ID，也就是 URL 中 project-ide 参数后的一串字符，例如 `https://www.coze.cncom/space/739174157340921****/project-ide/743996105122521****/workflow/744102227704147****` 中，应用的 ID 为 `743996105122521****`。 <br> 仅运行扣子应用中的工作流时，才需要设置 app_id。智能体绑定的工作流、空间资源库中的工作流无需设置 app_id。 <br>  |

**返回结果**


在流式响应中，开发者需要注意是否存在丢包现象。

* 事件 ID（id）默认从 0 开始计数，以包含 `event: Done` 的事件为结束标志。开发者应根据 id 确认响应消息整体无丢包现象。
* Message 事件的消息 ID 默认从 0 开始计数，以包含 `node_is_finish : true` 的事件为结束标志。开发者应根据 node_seq_id 确认 Message 事件中每条消息均完整返回，无丢包现象。

| **参数名** | **参数类型** | **参数描述** |
| --- | --- | --- |
| id | Integer | 此消息在接口响应中的事件 ID。以 0 为开始。 |
| event | String  | 当前流式返回的数据包事件。包括以下类型： <br>  <br> * Message：工作流节点输出消息，例如消息节点、结束节点的输出消息。可以在 data 中查看具体的消息内容。 <br> * Error：报错。可以在 data 中查看 error_code 和 error_message，排查问题。 <br> * Done：结束。表示工作流执行结束，此时 data 中包含 debug URL。 <br> * Interrupt：中断。表示工作流中断，此时 data 字段中包含具体的中断信息。 <br> * PING：心跳信号。表示工作流执行中，消息内容为空，用于维持连接。 |
| data | Object | 事件内容。各个 event 类型的事件内容格式不同。 |

**Message 事件**

Message 事件中，data 的结构如下：
| **参数名** | **参数类型** | **参数描述** |
| --- | --- | --- |
| content | String  | 流式输出的消息内容。 |
| node_title | String | 输出消息的节点名称，例如消息节点、结束节点。 |
| node_seq_id | String | 此消息在节点中的消息 ID，从 0 开始计数，例如消息节点的第 5 条消息。 |
| node_is_finish | Boolean | 当前消息是否为此节点的最后一个数据包。 |
| ext | Map[String]String | 额外字段。 |
| cost | String  | 预留字段，无需关注。 |

**Interrupt 事件**

Interrupt 事件中，data 的结构如下：
| **参数名** | **参数类型** | **参数描述** |
| --- | --- | --- |
| interrupt_data | Object | 中断控制内容。 |
| interrupt_data.event_id | String | 工作流中断事件 ID，恢复运行时应回传此字段。 |
| interrupt_data.type | Integer | 工作流中断类型，恢复运行时应回传此字段。 |
| node_title | String | 输出消息的节点名称，例如“问答”。 |

**Error 事件**

Error 事件中，data 的结构如下：
| **参数名** | **参数类型** | **参数描述** |
| --- | --- | --- |
| error_code | Integer | 调用状态码。  <br>  <br> * 0 表示调用成功。  <br> * 其他值表示调用失败。你可以通过 error_message 字段判断详细的错误原因。 |
| error_message | String  | 状态信息。API 调用失败时可通过此字段查看详细错误信息。 |

**示例**

**请求示例**

```Shell
curl --location --request POST 'https://{{host}}/v1/workflow/stream_run' \
--header 'Authorization: Bearer pat_fhwefweuk****' \
--header 'Content-Type: application/json' \
--data-raw '{
    "workflow_id": "73664689170551*****",
    "{\"user_id\":\"12345\",\"user_name\":\"George\"}"
}'
```


**响应示例**


* Message事件
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

* Error 事件
   ```Plain Text
   id: 0
   event: Error
   data: {"error_code":4000,"error_message":"Request parameter error"}
   ```

* Interrupt 事件
   ```Plain Text
   // 流式执行工作流，触发问答节点，Bot提出问题
   id: 0
   event: Message
   data: {"content":"请问你想查看哪个城市、哪一天的天气呢","content_type":"text","node_is_finish":true,"node_seq_id":"0","node_title":"问答"}
   
   id: 1
   event: Interrupt
   data: {"interrupt_data":{"data":"","event_id":"7404830425073352713/2769808280134765896","type":2},"node_title":"问答"}
   ```

# 恢复运行工作流

恢复运行已中断的工作流。

**接口说明**

执行包含问答节点的工作流时，智能体会以指定问题向用户提问，并等待用户回答。此时工作流为中断状态，开发者需要调用此接口回答问题，并恢复运行工作流。如果用户的响应和智能体预期提取的信息不匹配，例如缺少必选的字段，或字段数据类型不一致，工作流会再次中断并追问。如果询问 3 次仍未收到符合预期的回复，则判定为工作流执行失败。

**限制说明**


* 最多调用此接口恢复运行 3 次，如果第三次恢复运行时智能体仍未收到符合预期的回复，则判定为工作流执行失败。
* 恢复运行后，index 和节点 index 都会重置。
* 恢复运行工作流也会产生 token 消耗，且与**执行工作流（流式响应）**时消耗的 token 数量相同。


**基础信息**

| **请求方式** | POST |
| --- | --- |
| **请求地址** | https://&#123;&#123;host&#125;&#125;/v1/workflow/stream_resume |
| **接口说明** | 恢复运行已中断的工作流。 |

**Header**

| **参数** | **取值** | **说明** |
| --- | --- | --- |
| Authorization | Bearer <span style="color: #D83931"><em>$Access_Token</em></span> | 用于验证客户端身份的**个人访问令牌**。你可以在扣子平台中生成个人访问令牌，详细信息，参考【准备工作】。 |
| Content-Type | application/json | 解释请求正文的方式。 |

**Body**

| **参数** | **类型** | **是否必选** | **说明** |
| --- | --- | --- | --- |
| workflow_id | String  | 必选 | 待执行的 Workflow ID，此工作流应已发布。 <br> 进入 Workflow 编排页面，在页面 URL 中，`workflow` 参数后的数字就是 Workflow ID。例如 `https://www.coze.com/work_flow?space_id=42463***&workflow_id=73505836754923***`，Workflow ID 为 `73505836754923***`。 |
| event_id | String | Required | 工作流执行中断事件 ID。你可以从**执行工作流（流式响应）**的响应信息中获得中断事件 ID。 |
| resume_data | String  | Required | 恢复执行时，用户对智能体指定问题的回复。 <br> 如果是问答节点导致的中断，回复中应包含问答节点中的必选参数，否则工作流会再次中断并提问。 |
| interrupt_type | Integer | Required  | 中断类型，你可以从**执行工作流（流式响应）**的响应信息中获得中断时间的中断类型。 |

**返回结果**


在流式响应中，开发者需要注意是否存在丢包现象。

* 事件 ID（id）默认从 0 开始计数，以包含 `event: Done` 的事件为结束标志。开发者应根据 id 确认响应消息整体无丢包现象。
* Message 事件的消息 ID 默认从 0 开始计数，以包含 `node_is_finish : true` 的事件为结束标志。开发者应根据 node_seq_id 确认 Message 事件中每条消息均完整返回，无丢包现象。

| **参数名** | **参数类型** | **参数描述** |
| --- | --- | --- |
| id | Integer | 此消息在接口响应中的事件 ID。以 0 为开始。 |
| event | String  | 当前流式返回的数据包事件。包括以下类型： <br>  <br> * Message：工作流节点输出消息，例如消息节点、结束节点的输出消息。可以在 data 中查看具体的消息内容。 <br> * Error：报错。可以在 data 中查看 error_code 和 error_message，排查问题。 <br> * Done：结束。表示工作流执行结束，此时 data 中包含 debug URL。 <br> * Interrupt：中断。表示工作流中断，此时 data 字段中包含具体的中断信息。 <br> * PING：心跳信号。表示工作流执行中，消息内容为空，用于维持连接。 |
| data | Object | 事件内容。各个 event 类型的事件内容格式不同。 |

**Message 事件**

Message 事件中，data 的结构如下：
| **参数名** | **参数类型** | **参数描述** |
| --- | --- | --- |
| content | String  | 流式输出的消息内容。 |
| node_title | String | 输出消息的节点名称，例如消息节点、结束节点。 |
| node_seq_id | String | 此消息在节点中的消息 ID，从 0 开始计数，例如消息节点的第 5 条消息。 |
| node_is_finish | Boolean | 当前消息是否为此节点的最后一个数据包。 |
| ext | Map[String]String | 额外字段。 |
| cost | String  | 预留字段，无需关注。 |

**Interrupt 事件**

Interrupt 事件中，data 的结构如下：

| **参数名** | **参数类型** | **参数描述** |
| --- | --- | --- |
| interrupt_data | Object | 中断控制内容。 |
| interrupt_data.event_id | String | 工作流中断事件 ID，恢复运行时应回传此字段。 |
| interrupt_data.type | Integer | 工作流中断类型，恢复运行时应回传此字段。 |
| node_title | String | 输出消息的节点名称，例如“问答”。 |

**Error 事件**

Error 事件中，data 的结构如下：
| **参数名** | **参数类型** | **参数描述** |
| --- | --- | --- |
| error_code | Integer | 调用状态码。  <br>  <br> * 0 表示调用成功。  <br> * 其他值表示调用失败。你可以通过 error_message 字段判断详细的错误原因。 |
| error_message | String  | 状态信息。API 调用失败时可通过此字段查看详细错误信息。 |

**示例**

工作流执行中断场景下，以查看天气工作为例，完整的接口调用示例如下。

1. 调用接口**执行工作流（流式响应）**，要求查看天气。
   请求示例如下：
   ```Plain Text
   curl --location 'https://{{host}}/v1/workflow/stream_run' \
   --header 'Authorization: Bearer pat_vTG1****' \
   --header 'Content-Type: application/json' \
   --data '{
       "workflow_id": "739739507914235****",
       "parameters": "{\"BOT_USER_INPUT\":\"查看天气\"}"
   }'
   ```

2. 触发问答节点，工作流中断，响应信息中返回智能体提出的问题，要求用户提供城市和日期。
   返回示例如下：
   ```Plain Text
   id: 0
   event: Message
   data: {"content":"请问你想查看哪个城市、哪一天的天气呢","content_type":"text","node_is_finish":true,"node_seq_id":"0","node_title":"问答"}
   
   id: 1
   event: Interrupt
   data: {"interrupt_data":{"data":"","event_id":"7404831988202520614/6302059919516746633","type":2},"node_title":"问答"}
   ```

3. 调用接口恢复运行工作流，回复智能体城市和日期。
   请求示例如下：
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

4. 工作流执行完毕，完成天气查询，返回工作流输出消息。
   返回示例如下：
   ```Plain Text
   id: 0
   event: Message
   data: {"content":"{\"output\":[{\"condition\":\"中到大雨\",\"humidity\":72,\"predict_date\":\"2024-08-20\",\"temp_high\":35,\"temp_low\":26,\"weather_day\":\"中到大雨\",\"wind_dir_day\":\"西风\",\"wind_dir_night\":\"西风\",\"wind_level_day\":\"3\",\"wind_level_night\":\"3\"}]}","content_type":"text","cost":"0","node_is_finish":true,"node_seq_id":"0","node_title":"End","token":386}
   
   id: 1
   event: Done
   data: {}
   ```

# 执行对话流（流式响应）
执行已发布的对话流，响应方式为流式响应
## **接口说明**

* 对话流是基于对话场景的特殊工作流，专门用于处理对话类请求。对话流通过对话的方式和用户交互，并完成复杂的业务逻辑。在应用中添加对话流，将对话中的用户指令拆分为一个个步骤节点，并为其设计用户界面，你可以搭建出适用于移动端或网页端的对话式 AI 应用，实现自动化、智能化的对话流程。关于对话流的详细说明可参考[工作流与对话流](/guides/workflow_and_chatflow)。
* 此接口为流式响应模式，允许客户端在接收到完整的数据流之前就开始处理数据，例如在对话界面实时展示回复内容，减少客户端等待模型完整回复的时间。 
* 此接口支持包括问答节点、输入节点等可能导致对话中断的节点，**对话中断时只需再次调用对话流，在 additional_messages 中指定输入内容，即可继续对话**。

**如果对话流的输入中包含文件、图片等多模态内容，需要先将多模态内容上传到第三方存储工具中，并获取一个公开可访问的 URL 地址，将此 URL 作为对话流的输入**。
调用接口后，你可以从响应的 Done 事件中获得 debug_url，访问链接即可通过可视化界面查看对话流的试运行过程，其中包含每个执行节点的输入输出等详细信息，帮助你在线调试或排障。


此接口可用于调用空间资源库中的对话流，或扣子应用中的对话流。调用这两种对话流时，入参不同：
| **入参** | **资源库对话流** |  | **扣子应用对话流** |
| --- | --- | --- | --- |
|  | **在智能体中执行** | **在扣子应用中执行** |  |
| workflow_id | 必选 | 必选 | 必选 |
| app_id | 不传 | 必选 | 必选 |
| bot_id | 必选 | 不传 | 不传 |
| conversation_id | 可选 | 可选 | 可选 |

## **基础信息**
| **请求方式** | POST |
| --- | --- |
| **请求地址** | &#123;&#123;host&#125;&#125;/v1/workflows/chat |
| **接口说明** | 执行已发布的对话流，响应方式为流式响应。 |
### Header
| **参数** | **取值** | **说明** |
| --- | --- | --- |
| Authorization | Bearer *$Access_Token* | 用于验证客户端身份的**个人访问令牌**。你可以在扣子平台中生成个人访问令牌，详细信息，参考【准备工作】。 |
| Content-Type | application/json | 解释请求正文的方式。 |
### **Body**
| **参数** | **类型** | **是否必选** | **说明** |
| --- | --- | --- | --- |
| workflow_id | String  | 必选 | 待执行的 Workflow ID，此工作流应已发布。 <br> 进入 Workflow 编排页面，在页面 URL 中，`workflow` 参数后的数字就是 Workflow ID。例如 `https://www.coze.com/work_flow?space_id=42463***&workflow_id=73505836754923***`，Workflow ID 为 `73505836754923***`。 |
| additional_messages | array&lt;Object&gt;  <br>  <br> * Object  | 必选 | 对话中用户问题和历史消息。数组长度限制为 50，即最多传入 50 条消息。 |
| parameters | Object | 必选 | 设置对话流输入参数中的自定义参数。你可以根据实际业务需求，通过`parameters` 参数传入自定义参数以及对应的值。 <br>  <br> * 对话流的输入参数 USER_INPUT 应在 additional_messages 中传入，在 parameters 中的 USER_INPUT 不生效。 <br> * 如果 parameters 中未指定 CONVERSATION_NAME 或自定义输入参数，则使用参数默认值运行对话流；如果指定了这些参数，则使用指定值。 |
| app_id | String | 可选 | 需要关联的扣子应用 ID。调用对话流时，必须指定 app_id 或 bot_id，便于模型调用智能体或应用的数据库、变量等数据处理问题。 <br> 进入应用开发界面，开发页面 URL 中的 project-ide 参数后的数字就是 AppID，例如`https://www.coze.cn/space/74421656*****/project-ide/744208683**` ，扣子应用 ID 为`744208683**`**。** |
| bot_id <br>  | String  <br>  | 可选 <br>  | 需要关联的智能体ID。 部分工作流执行时需要指定关联的智能体，例如存在数据库节点、变量节点等节点的工作流。 <br> ![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/55746fa5540b488ea83a79064a223500~tplv-goo7wpa0wc-image.image) <br> 进入智能体的开发页面，开发页面 URL 中 `bot` 参数后的数字就是智能体ID。例如 `https://www.coze.com/space/341****/bot/73428668*****`，Bot ID 为 `73428668*****`。  <br> * 确保调用该接口使用的令牌开通了此智能体所在空间的权限。 <br> * 确保该智能体已发布为 API 服务。 <br>  |
| conversation_id | String | 可选 | 对话流对应的会话 ID，对话流产生的消息会保存到此对话中。会话默认为开始节点设置的 CONVERSATION_NAME，也可以通过 conversation_id 参数指定会话。 <br>  <br> * 指定 conversation_id 时，parameters 中设置的 CONVERSATION_NAME 不生效。 <br> * 会话的创建者必须和执行对话流的用户一致，即 API 访问令牌的创建者一致，否则无法执行对话流。 <br> * 会话与 app_id、渠道匹配，不同渠道的会话隔离。 <br> * 指定 bot_id 时，如果没有传入 conversation_id ，扣子会创建一个新的会话。不支持同时指定 bot_id 和 app_id |
| workflow_version | String | 可选 | 工作流的版本号，仅当运行的工作流属于资源库工作流时有效。未指定版本号时默认执行最新版本的工作流。 <br>  |
additional_messages object 
| 参数 | 类型 | 是否必须 | 说明 |
| --- | --- | --- | --- |
| content | String | 必须 | 消息的内容，仅支持纯文本。 <br> 暂不支持多模态（文本、图片、文件混合输入）、卡片等类型的内容。 |
| content_type | String | 必须 | 消息内容的类型。 <br> content_type 固定为 text，表示普通文本。 <br> 指定 content 时，应同时设置 content_type |
| role | String | 必须 | * **user**：代表该条消息内容是用户发送的。 <br> * **assistant**：代表该条消息内容是模型发送的。 |
| type |  |  | 消息类型。默认为 **question。** <br>  <br> * **question**：用户输入内容。 <br> * **answer**：模型返回给用户的消息内容，支持增量返回。如果对话流绑定了消息节点，可能会存在多 answer 场景，此时可以用流式返回的结束标志来判断所有 answer 完成。 |
## 返回参数

* 流式响应允许客户端在接收到完整的数据流之前就开始处理数据，例如在对话界面实时展示回复内容，减少客户端等待模型完整回复的时间。 

流式响应的整体流程如下：
### **流式响应流程**

   ```Python
   ######### 整体概览 （chat, MESSAGE 两级） 
   # chat - 开始
   # chat - 处理中 
   #   MESSAGE - 知识库召回 
   #   MESSAGE - function_call 
   #   MESSAGE - tool_output 
   #   MESSAGE - answer is normal text 
   #   MESSAGE - 多 answer 的情况下，会继续有 message.delta 
   # chat - 完成 
   # 流结束 event: done 
   ######### 
   ```


#### **流式响应事件列表**

   | **事件（event）名称** | **说明** |
   | --- | --- |
   | conversation.chat.created | 创建对话的事件，表示对话开始。 |
   | conversation.chat.in_progress | 服务端正在处理对话。 |
   | conversation.message.delta | 增量消息，通常是 type=answer 时的增量消息。 |
   | conversation.message.completed | message 已回复完成。此时流式包中带有所有 message.delta 的拼接结果，且每个消息均为 completed 状态。 |
   | conversation.chat.completed | 对话完成。 |
   | conversation.chat.failed | 此事件用于标识对话失败。 |
   | conversation.chat.requires_action | 对话中断，需要使用方上报工具的执行结果。通常是触发了问答节点或输入节点，需要再次调用此接口继续对话。 |
   | error | 流式响应过程中的错误事件。 |
   | done | 本次会话的流式返回正常结束。 |

#### **流式响应示例**

   ```Python
   # chat - 开始
   event: conversation.chat.created
   data: {"id":"120","conversation_id":"456","created_at":1733407180,"last_error":{"code":0,"msg":""},"status":"created","usage":{"token_count":0,"output_count":0,"input_count":0},"section_id":"789"}
   # chat - 处理中 
   event: conversation.chat.in_progress
   data: {"id":"121","conversation_id":"456","created_at":1733407180,"last_error":{"code":0,"msg":""},"status":"in_progress","usage":{"token_count":0,"output_count":0,"input_count":0},"section_id":"789"}
   # MESSAGE - answer is normal text 
   event: conversation.message.delta
   data: {"id":"122","conversation_id":"456","role":"assistant","type":"answer","content":"中午吃啥了","content_type":"text","chat_id":"567","section_id":"789","created_at":1733407182}
   
   # MESSAGE - 消息结束
   event: conversation.message.completed
   data: {"id":"124","conversation_id":"456","role":"assistant","type":"answer","content":"中午吃啥了","content_type":"text","chat_id":"567","section_id":"789","created_at":1733407182}
   
   event: conversation.message.completed
   data: {"id":"125","conversation_id":"456","role":"assistant","type":"verbose","content":"{\"msg_type\":\"interrupt\",\"data\":\"\",\"from_module\":null,\"from_unit\":null}","content_type":"text","chat_id":"567","section_id":"789","created_at":1733407182,"updated_at":1733407182}
   
   event: conversation.message.completed
   data: {"id":"130","conversation_id":"456","role":"assistant","type":"verbose","content":"{\"msg_type\":\"generate_answer_finish\",\"data\":\"{\\\"finish_reason\\\":1,\\\"FinData\\\":\\\"\\\"}\",\"from_module\":null,\"from_unit\":null}","content_type":"text","chat_id":"567","section_id":"789","created_at":1733407182,"updated_at":1733407182}
   # chat-需要操作（中断，通常为问答节点或者输入节点触发）
   event: conversation.chat.requires_action
   data: {"id":"131","conversation_id":"456","created_at":1733407180,"completed_at":1733407182,"last_error":{"code":0,"msg":""},"status":"requires_action","usage":{"token_count":0,"output_count":0,"input_count":0},"required_action":{"type":"submit_tool_outputs","submit_tool_outputs":{"tool_calls":[{"id":"","type":"reply_message","function":null,"require_info":null}]}},"section_id":"789"}
   
   event: done
   data: {"debug_url":"http://{{HOST}}/work_flow?execute_id=74449256856****\u0026space_id=7442165654356*****\u0026workflow_id=744224337778*****"}
   ```




* #### **事件消息体结构**
   | **参数** | **类型** | **说明** |
   | --- | --- | --- |
   | event | String | 当前流式返回的数据包事件。在流式响应中，服务端不会一次性发送所有数据，而是以数据流的形式逐条发送数据给客户端，数据流中包含对话过程中触发的各种事件（event），直至处理完毕或处理中断。处理结束后，服务端会通过 conversation.message.completed 事件返回拼接后完整的模型回复信息。各个事件的说明可参考下表。  |
   | data | Object | 消息内容。其中，chat 事件和 message 事件的格式不同。 <br>  <br> * chat 事件中，data 为 **Chat** **Object**。 <br> * message 事件中，data 为 **Message** **Object**。 |
* ##### **Chat Object** 
   | 参数 | 类型 | 是否可选 | 说明 |
   | --- | --- | --- | --- |
   | id | String | 必填 | 对话 ID，即对话的唯一标识。 |
   | conversation_id | String | 必填 | 会话 ID，即会话的唯一标识。 |
   | bot_id | String | 必填 | 要进行会话聊天的智能体 ID。 |
   | status <br>  | String | 必填 | 对话的运行状态。取值为： <br>  <br> * created：对话已创建。 <br> * in_progress：智能体正在处理中。 <br> * completed：智能体已完成处理，本次对话结束。 <br> * failed：对话失败。 <br> * requires_action：对话中断，需要进一步处理。 <br> * canceled：对话已取消。 |
   | required_action | Object | 选填 | 需要运行的信息详情。 |
   | usage | Object <br> ```JSON <br> { <br> "token_count":123, // token总的数量 <br> "output_count":100, // 输出消耗token <br> "input_count":23 // 输入 token  <br> } <br> ``` <br>  <br>  | 选填 | Token 消耗的详细信息。实际的 Token 消耗以对话结束后返回的值为准。 |

##### Message Object
| 参数 | 类型 | 说明 |
| --- | --- | --- |
| id | String | Message ID，即消息的唯一标识。 |
| conversation_id | String | 此消息所在的会话 ID。 |
| bot_id | String | 编写此消息的智能体ID。此参数仅在对话产生的消息中返回。 |
| chat_id | String | Chat ID。此参数仅在对话产生的消息中返回。 |
| role | String | 发送这条消息的实体。取值： <br>  <br> * **user**：代表该条消息内容是用户发送的。 <br> * **assistant**：代表该条消息内容是智能体发送的。 |
| content | String | 消息的内容，支持纯文本、多模态（文本、图片、文件混合输入）等多种类型的内容。 |
| content_type | String | 消息内容的类型，取值包括： <br>  <br> * text：文本。 <br> * object_string：多模态内容，即文本和文件的组合、文本和图片的组合。 |
| type | String | 消息类型。 <br>  <br> * **question**：用户输入内容。 <br> * **answer**：智能体返回给用户的消息内容，支持增量返回。如果对话流绑定了 messge 节点，可能会存在多 answer 场景，此时可以用流式返回的结束标志来判断所有 answer 完成。 |







