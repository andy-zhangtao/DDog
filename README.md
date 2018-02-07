# DDog
A tool for auto generate coredns configure file

## API Reference

## How to Run?
推荐使用封装好的镜像:`vikings/ddog`

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

### v0.6.4
* DDog修复以下issue:
 - 首次创建容器配置时,respon body为空的问题

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