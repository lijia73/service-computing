# 容器化技术与容器服务

## 实验环境

- Ubuntu 20.04 

## 实验内容

- 准备Docker环境
- 构建Docker镜像
- MySQL与容器化
- Docker Compose
- Docker 网络
	- 容器网络管理
	- 自定义网络
- Docker 仓库（Registry）
	- 搭建私有容器仓库
	- 阿里云容器镜像服务实践

## 准备知识
### Docker的核心概念
- `镜像`(image) – 类比`执行程序`  
一个`可执行包`: 包含运行应用程序所需的一切——`代码`、`运行环境`、`库`、`环境变量`和`配置文件`  

- `容器`(Container) – 类比`进程`  
是`镜像`的`runtime`实例——`镜像`在执行时在内存中会变成什么样子(即带有`状态`的`镜像`或`用户进程`)  

- `仓库`（Registry）– 类比 `repository`  
存放不同版本`镜像`的地方  

- `主机`（Host / Node）  
运行容器的`服务器`  

- `服务`（Service）  
一个镜像之上运行的一组可伸缩的`容器集合`，运行在一个容器集群中提供`同一功能服务`  

- `栈`（Stack）/ `命名空间`（Namaspace） / `应用`（Application）  
被编排的、可伸缩的一组`相互依赖的服务`，构成的一个`应用`  
  
https://gitee.com/li-jia666/service-computing/raw/master/Containerization/

## 实验过程

