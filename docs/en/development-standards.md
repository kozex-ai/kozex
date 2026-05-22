# Project architecture
Kozex is designed based on the principles of Domain-Driven Design (DDD). The overall project structure is as follows:
```Python
├── backend/              # Backend services
│   ├── api/              # API processors and routes
│   ├── application/      # Application layer, combining domain objects and infrastructure implementations
│   ├── conf/             # Configuration files
│   ├── crossdomain/      # Cross-domain anti-corruption layer
│   ├── domain/           # Domain layer, containing core business logic
│   ├── infra/            # Infrastructure implementation layer
│   ├── pkg/              # Utility methods with no external dependencies
│   └── types/            # Type definitions
├── common/               # Common components
├── docker/               # Docker configurations
├── frontend/             # Frontend application
│   ├── apps/             # Applications
│   ├── config/           # Configuration files
│   ├── infra/            # Infrastructure
│   └── packages/         # Packages
├── idl/                  # Interface Definition Language files
```


* **API layer (/backend/api):** Implements HTTP endpoints, uses the Hertz server, handles request/response processing, and includes middleware components
* **App layer (/backend/application):** Combines various domain objects and infrastructure implementations, provides API services, and encapsulates APIs
* **Domain layer (/backend/domain):** Contains core business logic, defines domain entities and value objects, and implements business rules and workflows
* **Cross-domain anti-corruption layer (/backend/crossdomain):** Defines cross-domain APIs, implements cross-domain APIs, and prevents direct dependencies between domains
* **Infrastructure contract layer (/backend/infra/contract):** Defines APIs for all external dependencies, serves as the boundary between domain logic and infrastructure, and includes APIs for storage systems, caching mechanisms, message queues, configuration management, etc.
* **Infrastructure implementation layer (/backend/infra/impl):** Implements APIs defined by the infrastructure contract layer, providing specific implementations such as database operations, caching mechanisms, message queue operations, etc.
* **Toolkit (/backend/pkg):** Utility methods with no external dependencies, which can be directly used by any layer
* **Docker**: Docker configuration
* **Frontend app (frontend):** Frontend app implementation
* **idl**: Backend API definition

# Code development and testing

In Kozex, both frontend and backend code must adhere to the coding styles and standards provided in this document.
## **Code structure**
```Python
├── backend/              # Backend services
│   ├── api/              # API processors and routes
│   │   ├── handler/      # Processors
│   │   ├── internal/     # Internal tools
│   │   ├── middleware/   # Middleware components
│   │   ├── model/        # API model definitions
│   │   └── router/       # Route definitions
│   ├── application/      # Application layer, combining domain objects and infrastructure implementations
│   │   ├── app/          # Application services
│   │   ├── conversation/ # Conversation application services
│   │   ├── knowledge/    # Knowledge application services
│   │   ├── memory/       # Memory application services
│   │   ├── modelmgr/     # Model management application services
│   │   ├── plugin/       # Plugin application services
│   │   ├── prompt/       # Prompt application services
│   │   ├── search/       # Search application services
│   │   ├── singleagent/  # Single agent application services
│   │   ├── user/         # User application services
│   │   └── workflow/     # Workflow application services
│   ├── conf/             # Configuration files
│   │   ├── model/        # Model configurations
│   │   ├── plugin/       # Plugin configurations
│   │   └── prompt/       # Prompt configurations
│   ├── crossdomain/      # Cross-domain anti-corruption layer
│   │   ├── contract/     # Cross-domain interface definitions
│   │   ├── impl/         # Cross-domain interface implementations
│   │   └── workflow/     # Workflow cross-domain implementations
│   ├── domain/           # Domain layer, containing core business logic
│   │   ├── agent/        # Agent domain logic
│   │   ├── app/          # Application domain logic
│   │   ├── conversation/ # Conversation domain logic
│   │   ├── knowledge/    # Knowledge domain logic
│   │   ├── memory/       # Memory domain logic
│   │   ├── modelmgr/     # Model management domain logic
│   │   ├── plugin/       # Plugin domain logic
│   │   ├── prompt/       # Prompt domain logic
│   │   ├── search/       # Search domain logic
│   │   ├── user/         # User domain logic
│   │   └── workflow/     # Workflow domain logic
│   ├── infra/            # Infrastructure implementation layer
│   │   ├── contract/     # Infrastructure interface definitions
│   │   └── impl/         # Infrastructure interface implementations
│   ├── pkg/              # Utility methods with no external dependencies
│   │   ├── ctxcache/     # Context cache tools
│   │   ├── errorx/       # Error handling tools
│   │   ├── goutil/       # Go language tools
│   │   ├── logs/         # Logging tools
│   │   └── safego/       # Safe Go tools
│   └── types/            # Type definitions
│       ├── consts/       # Constant definitions
│       ├── ddl/          # Data definition language
│       └── errno/        # Error code definitions
├── frontend/             # Frontend application
│   ├── apps/             # Applications
│   │   └── kozex/        # kozex application
│   ├── config/           # Configuration files
│   │   ├── eslint-config/ # ESLint configuration
│   │   ├── postcss-config/ # PostCSS configuration
│   │   ├── rsbuild-config/ # RSBuild configuration
│   │   ├── stylelint-config/ # StyleLint configuration
│   │   ├── tailwind-config/ # Tailwind configuration
│   │   ├── ts-config/    # TypeScript configuration
│   │   └── vitest-config/ # Vitest configuration
│   ├── infra/            # Infrastructure
│   │   ├── eslint-plugin/ # ESLint plugin
│   │   ├── idl/          # Interface definition language
│   │   ├── plugins/      # Plugins
│   │   └── utils/        # Utilities
│   └── packages/         # Packages
│       ├── agent-ide/    # Agent IDE
│       ├── arch/         # Architecture
│       ├── common/       # Common components
│       ├── components/   # UI components
│       ├── data/         # Data
│       ├── foundation/   # Foundation
│       ├── studio/       # Studio
│       └── workflow/     # Workflow
├── idl/                  # Interface definition language files
```

