After deploying the open-source version of Kozex, if you need to use the image upload functionality or knowledge base-related functionality, you should refer to this document to configure the required basic components for these features. These components typically rely on third-party services such as Volcengine. When configuring components, you need to provide authentication settings for the third-party service, such as the account, token, and other authentication configurations.
## Upload components
In multimodal chat scenarios, it is often necessary to upload images, files, and other information during the conversation. For example, sending an image message in the agent debugging area:

![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/6d18160e25a54359b3db9ceee90ec25a~tplv-goo7wpa0wc-image.image)

This functionality is provided by the upload component. Kozex currently supports the following three types of upload components. You can **choose any one** to use as the upload component.

* **(default) minio**: Images and files are automatically uploaded to the minio service on the local host and can be accessed via a specified port. However, the local host must be configured with a public domain name; otherwise, uploaded images and files will only have internal network access links and cannot be read or recognized by large models.
* **Volcengine TOS**: Images and files are automatically uploaded to the specified Volcengine TOS, and a publicly accessible URL is generated. If you select TOS, you must first enable TOS and configure the Volcengine token in kozex.
* **Volcengine ImageX**: Images and files are automatically uploaded to the specified Volcengine ImageX, and a publicly accessible URL is generated. If you choose ImageX, you must first enable ImageX and configure the Volcengine token in kozex.

The configuration methods for the upload component are as follows:

1. Set the upload component type.
   Open the `.env` file in the docker directory. The value of the configuration item "FILE_UPLOAD_COMPONENT_TYPE" specifies the type of upload component.
   * **storage** (default): Uses the storage component configured by `STORAGE_TYPE`. `STORAGE_TYPE` defaults to minio, but can also be configured as tos.
   * **imagex**: Indicates the use of the Volcengine ImageX component.
   ```Bash
   # This Upload component used in Agent / workflow File/Image With LLM  , support the component of imagex / storage
   # default: storage, use the settings of storage component
   # if imagex, you must finish the configuration of <VolcEngine ImageX>
   export FILE_UPLOAD_COMPONENT_TYPE="storage"
   ```

2. Add a secret key and other configurations for the upload component.
   Likewise, in the `.env` file in the docker directory, enter the following configuration based on the component type.
   | **Component types** | **Configure the authentication method** | **Example** |
   | --- | --- | --- |
   | minio <br> (Default) | 1. In the `docker/.env` file of the Kozex project, set FILE_UPLOAD_COMPONENT_TYPE to storage. <br> 2. In the Storage component area, set STORAGE_TYPE to minio. <br> 3. In the MiniO area, simply keep the default configuration. <br>  <br> If you choose to deploy Kozex on public clouds such as Volcengine, you must configure your cloud server in advance to allow access to ports 8888 and 8889. <br>  | ![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/64578ddbcbb94dafa875cf266e49a0fc~tplv-goo7wpa0wc-image.image) <br>  |
   | tos | 1. Enable the Volcengine TOS product. <br> 2. In the `docker/.env` file of the Kozex project, FILE_UPLOAD_COMPONENT_TYPE is set to storage. <br> 3. In the Storage component area, set STORAGE_TYPE to tos. <br> 4. In the TOS area, enter the following parameters: <br>    * TOS_ACCESS_KEY: Volcengine access key AK. For information on how to obtain it, refer to [Obtain a Volcengine API token](https://www.volcengine.com/docs/6257/64983). <br>    * TOS_SECRET_KEY: Volcengine secret key SK. To learn how to obtain it, refer to [Obtain Volcengine API token](https://www.volcengine.com/docs/6257/64983). <br>    * TOS_ENDPOINT: The Endpoint of the TOS service. For information on how to obtain it, see [Regions and access domain names](https://www.volcengine.com/docs/6349/107356). <br>    * TOS_REGION: The region where the TOS service is located. For information on how to obtain the region value, see [Regions and access endpoints](https://www.volcengine.com/docs/6349/107356). |  <br> ![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/9230b5a67296447ca807ea97ae6fdb62~tplv-goo7wpa0wc-image.image) <br>  |
   | imagex | 1. Activate the Volcengine veImageX product and create the **Asset Hosting** service. You can refer to the [Volcengine veImageX official documentation](https://www.volcengine.com/docs/508/8084). Please note that when creating the **Asset Hosting** service, you need to provide a domain name. We recommend obtaining a publicly available domain name in advance. <br> 2. In the `docker/.env` file of the Kozex project, FILE_UPLOAD_COMPONENT_TYPE is set to imagex. <br> 3. In the VolcEngine ImageX area, fill in the following parameters: <br>    * VE_IMAGEX_AK: Volcengine Access Key (AK). For information on how to obtain it, refer to [Obtain a Volcengine API token](https://www.volcengine.com/docs/6257/64983). <br>    * VE_IMAGEX_SK: Volcengine secret key (SK). Refer to [Obtain a Volcengine API token](https://www.volcengine.com/docs/6257/64983) for how to obtain it. <br>    * VE_IMAGEX_SERVER_ID: The service ID displayed on the **Service Management** page of the Volcengine veImageX product. <br>    * VE_IMAGEX_DOMAIN: The domain name specified when creating the service. <br>    * VE_IMAGEX_TEMPLATE: The name of the template for image processing configuration. |  <br> ![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/30ef93d304d542159141c4a00777c363~tplv-goo7wpa0wc-image.image) <br>  |
3. Run the following command to restart the service so that the above configuration takes effect.
   ```Shell
   docker compose --profile '*' up -d --force-recreate --no-deps coze-server
   ```

For other component configurations, please visit http://localhost:8888/admin/ to configure.
