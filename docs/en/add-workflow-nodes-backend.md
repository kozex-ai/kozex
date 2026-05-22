This article introduces how to add custom nodes to the Workflow in the **backend business logic**.
A node is a basic unit of a workflow and contains complete, relatively independent business logic. A workflow consists of at least two nodes: start and end. The complete process of a node—from being created by dragging it onto the canvas, to being passed to the backend and saved as part of a workflow, and then executed after instantiation—is illustrated in the following diagram:
![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/23860589ab064f9cb0fe1071c10db2a5~tplv-goo7wpa0wc-image.image)
A node on the canvas must belong to a node type, such as a large model node, a plugin node, and so on. Add a new node type to encapsulate business-related specialized logic, which can extend the capabilities of the workflow. In general, to add a basic node type, you need to do the following:

* Frontend: Refer to the document titled "Add new workflow node types(frontend)".
* Backend:
   * Add one Go file and define two structs in it: a node Config struct (which implements the NodeAdaptor and NodeBuilder APIs), and a node execution struct (which implements the InvokableNode and other APIs).
   * Configure the metadata for this new node type (such as name, style, and general execution parameters like timeout), and register the NodeAdaptor for this new type.

## Introduction to the Workflow operation mechanism
Kozex is built on the workflow capabilities of [Eino](https://github.com/cloudwego/eino) and is a directed acyclic graph (DAG) that includes both control flow and data flow.

* Control flow: Except for the 'Start' node, the prerequisites for a node to begin execution are that all of its incoming edges have produced results, and at least one of those results is 'successful'.
   * An edge can have three states: not completed, successful, and skipped. The initial state is "Incomplete." Both "Success" and "Skip" are considered as having produced a result.
   * When the source node of an edge completes execution, or when the source node contains a branch and the selection result of that branch is the current edge, the status of the current edge changes to "success".
   * When the starting node of the edge contains a branch and the selection result of that branch is not the current edge, the state of the current edge changes to "Skipped".
   * The 'skip' state can be propagated; that is, if the results of all incoming edges to a node are 'skip', then the results of all outgoing edges from that node are automatically set to 'skip'.

<div style="text-align: center"><img src="https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/335381ec3cd64cbfb991d2177fdc1077~tplv-goo7wpa0wc-image.image" width="2510px" /></div>

In the above diagram, the execution sequence is:

   1. If execution of the large model node fails, the process enters the exception branch. Edge "1" is not selected and is marked as "skipped", while edge "2" is selected and is marked as "successful".
   1. At this point, the two incoming edges to the "End" node, "1" and "4", are in the states "skipped" and "not completed", respectively, so the "End" node cannot start execution.
   1. Text processing has been completed. Edge 3 is marked as "Success".
   1. When the code node has finished executing, edge 4 is marked as "success".
   1. The two incoming edges of the "end" node, "1" and "4", correspond to "skip" and "success", respectively. The "end" node then begins execution.
* Data flow: The input data for a node is a `map[string]any` formed by merging the output fields from "any" predecessor node that the node references, configured static values, and configured variable values. As shown in the figure

![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/17a71d26e7d9496589e5651c37b4b52a~tplv-goo7wpa0wc-image.image)
## Write node logic

1. Create a new directory under backend/domain/workflow/internal/nodes, using your node's name as the directory name. For example, custom.
2. In this directory, create a new Go file, such as custom.go.
3. In this Go file, add a new struct such as YourNode to represent your node. You can define **any** fields required at runtime by the node within the struct, and these fields do not need to be exported.
4. This struct must implement one or more specific APIs, which represent "node execution logic."

For example:
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


5. The "Invoke" in the code above is a function with a fixed signature. This function receives input mapped from "predecessor nodes," executes its own business logic, and finally provides output for use by "successor nodes." Input and Output are not streams. **YourNode implements the InvokableNode API, which is defined in backend/domain/workflow/internal/nodes/node.go and does not require modification**:

```Go
// InvokableNode is a basic workflow node that can Invoke.
// Invoke accepts non-streaming input and returns non-streaming output.
// It does not accept any options.
type InvokableNode interface {
    Invoke(ctx context.Context, input map[string]any) (
       output map[string]any, err error)
}
```


6. The input and output of Invoke are both map[string]any. So what is inside the map? These are the fields specifically configured on the canvas, as well as the values mapped from predecessor nodes. For example, a simple large model node configuration:

   ![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/ede0f156e39c45d99c2b6e4f8e626bf6~tplv-goo7wpa0wc-image.image)

The user has configured input and output fields. As a result, the input map passed to the YourNode.Invoke method is `map[string]any{"input": "xxx"}`, and the returned output map is `map[string]any{"output": "yyy"}`.
## Set the node name

1. In backend/domain/workflow/entity/node_meta.go, give your new node a unique name:
   ```Go
   // NodeType is an enumeration value of a node type on the Coze workflow server.
   NodeTypeMyAwesomeNode NodeType = "MyAwesomeNode"
   ```

2. In the same file, configure some general "metadata" attributes for your new node. Take JSON serialization nodes as an example:
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


Focus areas:

* ID: A unique int64. For custom node types, it is recommended to start from a large number, such as 1000 or higher. This ID is primarily used for interacting with the frontend.
* Key: A unique string that accurately describes the functional role of the node. The backend primarily uses Key as the unique identifier, including those stored in the database.
* Name, Category, Desc, Color, IconURL, EnUSName, EnUSDescription: Display information for the canvas.
* SupportBatch: Whether the node can be configured to use "batch processing" mode. The ultimate control resides in the node styles on the frontend.
* ExecutableMeta: The general configuration for node runtime, which will be explained in detail later in the "advanced features" section.

Note: The Icon images corresponding to the nodes should be placed in docker/volumes/minio/default_icon/workflow_icon/. The file names of the icons need to match the IconURI field here.
## From canvas to backend

1. Workflow nodes from the canvas are sent to the backend, where they are converted into the Node struct in backend/domain/workflow/entity/vo/canvas.go. For example, a Node structure for a 'start node':
   ```JSON
   {
     "blocks": [], // Child nodes of composite nodes (batch processing, loop); start nodes do not have child nodes
     "data": { // Actual runtime-related configuration data
       "nodeMeta": { // Some meta information used by the frontend
         "description": "工作流的起始节点，用于设定启动工作流需要的信息",
         "icon": "https://lf3-static.bytednsdoc.com/obj/eden-cn/dvsmryvd_avi_dvsm/ljhwZthlaukjlkulzlp/icon/icon-Start-v2.jpg",
         "subTitle": "",
         "title": "开始"
       },
       "outputs": [ // Configuration of all "output" fields of the start node
         {
           "name": "input", // Only one output field
           "required": false, // Not required
           "schema": {
             "type": "string" // The array element type is string
           },
           "type": "list" // The overall type is an array
         }
       ],
     },
     "edges": null, // Edge information within composite nodes; not present in the start node
     "id": "100001", // Node ID. The ID of the start node is fixed.
     "meta": { // Location information for frontend use
       "position": {
         "x": -90.5329099821747,
         "y": -323.84999999999985
       }
     },
     "type": "1" // Node type, corresponds to NodeMeta.ID
   },
   ```

2. Before this Node struct can be executed, it must first be converted into the NodeSchema struct in backend/domain/workflow/internal/schema/node_schema.go on the backend.
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

3. Why convert? This is because the Node structure in the frontend schema is designed for the editing state of the workflow, while the NodeSchema structure in the backend schema is intended for the runtime state of the workflow.
4. How do I convert? Implement the NodeAdaptor interface:
   ```Go
   // NodeAdaptor provides conversion from frontend Node to backend NodeSchema.
   type NodeAdaptor interface {
       Adapt(ctx context.Context, n *vo.Node, opts ...AdaptOption) (
          *schema.NodeSchema, error)
   }
   ```

   This API is implemented by the Config Type of the node. For example, the JSON serialization node:
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

   Each node type must have a Config struct. The NodeAdaptor interface must be implemented.
   In Config, you can define fields of **any serializable** type (counter-examples: functions cannot be serialized, and the serialization of interfaces depends on specific types, so it is not recommended to use them in Config). These fields must be **exported**.
   These fields are essentially derived from the form configurations within nodes on the front-end canvas and serve as the primary source of information for node instantiation (with the remaining information coming from the node metadata in NodeMeta).
5. Register a NodeAdaptor: After obtaining a frontend node, how do you find the corresponding Config struct that implements the NodeAdaptor API? It is necessary to register the mapping from NodeType to NodeAdaptor. Modify the RegisterAllNodeAdaptors function in domain/workflow/internal/canvas/adaptor/to_schema.go:
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


## Generate node instances
Now that we have a unified backend NodeSchema, the final step is to actually instantiate this node and execute it.

1. When the workflow is running, the system obtains the NodeSchema and calls the New function in backend/domain/workflow/internal/compose/node_builder.go to ultimately convert the NodeSchema into eino's Lambda:
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

   Lambda is a fundamental component of the [Eino](https://github.com/cloudwego/eino) framework. It means "an arbitrary user-defined function". See [Eino Lambda User Guide](https://www.cloudwego.io/zh/docs/eino/core_modules/components/lambda_guide/).

2. During the instantiation process above, NodeSchema.Configs implements the NodeBuilder API:
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

   Using JSON-serialized nodes as an example, the simplest implementation of the NodeBuilder API is as follows:
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

   Since the JSON serialization node is completely stateless and can perform serialization as soon as it receives the input, there is no need to define any member variables in the Serializer type. In fact, you can define "any field that is required at runtime, remains stable across multiple runs, and does not change based on the current input"****, such as a question node:
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

3. Taking the JSON serialization node as an example, the build produces a *Serializer type, which is the concrete implementation class of the node. The Serializer implements the Invoke method, which corresponds to the InvokableNode API in backend/domain/workflow/internal/nodes/node.go:
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

   The semantics of Invoke are: non-streaming input and non-streaming output. Both the input and output types are `map[string]any`.
   In the same Go file, three other execution APIs are defined:
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

   Invoke, Stream, Collect, and Transform are the four basic streaming interaction paradigms of the [Eino](https://github.com/cloudwego/eino) framework. If you are interested in learning more, you can see [Eino Stream Programming Key Points](https://www.cloudwego.io/zh/docs/eino/core_modules/chain_and_graph_orchestration/stream_programming_essentials/).

4. For the same node type, it is possible to implement only one API or multiple APIs, depending on the specific business requirements of the node. For example,
   1. At present, most nodes have only implemented the InvokableNode API, because most nodes can only handle non-streaming input and produce non-streaming output.
   2. The "batch" and "loop" nodes implement the InvokableNodeWOpt API, as these two composite nodes need to pass information through to their internal nodes. This passthrough is achieved by passing NodeOption into the Invoke method.
   3. The large model node additionally implements the StreamableNodeWOpt API, as it can generate true streaming output. It may also need to transparently pass information to the internal workflow tool, and this requires the use of NodeOption.
   4. The output node additionally implements the TransformableNode API because it needs to receive streaming input and produce streaming output (displayed on the screen in a typewriter-style effect).
   At this point, it should already be possible to implement a basic new node type, drag it onto the canvas, and perform a test run. The following sections will provide individual introductions to advanced features related to nodes, which you can explore in more detail as needed.

## Advanced features
### Perceptual input and output types and their mappings
Sometimes, when a node is executed, it needs to know the "source" of its output fields. For example: a large model node (JSON-structured output), an input node (fields that require input), or a question node (information to be extracted from the answer). You can store the OutputTypes from the NodeSchema in the Node within NodeBuilder for use at runtime, for example, in the Build method of an input node:
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

Similarly, the InputTypes (the types of input fields) from the NodeSchema can be stored in the Node for use at runtime.
On the other hand, some nodes need to know the "source" of an output field. For example, batch processing and loop nodes need to know which field of which internal node a particular output field comes from. At this point, save NodeSchema.OutputSources in the Node, as in the Build method of the loop node:
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

Similarly, the InputSources of NodeSchema (the sources of input fields) can be stored in Node for use during node runtime.
### Customize frontend schemas
Most nodes that are even slightly complex require a customized frontend schema. For example, a large model node needs skill information, a condition node requires branch selection information, a sub-workflow node needs the sub-workflow ID and version, and a batch processing node requires the batch size and concurrency, and so on. All of this information needs to be configured by the user on the canvas, but it does not fall under 'user-defined fields'; instead, it is a fixed configuration specific to each node. These custom configurations are uniformly defined in the Inputs type of the frontend schema (backend/domain/workflow/entity/vo/canva.go):
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

For new node types, you can also extend the above Inputs structure to add fields specific to the new node type, and then parse these fields in NodeAdaptor to map them to the standard fields of NodeSchema. Because no matter how unique they are, **these pieces of information can ultimately be classified into three types: "input", "output", and "configuration"**.
These type-specific fields are highly likely to require customized frontend form support, so frontend adaptation must be handled in parallel.
### Customize input and output display fields
For certain node types, the 'input field information' displayed on the canvas during trial runs does not exactly match the `map[string]any` input at runtime. For example, the display of the canvas input for a condition:
![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/0e8f95b6c76a4811997237c19a62ad6f~tplv-goo7wpa0wc-image.image)

When displaying the left value of "Condition 1", also explicitly show that this field originates from the "Batch Processing" node. At this point, this node type needs to implement the `CallbackInputConverted` API to distinguish between the displayed input information and the actual input information:
```Go
// CallbackInputConverted converts node input to a form better suited for UI.
// The converted input will be displayed on canvas when test run, 
// and will be returned when querying the node's input through OpenAPI.
type CallbackInputConverted interface {
    ToCallbackInput(ctx context.Context, in map[string]any) (map[string]any, error)
}
```

For example, the condition node implements `CallbackInputConverted`; refer to backend/domain/workflow/internal/nodes/selector/callbacks.go.
On the other hand, some field types need to display information on the canvas that is not exactly the same as the original output. For example, for a question node, the original output is the "final answer," while the test run output shown on the canvas includes the entire question and answer process for all rounds. Or, for example, in the case of code nodes and similar cases, if formatted output fails, an additional "warning message" will be displayed. At this point, this node type needs to implement the `CallbackOutputConverted` API to distinguish between the displayed output information and the actual output information:
```Go
// CallbackOutputConverted converts node input to a form better suited for UI.
// The converted output will be displayed on canvas when test run,
// and will be returned when querying the node's output through OpenAPI.
type CallbackOutputConverted interface {
    ToCallbackOutput(ctx context.Context, out map[string]any) (*StructuredCallbackOutput, error)
}
```

For example, the question node implements `CallbackOutputConverted`; see backend/domain/workflow/internal/nodes/qa/question_answer.go.
Kozex uses the callback mechanism of the [Eino](https://github.com/cloudwego/eino) framework to achieve the effect of passing the modified input/output through a bypass. For information about the Eino Callback mechanism, refer to [Eino Callback User Manual](https://www.cloudwego.io/zh/docs/eino/core_modules/chain_and_graph_orchestration/callback_manual/).

### Exception handling strategies
Some nodes already support general exception handling strategies, including timeout, retry, fallback to default output after an exception, and executing an exception branch after an exception. These nodes include nodes such as large model nodes and code nodes. The configuration interface is as follows:
![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/8622ae4ae8a9482eaa97dbf3444e953a~tplv-goo7wpa0wc-image.image)

In the diagram, the "Code" node is configured with a 60-second timeout, does not retry, and executes the exception flow when an exception occurs. These configurations are all included in the SettingOnError structure of the frontend schema:
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

To enable the above exception handling strategy in a new node type, you only need to add an "Exception Handling" section to the node's frontend form and ensure that the schema sent to the backend includes SettingOnError. By doing this, **the exception execution policy will take effect automatically**.
### Nodes with the "branch" functionality
Some nodes have a "branch selection" function, such as condition, intent recognition, and question answering. On the canvas, this is shown by there being more than one "port" after this node:
![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/54cfd18911084eb68f0f6274b402e882~tplv-goo7wpa0wc-image.image)

In the front-end schema definition, a "port" refers to an element included on an edge (Edge). For example:
```JSON
{
  "sourceNodeID": "133234",
  "targetNodeID": "163493",
  "sourcePortID": "branch_0"
}
```

This frontend schema fragment indicates a connection that starts from node 133234 and leads to node 163493, with the starting port being "branch_0". This branch selection functionality, as implemented in Kozex, makes choices among multiple 'ports' using Eino's [Branch](https://www.cloudwego.io/zh/docs/eino/core_modules/chain_and_graph_orchestration/orchestration_design_principles/#branch) mechanism, based on the actual output of the node—for example, the branch selected by a condition, the intent identified by intent recognition, or the option chosen in a question-and-answer scenario. In the workflow topology, this is reflected by an additional Branch appearing after a Node.
These "ports" are divided into three categories:

1. Typical ports include, for example, each branch of a condition, each intent in intent recognition, and each static option in question answering (for dynamic options, they are grouped together as a single port).
2. Default port: for example, the else branch of a condition, the default intent in intent recognition, or the "Other" option in Q&A.
3. Exception port: If the node enables the general exception handling policy described above, this port will be available. The framework handles this, so developers do not need to be concerned with it.

If you want to enable the Branch feature for a new node, you need to ensure the following in the front-end schema:

1. For standard ports, the sourcePortID follows the format "branch_%d". For example, in intent recognition, the first intent is branch_0, and the second is branch_1.
2. Default port, sourcePortID = "default"
3. Exception port: sourcePortID = "branch_error"

In addition, the Config Type of the node needs to implement two APIs:

1. BranchBuilder: Receives the actual output of a node, calculates the corresponding port, and constructs a Branch.

```Go
// BranchBuilder builds the extractor function that maps node output to port index.
type BranchBuilder interface {
    BuildBranch(ctx context.Context) (extractor func(ctx context.Context,
       nodeOutput map[string]any) (int64, bool /*if is default branch*/, error), hasBranch bool)
}
```

Using the implementation of the condition as an example:
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

When the framework calls BuildBranch, Config has already completed the NodeAdaptor process, and all its internal fields are available for use.

2. BranchAdaptor: Determine which ports "should" exist based on the frontend schema (vo.Node) for validation:

```Go
// BranchAdaptor provides validation and conversion from frontend port to backend port.
type BranchAdaptor interface {
    ExpectPorts(ctx context.Context, n *vo.Node) []string
}
```

Using the implementation of intent recognition as an example:
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

Neither BranchBuilder nor BranchAdaptor needs to be aware of "exception branches", as these are entirely handled by the Kozex framework.
### The input and output assignment fallback strategy
The node input and node output of a Workflow are both `map[string]any`. There are two possible scenarios:

* The node has declared an input field, but the actual input map does not contain the corresponding key for this field.
* The node declared an output field, but the actual output map does not contain the key corresponding to this field.

However, the underlying [Eino](https://github.com/cloudwego/eino) Workflow engine requires that, when field mapping occurs, the upstream field must exist in the upstream output map. Therefore, for node types with user-customizable output fields, it is recommended to configure `PostFillNil = true`. Kozex will then automatically assign nil to any missing output fields in the map.
On the other hand, the business logic of certain nodes will automatically set nil keys in the input map to zero values. For these nodes, you can configure `PreFillZero = true`, and Kozex will automatically replace input fields that have nil values with the corresponding zero value for their type.
For example, in a text processing node, Nil values in the input must be converted to zero values, but the output field can never be Nil:
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

The variable merge node is designed to retain the original nil values in the output; however, some fields may be missing from the output if all values in a group are nil.
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

### The canvas typewriter effect
For output nodes and end nodes (in return text mode), if 'streaming output' is selected, the streaming output content will be displayed on the canvas in real time with a typewriter effect during test runs.
If you want your new node to enable this feature—that is, to display streaming output on the canvas in real time during trial runs—you need to set IncrementalOutput = true in NodeMeta:
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

In addition to this configuration, there are two prerequisites:

1. The node can indeed return streaming data, that is, it implements the StreamableNode or TransformableNode API.
2. In backend/domain/workflow/internal/execute/callback.go, modify the `func (n *NodeHandler) OnEndWithStreamOutput(ctx context.Context, info *callbacks.RunInfo, output *schema.StreamReader[callbacks.CallbackOutput]) context.Context` method to ensure that your new node type can emit the **`*NodeStreamingOutput *`** event type (currently, it is not easy to change, but we will improve it in the future).

### Composite nodes
Nodes such as batch processing and loop are referred to as composite nodes (CompositeNode) and have the following characteristics:

* In the frontend schema, a Node contains at least one Block (internal node).
* Internal nodes may reference the outputs from composite nodes and the outputs from their predecessor nodes.
* Internal nodes may be interrupted and resumed, which causes composite nodes to also be interrupted and resumed.

To add a composite node type, set `IsComposite = true` in NodeMeta, for example:
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

Implementing a new composite node is a major undertaking. You can refer to the implementation in backend/domain/workflow/internal/nodes/batch/batch.go and backend/domain/workflow/internal/nodes/loop/loop.go. Main areas of concern:

* Internal workflow scheduling
* Interruption recovery for internal workflows

The new composite nodes need to be adapted together with the frontend.
### The batch processing mode of regular nodes
Currently, three node types support "batch processing mode": large models, plugins, and sub-workflows.
After enabling "batch processing mode," additional batch processing-related configurations will appear in the vo.Node section of the node's frontend schema:
```Go
type NodeBatch struct {
    BatchEnable    bool     `json:"batchEnable"`
    BatchSize      int64    `json:"batchSize"`
    ConcurrentSize int64    `json:"concurrentSize"`
    InputLists     []*Param `json:"inputLists,omitempty"`
}
```

To enable batch mode for a node type, set `SupportBatch = true` in NodeMeta:
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

Frontend adaptation is also required.
### Node token consumption statistics
The workflow can track the model's token consumption at both the node level and the overall workflow level.
If your new node type needs to track token consumption—for example, if you need to display the node's token consumption on the trial run page, add the node's token consumption to the overall workflow consumption, or include the node's token consumption in the messages returned by the OpenAPI—you must set MayUseChatModel = true in node_meta.go:
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

Once MayUseChatModel is set to true, if a large model is actually used, statistics will be collected automatically. If the configuration is set but the large model is not used, the token count will be zero.
### Pass a runtime Option
In node.go, the framework provides four additional run APIs that accept NodeOption, which nodes can choose to implement:
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

NodeOption can inject additional information into a node at runtime, apart from input and state, for example:

* The sub-workflow node passes Eino's runtime option to the internal sub-workflow.
* Batch processing nodes and loop nodes pass information about resuming after an interruption to their internal nodes.
* StreamWriter for passing the message stream from the large model node to the internal workflow tool.

To implement a custom NodeOption, you can refer to the implementation of the LLM node:
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

Obtain and use this NodeOption in the LLM node:
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

As shown in the code above, the general NodeOption can define and hold node-specific options. This mechanism is derived from [Eino](https://github.com/cloudwego/eino)'s [CallOption mechanism](https://www.cloudwego.io/zh/docs/eino/core_modules/chain_and_graph_orchestration/call_option_capabilities/).
It is important to note that for the same node type, if multiple APIs are implemented, these APIs must handle NodeOption consistently; that is, either all require NodeOption or none do.
