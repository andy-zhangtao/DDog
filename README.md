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
- DDOG_ETCD_ENDPOINT：ETCD链接信息
   
以下变量为可选项
- DDOG_MONGO_NAME: mongo用户名 
- DDOG_MONGO_PASSWD: mongo口令
- DDOG_DEBUG: 是否输出调试信息,默认为false