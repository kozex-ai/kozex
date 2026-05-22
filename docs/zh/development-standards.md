# 项目架构
Kozex 基于领域驱动设计（DDD）原则架构设计，项目整体结构如下：
```Python
├── backend/              # 后端服务
│   ├── api/              # API 处理器和路由
│   ├── application/      # 应用层，组合领域对象和基础设施实现
│   ├── conf/             # 配置文件
│   ├── crossdomain/      # 跨领域防腐层
│   ├── domain/           # 领域层，包含核心业务逻辑
│   ├── infra/            # 基础设施实现层
│   ├── pkg/              # 无外部依赖的工具方法
│   └── types/            # 类型定义
├── common/               # 公共组件
├── docker/               # Docker 配置
├── frontend/             # 前端应用
│   ├── apps/             # 应用程序
│   ├── config/           # 配置文件
│   ├── infra/            # 基础设施
│   └── packages/         # 包
├── idl/                  # 接口定义语言文件
```


* **API 层 ( /backend/api )：​**实现 HTTP 端点，使用 Hertz 服务器、处理请求/响应处理、包含中间件组件
* **应用层 ( /backend/application )：​**组合各种领域对象和基础设施实现、提供 API 服务、封装 API 
* **领域层 ( /backend/domain )：​**包含核心业务逻辑、定义领域实体和值对象、实现业务规则和工作流
* **跨领域防腐层 ( /backend/crossdomain )：​**定义跨领域接口、实现跨领域接口、防止领域间直接依赖
* **基础设施契约层 ( /backend/infra/contract )：​**定义所有外部依赖的接口、作为领域逻辑和基础设施之间的边界、包括存储系统、缓存机制、消息队列、配置管理等接口
* **基础设施实现层 ( /backend/infra/impl )：​**实现基础设施契约层定义的接口、提供具体实现，如数据库操作、缓存机制、消息队列操作等
* **工具包 ( /backend/pkg )：​**无外部依赖的工具方法、可以被任何层直接使用
* **docker**：Docker 配置
* **前端应用(frontend)：​**前端应用实现
* **idl**:后端接口定义


# 代码开发与测试
在 Kozex 中，前后端代码均需要遵守本文档提供的编码风格与规范。
## **代码结构**
```Python
├── backend/              # 后端服务
│   ├── api/              # API 处理器和路由
│   │   ├── handler/      # 处理器
│   │   ├── internal/     # 内部工具
│   │   ├── middleware/   # 中间件组件
│   │   ├── model/        # API 模型定义
│   │   └── router/       # 路由定义
│   ├── application/      # 应用层，组合领域对象和基础设施实现
│   │   ├── app/          # 应用服务
│   │   ├── conversation/ # 会话应用服务
│   │   ├── knowledge/    # 知识应用服务
│   │   ├── memory/       # 内存应用服务
│   │   ├── modelmgr/     # 模型管理应用服务
│   │   ├── plugin/       # 插件应用服务
│   │   ├── prompt/       # 提示词应用服务
│   │   ├── search/       # 搜索应用服务
│   │   ├── singleagent/  # 单一代理应用服务
│   │   ├── user/         # 用户应用服务
│   │   └── workflow/     # 工作流应用服务
│   ├── conf/             # 配置文件
│   │   ├── model/        # 模型配置
│   │   ├── plugin/       # 插件配置
│   │   └── prompt/       # 提示词配置
│   ├── crossdomain/      # 跨领域防腐层
│   │   ├── contract/     # 跨领域接口定义
│   │   ├── impl/         # 跨领域接口实现
│   │   └── workflow/     # 工作流跨领域实现
│   ├── domain/           # 领域层，包含核心业务逻辑
│   │   ├── agent/        # 代理领域逻辑
│   │   ├── app/          # 应用领域逻辑
│   │   ├── conversation/ # 会话领域逻辑
│   │   ├── knowledge/    # 知识领域逻辑
│   │   ├── memory/       # 内存领域逻辑
│   │   ├── modelmgr/     # 模型管理领域逻辑
│   │   ├── plugin/       # 插件领域逻辑
│   │   ├── prompt/       # 提示词领域逻辑
│   │   ├── search/       # 搜索领域逻辑
│   │   ├── user/         # 用户领域逻辑
│   │   └── workflow/     # 工作流领域逻辑
│   ├── infra/            # 基础设施实现层
│   │   ├── contract/     # 基础设施接口定义
│   │   └── impl/         # 基础设施接口实现
│   ├── pkg/              # 无外部依赖的工具方法
│   │   ├── ctxcache/     # 上下文缓存工具
│   │   ├── errorx/       # 错误处理工具
│   │   ├── goutil/       # Go 语言工具
│   │   ├── logs/         # 日志工具
│   │   └── safego/       # 安全 Go 工具
│   └── types/            # 类型定义
│       ├── consts/       # 常量定义
│       ├── ddl/          # 数据定义语言
│       └── errno/        # 错误码定义
├── frontend/             # 前端应用
│   ├── apps/             # 应用程序
│   │   └── kozex/        # kozex 应用
│   ├── config/           # 配置文件
│   │   ├── eslint-config/# ESLint 配置
│   │   ├── postcss-config/# PostCSS 配置
│   │   ├── rsbuild-config/# RSBuild 配置
│   │   ├── stylelint-config/# StyleLint 配置
│   │   ├── tailwind-config/# Tailwind 配置
│   │   ├── ts-config/    # TypeScript 配置
│   │   └── vitest-config/# Vitest 配置
│   ├── infra/            # 基础设施
│   │   ├── eslint-plugin/# ESLint 插件
│   │   ├── idl/          # 接口定义语言
│   │   ├── plugins/      # 插件
│   │   └── utils/        # 工具
│   └── packages/         # 包
│       ├── agent-ide/    # 代理 IDE
│       ├── arch/         # 架构
│       ├── common/       # 公共组件
│       ├── components/   # UI 组件
│       ├── data/         # 数据
│       ├── foundation/   # 基础
│       ├── studio/       # Studio
│       └── workflow/     # 工作流
├── idl/                  # 接口定义语言文件
```

