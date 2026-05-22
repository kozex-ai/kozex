* For information on how to develop and debug frontend modules, refer to: [Development Standards - Code development and testing](development-standards.md#code-development-and-testing)
* The MR related to this document: https://github.com/kozex-ai/kozex/pull/215

## Background
This document uses the addition of a **JSON serialization node** as an example. It shows how to serialize the return variable of a preceding node into a string and demonstrates how to add a node type in the Kozex frontend interface.
| ![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/240fbbeec7d24fee9976828d9d5245bb~tplv-goo7wpa0wc-image.image) <br> The node panel can add JSON serialization nodes. | ![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/44331a1a27d04904a8fa9e5e394c7e2d~tplv-goo7wpa0wc-image.image) <br> JSON serialization node configuration | ![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/a5d73159c95c4911951f1f2e5cc9a58a~tplv-goo7wpa0wc-image.image) <br> Test run of a JSON serialization node |
| --- | --- | --- |
## Key terminology
Let us first define some common concepts.

* **Node type**: In node registration logic, each node has its own type. This type should be agreed upon with the backend and must not conflict with existing nodes.
* **Node instance**: After a node type is added to the canvas, a workflow node instance is generated.
* **Stage node**: A node display on the stage canvas that provides a summary of key node information and, during trial runs, displays a trial run result bar.
* **Node form**: After clicking a stage node, the node form that appears in the side drawer displays all configuration items for the node instance.
* **Dynamic ports**: By default, node ports are static inputs and outputs. However, in some models, the ports change dynamically based on the node configuration. For example, if an intent recognition node has multiple options, it will have multiple output ports.
* **Node Registry**: Node registration configuration
* **Form Meta**: Node form metadata configuration
* **VO**: View Object, a display layer object used directly to present the UI
* **DTO**: Data Transfer Object, an object transferred from the backend to the frontend

## Requirements confirmation

1. Functionality confirmation
   1. Does it support single-node debugging?
   2. Does it support exception settings?
   3. Constraints
      1. Maximum number of items that can be added
      2. Types of variable selection
      3. Input length and character restrictions
      4. Other
2. Design confirmation
   1. Node card: contains dynamic ports
   2. Node form, linkage between form fields
   3. The node form in read-only mode
3. Backend confirmation
   1. Node schema and type
   2. Can a node support single-node debugging

A separate UI analysis will be conducted below.
## Interface analysis
| **Interface examples** | **Analysis** |
| --- | --- |
| Stage node <br> ![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/89055c5afc6d458a98d0845d1a1d8bde~tplv-goo7wpa0wc-image.image) | * The header area already contains components, so no changes are needed. <br> * The input area already contains a component; no changes are needed. <br> * Output area: components already exist; no changes needed. <br>  <br>  |
| Node form configurations <br> ![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/5b4bceb1b173497d84f3826929d77442~tplv-goo7wpa0wc-image.image) | * In the input area, components are already available; simply configure them as needed. <br> * In the output area, existing components are available; simply configure them as needed. |
| Single node test run form <br> ![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/db450559658a4e3c9dafa80d19563f05~tplv-goo7wpa0wc-image.image) | When the input contains variables, the single-node trial run needs to extract the variables to generate a form. <br>  <br> * The existing logic generates different components based on the variable type of the input, so no changes are required. <br>  <br>  |
## API conventions
Items the front end and back end must agree upon at a minimum:

* Node type: here we assume that the node type for JSON serialization is '58'.
* The format of a node in the Schema is the same as the format in which the save API stores it on the backend. Since the content of this node is relatively simple, you can just reuse the existing format.

```TypeScript
import { type NodeDataDTO } from '@flow-workflow/base';
```

Alright, everything is ready; you can start writing code now.
## Node development
Development centered on a node mainly includes the following parts:

![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/6f09d384431d4b0591ede2ac0336605b~tplv-goo7wpa0wc-image.image)

### Node definition

1. Go to `packages/workflow/base/src/types/node-type.ts` to update the node type. Since this type needs to be confirmed with the backend, it must be entered manually for now. In the future, it will be collected automatically via a command.
   ```C++
   /**
    * Definition of basic node types
    */
   export enum StandardNodeType {
     // ...
     JsonStringify = '58',
   }
   ```

2. Go to `frontend/packages/workflow/adapter/base/src/utils/get-enabled-node-types.ts`, add custom display logic, and add the newly defined type to enabledNodeTypes and typeEnable.
   ```TypeScript
   export const getEnabledNodeTypes = (_params: {
     loopSelected: boolean;
     isProject: boolean;
     isSupportImageflowNodes: boolean;
     isSceneFlow: boolean;
     isBindDouyin: boolean;
   }) => {
       const nodesMap = {
         // ...
         [StandardNodeType.JsonStringify]: true,
       };
   };
   ```

3. Use a scaffolding command to quickly generate a node skeleton.
   ```TypeScript
   > cd packages/workflow/playground
   > rushx create:node
   ```

   ![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/92a632df08d4489fb4b6ce05ccf2a81d~tplv-goo7wpa0wc-image.image)

It will be added to ` packages/workflow/playground/src/nodes-registries`. Now, let's analyze the directory structure:
```TypeScript
src/node - registries/json - stringify
├── components          // (Optional) Stores general - purpose node components
├── constants.ts        // (Optional) Stores node constant configurations
├── data - transformer.ts // (Optional) Stores node data transformation logic, front - end form data <-> back - end Schema data
├── form - meta.tsx       // (Required) Stores node form unit data configurations
├── form.tsx            // (Required) Stores node form rendering logic
├── hooks               // (Optional) Stores node custom business logic
├── index.ts            // (Required) Entry
├── node - content.tsx    // (Required) Stage node card component
├── node - registry.ts    // (Required) Node registration configuration
├── node - test.ts        // (Required) Single - node test configuration
├── types.ts            // (Required) Node types
└── utils               // (Optional) Node utility functions
```

Let's focus on the node definition component `node-registry.ts`.
```TypeScript
import {
  DEFAULT_NODE_META_PATH,
  DEFAULT_OUTPUTS_PATH,
} from '@coze-workflow/nodes';
import {
  StandardNodeType,
  type WorkflowNodeRegistry,
} from '@coze-workflow/base';

import { JSON_STRINGIFY_FORM_META } from './form-meta';
import { INPUT_PATH } from './constants';
import { test, type NodeTestMeta } from './node-test';

export const JSON_STRINGIFY_NODE_REGISTRY: WorkflowNodeRegistry<NodeTestMeta> = {
  type: StandardNodeType.JsonStringify,
  meta: {
    nodeDTOType: StandardNodeType.JsonStringify,
    size: { width: 360, height: 130.7 },
    nodeMetaPath: DEFAULT_NODE_META_PATH,
    outputsPath: DEFAULT_OUTPUTS_PATH,
    inputParametersPath: INPUT_PATH,
    test,
  },
  formMeta: JSON_STRINGIFY_FORM_META,
};
```

This file does not need to be modified. In it, `formMeta` defines the form for this node. Now, let's run it:
```TypeScript
cd apps/kozex
npm run dev
```

You should now be able to see the JSON serialization node in the node addition panel.
![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/7f478f62bf594b7da5f32dda0d7740e7~tplv-goo7wpa0wc-image.image)
After we add a node, if we click to expand the configuration panel, we will notice some issues:

1. The input parameter should support only one variable.
2. The output parameter format is incorrect
   ![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/a451f33076c24bb2bcf638d6398ef162~tplv-goo7wpa0wc-image.image)

Next, we will solve these problems.
### Form definitions
Let's take a look at `form-meta.tsx`. Of particular importance is the `FormMetaV2` type, which is the core API definition for form metadata. The underlying form engine uses the proprietary form engine developed by [FlowGram.ai](https://flowgram.ai/). For more details, see: https://flowgram.ai/guide/advanced/form.html
```TypeScript
import {
  ValidateTrigger,
  type FormMetaV2,
} from '@flowgram - adapter/free - layout - editor';

import { createValueExpressionInputValidate } from '@/node - registries/common/validators';
import {
  fireNodeTitleChange,
  provideNodeOutputVariablesEffect,
} from '@/node - registries/common/effects';

import { type FormData } from './types';
import { FormRender } from './form';
import { transformOnInit, transformOnSubmit } from './data - transformer';

export const JSON_STRINGIFY_FORM_META: FormMetaV2<FormData> = {
  // Node form rendering
  render: () => <FormRender />,

  // Validation trigger timing
  validateTrigger: ValidateTrigger.onChange,

  // Validation rules
  validate: {
    // Required
    'inputs.inputParameters.0.input': createValueExpressionInputValidate({
      required: true,
    }),
  },

  // Side - effect management
  effect: {
    nodeMeta: fireNodeTitleChange,
    outputs: provideNodeOutputVariablesEffect,
  },

  // Node back - end data -> front - end form data
  formatOnInit: transformOnInit,

  // Front - end form data -> node back - end data
  formatOnSubmit: transformOnSubmit,
};
```

#### Form rendering
Form rendering uses the component from `form.tsx`, which has input and output components built in by default.
```JavaScript
export const FormRender = () => (
  <NodeConfigForm>
    <InputsParametersField
      name={INPUT_PATH}
      title={I18n.t('node_http_request_params')}
      tooltip={I18n.t('node_http_request_params_desc')}
      defaultValue={[]}
    />

    <OutputsField
      title={I18n.t('workflow_detail_node_output')}
      tooltip={I18n.t('node_http_response_data')}
      id="jsonStringify-node-outputs"
      name="outputs"
      topLevelReadonly={true}
      customReadonly
    />
  </NodeConfigForm>
);
```

There is an issue above: the input currently only allows a single parameter. This should be modified here. We can implement our own input component in the components directory. This input component will support only one variable input and will not allow adding or removing inputs.
> packages/workflow/playground/src/node-registries/json-stringify/components/inputs/index.tsx

```TypeScript
import {
  FieldArray,
  type FieldArrayRenderProps,
} from '@flowgram-adapter/free-layout-editor';
import type { ViewVariableType, InputValueVO } from '@coze-workflow/base';
import { I18n } from '@coze-arch/i18n';

import { useReadonly } from '@/nodes-v2/hooks/use-readonly';
import { ValueExpressionInputField } from '@/node-registries/common/fields';
import { FieldArrayItem, FieldRows, Section, type FieldProps } from '@/form';

interface InputsFieldProps extends FieldProps<InputValueVO[]> {
  title?: string;
  paramsTitle?: string;
  expressionTitle?: string;
  disabledTypes?: ViewVariableType[];
  onAppend?: () => InputValueVO;
  inputPlaceholder?: string;
  literalDisabled?: boolean;
  showEmptyText?: boolean;
  nthCannotDeleted?: number;
}

export const InputsField = ({
  name,
  defaultValue,
  title,
  tooltip,
  disabledTypes,
  inputPlaceholder,
  literalDisabled,
  showEmptyText = true,
}: InputsFieldProps) => {
  const readonly = useReadonly();
  return (
    <FieldArray<InputValueVO> name={name} defaultValue={defaultValue}>
      {({ field }: FieldArrayRenderProps<InputValueVO>) => {
        const { value = [] } = field;
        const length = value?.length ?? 0;
        const isEmpty = !length;
        return (
          <Section
            title={title}
            tooltip={tooltip}
            isEmpty={showEmptyText && isEmpty}
            emptyText={I18n.t('workflow_inputs_empty')}
          >
            <FieldRows>
              {field.map((item, index) => (
                <FieldArrayItem key={item.key} disableRemove hiddenRemove>
                  <div style={{ flex: 3 }}>
                    <ValueExpressionInputField
                      name={`${name}.${index}.input`}
                      disabledTypes={disabledTypes}
                      readonly={readonly}
                      inputPlaceholder={inputPlaceholder}
                      literalDisabled={literalDisabled}
                    />
                  </div>
                </FieldArrayItem>
              ))}
            </FieldRows>
          </Section>
        );
      }}
    </FieldArray>
  );
};
```

Then replace the original input section:
> packages/workflow/playground/src/node-registries/json-stringify/form.tsx

```JavaScript
import { I18n } from '@coze-arch/i18n';

import { NodeConfigForm } from '@/node-registries/common/components';

import { OutputsField } from '../common/fields';
import { INPUT_PATH } from './constants';
import { InputsField } from './components/inputs';

export const FormRender = () => (
  <NodeConfigForm>
    <InputsField
      name={INPUT_PATH}
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      defaultValue={[{ name: 'input' } as any]}
      title={I18n.t('workflow_250429_01')}
      tooltip={I18n.t('workflow_250429_03')}
      required={false}
      layout="horizontal"
    />

    <OutputsField
      title={I18n.t('workflow_detail_node_output')}
      tooltip={I18n.t('node_http_response_data')}
      id="jsonStringify-node-outputs"
      name="outputs"
      topLevelReadonly={true}
      customReadonly
    />
  </NodeConfigForm>
);
```

#### Form validation
Form validation is defined in the validate field in `form-meta.tsx`. The key represents the field path and can be configured to support regular expressions (*), while the value is the validation function. For more details, see https://flowgram.ai/guide/advanced/form.html#%E6%A0%A1%E9%AA%8C
```JavaScript
export const JSON_STRINGIFY_FORM_META: FormMetaV2<FormData> = {

  // Validation trigger timing
  validateTrigger: ValidateTrigger.onChange,

  // Validation rules
  validate: {
    // Validate that the first item is required
    'inputs.inputParameters.0.input': createValueExpressionInputValidate({
      required: true,
    }),
  },
};
```

There is currently no need to modify the validation logic here; it also checks that the first item is required.
#### Data transformations
Because the backend schema format is inconsistent with the data format required by the frontend form engine, data conversion between the two is necessary.

* The `formatOnInit` defined in `form-meta.tsx` is typically implemented in a separate file `data-transformer.tsx`.
* Frontend-to-backend transformation: the `formatOnSubmit` defined in `form-meta.tsx`, typically implemented in a separate file, `data-transformer.tsx`.

```TypeScript
import { type NodeDataDTO } from '@coze - workflow/base';

import { type FormData } from './types';
import { OUTPUTS } from './constants';

/**
 * Node backend data -> Front - end form data
 */
export const transformOnInit = (value: NodeDataDTO) => ({
 ...(value?? {}),
  outputs: value?.outputs?? OUTPUTS,
});

/**
 * Front - end form data -> Node backend data
 * @param value
 * @returns
 */
export const transformOnSubmit = (value: FormData): NodeDataDTO =>
  value as unknown as NodeDataDTO;
```

Since there is no special conversion logic here, no modification is needed.
Note that there is an additional layer of conversion logic for input variables (inputsParameters) and output variables (outputs), which automatically converts the variable formats. For more details, see: frontend/packages/workflow/nodes/src/workflow-json-format.ts

#### Side effect management
Side effect management is primarily used to handle linkage between fields. For example, when the value of field A changes, the value of field B needs to be updated. This component is not currently in use. For more information, see: https://flowgram.ai/guide/advanced/form.html#%E5%89%AF%E4%BD%9C%E7%94%A8-effect.
#### Variable synchronization
The generation logic for node output variables is as follows:

![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/28c0185e9191491cbe5ea4c2bb8f8a6d~tplv-goo7wpa0wc-image.image)

Generally, use provideNodeOutputVariablesEffect by default.
```TypeScript
import {
  provideNodeOutputVariablesEffect,
} from '@/node-registries/common/effects';

export const FORM_META: FormMetaV2<FormData> = {
  ...
  effect: {
    outputs: provideNodeOutputVariablesEffect,
  },
};
```

Since there is no linkage or variable creation or destruction, there is no need to modify the side effects here.
#### Output variable modifications
To modify output variables, simply change the constants inside `frontend/packages/workflow/playground/src/node-registries/json-stringify/constants.ts`.
```JavaScript
import { nanoid } from 'nanoid';
import { ViewVariableType } from '@coze-workflow/variable';

// Input parameter path. Functions such as trial operation rely on this path to extract parameters.
export const INPUT_PATH = 'inputs.inputParameters';

// Define fixed output parameters.
export const OUTPUTS = [
  {
    key: nanoid(),
    name: 'output',
    type: ViewVariableType.String,
  },
];

export const DEFAULT_INPUTS = [{ name: 'input' }];
```

At this point, we delete the original JSON serialization node and create a new serialization node to test it.
<div style="text-align: center"></div>

![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/376d08d73ffc462a9005fefb6f84dde5~tplv-goo7wpa0wc-image.image)

The relevant capabilities have been largely implemented.
### Stage nodes
![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/171f349aa0eb43d7a8907bcc7a532489~tplv-goo7wpa0wc-image.image)

A node is a card displayed on the canvas, implemented in `node-content.tsx`.
```JavaScript
import { InputParameters, Outputs } from '../common/components';

export function JsonStringifyContent() {
  return (
    <>
      <InputParameters />
      <Outputs />
    </>
  );
}
```

Note that you need to export this component in the current directory's `index.ts` file, and then have another component reference it. All of this is handled automatically by the scaffolding. Since the stage node for JSON serialization only displays input and output, no changes are needed here.
### Node registration
The generated node registry (in this case, JSON_STRINGIFY_NODE_REGISTRY) needs to be registered with the node list. This process is also automated by the scaffold. For details, see: `packages/workflow/playground/src/nodes-v2/constants.ts`.
### Test run
#### Single-node trial run
You need to define the form extraction logic for single-node test runs. This is specified in the node-test.ts file. If you set it to true, the default form extraction logic for test runs will be used.
```TypeScript
import type { NodeTestMeta } from '@/test-run-kit';

const test: NodeTestMeta = true;

export { test, type NodeTestMeta };
```

We assign a person object to the variable and try running it:
![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/dade1349b27c4a11bf47eb3aafd88e13~tplv-goo7wpa0wc-image.image)
We found a problem. For a single-node trial run, the person variable should be extracted and used as an input. Let's make some modifications.
> frontend/packages/workflow/playground/src/node-registries/json-stringify/node-test.ts

```TypeScript
import { FlowNodeFormData } from '@flowgram-adapter/free-layout-editor';

import {
  type NodeTestMeta,
  generateParametersToProperties,
} from '@/test-run-kit';

export const test: NodeTestMeta = {
  generateFormInputProperties(node) {
    const formData = node
      .getData(FlowNodeFormData)
      .formModel.getFormItemValueByPath('/');
    const parameters = formData?.inputs?.inputParameters;

    return generateParametersToProperties(parameters, {
      node,
    });
  },
};
export type { NodeTestMeta };
```

The above logic primarily extracts input parameters and then calls **generateParametersToProperties** with these parameters to generate the trial run form.
Try running it again:
![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/b8f327551c854e62a0f31b12ec12c6f6~tplv-goo7wpa0wc-image.image)
As you can see, there are no issues now.
#### A full process trial run
A full-process trial run is executed from the start node, so it does not concern the current node.
### List of modifications
To summarize, the following files need to be modified based on the scaffolding:
| **Edit files** | **Function** |
| --- | --- |
| constants.ts | Definition of output variable types |
| form.tsx | Replace the input component with the new version |
| node-test.ts | Single-node trial run adjustments |
| components/inputs/index.tsx | Add specialized input components |

For the final changes, you can refer to this MR: https://github.com/kozex-ai/kozex/pull/215
## Node integration testing
After completing the above steps, you can begin integration testing with the backend and adjust related features as needed.
