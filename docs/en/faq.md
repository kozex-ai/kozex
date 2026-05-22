## 1. "Error response from daemon: Ports are not available: exposing port TCP http://0.0.0.0:2379 -＞ http://127.0.0.1:0" when deploying locally on Windows?
```Plain Text
# Check port usage
netstat -ano | findstr :2379

net stop winnat
net start winnat
```


## 2. "Something error: Internal server error" during Agent conversation debugging?
You can query the specific error logs with the following command:
```Plain Text
docker logs coze-server | grep -i 'node execute failed'
```

![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/d1d490a6cac74051a222f51a110fad08~tplv-goo7wpa0wc-image.image)


## 3. Why does the model report an error when debugging after uploading an image/file in an Agent conversation or a workflow's large model node?

The image/file links accessed by the large model must be publicly accessible URLs. The image storage component needs to be deployed on the public network. For details, refer to: [Upload Component Configuration](basic-component-configuration.md#upload-components)

## 4. How to add Python third-party libraries to a workflow code node
In the `Kozex` project, the code node comes with two third-party dependency libraries by default: `httpx` and `numpy`. Kozex also allows developers to add other third-party `Python` libraries on their own. The detailed steps are as follows:

1. Modify configuration files.
   In the `./scripts/setup/[python.sh](python.sh)` script and the `./backend/Dockerfile` file, you can find the `third-party libraries` comment. Simply add the corresponding `pip install` command for the dependency directly below the third-party libraries comment in both files.
   For example, add version 2.0.0 of `torch`:
   ```Bash
   # If you want to use other third - party libraries, you can install them here.
   pip install torch==2.0.0
   ```

2. Add the package names of third-party modules in `./backend/conf/workflow/config.yaml`. For example, to add `torch`:
   ```Go
   NodeOfCodeConfig:
       SupportThirdPartModules:
           - httpx
           - numpy
           - torch
   ```

3. Modify the coze-server command in `./docker/docker-compose.yml`.
   ```YAML
   coze-server: 
       build:  # Uncomment the build instruction.
         context: ../
         dockerfile: backend/Dockerfile
   ```

4. Execute the following commands to restart and compile the coze-server service.
   ```Bash
   docker compose --profile "*" up -d --build coze-server
   ```
