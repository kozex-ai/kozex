## 1、windows本地部署，"Error response from daemon: Ports are not available: exposing port TCP [0.0.0.0:2379](http://0.0.0.0:2379/) -＞ [127.0.0.1:0](http://127.0.0.1:0/):2379" 如何解决？

```bash
# 查看占用
netstat -ano | findstr :2379

net stop winnat
net start winnat
```


## 2、Agent 对话调试 "Something error:Internal server error" 错误？

可以通过以下命令查询具体错误日志

```bash
docker logs coze-server | grep -i 'node execute failed'
```

![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/d1d490a6cac74051a222f51a110fad08~tplv-goo7wpa0wc-image.image)

## 3、Agent对话或者workflow大模型节点上传图片/文件后调试时模型报错？

大模型访问的图片/文件链接需要是公网可访问的链接，图片存组件需要在公网部署，具体可以参考：[上传组件配置](basic-component-configuration.md#上传组件)

## 4、工作流代码节点如何添加 Python 第三方库
在 `Kozex` 项目中，代码节点默认内置了两个第三方依赖库：`httpx` 和 `numpy`。Kozex 也支持开发者自行添加其他的 `Python` 第三方库。详细操作步骤如下：

1. 修改配置文件。
   在`./scripts/setup/[python.sh](python.sh)`脚本与`./backend/Dockerfile` 文件中，可找到 `third - party libraries` 注释，在这两个文件的第三方库注释下方直接添加依赖库对应的 `pip install` 命令即可。
   例如添加 2.0.0 版本的 `torch`：
   ```Bash
   # If you want to use other third - party libraries, you can install them here.
   pip install torch==2.0.0
   ```

2. 在 `./backend/conf/workflow/config.yaml` 中添加第三方模块的包名称。例如添加 torch：
   ```yaml
   NodeOfCodeConfig:
    SupportThirdPartModules:
        - httpx
        - numpy
        - torch

   ```

3. 在 `./docker/docker-compose.yml` 中修改 coze-server 命令。
   ```YAML
   coze-server: 
       build:  # 将 build 指令注释放开
         context: ../
         dockerfile: backend/Dockerfile
   ```

4. 执行以下命令重启并编译 coze-sever 服务。
   ```Bash
   docker compose --profile "*" up -d --build coze-server
   ```