### 准备Docker环境
1. 检查Ubuntu的内核
	docker需要ubuntu的内核是高于3.10

	`uname -r`

	![](https://gitee.com/li-jia666/service-computing/raw/master/Containerization/img/1.png)

2. 安装docker

	` sudo apt-get install docker.io`

	![](https://gitee.com/li-jia666/service-computing/raw/master/Containerization/img/2.png)

3. 查看版本

	`docker --version`

	![](https://gitee.com/li-jia666/service-computing/raw/master/Containerization/img/3.png)

### 构建Docker镜像
1. 编写 Dockerfile

	以下是一个简单的镜像构建文件。
	```
	FROM ubuntu
	ENTRYPOINT ["top", "-b"]
	CMD ["-c"]
	```

	在该文件中，我们指明镜像基于`ubuntu:latest`镜像，镜像启动后运行`top -b -c`命令。其中，`ENTRYPOINT`描述了容器的入口点，这个入口点程序是容器的初始进程（PID 1），因此在退出后容器就会退出，这就规定了容器的生命周期和 top 命令一致。

	![](https://gitee.com/li-jia666/service-computing/raw/master/Containerization/img/4.png)

	![](https://gitee.com/li-jia666/service-computing/raw/master/Containerization/img/5.png)
	
	
2. 构建镜像
   `docker build . -t hello`

   ![](https://gitee.com/li-jia666/service-computing/raw/master/Containerization/img/6.png)


3. 运行镜像
   `docker run -it --rm hello -H `

   其中，`-it`表示通过终端与进程（容器）交互，stdin，stdout，stderr定向到 TTY，`-rm`表示容器运行完毕后删除此容器。

   ![](https://gitee.com/li-jia666/service-computing/raw/master/Containerization/img/7.png)

### MySQL与容器化

1. 启动服务器
	我们可执行以下命令运行一个 MySQL 容器，在命令中我们通过 -e 设置容器的环境变量参数，设定了 MySQL 数据库密码。
	`docker run -p 3386:3306 --name testmysql -e MYSQL_ROOT_PASSWORD=123456 -d mysql:5.7`

	![](https://gitee.com/li-jia666/service-computing/raw/master/Containerization/img/8.png)

	显示运行中容器

	`docker ps`

	![](https://gitee.com/li-jia666/service-computing/raw/master/Containerization/img/9.png)

	结果中的`Up 51 seconds`表明容器已正常运行。

2. 数据库持久化
	上述执行的命令是不具备持久化能力的，因为一旦容器停止运行，所存储的数据都丢失了。

	因此，我们需要创建`Volumes`，并与镜像链接，实现数据持久化。首先，我们需要删除原先创建的镜像，具体命令如下。

	`docker rm $(docker ps -a -q) -f -v`

	创建卷并挂载

	`docker volume create mydb`
	`docker run --name testmysql -e MYSQL_ROOT_PASSWORD=123456 -v mydb:/var/lib/mysql -d mysql:5.7`

	![](https://gitee.com/li-jia666/service-computing/raw/master/Containerization/img/10.png)

3. 启动客户端容器链接服务器

	`docker run --name myclient --link testmysql:mysql -it mysql:5.7 bash`

	在该命令中，我们使用`--link`参数，使得`myclient`和 `testmysql`容器链接在一起，命令执行结果如下图所示。

	客户端容器内可以使用 mysql 这个别名访问服务器。

	` mysql -hmysql -uroot -p123456`

	![](https://gitee.com/li-jia666/service-computing/raw/master/Containerization/img/11.png)

### Docker Compose 
Docker Compose 是一个命令行工具，它允许你定义和编排多容器 Docker 应用。它使用 YAML 文件来配置应用服务，网络和卷。

使用 Compose， 你可以定义一个可以运行在任何系统上的可移植应用环境。

Compose 通常被用来本地开发，单机应用部署，和自动测试。

1. 安装 Docker Compose
	`sudo apt install docker-compose`

	验证安装成功:
	`docker-compose --version`

	![](https://gitee.com/li-jia666/service-computing/raw/master/Containerization/img/12.png)

2. 编写：`stack.yml`
	`mkdir comptest && cd comptest`,`vi stack.yml`
	我们有服务器，db 和 adminer。当 docker-compose 运行，每一个服务运行一个镜像，创建一个独立的容器。其中**adminer**是一个数据库管理工具，提供了一个 Web 页面，用于管理数据库。在 **mysql** 中，我们指明了所使用的镜像、密码配置和重启策略；在 **adminer** 中，我们指明了所使用的镜像、重启策略和端口映射。
	配置文件内容如下：
```yml
version: '3.1'
services:
 db:
  image: mysql:5.7
  command: --default-authentication-plugin=mysql_native_password
  restart: always
  environment:
   MYSQL_ROOT_PASSWORD: 123456
 adminer:
  image: adminer
  restart: always
  ports:
   - 8080:8080 
```

3. 启动服务：`sudo docker-compose -f stack.yml up` 
	![](https://gitee.com/li-jia666/service-computing/raw/master/Containerization/img/13.png)

	在成功运行上述命令后，我们即可在浏览器中访问`:8080`端口，进入 Adminer 管理页面。

	![](https://gitee.com/li-jia666/service-computing/raw/master/Containerization/img/14.png)

	输入 MySQL 账户名和密码后，我们即可进入到数据库管理页面。

	![](https://gitee.com/li-jia666/service-computing/raw/master/Containerization/img/15.png)

### Docker 网络
#### 容器网络管理
1. 管理容器网络
	`sudo docker network ls`

	![](https://gitee.com/li-jia666/service-computing/raw/master/Containerization/img/16.png)

   - 容器`默认`使用网络：`Docker0`（桥接）
   - 容器支持的网络与类型
      - `bridge` （本机内网络）
      - `host`     （主机网卡） 
      - `overlay` （跨主机网络） 
      - `none`
      - `Custom`（网络插件）
  
	注意：docker-compose为每个应用建立自己的网络。
#### 自定义网络
2. 备制支持 `ifconfig` 和 `ping` 命令的 `ubuntu 容器`  

	运行`ubuntu`容器:`sudo docker run --name unet -it --rm ubuntu bash`  

    - 更新容器：`apt-get update`  
    ![](https://gitee.com/li-jia666/service-computing/raw/master/Containerization/img/17.png)
	- 安装网络工具包`net-tools`：`apt-get install net-tools`  
    ![](https://gitee.com/li-jia666/service-computing/raw/master/Containerization/img/18.png)
	- 安装`ping`的依赖包：`apt-get install iputils-ping -y`  
    ![](https://gitee.com/li-jia666/service-computing/raw/master/Containerization/img/19.png)
	- `ifconfig`命令  
    ![](https://gitee.com/li-jia666/service-computing/raw/master/Containerization/img/20.png)
	- `ping`命令：`ping 172.17.0.3`  
	![](https://gitee.com/li-jia666/service-computing/raw/master/Containerization/img/21.png)
	
	
	

3. 启动另一个命令窗口，由容器制作镜像：`sudo docker commit unet ubuntu:net` 
	![](https://gitee.com/li-jia666/service-computing/raw/master/Containerization/img/22.png)

	查看容器和镜像：

	![](https://gitee.com/li-jia666/service-computing/raw/master/Containerization/img/25.png)
	
4. 创建自定义网络：`sudo docker network create mynet`  
	![](https://gitee.com/li-jia666/service-computing/raw/master/Containerization/img/23.png)
	再次查看网络`sudo docker network ls`，可以看到`mynet`这个自定义网络已经创建成功了。
	![](https://gitee.com/li-jia666/service-computing/raw/master/Containerization/img/24.png)

5. 在两个窗口创建 u1,u2 容器网络：
   	```
  	docker run --name u1 -it -p 8080:80 --net mynet --rm ubuntu:net bash
	docker run --name u2 --net mynet -it --rm ubuntu:net bash
	```

	假设我们容器想要使用前面自定义的网络，在启动时通过--network或者缩写--net指定即可。

	另开一个窗口，使用如下命令：

	```
	docker info
	```

	![](https://gitee.com/li-jia666/service-computing/raw/master/Containerization/img/26.png)

	```
	docker network connect bridge u1
	docker network disconnect mynet u1
  	```
   
	![](https://gitee.com/li-jia666/service-computing/raw/master/Containerization/img/27.png)
	可以看到u1已经被从网络`mynet`删除了。

### Docker 仓库（Registry）
	
#### 搭建私有容器仓库

1. 运行一个本地仓库
   运行一个仓库容器
   
   `docker run -d -p 5000:5000 --restart=always --name registry registry:2`  

   ![](https://gitee.com/li-jia666/service-computing/raw/master/Containerization/img/28.png)

2. 从Docker Hub 复制一个**镜像**到仓库中

   - 从docker hub拉取ubuntu:16.04镜像
	```
  	docker pull ubuntu:16.04
	```
	![](https://gitee.com/li-jia666/service-computing/raw/master/Containerization/img/29.png)
   - 标记该镜像为`localhost:5000/my-ubuntu`，当标记第一部分为主机名和端口时，docker将它们视为push时仓库的地址。
    ```
  	docker tag ubuntu:16.04 localhost:5000/my-ubuntu
	```

   - 将镜像push到正运行在localhost:5000的本地仓库中 
	```
	docker push localhost:5000/my-ubuntu
	```
	![](https://gitee.com/li-jia666/service-computing/raw/master/Containerization/img/30.png)
	- 删除本地缓存的`ubuntu:16.04` 和 `localhost:5000/my-ubuntu`镜像，以测试从仓库中pull镜像。这些操作不回移除仓库中的`localhost:5000/my-ubuntu`镜像。
	```
	sudo docker image remove ubuntu:16.04
	sudo docker image remove localhost:5000/my-ubuntu  
	```

	![](https://gitee.com/li-jia666/service-computing/raw/master/Containerization/img/31.png)

	- 从本地仓库拉取镜像`localhost:5000/my-ubuntu`
	```
	sudo docker pull localhost:5000/my-ubuntu
	```

	![](https://gitee.com/li-jia666/service-computing/raw/master/Containerization/img/32.png)

3. 停用本地仓库
	停用仓库使用命令`docker container stop`，例如
	```
	docker container stop registry
	```

	删除容器使用命令`docker container rm`，例如
	```
	docker container rm -v registry
	```

#### 阿里云容器镜像服务实践

1. 注册账号
	访问 https://cr.console.aliyun.com
	
	开通**容器镜像服务**：设置registry密码
	
	![](https://gitee.com/li-jia666/service-computing/raw/master/Containerization/img/33.png)
	
	创建镜像仓库的命名空间

	![](https://gitee.com/li-jia666/service-computing/raw/master/Containerization/img/35.png)

	创建镜像仓库

	![](https://gitee.com/li-jia666/service-computing/raw/master/Containerization/img/36.png)

	选择代码源为本地仓库

	![](https://gitee.com/li-jia666/service-computing/raw/master/Containerization/img/37.png)


2. 测试上传`hello-world`镜像  
	- 登陆：`sudo docker login --username=卷耳多多多 registry.cn-shenzhen.aliyuncs.com`
	![](https://gitee.com/li-jia666/service-computing/raw/master/Containerization/img/34.png)
    

	- 标签：`sudo docker tag hello-world registry.cn-shenzhen.aliyuncs.com/gail/repo:hello-world`  

		使用命令`docker images`显示本地镜像库内容，如果没有该镜像，执行以下命令
		`sudo docker run hello-world`

	- 上传：`sudo docker push registry.cn-shenzhen.aliyuncs.com/gail/repo:hello-world`
		![](https://gitee.com/li-jia666/service-computing/raw/master/Containerization/img/38.png)

		这时在阿里云仓库上可以看到上传的镜像：
		
		![](https://gitee.com/li-jia666/service-computing/raw/master/Containerization/img/39.png)
	- 下载：`sudo docker pull registry.cn-shenzhen.aliyuncs.com/gail/repo:hello-world`  
	![](https://gitee.com/li-jia666/service-computing/raw/master/Containerization/img/40.png)

	下载成功后可以看到本地镜像多了registry.cn-shenzhen.aliyuncs.com/gail/repo:hello-world

	![](https://gitee.com/li-jia666/service-computing/raw/master/Containerization/img/41.png)

	- 删除：`sudo docker rmi registry.cn-shenzhen.aliyuncs.com/gail/repo:hello-world`  

	删除成功后可以看到本地镜像少了registry.cn-shenzhen.aliyuncs.com/gail/repo:hello-world

	![](https://gitee.com/li-jia666/service-computing/raw/master/Containerization/img/42.png)
	
	- 运行：`sudo docker run --rm registry.cn-shenzhen.aliyuncs.com/gail/repo:hello-world`

	![](https://gitee.com/li-jia666/service-computing/raw/master/Containerization/img/43.png)
	

	- 退出：`sudo docker logout registry.cn-shenzhen.aliyuncs.com`  

	![](https://gitee.com/li-jia666/service-computing/raw/master/Containerization/img/44.png)

	