如果你基于 Kozex 进行二次开发，需要新增 API 接口，应编写 IDL 定义和对应服务端代码。
## **步骤一：编写 IDL 接口定义**
### 1 定义 IDL 接口
在 idl/ 目录下定义你的 API 接口。可以创建新的 thrift 文件或在现有 thrift 文件中添加新接口的相关定义。
例如，在 `idl/example/example_service.thrift` 中新增以下接口定义：
```Thrift
namespace go example
include "../base.thrift"


// 请求结构体
struct CreateUserRequest {
    1: required string name
    2: optional string email
    3: optional i32 age
    255: optional base.Base Base (api.none="true")
}

struct CreateUserResponse {
    1: CreateUserData data
    253: required i64 code
    254: required string msg
    255: optional base.BaseResp BaseResp (api.none="true")
}

struct CreateUserData {
    1: i64 user_id (agw.js_conv="str", api.js_conv="true")
    2: string name
}

// 服务定义
service ExampleService {
    CreateUserResponse CreateUser(1: CreateUserRequest request)(
        api.post='/api/example/user/create', 
        api.category="user"
    )
}
```

### 2 更新主 IDL 文件
如果是新增模块，在 `idl/api.thrift` 中添加新服务的引用；如果是在现有模块中添加接口，则不需要修改 `api.thrift`。
```Thrift
include "./example/example_service.thrift"

// 在service列表中添加
service ExampleService extends example_service.ExampleService {}
```

## 步骤二：生成并实现服务端代码

1. 进入 backend 目录，执行以下命令，生成服务端代码。

```Bash
cd backend
hz update -idl ../idl/api.thrift -enable_extends
```


2. 编写业务逻辑代码。
   生成的 handler 文件位于 `backend/api/handler/coze/example_service.go` 中，你需要自行实现具体的业务逻辑。
   ```Go
   // CreateUser .
   // @router /api/example/user/create [POST]
   func CreateUser(ctx context.Context, c *app.RequestContext) {
       var err error
       var req example.CreateUserRequest
       err = c.BindAndValidate(&req)
       if err != nil {
           invalidParamRequestResponse(c, err.Error())
           return
       }
       // 调用业务逻辑层
       resp, err := userService.CreateUser(ctx, &req)
       if err != nil {
           internalServerErrorResponse(ctx, c, err)
           return
       }
       c.JSON(consts.StatusOK, resp)
   }
   ```


## 步骤三：生成前端 Client 或配置鉴权
如果需要对应的前端页面来调用新的 API 接口，则应生成前端 Client 并编写前端代码；对于新增的 OpenAPI 接口，如果无需开发对应的前端页面，则需要开启 OpenAPI 鉴权逻辑。
### 生成前端 Client

1. 执行以下命令，通过服务端定义的 idl 生成前端 Client 代码。
   ```Bash
   cd frontend/packages/arch/api-schema 
   npm run update
   ```

2. 编写前端代码。
   编写前端的页面展示逻辑，并调用 Client 发起请求。
   完成以上操作后，可以开始前后端的联调。建议测试通过后再正式部署到线上使用。

### 开启 OpenAPI 鉴权

1. 打开文件 `backend/api/middleware/openapi_auth.go`。
2. 在 needAuthPath 函数中新增一行代码，将需要鉴权的 path 添加到列表中。
   ![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/4dceecdda0ee4b61a5ae124455a7e2ea~tplv-goo7wpa0wc-image.image)


