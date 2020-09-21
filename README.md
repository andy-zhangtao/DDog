# DDog
A tool for auto generate coredns configure file

## How to Build?

** All build scripts in `buildscript` dir **

The belowing commands exec in `buildscript` dir. 

+ k8smonitor

```
./buildK8sMonitorAgent.sh k8smonitor
```

+ deployagent

```
./buildDeployAgent.sh deployagent
```

+ caas-deploy-agent

```
./buildGraphAgent.sh caas-deploy-agent
```

+ destoryagent

```
./buildDestoryAgent.sh destoryagent
```

+ monitoragent

```
./buildMonitorAgent.sh monitoragent
```

+ replicaagent

```
./buildReplicaAgent.sh replicaagent
```

## API Reference

## How to Run?
推荐使用封装好的镜像:`vikings/ddog`

* MakeFile Command
    - make client 构建Agent客户端
    - make runclient 构建Agent客户端并且同时运行
    - make build 构建DDog主程序
    - make run 构建DDog主程序并且运行
    - make agent-release 构建Agent可发布版本
    - make srv-release 构建DDog可发布版本
    - make release 同时构建Agent DDog可发布版本
## ENV Reference
以下变量不允许为空

- DDOG_MONGO_DB: mongo数据库名称
- DDOG_REGION: 集群所在区域   
- DDOG_MONGO_ENDPOINT: mongo链接信息
- DDOG_NAME_SPACE: 默认命名空间
- DDOG_NSQD_ENDPOINT: NSQ链接地址
   
以下变量为可选项
- DDOG_MONGO_NAME: mongo用户名 
- DDOG_MONGO_PASSWD: mongo口令
- DDOG_DEBUG: 是否输出调试信息,默认为false

# Change Log

### v0.6.10
* 优化服务状态查询服务. 当部署状态为失败时,返回500错误

### v0.6.9
* 将部署功能从主服务中剥离出来，单独成立一个DeployAgent服务

### v0.6.8
* 每个服务都增加一个sidecar,用来侦测服务状态

### v0.6.7
* 增加资源限额. 默认CPU:0.5,CPU上限:1.5,Memory:300,Memory上限:800

### v0.6.6
* 调整创建服务时的状态轮询策略,当遇到多次失败后，将此服务置为失败

### v0.6.5
* 修复调用svcconf check接口，预期结果不幂等的问题
* 调整升级规则，当直接调用Svcconf Deploy接口时按照直接升级来处理

### v0.6.4
* DDog修复以下issue:
 - 首次创建容器配置时,respon body为空的问题
 - 修复注册集群元数据时响应超时的问题

### v0.6.3
* 修复升级服务后无法获取服务状态的问题

### v0.6.2
* DDog在创建服务时启用健康检测和就绪检测

### v0.6.1
* DDog将所有log替换为logrus

### v0.6.0
* 使用logrus日志框架替换原生log框架
* 使用Nsq作为任务分发工具
* 将删除服务功能由同步改为异步
* 修复DDog直接直接升级时不会删除旧服务的问题

### v0.5.1
* 在创建服务之前会尝试删除当前正在运行的服务

### v0.1.2

* 剥离服务扫描功能
* 配合Docker Logging Plugin(logchain)，增加默认环境变量
* 去掉健康检测和就绪检测
