## Requirements
Before installing Kozex, ensure that your software and hardware environment meets the following requirements:
| <span style="background-color: rgb(255, 255, 255)"><span style="color: #1F2328"><strong>Requirements</strong></span></span> | <span style="background-color: rgb(255, 255, 255)"><span style="color: #1F2328"><strong>Description</strong></span></span> |
| --- | --- |
| CPU | 2 Core |
| RAM | 4 GiB |
| Docker | Install Docker and Docker Compose in advance, and start the Docker service. For detailed instructions, refer to the Docker documentation: <br>  <br> * **macOS**: It is recommended to use Docker Desktop for installation. Refer to the [Docker Desktop For Mac](https://docs.docker.com/desktop/setup/install/mac-install/) installation guide. <br> * **Linux**: Refer to the [Docker installation guide](https://docs.docker.com/engine/install/) and the [Docker Compose](https://docs.docker.com/compose/install/) installation guide. <br> * **Windows**: It is recommended to install using Docker Desktop. Refer to the [Docker Desktop For Windows](https://docs.docker.com/desktop/setup/install/windows-install/) installation guide. |

## Install Kozex

### Step 1: Clone the source code
Run the following commands in your local project to clone the latest version of the Kozex source code.

```Bash
Clone code
git clone https://github.com/kozex-ai/kozex.git
```


### Step 2: Deploy and start services

The initial deployment and startup of Kozex require retrieving images and building local images. This process may take some time, so please be patient. If you see the message "Container coze-server Started", it means that the Kozex service has started successfully.

```Bash
cd kozex
# start service
# for macOS or Linux
make web  
# for windows
cp .env.example .env
docker compose -f ./docker/docker-compose.yml up
```


### Step 3: Register

Register an account by visiting `http://localhost:8888/sign`, entering your username and password, and clicking the Register button.

### Step 4: Configure a model

Configure the model at `http://localhost:8888/admin/#model-management` by adding a new model.

### Step 5: Access the service

Visit Kozex at `http://localhost:8888/`.

![Image](https://p9-arcosite.byteimg.com/tos-cn-i-goo7wpa0wc/c05109bc278e4a7b87164c2db2738b41~tplv-goo7wpa0wc-image.image)

## What's next
<span style="background-color: rgb(255, 255, 255)"><span style="color: #1F2328">After successfully deploying Kozex, if you need to use functions such as plugins and knowledge bases, you also need to configure the relevant components at `http://localhost:8888/admin/`.</span></span>


## Security risks in public networks
If Kozex is to be deployed in a public network environment, it is recommended to pay attention to the following security risks:
- It is recommended to turn off the registration function or enable the email whitelist; otherwise, any user can use it via the link.
- It is recommended to enable the default sandbox environment for the workflow code nodes to enhance security. For detailed configuration instructions.
- It is recommended to configure the network for the deployment environment according to business requirements (such as intranet access) to avoid SSRF risks.
- By default, the Kozex server only listens to [localhost](http://localhost/). When deployed in a public network environment, it is recommended to listen to `0.0.0.0` only when necessary or add additional security measures to prevent the service from being directly exposed to the public network.
