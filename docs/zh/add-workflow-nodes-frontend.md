* 关于如何开发调试前端模块，参考：[开发调试](development-standards.md#开发调试)
* 本文档涉及到的 MR：https://github.com/kozex-ai/kozex/pull/215

## 背景
本文档以添加一个 **JSON 序列化节点**为例，将一个前序节点的返回变量序列化成字符串，演示如何在 Kozex 前端界面添加一个节点类型。
| ![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/240fbbeec7d24fee9976828d9d5245bb~tplv-goo7wpa0wc-image.image) <br> 节点面板可以添加 JSON 序列化节点 | ![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/44331a1a27d04904a8fa9e5e394c7e2d~tplv-goo7wpa0wc-image.image) <br> JSON 序列化节点配置 | ![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/a5d73159c95c4911951f1f2e5cc9a58a~tplv-goo7wpa0wc-image.image) <br> JSON 序列化节点试运行 |
| --- | --- | --- |
## 基本概念
我们先约定一些常见的概念。

* **节点类型**：节点注册逻辑，每个节点有自己的类型（type），这个类型和后端约定好即可，不能和已有节点冲突
* **节点实例**：节点类型添加到画布后，会生成一个工作流节点实例
* **舞台节点**：舞台画布上的节点展示，包含节点的重要信息摘要，以及在试运行时会展示试运行结果 bar
* **节点表单**：点击舞台节点后，出现在侧边抽屉的节点表单，展示了节点实例的所有配置项
* **动态端口**：默认的节点端口是静态的输入和输出，而有些模型的端口根据节点配置动态变化的，如意图识别节点有多个选项就会有多个输出的端口
* **Node Registry**：节点注册配置
* **Form Meta**：节点表单元数据配置
* **VO**：View Object，显示层对象，直接用来展示 UI
* **DTO**：Data Transfer Object，数据传输对象，后端传输到前端的对象

## 需求确认

1. 功能确认
   1. 是否支持单节点调试
   2. 是否支持异常设置
   3. 限制条件
      1. 添加条目的上限
      2. 变量选择的类型
      3. 输入的长度和字符限制
      4. 其它
2. 设计确认
   1. 节点卡片，是否包含动态端口
   2. 节点表单，表单字段之间的联动关系
   3. 只读态下节点表单
3. 后端确认
   1. 节点 schema 和 type
   2. 节点是否能支持单节点调试

下边单独进行 UI 分析。
## 界面分析
| **界面示例** | **分析** |
| --- | --- |
| 舞台节点 <br> ![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/89055c5afc6d458a98d0845d1a1d8bde~tplv-goo7wpa0wc-image.image) | * 头部区域，已有组件，这个无需改动 <br> * 输入区域，已有组件，这个无需改动 <br> * 输出区域，已有组件，这个无需改动 <br>  <br>  |
| 节点表单配置 <br> ![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/5b4bceb1b173497d84f3826929d77442~tplv-goo7wpa0wc-image.image) | * 输入区域，已有组件，做相关配置即可 <br> * 输出区域，已有组件，做相关配置即可 |
| 单节点试运行表单 <br> ![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/db450559658a4e3c9dafa80d19563f05~tplv-goo7wpa0wc-image.image) | 当输入中包含变量时，单节点试运行需要提取变量生成表单。 <br>  <br> * 已有逻辑，输入会根据变量类型生成不同的组件，这个无需改动 <br>  <br>  |
## 接口约定
前后端至少需要约定如下内容：

* 节点类型，这里我们假设 JSON 序列化的节点类型为 '58'
* 节点在 Schema 中的格式，即 save 接口保存在后端的格式，由于这个节点内容比较简单，复用已有格式即可

```TypeScript
import { type NodeDataDTO } from '@flow-workflow/base';
```

好了万事俱备，可以开始写代码了。
## 节点开发
围绕一个节点开发，主要包括如下几个部分：
![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/cd09e3a0aaa84e27829eec3d2d2df2ca~tplv-goo7wpa0wc-image.image)
### 节点定义

1. 去 `packages/workflow/base/src/types/node-type.ts` 更新节点类型（由于这个 type 需要和后端确认，所以需要手动填一下，后续会升级通过命令收集）。
   ```C++
   /**
    * 节点基础类型定义
    */
   export enum StandardNodeType {
     // ...
     JsonStringify = '58',
   }
   ```

2. 去 `frontend/packages/workflow/adapter/base/src/utils/get-enabled-node-types.ts` 添加自定义显示逻辑，在 enabledNodeTypes 和 typeEnable 中增加刚刚定义的类型。
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

3. 使用脚手架命令，快速生成一个节点的骨架。
   ```TypeScript
   > cd packages/workflow/playground
   > rushx create:node
   ```

   ![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/92a632df08d4489fb4b6ce05ccf2a81d~tplv-goo7wpa0wc-image.image)

会被添加到` packages/workflow/playground/src/nodes-registries`，现在我们分析下目录结构：
```TypeScript
src/node-registries/json-stringify
├── components          //（可选）存放节点通用组件
├── constants.ts        //（可选）存放节点常量配置
├── data-transformer.ts //（可选）存放节点数据转换逻辑，前端表单数据 <-> 后端 Schema 数据
├── form-meta.tsx       //（必选）存放节点表单元数据配置
├── form.tsx            //（必选）存放节点表单渲染逻辑
├── hooks               //（可选）存放节点自定义业务逻辑
├── index.ts            //（必选）入口
├── node-content.tsx    //（必选）舞台节点卡片组件
├── node-registry.ts    //（必选）节点注册配置
├── node-test.ts        //（必选）单节点测试配置
├── types.ts            //（必选）节点类型
└── utils               //（可选）节点实用函数
```

我们重点看下，节点定义组件 `node-registry.ts`。
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

该文件无需修改，其中 `formMeta` 定义了该节点的表单定义，此时我们来运行下：
```TypeScript
cd apps/kozex
npm run dev
```

应该在节点添加面板可以看到 JSON 序列化这个节点了。
![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/7f478f62bf594b7da5f32dda0d7740e7~tplv-goo7wpa0wc-image.image)
我们添加节点后，点击展开配置面板，会发现存在一些问题：

1. 入参应该只支持 1 个变量
2. 出参格式不对
   ![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/a451f33076c24bb2bcf638d6398ef162~tplv-goo7wpa0wc-image.image)

我们接下来解决这些问题。
### 表单定义
我们来看下 `form-meta.tsx`，其中比较重要的是 `FormMetaV2` 类型，这是表单元数据的核心接口定义。底层表单引擎使用了 [FlowGram.ai](https://flowgram.ai/) 自研表单引擎，详情可以参见：https://flowgram.ai/guide/advanced/form.html
```TypeScript
import {
  ValidateTrigger,
  type FormMetaV2,
} from '@flowgram-adapter/free-layout-editor';

import { createValueExpressionInputValidate } from '@/node-registries/common/validators';
import {
  fireNodeTitleChange,
  provideNodeOutputVariablesEffect,
} from '@/node-registries/common/effects';

import { type FormData } from './types';
import { FormRender } from './form';
import { transformOnInit, transformOnSubmit } from './data-transformer';

export const JSON_STRINGIFY_FORM_META: FormMetaV2<FormData> = {
  // 节点表单渲染
  render: () => <FormRender />,

  // 验证触发时机
  validateTrigger: ValidateTrigger.onChange,

  // 验证规则
  validate: {
    // 必填
    'inputs.inputParameters.0.input': createValueExpressionInputValidate({
      required: true,
    }),
  },

  // 副作用管理
  effect: {
    nodeMeta: fireNodeTitleChange,
    outputs: provideNodeOutputVariablesEffect,
  },

  // 节点后端数据 -> 前端表单数据
  formatOnInit: transformOnInit,

  // 前端表单数据 -> 节点后端数据
  formatOnSubmit: transformOnSubmit,
};
```

#### 表单渲染
表单渲染使用了 `form.tsx` 中的组件，该组件默认内置了输入组件，以及输出组件。
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

上边有个问题，就是输入只能输入一个参数，应该在这里进行修改，这里我们可以在 components 目录实现一个自己的输入组件，这个输入组件就支持一个变量输入，无法新增和删除。
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

然后替换原来的输入部分：
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

#### 表单校验
表单校验在 `form-meta.tsx` 中的 validate 字段定义，key 为字段路径，支持正则（*）配置，value 为校验函数，具体可以参见：https://flowgram.ai/guide/advanced/form.html#%E6%A0%A1%E9%AA%8C
```JavaScript
export const JSON_STRINGIFY_FORM_META: FormMetaV2<FormData> = {

  // 验证触发时机
  validateTrigger: ValidateTrigger.onChange,

  // 验证规则
  validate: {
    // 验证第一项为必填
    'inputs.inputParameters.0.input': createValueExpressionInputValidate({
      required: true,
    }),
  },
};
```

校验逻辑这里暂无需修改，也是验证第一项为必填。
#### 数据转化
由于后端 Schema 格式以及前端表单引擎需要的数据格式不一致，因此两者存在一个数据转化。

* 后端转前端：定义在 `form-meta.tsx` 的 `formatOnInit`，一般在单独的文件 `data-transformer.tsx` 实现
* 前端转后端：定义在 `form-meta.tsx` 的 `formatOnSubmit`，一般在单独的文件 `data-transformer.tsx` 实现

```TypeScript
import { type NodeDataDTO } from '@coze-workflow/base';

import { type FormData } from './types';
import { OUTPUTS } from './constants';

/**
 * 节点后端数据 -> 前端表单数据
 */
export const transformOnInit = (value: NodeDataDTO) => ({
  ...(value ?? {}),
  outputs: value?.outputs ?? OUTPUTS,
});

/**
 * 前端表单数据 -> 节点后端数据
 * @param value
 * @returns
 */
export const transformOnSubmit = (value: FormData): NodeDataDTO =>
  value as unknown as NodeDataDTO;
```

由于这里没有特殊的转化逻辑，也无需修改。
需要注意的是，对于输入变量（inputsParamters）以及输出变量（outputs）还有一层转换逻辑，会对变量格式进行自动换转，详情参见：frontend/packages/workflow/nodes/src/workflow-json-format.ts

#### 副作用管理
副作用管理主要用于处理联动，例如字段 A 变更，需要更新字段 B 的值，这个组件暂时没有用到，详情可以参见：https://flowgram.ai/guide/advanced/form.html#%E5%89%AF%E4%BD%9C%E7%94%A8-effect。
#### 变量同步
节点输出变量的生成逻辑如下：
![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/b9545158b99b4def89b173230b3517fc~tplv-goo7wpa0wc-image.image)
一般默认使用 provideNodeOutputVariablesEffect 即可。
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

由于没有联动关系，或者变量的生成与销毁，副作用这里无需修改。
#### 输出变量修改
输出变量修改 `frontend/packages/workflow/playground/src/node-registries/json-stringify/constants.ts` 里边的常量即可。
```JavaScript
import { nanoid } from 'nanoid';
import { ViewVariableType } from '@coze-workflow/variable';

// 入参路径，试运行等功能依赖该路径提取参数
export const INPUT_PATH = 'inputs.inputParameters';

// 定义固定出参
export const OUTPUTS = [
  {
    key: nanoid(),
    name: 'output',
    type: ViewVariableType.String,
  },
];

export const DEFAULT_INPUTS = [{ name: 'input' }];
```

此时我们删掉原来的 JSON 序列化节点，新建一个序列化节点测试下：
![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/62bdbeb293e144f7b4ef3399d1416bd3~tplv-goo7wpa0wc-image.image)
可以看到基本实现了相关能力。
### 舞台节点
![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/171f349aa0eb43d7a8907bcc7a532489~tplv-goo7wpa0wc-image.image)

节点为在画布上展示的节点卡片，实现在 `node-content.tsx` 中。
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

需要注意的是，需要将该组件在当前目录 `index.ts` 导出，然后另一个组件引用，这一切脚手架都自动完成了，由于 JSON 序列化的舞台节点也只展示输入输出，这里也无需修改。
### 节点注册
需要将生成的节点 registry（这里是 JSON_STRINGIFY_NODE_REGISTRY）注册到节点列表中去，这个脚手架也自动完成好了，具体参见：`packages/workflow/playground/src/nodes-v2/constants.ts`。
### 试运行
#### 单节点试运行
单节点试运行需要定义表单提取逻辑，这个在文件 node-test.ts 中定义，如果设置成 true，会使用默认的试运行表单提取逻辑。
```TypeScript
import type { NodeTestMeta } from '@/test-run-kit';

const test: NodeTestMeta = true;

export { test, type NodeTestMeta };
```

我们配置变量为一个 person 对象，试着运行一下：
![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/dade1349b27c4a11bf47eb3aafd88e13~tplv-goo7wpa0wc-image.image)
发现是存在问题的，单节点试运行应该需要提取 person 变量作为一个输入，我们来改造下。
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

上边逻辑主要是提取入参参数，然后针对入参参数调用 **generateParametersToProperties** 生成试运行表单。
再运行试一下：
![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/b8f327551c854e62a0f31b12ec12c6f6~tplv-goo7wpa0wc-image.image)
可以看到没有问题了。
####  全流程试运行
全流程试运行由于是针对开始节点来运行，因此与当前节点没有关系。
### 修改清单
总结一下，在脚手架的基础上需要修改如下文件：
| **修改文件** | **作用** |
| --- | --- |
| constants.ts | 输出变量类型定义 |
| form.tsx | 替换新版本输入组件 |
| node-test.ts | 单节点试运行调整 |
| components/inputs/index.tsx | 新增特化的输入组件 |

最终变更文件，可以参考这个 MR：https://github.com/kozex-ai/kozex/pull/215
## 节点联调
完成上述操作后，可以和后端联调了，按需修改相关功能即可。


