To perform secondary development based on Kozex and add new API endpoints, you should write IDL definitions and the corresponding server-side code.
## **Step 1: Write the IDL API definition**
### Step 1: Define IDL API
Define your API in the `idl/` directory. You can create new thrift files or add relevant definitions for new APIs to existing thrift files.
For example, add the following API definitions in `idl/example/example_service.thrift`:
```Thrift
namespace go example
include "../base.thrift"


// Request structure
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

// Service definition
service ExampleService {
    CreateUserResponse CreateUser(1: CreateUserRequest request)(
        api.post='/api/example/user/create', 
        api.category="user"
    )
}
```

### Step 2: Update main IDL file
If it is a new module, add a reference to the new service in `idl/api.thrift`; if you are adding an API to an existing module, no modification to `api.thrift` is required.
```Thrift
include "./example/example_service.thrift"

// Add to the service list
service ExampleService extends example_service.ExampleService {}
```

## Step 2: Generate and implement server code

1. Go to the backend directory and run the following command to generate the server code.

```Bash
cd backend
hz update -idl ../idl/api.thrift -enable_extends
```


2. Write business logic code.
   The generated handler file is located in `backend/api/handler/coze/example_service.go`. You need to implement the specific business logic.
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
       // Call the business logic layer
       resp, err := userService.CreateUser(ctx, &req)
       if err != nil {
           internalServerErrorResponse(ctx, c, err)
           return
       }
       c.JSON(consts.StatusOK, resp)
   }
   ```


## Step 3: Generate a frontend client or configure authentication
If a corresponding frontend page is needed to call the new API, generate the frontend client and write the frontend code. For newly added OpenAPI APIs, if there is no need to develop a corresponding frontend page, OpenAPI authentication logic must be enabled.
### Generate a frontend client

1. Run the following command to generate frontend client code using the IDL defined by the server.
   ```Bash
   cd frontend/packages/arch/api-schema 
   npm run update
   ```

2. Write front-end code.
   Write the display logic for the frontend page(s), and invoke the Client to initiate requests.
   After completing the above operations, you can begin front-end and back-end joint debugging. It is recommended to formally deploy to the production environment only after testing has been successfully completed.

### Enable OpenAPI authentication

1. Open the file `backend/api/middleware/openapi_auth.go`.
2. Add a line of code in the needAuthPath function to add the path that requires authentication to the list.
   ![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/4dceecdda0ee4b61a5ae124455a7e2ea~tplv-goo7wpa0wc-image.image)
