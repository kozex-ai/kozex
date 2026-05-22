本文介绍如何在**后端业务逻辑中**为 Workflow 增加自定义节点。
节点是工作流的组成单元，包含完整的、相对独立的业务逻辑。一个工作流由至少两个节点组成：开始和结束。一个节点从在画布中被拖拽生成开始，到作为工作流的一部分传入后端保存，再到实例化后执行，完整过程如下图所示：
![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/c838039de41040f3b7654deb7ebb70aa~tplv-goo7wpa0wc-image.image)

画布上的一个节点，一定属于一个“节点类型”，比如“大模型节点”，“插件节点”等。新增一个节点类型，封装业务相关的专业逻辑，能拓展工作流的能力范围。一般来说，新增一个基本的节点类型，需要做的事情：

* 前端：参考【新增工作流节点类型】文档。
* 后端：
   * 新增 1 个 go 文件，在其中定义 2 个结构体，分别为节点 Config 结构体（实现 NodeAdaptor 和 NodeBuilder 接口）和节点执行结构体（实现 InvokableNode 等接口）。
   * 配置这个新节点类型的“元信息”（比如名称、样式、通用执行参数如超时等），并注册这个新类型的 NodeAdaptor。

## Workflow 运行机制简介
Kozex 基于 [Eino](https://github.com/cloudwego/eino) 的 workflow 能力搭建，是一个有向无环图（DAG），包含“控制流”和“数据流”两个方面：

* 控制流：除“开始”节点外，一个节点可以开始执行的前提条件，是它的所有入边“得出结果”，且至少一条入边的结果是“成功”。
   * 一条边的状态有三种：未完成，成功，跳过。初始状态是“未完成”。“成功”和“跳过”都认为已“得出结果”。
   * 当边的起始节点完成执行，或边的起始节点包含分支、且该分支的选择结果是当前边时，当前边状态变为“成功”。
   * 当边的起始节点包含分支且该分支的选择结果不是当前边时，当前边的状态变为“跳过”。
   * “跳过”状态可以传导，即若一个节点的所有入边的结果都是“跳过”，则该节点的所有出边的结果自动改为“跳过”。

<div style="text-align: center"><img src="https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/335381ec3cd64cbfb991d2177fdc1077~tplv-goo7wpa0wc-image.image" width="2510px" /></div>

  在上图中，执行时序是：

   1. 大模型节点执行失败，进入异常分支，边“1”未被选择，记为“跳过”，边“2”被选择，记为成功。
   1. 此时“结束”节点的两个入边“1、4”，分别是“跳过，未完成”，“结束”节点不能开始执行。
   1. 文本处理执行完成，边 3 记为“成功”。
   1. 代码节点执行完成，边 4 记为“成功”。
   1. 结束节点的两个入边“1、4”，分别为“跳过、成功”，“结束”节点开始执行。
* 数据流：一个节点的输入数据，是该节点引用的“任意”前驱节点某个输出字段 + 配置的静态固定值(static values) + 配置的变量(variable)值合并而成的 `map[string]any`。如图所示

![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/17a71d26e7d9496589e5651c37b4b52a~tplv-goo7wpa0wc-image.image)
## 写节点逻辑

1. 在 backend/domain/workflow/internal/nodes 下新建一个目录，名字是你的节点名称，如 custom
2. 在这个目录下，新建一个 go 文件，如 custom.go
3. 在这个 go 文件中，新增一个结构体如 YourNode，代表你的节点。结构体内可以定义“**任意**”的节点运行时需要的字段，这些字段不需要导出。
4. 这个结构体必须要实现一个或多个固定的接口，代表“节点执行逻辑”。

例如：
```Go
package custom

import "context"

// YourNode is the definition of your custom node.
type YourNode struct{
    // define any fields here which are needed during node execution
}

// Invoke is the execution method of your custom node.
func (c *YourNode) Invoke(ctx context.Context, input map[string]any) (
    output map[string]any, err error) {
    // your business logic
    return
}
```


5. 上面代码中的 " Invoke" 是一个有固定签名的函数。这个函数接收“前驱节点”映射而来的 input，运行自身的业务逻辑，再最终给出 output 供“后置节点”使用。Input 和 Output 都不是流。**YourNode 实现了 InvokableNode 接口，定义在 backend/domain/workflow/internal/nodes/node.go 中，不需要修改**：

```Go
// InvokableNode is a basic workflow node that can Invoke.
// Invoke accepts non-streaming input and returns non-streaming output.
// It does not accept any options.
type InvokableNode interface {
    Invoke(ctx context.Context, input map[string]any) (
       output map[string]any, err error)
}
```


6. Invoke 的 input 和 output 都是 map[string]any，那 map 里面是什么？是画布上具体配置的字段，以及前驱节点映射而来的值。比如简单的大模型节点配置：
   ![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/ede0f156e39c45d99c2b6e4f8e626bf6~tplv-goo7wpa0wc-image.image)

用户配置了一个“输入字段”和一个“输出字段”，因此传到 YourNode.Invoke 方法中的 input map 是 `map[string]any{"input": "xxx"}`，返回的 output map 是 `map[string]any{"output": "yyy"}`
## 定节点名称

1. 在 backend/domain/workflow/entity/node_meta.go 中，给你的新节点一个独占的名字：
   ```Go
   // NodeType 是 coze workflow 服务端对一个节点类型的枚举值
   NodeTypeMyAwesomeNode NodeType = "MyAwesomeNode"
   ```

2. 在同一个文件中，为你的新节点配置一些通用“元信息”属性。以 JSON 序列化节点为例：
   ```Go
   NodeTypeJsonSerialization: {
       // ID is the unique identifier of this node type. Used in various front-end APIs. 
       ID:         58,
       
       // Key is the unique NodeType of this node. Used in backend code as well as saved in DB.
       Key:        NodeTypeJsonSerialization,
       
       // DisplayKey is the string used in frontend to identify this node.
       // Example use cases: 
       // - used during querying test-run results for nodes
       // - used in returned messages from streaming openAPI Runs.
       // If empty, will use Key as DisplayKey.
       DisplayKey: "ToJSON",
   
       // Name is the node in ZH_CN, will be displayed on Canvas.
       Name: "JSON 序列化",
   
       // Category is the category of this node, determines which category this node will be displayed in.
       Category: "utilities",
   
       // Desc is the desc in ZH_CN, will be displayed as tooltip on Canvas.
       Desc: "用于把变量转化为JSON字符串",
   
       // Color is the color of the upper edge of the node displayed on Canvas.
       Color: "F2B600",
   
       // IconURI is the resource identifier for the icon displayed on the Canvas. It's resolved into a full URL by the backend to support different deployment environments.
       IconURI: "default_icon/workflow_icon/icon-json-stringify.jpg",
       
       // SupportBatch indicates whether this node can set batch mode.
       // NOTE: ultimately it's frontend that decides which node can enable batch mode.
       SupportBatch: false,
   
       // ExecutableMeta configures certain common aspects of request-time behaviors for this node.
       ExecutableMeta: ExecutableMeta{
          // DefaultTimeoutMS configures the default timeout for this node, in milliseconds. 0 means no timeout.
          DefaultTimeoutMS: 60 * 1000, // 1 minute
          // PreFillZero decides whether to pre-fill zero value for any missing fields in input.
          PreFillZero: true,
          // PostFillNil decides whether to post-fill nil value for any missing fields in output.
          PostFillNil: true,
       },
       // EnUSName is the name in EN_US, will be displayed on Canvas if language of kozex is set to EnUS.
       EnUSName:        "JSON serialization",
       // EnUSDescription is the description in EN_US, will be displayed on Canvas if language of kozex is set to EnUS.
       EnUSDescription: "Convert variable to JSON string",
   },
   ```


关注点：

* ID：唯一的一个 int64，自定义节点类型建议从一个大数开始，比如 1000+。与前端交互主要用这个 ID。
* Key：唯一的一个 string，准确的描述出节点的功能定位。后端主要用 Key 做唯一标识，包括数据库中保存的。
* Name, Category, Desc, Color, IconURI, EnUSName, EnUSDescription：给画布用的展示信息。
* SupportBatch: 节点是否可配置“批处理”模式。最终控制权在前端的节点样式。
* ExecutableMeta：节点运行时的通用配置，在后面“进阶功能”中会展开说明。

注：节点对应的Icon图片需放置在docker/volumes/minio/default_icon/workflow_icon/中，图标的文件名和此处的IconURI字段需要匹配。
## 从画布到后端

1. 画布中的 workflow 节点传到后端，变成了 backend/domain/workflow/entity/vo/canvas.go 中的 Node 结构体。比如一个“开始节点”的 Node 结构体：
   ```JSON
   {
     "blocks": [], // 复合节点（批处理、循环）的子节点，开始节点没有
     "data": { // 真正运行相关的配置数据
       "nodeMeta": { // 前端用的一些元信息
         "description": "工作流的起始节点，用于设定启动工作流需要的信息",
         "icon": "https://lf3-static.bytednsdoc.com/obj/eden-cn/dvsmryvd_avi_dvsm/ljhwZthlaukjlkulzlp/icon/icon-Start-v2.jpg",
         "subTitle": "",
         "title": "开始"
       },
       "outputs": [ // 开始节点的所有“输出”字段的配置
         {
           "name": "input", // 只有一个输出字段
           "required": false, // 不是必填的
           "schema": {
             "type": "string" // 数组元素类型是 string
           },
           "type": "list" // 整体是数组类型
         }
       ],
     },
     "edges": null, // 复合节点内部的边信息，开始节点没有
     "id": "100001", // 节点 ID，开始节点的 ID 是固定的
     "meta": { // 位置信息，前端使用
       "position": {
         "x": -90.5329099821747,
         "y": -323.84999999999985
       }
     },
     "type": "1" // 节点类型，对应 NodeMeta.ID
   },
   ```

2. 这个 Node 结构体在能够执行之前，先要转换成后端的 backend/domain/workflow/internal/schema/node_schema.go 中的 NodeSchema 结构体。
   ```Go
   // NodeSchema is the universal description and configuration for a workflow Node.
   // It should contain EVERYTHING a node needs to instantiate.
   type NodeSchema struct {
       // Key is the node key within the Eino graph.
       // A node may need this information during execution,
       // e.g.
       // - using this Key to query workflow State for data belonging to current node.
       Key vo.NodeKey `json:"key"`
   
       // Name is the name for this node as specified on Canvas.
       // A node may show this name on Canvas as part of this node's input/output.
       Name string `json:"name"`
   
       // Type is the NodeType for the node.
       Type entity.NodeType `json:"type"`
   
       // Configs are node specific configurations, with actual struct type defined by each Node Type.
       // Will not hold information relating to field mappings, nor as node's static values.
       // In a word, these Configs are INTERNAL to node's implementation, NOT related to workflow orchestration.
       // Actual type of these Configs should implement two interfaces:
       // - NodeAdaptor: to provide conversion from vo.Node to NodeSchema
       // - NodeBuilder: to provide instantiation from NodeSchema to actual node instance.
       Configs any `json:"configs,omitempty"`
   
       // InputTypes are type information about the node's input fields.
       InputTypes map[string]*vo.TypeInfo `json:"input_types,omitempty"`
       // InputSources are field mapping information about the node's input fields.
       InputSources []*vo.FieldInfo `json:"input_sources,omitempty"`
   
       // OutputTypes are type information about the node's output fields.
       OutputTypes map[string]*vo.TypeInfo `json:"output_types,omitempty"`
       // OutputSources are field mapping information about the node's output fields.
       // NOTE: only applicable to composite nodes such as NodeTypeBatch or NodeTypeLoop.
       OutputSources []*vo.FieldInfo `json:"output_sources,omitempty"`
   
       // ExceptionConfigs are about exception handling strategy of the node.
       ExceptionConfigs *ExceptionConfig `json:"exception_configs,omitempty"`
       // StreamConfigs are streaming characteristics of the node.
       StreamConfigs *StreamConfig `json:"stream_configs,omitempty"`
   
       // SubWorkflowBasic is basic information of the sub workflow if this node is NodeTypeSubWorkflow.
       SubWorkflowBasic *entity.WorkflowBasic `json:"sub_workflow_basic,omitempty"`
       // SubWorkflowSchema is WorkflowSchema of the sub workflow if this node is NodeTypeSubWorkflow.
       SubWorkflowSchema *WorkflowSchema `json:"sub_workflow_schema,omitempty"`
   
       // FullSources contains more complete information about a node's input fields' mapping sources,
       // such as whether a field's source is a 'streaming field',
       // or whether the field is an object that contains sub-fields with real mappings.
       // Used for those nodes that need to process streaming input.
       // Set InputSourceAware = true in NodeMeta to enable.
       FullSources map[string]*SourceInfo
   
       // Lambda directly sets the node to be an Eino Lambda.
       // NOTE: not serializable, used ONLY for internal test.
       Lambda *compose.Lambda
   }
   ```

3. 为什么要转化？因为前端 schema 的 Node 结构体中“针对工作流的编辑态”的，后端 schema 的 NodeSchema 结构体是“针对工作流的运行态”的。
4. 如何转化？需要实现 NodeAdaptor 接口：
   ```Go
   // NodeAdaptor provides conversion from frontend Node to backend NodeSchema.
   type NodeAdaptor interface {
       Adapt(ctx context.Context, n *vo.Node, opts ...AdaptOption) (
          *schema.NodeSchema, error)
   }
   ```

   这个接口由节点的 Config Type 实现。比如 JSON 序列化节点：
   ```Go
   // SerializationConfig is the Config type for NodeTypeJsonSerialization.
   // Each Node Type should have its own designated Config type,
   // which should implement NodeAdaptor and NodeBuilder.
   // NOTE: we didn't define any fields for this type,
   // because this node is simple, we doesn't need to extract any SPECIFIC piece of info
   // from frontend Node. In other cases we would need to do it, such as LLM's model configs.
   type SerializationConfig struct {
       // you can define ANY number of fields here,
       // as long as these fields are SERIALIZABLE and EXPORTED.
       // to store specific info extracted from frontend node.
       // e.g.
       // - LLM model configs
       // - conditional expressions
       // - fixed input fields such as MaxBatchSize
   }
   
   // Adapt provides conversion from Node to NodeSchema.
   // NOTE: in this specific case, we don't need AdaptOption.
   func (s *SerializationConfig) Adapt(_ context.Context, n *vo.Node, _ ...nodes.AdaptOption) (*schema.NodeSchema, error) {
       ns := &schema.NodeSchema{
          Key:     vo.NodeKey(n.ID),
          Type:    entity.NodeTypeJsonSerialization,
          Name:    n.Data.Meta.Title,
          Configs: s, // remember to set the Node's Config Type to NodeSchema as well
       }
   
       // this sets input fields' type and mapping info
       if err := convert.SetInputsForNodeSchema(n, ns); err != nil {
          return nil, err
       }
   
       // this set output fields' type info
       if err := convert.SetOutputTypesForNodeSchema(n, ns); err != nil {
          return nil, err
       }
   
       return ns, nil
   }
   ```

   Config 结构体，每个节点类型必须有。必须实现 NodeAdaptor 接口。
   Config 内可以定义**“任意”的可序列化**（反例：function 无法被序列化，interface 的序列化依赖具体类型，不建议在 Config 中使用）的字段，这些字段需要**导出**。
   这些字段本质上是由前端画布中的节点中各表单配置转换而来，是节点实例化的“主要信息来源”（剩下的来自 NodeMeta 中的节点元信息）。
5. 注册 NodeAdaptor：拿到一个 frontend node，怎么找到对应的实现 NodeAdaptor 接口的 Config 结构体？需要注册 NodeType 到 NodeAdaptor 的映射关系。修改 domain/workflow/internal/canvas/adaptor/to_schema.go 中的 RegisterAllNodeAdaptors 函数：
   ```Go
   // RegisterAllNodeAdaptors register all NodeType's NodeAdaptor.
   func RegisterAllNodeAdaptors() {
       // omitted multiple registrations ...
       
       // register a generator function so that each time a NodeAdaptor is needed,
       // we can provide a brand new Config instance.
       nodes.RegisterAdaptor(entity.NodeTypeJsonSerialization, func() nodes.NodeAdaptor {
          return &json.SerializationConfig{}
       })
       
       // omitted multiple registrations...
   }
   ```


## 生成节点实例
现在我们已经有了统一的后端 NodeSchema，最后一步是真正实例化这个节点并执行。

1. Workflow 运行时，拿到 NodeSchema，会调用 backend/domain/workflow/internal/compose/node_builder.go 中的 New 函数，把 NodeSchema 最终转换成 eino 的 Lambda：
   ```Go
   // New instantiates the actual node type from NodeSchema.
   func New(ctx context.Context, s *schema.NodeSchema,
       // omitted multiple lines...
   
       // if NodeSchema's Configs implements NodeBuilder, will use it to build the node
       nb, ok := s.Configs.(schema.NodeBuilder)
       if ok {
          opts := []schema.BuildOption{
             schema.WithWorkflowSchema(sc),
             schema.WithInnerWorkflow(inner),
          }
   
          // build the actual InvokableNode, etc.
          n, err := nb.Build(ctx, s, opts...)
          if err != nil {
             return nil, err
          }
   
          // wrap InvokableNode, etc. within NodeRunner, converting to eino's Lambda
          return toNode(s, n), nil
       }
       
       // omitted multiple lines...
   }
   ```

   Lambda 是 [Eino](https://github.com/cloudwego/eino) 框架的一个基础组件，意思是“用户定制的任意函数”，参见：[Eino Lambda 使用说明](https://www.cloudwego.io/zh/docs/eino/core_modules/components/lambda_guide/)

2. 在上面的实例化过程中，NodeSchema.Configs 实现了 NodeBuilder 接口：
   ```Go
   // NodeBuilder takes a NodeSchema and several BuildOption to build an executable node instance.
   // The result 'executable' MUST implement at least one of the execute interfaces:
   // - nodes.InvokableNode
   // - nodes.StreamableNode
   // - nodes.CollectableNode
   // - nodes.TransformableNode
   // - nodes.InvokableNodeWOpt
   // - nodes.StreamableNodeWOpt
   // - nodes.CollectableNodeWOpt
   // - nodes.TransformableNodeWOpt
   // NOTE: the 'normal' version does not take NodeOption, while the 'WOpt' versions take NodeOption.
   // NOTE: a node should either implement the 'normal' versions, or the 'WOpt' versions, not mix them up.
   type NodeBuilder interface {
       Build(ctx context.Context, ns *NodeSchema, opts ...BuildOption) (
          executable any, err error)
   }
   ```

   以 JSON 序列化节点为例，最简单的 NodeBuilder 接口实现：
   ```Go
   func (s *SerializationConfig) Build(_ context.Context, _ *schema.NodeSchema, _ ...schema.BuildOption) (
       any, error) {
       return &Serializer{}, nil
   }
   
   // Serializer is the actual node implementation.
   type Serializer struct {
       // here can holds ANY data required for node execution
   }
   
   // Invoke implements the InvokableNode interface.
   func (js *Serializer) Invoke(_ context.Context, input map[string]any) (map[string]any, error) {
       // Directly use the input map for serialization
       if input == nil {
          return nil, fmt.Errorf("input data for serialization cannot be nil")
       }
   
       originData := input[InputKeySerialization]
       serializedData, err := sonic.Marshal(originData) // Serialize the entire input map
       if err != nil {
          return nil, fmt.Errorf("serialization error: %w", err)
       }
       return map[string]any{OutputKeySerialization: string(serializedData)}, nil
   }
   ```

   因为 JSON 序列化节点完全无状态，只要拿到 input 之后序列化就可以，因此 Serializer type 中不需要定义任何成员变量。实际上你可以在里面定义**任何“​**运行时需要，且在多次运行期间保持稳定，不因为本次运行的 input 而改变**”​**的字段，比如问答节点：
   ```Go
   type QuestionAnswer struct {
       model    model.BaseChatModel
       nodeMeta entity.NodeTypeMeta
   
       questionTpl string
       answerType  AnswerType
   
       choiceType   ChoiceType
       fixedChoices []string
   
       needExtractFromAnswer     bool
       additionalSystemPromptTpl string
       maxAnswerCount            int
   
       nodeKey      vo.NodeKey
       outputFields map[string]*vo.TypeInfo
   }
   ```

3. JSON 序列化节点为例，Build 出来的是 *Serializer 类型，是节点具体的实现类。*Serializer 实现了 Invoke 方法，对应的是 backend/domain/workflow/internal/nodes/node.go 中的 InvokableNode 接口：
   ```Go
   // InvokableNode is a basic workflow node that can Invoke.
   // Invoke accepts non-streaming input and returns non-streaming output.
   // It does not accept any options.
   // Most nodes implement this, such as NodeTypePlugin.
   type InvokableNode interface {
       Invoke(ctx context.Context, input map[string]any) (
          output map[string]any, err error)
   }
   ```

   Invoke 的语义是：非流式输入，非流式输出。输入输出类型都是 `map[string]any`。
   在同一个 go 文件中，还定义了另外三种执行接口：
   ```Go
   // StreamableNode is a workflow node that can Stream.
   // Stream accepts non-streaming input and returns streaming output.
   // It does not accept and options
   // Currently NO Node implement this.
   // A potential example would be streamable plugin for NodeTypePlugin.
   type StreamableNode interface {
       Stream(ctx context.Context, in map[string]any) (
          *einoschema.StreamReader[map[string]any], error)
   }
   
   // CollectableNode is a workflow node that can Collect.
   // Collect accepts streaming input and returns non-streaming output.
   // It does not accept and options
   // Currently NO Node implement this.
   // A potential example would be a new condition node that makes decisions 
   // based on streaming input.
   type CollectableNode interface {
       Collect(ctx context.Context, in *einoschema.StreamReader[map[string]any]) (
          map[string]any, error)
   }
   
   // TransformableNode is a workflow node that can Transform.
   // Transform accepts streaming input and returns streaming output.
   // It does not accept and options
   // e.g.
   // NodeTypeVariableAggregator implements TransformableNode.
   type TransformableNode interface {
       Transform(ctx context.Context, in *einoschema.StreamReader[map[string]any]) (
          *einoschema.StreamReader[map[string]any], error)
   }
   ```

   Invoke, Stream, Collect, Transform 是 [Eino](https://github.com/cloudwego/eino) 框架的四种基本流式交互范式。有兴趣深入了解的话，可以看 [Eino 流式编程要点](https://www.cloudwego.io/zh/docs/eino/core_modules/chain_and_graph_orchestration/stream_programming_essentials/)。

4. 对同一个节点类型来说，可以只实现一个接口，也可以实现多个，由节点的具体业务决定。比如，
   1. 目前大部分节点只实现了 InvokableNode 接口，因为大部分节点只能处理非流式输入和产生非流式输出。
   2. “批处理”和“循环”节点，实现的是 InvokableNodeWOpt 接口，因为这两个“复合节点”需要给内部的节点透传信息，这个透传是通过 Invoke 方法中传入 NodeOption 实现的。
   3. 大模型节点，额外实现了 StreamableNodeWOpt 接口，因为大模型节点可以产生真正的流式输出，并且可能需要给内部的 workflow tool 透传信息，也需要借助 NodeOption。
   4. 输出节点，额外实现了 TransformableNode 接口，因为需要接收流式输入，并产生流式输出（打字机上屏）。
   到目前为止，应当已经可以实现一个基础的新节点类型，在画布上拖出来，并能够试运行。下面会展开单独介绍下节点相关的进阶功能，可以按需深入了解。

## 进阶功能
### 感知输入输出类型及映射关系
有时节点执行时需要知道输出字段的“来源”，比如：大模型节点（JSON 结构化输出），输入节点（需要输入的字段），问答节点（需要从回答中提取的信息）。可以在 NodeBuilder 中将 NodeSchema 的 OutputTypes 保存在 Node 中，供运行时使用，比如输入节点的 Build 方法：
```Go
func (c *Config) Build(_ context.Context, ns *schema.NodeSchema, _ ...schema.BuildOption) (any, error) {
    // omitted multiple lines...

    return &InputReceiver{
       outputTypes:   ns.OutputTypes, // so the node can refer to its output types during execution
       nodeMeta:      *nodeMeta,
       nodeKey:       ns.Key,
       interruptData: interruptDataStr,
    }, nil
}
```

同理，可以将 NodeSchema 的 InputTypes（输入字段的类型）保存在 Node 中，供节点运行时使用。
另一方面，有的节点需要知道输出字段的“来源”，比如批处理和循环节点，需要知道某个输出字段是来自于内部的某个节点的某个字段。这时，将 NodeSchema.OutputSources 保存在 Node 中，如循环节点的 Build 方法：
```Go
func (c *Config) Build(_ context.Context, ns *schema.NodeSchema, opts ...schema.BuildOption) (any, error) {
    // omitted multiple lines...

    b := &Batch{
       outputs:       make(map[string]*vo.FieldSource),
       innerWorkflow: bo.Inner,
       key:           ns.Key,
       inputArrays:   inputArrays,
    }

    for i := range ns.OutputSources {
       source := ns.OutputSources[i]
       path := source.Path
       if len(path) != 1 {
          return nil, fmt.Errorf("invalid path %q", path)
       }

       // from which inner node's which field does the batch's output fields come from
       b.outputs[path[0]] = &source.Source
    }

    return b, nil
}
```

同理，可以将 NodeSchema 的 InputSources（输入字段的来源）保存在 Node 中，供节点运行时使用。
### 定制前端 schema
大部分稍微复杂的节点，都需要定制的前端 schema，比如大模型节点的“技能信息”，选择器节点的“选择支信息”，子工作流节点的“子工作流 ID 和 version”，批处理节点的“批次大小”和“并发数”等。这些都是需要用户在画布上配置的信息，但是不属于“用户自定义的字段”，而是节点特定的固定配置。这些定制的配置，统一定义在前端 schema 的 Inputs 类型中(backend/domain/workflow/entity/vo/canva.go)：
```C++
type Inputs struct {
    // InputParameters are the fields defined by user for this particular node.
    InputParameters []*Param `json:"inputParameters"`

    // SettingOnError configures common error handling strategy for nodes.
    // NOTE: enable in frontend node's form first.
    SettingOnError *SettingOnError `json:"settingOnError,omitempty"`

    // NodeBatchInfo configures batch mode for nodes.
    // NOTE: not to be confused with NodeTypeBatch.
    NodeBatchInfo *NodeBatch `json:"batch,omitempty"`

    // LLMParam may be one of the LLMParam or IntentDetectorLLMParam or SimpleLLMParam.
    // Shared between most nodes requiring an ChatModel to function.
    LLMParam any `json:"llmParam,omitempty"`

    *OutputEmitter      // exclusive configurations for NodeTypeEmitter and NodeTypeExit in Answer mode
    *Exit               // exclusive configurations for NodeTypeExit
    *LLM                // exclusive configurations for NodeTypeLLM
    *Loop               // exclusive configurations for NodeTypeLoop
    *Selector           // exclusive configurations for NodeTypeSelector
    *TextProcessor      // exclusive configurations for NodeTypeTextProcessor
    *SubWorkflow        // exclusive configurations for NodeTypeSubWorkflow
    *IntentDetector     // exclusive configurations for NodeTypeIntentDetector
    *DatabaseNode       // exclusive configurations for various Database nodes
    *HttpRequestNode    // exclusive configurations for NodeTypeHTTPRequester
    *Knowledge          // exclusive configurations for various Knowledge nodes
    *CodeRunner         // exclusive configurations for NodeTypeCodeRunner
    *PluginAPIParam     // exclusive configurations for NodeTypePlugin
    *VariableAggregator // exclusive configurations for NodeTypeVariableAggregator
    *VariableAssigner   // exclusive configurations for NodeTypeVariableAssigner
    *QA                 // exclusive configurations for NodeTypeQuestionAnswer
    *Batch              // exclusive configurations for NodeTypeBatch
    *Comment            // exclusive configurations for NodeTypeComment
    *InputReceiver      // exclusive configurations for NodeTypeInputReceiver
}
```

对于新的节点类型，也可以扩展上面的 Inputs 结构体，增加新节点类型的特定字段，并在 NodeAdaptor 中解析这些字段，转化到 NodeSchema 的标准字段中。因为无论如何特殊，**这些信息终归是可以归结为“输入”，“输出”，“配置”三个类型**。
这些类型特定的字段，大概率是需要定制化的前端表单支持，因此需要同步处理前端的适配。
### 定制输入输出显示字段
有的节点类型，在试运行时，画布上展示的“输入字段信息”，与运行时输入的 `map[string]any` 并不完全匹配。比如选择器的画布输入展示：
![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/0e8f95b6c76a4811997237c19a62ad6f~tplv-goo7wpa0wc-image.image)

在展示“条件 1 的左值”时，把这个字段的来源是“批处理”这个节点，也一起显式出来。这时，这个节点类型需要实现 `CallbackInputConverted` 接口，把“展示的输入信息和实际的输入信息”区分开：
```Go
// CallbackInputConverted converts node input to a form better suited for UI.
// The converted input will be displayed on canvas when test run, 
// and will be returned when querying the node's input through OpenAPI.
type CallbackInputConverted interface {
    ToCallbackInput(ctx context.Context, in map[string]any) (map[string]any, error)
}
```

举例：选择器节点实现了 `CallbackInputConverted`，参考 backend/domain/workflow/internal/nodes/selector/callbacks.go。
另一方面，有的字段类型需要在画布上展示与原始输出不完全一致的信息，比如问答节点，原始输出是“最终答案”，画布上展示的试运行输出，则是包含“所有轮次的问答过程”。或者比如代码节点等，在格式化输出失败时，会额外展示一个“警告信息”。这时，这个节点类型需要实现 `CallbackOutputConverted` 接口，把“展示的输出信息和实际的输出信息”区分开：
```Go
// CallbackOutputConverted converts node input to a form better suited for UI.
// The converted output will be displayed on canvas when test run,
// and will be returned when querying the node's output through OpenAPI.
type CallbackOutputConverted interface {
    ToCallbackOutput(ctx context.Context, out map[string]any) (*StructuredCallbackOutput, error)
}
```

举例：问答节点实现了`CallbackOutputConverted`，参考 backend/domain/workflow/internal/nodes/qa/question_answer.go。
Kozex 利用 [Eino](https://github.com/cloudwego/eino) 框架的 Callback 机制实现“通过旁路把修改后的 input/output 传递出去”的效果。关于 Eino 的 Callback 机制，可以参考 [Eino Callback 用户手册](https://www.cloudwego.io/zh/docs/eino/core_modules/chain_and_graph_orchestration/callback_manual/)。

### 异常处理策略
部分节点已支持通用的异常处理策略，包括超时、重试、异常后降级到默认输出，异常后执行异常分支等。这些节点包括大模型节点，代码节点等，配置界面如下：
![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/8622ae4ae8a9482eaa97dbf3444e953a~tplv-goo7wpa0wc-image.image)

图中，“代码”节点配置了 60s 超时时间，不重试，异常时执行异常流程。这些配置统一出现在前端 schema 的 SettingOnError 结构体中：
```Go
type ErrorProcessType int

const (
    ErrorProcessTypeThrow             ErrorProcessType = 1 // throws the error as usual
    ErrorProcessTypeReturnDefaultData ErrorProcessType = 2 // return DataOnErr configured in SettingOnError
    ErrorProcessTypeExceptionBranch   ErrorProcessType = 3 // executes the exception branch on error
)

// SettingOnError contains common error handling strategy.
type SettingOnError struct {
    // DataOnErr defines the JSON result to be returned on error.
    DataOnErr   string            `json:"dataOnErr,omitempty"`
    // Switch defines whether ANY error handling strategy is active.
    // If set to false, it's equivalent to set ProcessType = ErrorProcessTypeThrow 
    Switch      bool              `json:"switch,omitempty"`
    // ProcessType determines the error handling strategy for this node.
    ProcessType *ErrorProcessType `json:"processType,omitempty"`
    // RetryTimes determines how many times to retry. 0 means no retry.
    // If positive, any retries will be executed immediately after error.
    RetryTimes  int64             `json:"retryTimes,omitempty"`
    // TimeoutMs sets the timeout duration in millisecond.
    // If any retry happens, ALL retry attempts accumulates to the same timeout threshold.
    TimeoutMs   int64             `json:"timeoutMs,omitempty"`
    // Ext sets any extra settings specific to NodeType
    Ext         *struct {
       // BackupLLMParam is only for LLM Node, marshaled from SimpleLLMParam.
       // If retry happens, the backup LLM will be used instead of the main LLM.
       BackupLLMParam string `json:"backupLLMParam,omitempty"`
    } `json:"ext,omitempty"`
}
```

如果希望在一个新的节点类型中启用上述异常执行策略，只需要在节点的前端表单中加入“异常处理”部分，确保在传给后端的 schema 中包含 SettingOnError。只要这样做，**异常执行策略会自动生效**。
### 带 "branch" 功能的节点
有的节点带有“分支选择”功能，如选择器、意图识别、问答，在画布上，表现在这个节点后面有不止一个“端口”：
![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/54cfd18911084eb68f0f6274b402e882~tplv-goo7wpa0wc-image.image)

“端口”在前端 schema 定义中，包含在边（Edge）上，称为 Port，举例：
```JSON
{
  "sourceNodeID": "133234",
  "targetNodeID": "163493",
  "sourcePortID": "branch_0"
}
```

这个前端 schema 片段的含义是，一条从 133234 节点出发，连接到 163493 节点，且开始的“端口”是“branch_0”。这种分支选择功能，在 Kozex 实现中，是根据节点的实际输出（比如选择器选的分支，意图识别的意图，问答的选项等），在多个“端口”中通过 Eino 的 [Branch](https://www.cloudwego.io/zh/docs/eino/core_modules/chain_and_graph_orchestration/orchestration_design_principles/#branch) 机制做选择。表现在 workflow 拓扑上，是 Node 后面额外多了个 Branch。
这些“端口”分为三类：

1. 常规的端口：比如选择器的各个分支，意图识别的各意图，问答的各静态选项（动态选项的话，合在一起是一个端口）。
2. 默认端口：比如选择器的 else 分支，意图识别的默认意图，问答的“其他”选项。
3. 异常端口：如果节点开启了上面说的通用异常处理策略，会有这个端口。框架来处理，开发者不需要关注。

一个新的节点要启用 branch 功能，则在前端 schema 中，需要确保：

1. 常规的端口，sourcePortID 符合 "branch_%d" 的格式，比如意图识别第一个意图是 branch_0，第二个是 branch_1。
2. 默认端口，sourcePortID = "default"
3. 异常端口：sourcePortID = "branch_error"

另外，节点的 Config Type 需要实现两个接口：

1. BranchBuilder: 接收节点的实际输出，计算对应的端口，用于构建 Branch。

```Go
// BranchBuilder builds the extractor function that maps node output to port index.
type BranchBuilder interface {
    BuildBranch(ctx context.Context) (extractor func(ctx context.Context,
       nodeOutput map[string]any) (int64, bool /*if is default branch*/, error), hasBranch bool)
}
```

以选择器的实现为例：
```Go
func (c *Config) BuildBranch(_ context.Context) (
    func(ctx context.Context, nodeOutput map[string]any) (int64, bool, error), bool) {
    return func(ctx context.Context, nodeOutput map[string]any) (int64, bool, error) {
       choice := nodeOutput[SelectKey].(int64)
       if choice < 0 || choice > int64(len(c.Clauses)+1) {
          return -1, false, fmt.Errorf("selector choice out of range: %d", choice)
       }

       if choice == int64(len(c.Clauses)) { // default
          return -1, true, nil
       }

       return choice, false, nil
    }, true
}
```

框架调用 BuildBranch 时，Config 已经完成了 NodeAdaptor 的过程，内部各字段均可使用。

2. BranchAdaptor: 根据前端 schema (vo.Node)，计算“应当”有哪些端口，用于校验：

```Go
// BranchAdaptor provides validation and conversion from frontend port to backend port.
type BranchAdaptor interface {
    ExpectPorts(ctx context.Context, n *vo.Node) []string
}
```

以意图识别的实现为例：
```Go
func (c *Config) ExpectPorts(ctx context.Context, n *vo.Node) []string {
    expects := make([]string, len(n.Data.Inputs.Intents)+1)
    expects[0] = schema2.PortDefault
    for i := 0; i < len(n.Data.Inputs.Intents); i++ {
       expects[i+1] = fmt.Sprintf(schema2.PortBranchFormat, i)
    }
    return expects
}
```

无论是 BranchBuilder 还是 BranchAdaptor，都不需要感知“异常分支”，完全由 Kozex 框架处理。
### 输入输出赋值兜底策略
Workflow 的节点输入和节点输出都是 `map[string]any`，可能有两种情况发生：

* 节点声明了一个输入 field，但是实际输入的 map 中没有这个 field 对应的 key；
* 节点声明了一个输出 field，但是实际输出的 map 中没有这个 field 对应的 key；

但是，底层的 [Eino](https://github.com/cloudwego/eino) Workflow 引擎，要求在发生字段映射时，上游的字段在上游输出的 map 中务必存在。因此，对于用户可自定义输出字段的节点类型，建议配置 `PostFillNil = true`，Kozex 会自动将缺失的输出字段赋 nil 到 map 中。
另一方面，部分节点的业务逻辑，会自动将输入 map 中的 nil key，设置成零值，针对这些节点，可配置 `PreFillZero = true`，Kozex 会自动将值为 nil 的输入字段替换为对应类型的零值。
比如文本处理节点，需要输入 Nil 值转零值，但是输出字段一定不可能是 Nil：
```YAML
NodeTypeTextProcessor: {
    ID:           15,
    Key:          NodeTypeTextProcessor,
    DisplayKey:   "Text",
    Name:         "文本处理",
    
    // omitted multiple lines...
    ExecutableMeta: ExecutableMeta{
        // omitted multiple lines...
       PreFillZero:        true,
    },
},
```

变量聚合节点，就是希望保留输出中的原始 nil，但是输出的时候可能缺失字段（某个 group 全是 nil）：
```YAML
NodeTypeVariableAggregator: {
    ID:           32,
    Key:          NodeTypeVariableAggregator,
    Name:         "变量聚合",

    // omitted multiple lines...
    ExecutableMeta: ExecutableMeta{
       PostFillNil:        true,
       // omitted multiple lines...
    },
},
```

### 画布打字机效果
输出节点和结束节点（返回文本模式），如果选择了“流式输出”，则在试运行时，会以打字机效果把流式输出内容实时展示在画布上。
如果你的新节点希望开启这个效果，即“试运行时，把流式输出内容实时展示在画布上”，需要在 NodeMeta 中配置 IncrementalOutput = true:
```YAML
NodeTypeOutputEmitter: {
    ID:           13,
    Key:          NodeTypeOutputEmitter,
    DisplayKey:   "Message",
    Name:         "输出",
    
    // omitted multiple lines...
    
    ExecutableMeta: ExecutableMeta{
        // omitted multiple lines...
    
       IncrementalOutput:    true,
    },
},
```

除了这个配置，还有两个前提：

1. 节点确实能返回流式数据，即实现了 StreamableNode 或 TransformableNode 接口。
2. 在 backend/domain/workflow/internal/execute/callback.go 中，修改 `func (n *NodeHandler) OnEndWithStreamOutput(ctx context.Context, info *callbacks.RunInfo, output *schema.StreamReader[callbacks.CallbackOutput]) context.Context` 方法，确保你的新节点类型能够发出 **`*NodeStreamingOutput*` ****这个事件类型（目前还不太好改，后面我们会改进）。

### 复合节点
批处理和循环这种节点，我们称为复合节点（CompositeNode），具有如下特点：

* 在前端 schema 中，Node 内含有至少一个 Block（内部节点）。
* 内部的节点可能引用复合节点及其前驱节点的输出。
* 内部的节点可能中断和恢复，导致复合节点也需要中断和恢复。

如果要新增一个复合节点类型，在 NodeMeta 中设置 `IsComposite = true`，如：
```YAML
NodeTypeLoop: {
    ID:           21,
    Key:          NodeTypeLoop,
    DisplayKey:   "Loop",
    Name:         "循环",
    // omitted multiple lines...
    ExecutableMeta: ExecutableMeta{
       IsComposite:        true,
       // omitted multiple lines...
    },
},
```

实现新的复合节点是一个大工程，可以参考 backend/domain/workflow/internal/nodes/batch/batch.go 以及 backend/domain/workflow/internal/nodes/loop/loop.go 中的实现方式。主要的关注点：

* 内部工作流的调度
* 内部工作流的中断恢复

新的复合节点需要前端一起适配。
### 普通节点的批处理模式
目前三个节点类型支持启用“批处理模式”：大模型，插件和子工作流。
启用“批处理模式”后，节点的前端 schema 的 vo.Node 中，会额外出现批处理相关配置：
```Go
type NodeBatch struct {
    BatchEnable    bool     `json:"batchEnable"`
    BatchSize      int64    `json:"batchSize"`
    ConcurrentSize int64    `json:"concurrentSize"`
    InputLists     []*Param `json:"inputLists,omitempty"`
}
```

要为一个节点类型启用“批处理模式”，在 NodeMeta 中设置 `SupportBatch = true`:
```YAML
NodeTypeLLM: {
    ID:           3,
    Key:          NodeTypeLLM,
    DisplayKey:   "LLM",
    Name:         "大模型",
    SupportBatch: true,
    // omitted multiple lines...
},
```

需要前端一起适配。
### 节点消耗 Token 统计
Workflow 可以统计模型的 token 消耗，既包括节点维度的，也包括工作流整体维度的。
如果你的新节点类型，需要统计 token 消耗，比如需要在试运行页面展示节点的 token 消耗，或需要把节点的 token 消耗累加到工作流整体的消耗中，或者需要在 OpenAPI 返回的消息中统计这个新节点的 token 消耗，则需要在 node_meta.go 中配置 MayUseChatModel = true:
```YAML
NodeTypeLLM: {
    ID:           3,
    Key:          NodeTypeLLM,
    DisplayKey:   "LLM",
    Name:         "大模型",
    
    // omitted multiple lines...

    ExecutableMeta: ExecutableMeta{
        // omitted multiple lines...
    
       MayUseChatModel:    true,
    },
},
```

当配置了 MayUseChatModel = true 后，如果真的用到了大模型，则会自动进行统计。如果配置了，但是没有用到大模型，统计 token 数字就是 0。
### 传递运行时 Option
在 node.go 中，框架提供了 4 个额外带 NodeOption 的运行接口，节点可以有选择的实现：
```Go
// InvokableNodeWOpt is a workflow node that can Invoke.
// Invoke accepts non-streaming input and returns non-streaming output.
// It can accept NodeOption.
// e.g. NodeTypeLLM, NodeTypeSubWorkflow implement this.
type InvokableNodeWOpt interface {
    Invoke(ctx context.Context, in map[string]any, opts ...NodeOption) (
       map[string]any, error)
}

// StreamableNodeWOpt is a workflow node that can Stream.
// Stream accepts non-streaming input and returns streaming output.
// It can accept NodeOption.
// e.g. NodeTypeLLM implement this.
type StreamableNodeWOpt interface {
    Stream(ctx context.Context, in map[string]any, opts ...NodeOption) (
       *einoschema.StreamReader[map[string]any], error)
}

// CollectableNodeWOpt is a workflow node that can Collect.
// Collect accepts streaming input and returns non-streaming output.
// It accepts NodeOption.
// Currently NO Node implement this.
// A potential example would be a new batch node that accepts streaming input,
// process them, and finally returns non-stream aggregation of results.
type CollectableNodeWOpt interface {
    Collect(ctx context.Context, in *einoschema.StreamReader[map[string]any], opts ...NodeOption) (
       map[string]any, error)
}

// TransformableNodeWOpt is a workflow node that can Transform.
// Transform accepts streaming input and returns streaming output.
// It accepts NodeOption.
// Currently NO Node implement this.
// A potential example would be an audio processing node that
// transforms input audio clips, but within the node is a graph
// composed by Eino, and the audio processing node needs to carry
// options for this inner graph.
type TransformableNodeWOpt interface {
    Transform(ctx context.Context, in *einoschema.StreamReader[map[string]any], opts ...NodeOption) (
       *einoschema.StreamReader[map[string]any], error)
}
```

NodeOption 可以在运行时向节点注入 input 和 state 外的额外信息，比如：

* 子工作流节点，向内部的子工作流传递 Eino 的运行时 option。
* 批处理、循环节点，向内部的节点传递中断后恢复的信息。
* 大模型节点，向内部的 workflow tool 传递 message 流的 StreamWriter。

实现自定义的 NodeOption，可以参考 LLM 节点的实现：
```Go
type llmOptions struct {
    toolWorkflowSW *schema.StreamWriter[*entity.Message]
}

func WithToolWorkflowMessageWriter(sw *schema.StreamWriter[*entity.Message]) nodes.NodeOption {
    return nodes.WrapImplSpecificOptFn(func(o *llmOptions) {
       o.toolWorkflowSW = sw
    })
}
```

在 LLM 节点中获取并使用这个 NodeOption:
```Go
llmOpts := nodes.GetImplSpecificOptions(&llmOptions{}, opts...)
if llmOpts.toolWorkflowSW != nil {
    // omitted multiple lines...

    safego.Go(ctx, func() {
          // omitted multiple lines...

          llmOpts.toolWorkflowSW.Send(msg, nil)
       }
    })
}
```

如上面代码所示，通用的 NodeOption 中，可以定义和承载“节点特定的”具体 option，这个机制来自于 [Eino](https://github.com/cloudwego/eino) 的 [CallOption 机制](https://www.cloudwego.io/zh/docs/eino/core_modules/chain_and_graph_orchestration/call_option_capabilities/)。
需要注意的是，同一个节点类型，如果实现了多个接口，则这些接口对 NodeOption 的处理需要一致，即要么都需要 NodeOption，要么都不需要。