## 基础组件 

* **后端框架**： 
   * Hertz ( Cloudwego 高性能 HTTP 框架) 
* **数据存储**： 
   * MySQL：结构化数据存储 
   * Redis：缓存和临时数据 
   * MinIO：对象存储 
   * Imagex: 图像、文档等各类素材上传
   * RocketMQ: 消息队列 
* **容器化**：Docker 和 Docker Compose


## 开发调试
当你对 Kozex 进行二次开发保存后，需要本地开发调试时，可以按照以下命令启动本地服务，确保在项目根目录下执行，更多命令可以参考`kozex/Makefile`文件。

**环境要求**：
docker、go（≥1.24）、make、node(>=22) 、npm、rush

**调试步骤**：
1. 首次在本地启动调试时，需要执行以下命令：

   ```Bash
   # 确保 docker 是正常运行的
   make debug
   ```


2. 后面修改代码保存，直接重新编译启动coze-server服务，执行以下命令即可：

   ```Bash
   make server
   ```


3. 除上述外，如果还需要调试前端代码，重新打开Terminal 程序，启动前端工程：

   ```Plain Text
   bash scripts/setup_fe.sh
   cd frontend/apps/kozex
   npm run dev
   ```

之后，访问命令行输出的地址即可，如下图所示的 `http://localhost:8080`：
![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/7fac0074d1d9412494497afcaf31b464~tplv-goo7wpa0wc-image.image)