## Basic components

* **Backend framework**:
   * Hertz (Cloudwego high-performance HTTP framework)
* **Data storage**:
   * MySQL: Structured data storage
   * Redis: Cache and temporary data
   * MinIO: Object storage
   * Imagex: Upload various types of materials such as images and documents
   * RocketMQ: Message queue
* **Containerization**: Docker and Docker Compose

## Development and Debugging

When you have made changes to Kozex and saved them, you can start a local server for development and debugging by following these commands. Make sure to execute them from the project root directory. For more commands, please refer to the `kozex/Makefile` file.

**Prerequisites**:
docker, go (≥1.24), make, node (>=22), npm, rush

**Debugging Steps**:
1.  When starting debugging locally for the first time, you need to execute the following command:

    ```Bash
    # Make sure Docker is running
    make debug
    ```

2.  After modifying and saving the code, you can recompile and restart the coze-server service by executing the following command:

    ```Bash
    make server
    ```

3.  In addition to the above, if you also need to debug the front-end code, open a new Terminal and start the front-end project:

    ```Plain Text
    bash scripts/setup_fe.sh
    cd frontend/apps/kozex
    npm run dev
    ```

After that, you can access the address output in the command line, as shown in the figure below `http://localhost:8080`:
![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/7fac0074d1d9412494497afcaf31b464~tplv-goo7wpa0wc-image.image)