## 后端  规范
### **Go 规范**
Go 语言的代码规范可参考[Google规范](https://google.github.io/styleguide/go/best-practices.html)。建议使用 gofmt 等格式化工具进行代码格式化。 
### **IDL 规范**
| **类别** | **说明** |
| --- | --- |
| Service 定义  | * 服务命名采用驼峰命名  <br> * 一个Thrift文件只定义一个Service, extends聚合除外  |
| Method 定义  | * 接口命名采用驼峰命名  <br> * 接口只能拥有一个参数和一个返回值，且是自定义Struct类型  <br> * 入参须命名为{Method}Request，返回值命名为{Method}Response  <br> * 每个Request类型须包含Base字段，类型base.Base，字段序号为255，optional类型  <br> * 每个Response类型须包含BaseResp字段，类型base.BaseResp，字段序号为255  |
| Struct 定义  | * 结构体命名采用驼峰命名  <br> * 字段命名采用蛇形命名  <br> * 新增字段设置为optional，禁止required  <br> * 禁止修改现有字段的ID和类型  |
| 枚举定义  | * 推荐使用typedef来定义枚举值  <br> * 枚举值命名采用驼峰命名，类型和名字之间用下划线连接  |
| API定义  | * 使用Restful风格定义API  <br> * 参考现有模块的API定义，风格保持一致  |
| 注解定义  | * 可参考[Hertz](https://www.cloudwego.cn/zh/docs/hertz/tutorials/toolkit/annotation/#%E6%94%AF%E6%8C%81%E7%9A%84-api-%E6%B3%A8%E8%A7%A3)支持的注解  |
### **单测规范**
| **类别** | **规范说明** |
| --- | --- |
| UT 函数命名  <br>  | * 普通函数命名为Test{FunctionName}(t *testing.T)  <br> * 对象方法命名为Test{ObjectName}_{MethodName}(t *testing.T)  <br> * 基准测试函数命名为Benchmark{FunctionName}(b *testing.B)  <br> * 基准测试对象命名为Benchmark{ObjectName}_{MethodName}(b *testing.B)  |
| 文件命名  | 测试文件与被测试文件同名，后缀为_test.go，处于同一目录下  |
| 测试设计  | * 推荐使用 Table-Driven 的方式定义输入/输出，覆盖多种场景  <br> * 使用github.com/stretchr/testify简化断言逻辑  <br> * 使用github.com/uber-go/mock生成Mock对象，尽量避免Patch打桩的方式  |
## 前端开发规范
### **语言规范**
Kozex 前端部分整体使用 Typescript、React、Tailwind 等框架实现，仓库内已内置丰富的代码规范检测体系，在提交代码、合入 PR 时会执行自动检测，开发者只需确保提交时未出现 Lint 报错即可。
### **单测规范**
仓库中选用 Vitest 作为主要前端测试框架，测试文件需遵循如下规范：

* 在 pakcage 根目录下设立 `__tests__` 文件夹，与 `src` 同级。
* Package 内所有的测试用例，都保存在上一步创建的文件夹中。
* 为 `src` 目录每一个源码模块 `foo.ts`，创建对应同名测试模块 `foo.test.ts`，且测试代码的目录结构与源码保持一致，方便对应。
* 所有单测文件名均以 `test.ts` 结束。

最终形成如下结构：

<div style="display: flex;">
<div style="flex-shrink: 0;width: calc((100% - 16px) * 0.5000);">

![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/62fe5102cc0f4214bb2f07592a0b0ef9~tplv-goo7wpa0wc-image.image)


</div>
<div style="flex-shrink: 0;width: calc((100% - 16px) * 0.5000);margin-left: 16px;">

```JavaScript
infra/flags-devtool/
├── README.md
├── OWNERS
├── vitest.config.ts
├── package.json
├── src/
│   ├── index.tsx
│   ├── root.tsx
│   ├── indicator.tsx
│   ├── global.d.ts
│   ├── index.module.less
│   ├── hooks/
│   ├── utils/
│   └── config-panel/
├── stories/
├── setup/
└── __tests__/
    ├── root.test.tsx
    ├── index.test.tsx
    ├── indicator.test.tsx
    ├── hooks/
    ├── utils/
    └── config-panel/
```



</div>
</div>

上述示例中：

* `src/root.tsx` 相关单测代码集中在 `__tests__/root.test.tsx` 中.
* `src/config-panel/foo.tsx` 则集中在 `__tests__/config-panel/foo.test.tsx` 中，单测目录结构与源码目录结构保持一致。

另外，源码文件应与单测文件一一对应，若出现某些测试文件代码行数过多时，请不要拆解出多个单测文件，**而应该优先判断对应源码模块的逻辑是否过于复杂，是否应该做进一步模块拆解**。
其次，单元测试本质上就是"**在可控环境中，模拟触发代码逻辑，验证执行结果**"的过程，一个标准的单测用例通常包含如下要素：

* `arrange`：调用 `vi.mock` 等接口模拟上下文状态，构建"可控"的测试环境；
* `act`：调用测试目标代码，触发执行效果；
* `assert`：检测，验证 act 的响应效果是否符合预期，注意，单测中务必包含足够完整的 `assert`，否则无法达成验证效果的目标。

建议后续 UT 代码均 AAA(Arrange-Act-Assert) 结构组织代码，遵循如下结构要求：

* 除 `vi.importActual` 等特殊语句外，所有 import 语句均保存到文件开头；
* `import` 语句之后，放置全局 `vi.mock` 调用，原则上应 `mock` 掉所有下游模块；
   ![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/41057f6edf4b4bbca444a3633cb61708~tplv-goo7wpa0wc-image.image)
* Mock 语句之后放置 `describe` 测试套件函数，函数内原则上不可嵌套多个 `describe`；函数内应包含多个 `it` 用例；
* `it` 用例内部遵循 `arrange => act => asset` 顺序，例如：
   ![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/88021cf7063b4a8a98ac8745791ab686~tplv-goo7wpa0wc-image.image)

完整示例：
```JavaScript
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { exec } from 'shelljs';

import { ensureNotUncommittedChanges } from '@/utils/git';
// 导入被mock的模块，以便我们可以访问mock函数

import { env } from '@/ai-scripts';

// arrange
// Mock shelljs
vi.mock('shelljs', () => ({
  exec: vi.fn(),
}));

// Mock ../ai-scripts
vi.mock('@/ai-scripts', () => ({
  env: vi.fn(),
}));

describe('git utils', () => {
  it('应该在 BYPASS_UNCOMMITTED_CHECK 为 true 时直接返回 true', async () => {
    // arrange
    // mock
    vi.mocked(env).mockReturnValue('true');

    // act
    const result = await ensureNotUncommittedChanges('/fake/path');

    // assert
    expect(result).toBe(true);
    expect(exec).not.toHaveBeenCalled();
  });
});
```

### 依赖管理
请尽量遵循现有依赖定义，避免引入重复依赖项，例如：

* 目前仓库主要使用 Rspack、Rollup 作为构建工具，分别对应 Web 应用与 Node Package 场景，应避免引入 Webpack 等构建工具。
* 使用 Vitest 作为测试框架，请勿引入 Jest、Mocha 等。
* 使用 `@coze/coze-design` 作为 UI 框架，请勿引入其它框架。
* 其它依赖同理。

## 测试流程 
### **功能测试** 
本地部署后可以打开平台，测试各模块的功能是否正常。

# 故障排查
### 查看容器状态 
项目启动时，会运行如下几个容器： 
```Bash
coze-server: 后端服务
coze-mysql: mysql数据库
coze-redis: redis缓存
coze-milvus: 向量数据库
coze-elasticsearch: es
coze-nsqlookupd: nsq lookupd
coze-nsqadmin: nsqadmin
coze-nsqd: nsqd
coze-minio: 对象存储
coze-etcd: etcd
```

你可以参考以下步骤排查容器状态。 

1. 查看是否所有容器都已经启动并处于 healthy 状态。 
   ```Bash
   docker ps -a
   ```

2. 如果有服务处于 unhealty 状态，可以查看对应组件的日志并定位报错原因。 
   ```Bash
   docker logs coze-server #后端服务 
   docker logs coze-mysql #mysql数据库
   docker logs coze-redis #redis缓存
   docker logs coze-milvus #向量数据库
   docker logs coze-nsqlookupd #nsq lookupd
   docker logs coze-nsqadmin #nsqadmin
   docker logs coze-nsqd #nsqd
   docker logs coze-elasticsearch #es
   docker logs coze-minio #对象存储
   docker logs coze-etcd #etcd
   ```

3. 对于基础组件，可以参考 docker-compose.yaml 文件中的说明进入基础组件查看数据。 
   例如： 
   ```Bash
   # mysql，查询user表的数据
   docker exec -it coze-mysql mysql -u root -proot
   SHOW DATABASES;
   USE opencoze;
   SHOW TABLES;
   SELECT * FROM user LIMIT 10;
   ```


### 服务日志 
所有容器都正常启动之后，主要关注点就是后端服务的运行日志，如果出现接口报错等情况，我们可以查看服务容器对应的日志，查看是否有错误日志： 
```Bash
docker logs coze-server
```

页面如果接口报错，可以通过F12进入浏览器控制台，查看对应报错的请求，从相应头中获取LogID，位于x-log-id，拿到logid后再去容器内查看日志，如果接口正常，可能会搜索不到对应日志。 

![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/653fe58b6aa74ac59e2f8e57feb9c63f~tplv-goo7wpa0wc-image.image)

去容器内查看日志： 
```Bash
docker logs coze-server | grep {logid} 
```

### **错误码**
服务端的错误定义在 `backend/types/errno` 目录下，若请求返回了错误码，可以在项目中搜索错误码来简单定位问题，然后再通过以上日志查询的方式定位具体问题。