## Backend specifications
### **Go specifications**
Code specifications for the Go language can refer to the [Google specifications](https://google.github.io/styleguide/go/best-practices.html). It is recommended to use formatting tools like gofmt for code formatting.
### **IDL specifications**
| **Category** | **Note** |
| --- | --- |
| Service definition | * The service name follows camelCase naming convention. <br> * Each Thrift file only defines one Service, except for extends aggregation. |
| Method definition | * API naming adopts camelCase naming convention. <br> * The API can only have one parameter and one return value, both of which must be custom Struct types. <br> * The input parameter must be named {Method}Request, and the return value must be named {Method}Response <br> * Each Request type must include a Base field, with the type base.Base, field number 255, and the optional type <br> * Each Response type must include a BaseResp field, with the type base.BaseResp, field number 255 |
| Struct definition | * Struct names should use camel case naming <br> * Field names should use snake case naming <br> * New fields should be set as optional, required is prohibited <br> * Existing field IDs and types must not be modified |
| Enumeration definition | * It is recommended to use typedef to define enumeration values <br> * Enum values are named using camel case, with the type and name connected by underscores |
| API definition | * Define API using Restful style <br> * Refer to existing module's API definitions and keep the style consistent |
| Annotation definition | * You can refer to [Hertz](https://www.cloudwego.cn/zh/docs/hertz/tutorials/toolkit/annotation/#%E6%94%AF%E6%8C%81%E7%9A%84-api-%E6%B3%A8%E8%A7%A3)-supported annotations |
### **Unit testing specifications**
| **Category** | **Specification description** |
| --- | --- |
| UT function naming <br>  | * Ordinary functions should be named as Test{FunctionName}(t *testing.T) <br> * The object method is named Test{ObjectName}_{MethodName}(t *testing.T) <br> * The benchmark test function is named Benchmark{FunctionName}(b *testing.B) <br> * The benchmark test object is named Benchmark{ObjectName}_{MethodName}(b *testing.B) |
| File naming | Test files share the same name as the files being tested, with the suffix _test.go, and are located in the same directory |
| Test design | * It is recommended to use the Table-Driven approach to define inputs/outputs, covering various scenarios <br> * Use github.com/stretchr/testify to simplify assertion logic <br> * Use github.com/uber-go/mock to generate Mock objects, and try to avoid using the patch stubbing method |
## Frontend development specifications
### **Language specification**
The front-end part of Kozex is implemented using frameworks such as Typescript, React, and Tailwind. The repository has a built-in comprehensive code specification detection system, which performs automatic checks when committing code and merging PRs. Developers only need to ensure there are no Lint errors when committing.
### **Unit test specification**
The repository adopts Vitest as the main front-end testing framework. Test files must follow the following specifications:

* Set up the `__tests__` folder at the package root directory, on the same level as `src`.
* All test cases in the package are saved in the folder created in the previous step.
* For each source module `foo.ts` in the `src` directory, create a corresponding test module with the same name `foo.test.ts`. The directory structure of the test code should remain consistent with the source code for easier mapping.
* All unit test filenames should end with `test.ts`.

The final structure is as follows:

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

In the above example:

* `src/root.tsx` related unit test code is concentrated in `__tests__/root.test.tsx`.
* `src/config-panel/foo.tsx` is concentrated in `__tests__/config-panel/foo.test.tsx`, and the unit test directory structure is consistent with the source code directory structure.

Additionally, source files should correspond one-to-one with test files. If some test files have too many lines of code, **instead, you should prioritize determining whether the corresponding source module logic is overly complex and whether further modularization is needed**.
Secondly, the essence of unit testing is **"simulating and triggering code logic in a controlled environment to verify the execution results"**. A standard unit test case typically includes the following elements:

* `Arrange`: Use APIs such as `vi.mock` to simulate the context state and construct a "controlled" testing environment.
* `Act`: Invoke the test target code to trigger execution effects.
* `Assert`: Check and verify whether the response effects of the act step meet expectations. Note that unit tests must include sufficiently complete `assert`; otherwise, the goal of validation cannot be achieved.

It is recommended that all subsequent UT code should organize using the AAA (Arrange-Act-Assert) structure and follow the structural requirements below:

* Except for special statements like `vi.importActual`, all import statements should be placed at the beginning of the file.
* After the `import` statements, global `vi.mock` calls should be placed. In principle, all downstream modules should be `mocked`.
   ![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/41057f6edf4b4bbca444a3633cb61708~tplv-goo7wpa0wc-image.image)
* Place the `describe` test suite function after the mock statement. The function should not generally contain multiple nested `describe` functions and should include multiple `it` test cases within it.
* Within the `it` test case, follow the sequence of `arrange => act => assert`, for example:
   ![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/88021cf7063b4a8a98ac8745791ab686~tplv-goo7wpa0wc-image.image)

Complete example:
```JavaScript
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { exec } from 'shelljs';

import { ensureNotUncommittedChanges } from '@/utils/git';
// Import the mocked module so that we can access the mock functions

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

### Dependency management
Please try to follow the existing dependency definitions and avoid introducing duplicate dependencies, for example:

* The current repository mainly uses Rspack and Rollup as build tools for Web apps and Node Package scenarios, respectively. Avoid introducing tools like Webpack.
* Use Vitest as the testing framework, and do not introduce Jest, Mocha, etc.
* Use `@coze/coze-design` as the UI framework, and do not introduce other frameworks.
* The same principle applies to other dependencies.

## Testing process
### **Function test**
After local deployment, the platform can be opened to test whether the functions of each module are normal.

# Troubleshooting
### Check container status
When the project starts, the following containers will run:
```Bash
coze-server: kozex backend service
coze-mysql: kozex MySQL database
coze-redis: kozex Redis cache
coze-milvus: kozex Milvus vector database
coze-nsqlookupd: nsq lookupd
coze-nsqadmin: nsqadmin
coze-nsqd: nsqd
coze-elasticsearch: kozex Elasticsearch
coze-minio: kozex MinIO object storage
coze-etcd: kozex etcd
```

You can refer to the following steps to troubleshoot container status.

1. Check whether all containers have started and are in a healthy state.
   ```Bash
   docker ps -a
   ```

2. If any service is in an unhealthy state, you can check the logs of the corresponding component and identify the cause of the error.
   ```Bash
   docker logs coze-server #Backend services
   docker logs coze-mysql #MySQL database
   docker logs coze-redis #Redis cache
   docker logs coze-milvus #Vector database
   docker logs coze-nsqlookupd #nsq lookupd
   docker logs coze-nsqadmin #nsqadmin
   docker logs coze-nsqd #nsqd
   docker logs coze-elasticsearch #es
   docker logs coze-minio #Object Storage
   docker logs coze-etcd #etcd
   ```

3. For basic components, you can refer to the instructions in the docker-compose.yaml file to access the basic components and check the data.
   For example:
   ```Bash
   #mysql, query data from the user table
   docker exec -it coze-mysql mysql -u root -proot
   SHOW DATABASES;
   USE opencoze;
   SHOW TABLES;
   SELECT * FROM user LIMIT 10;
   ```


### Service logs
After all containers have started normally, the main focus is the running logs of the backend services. If an API error or similar issue occurs, we can check the logs of the corresponding service container to see if there are any error logs:
```Bash
docker logs coze-server
```

If the API reports an error on the page, you can press F12 to access the browser console, check the request corresponding to the error, obtain the LogID from the header located at x-log-id, and then use the logid within the container to view logs. If the API is functioning correctly, you may not be able to find the corresponding logs.

![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/653fe58b6aa74ac59e2f8e57feb9c63f~tplv-goo7wpa0wc-image.image)

View logs within the container:
```Bash
docker logs coze-server | grep {logid} 
```

### **Error codes**
The server error definitions are located under the `backend/types/errno` directory. If the request returns an error code, you can search for the error code in the project to quickly identify the issue, then use the aforementioned log querying method to pinpoint the specific problem.
