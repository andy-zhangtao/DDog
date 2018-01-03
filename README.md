# DDog
A tool for auto generate coredns configure file

## API Reference

## How to Run?
推荐使用封装好的镜像:`vikings/ddog`

## ENV Reference
以下变量不允许为空
- DDOG_DOMAIN: 需要DNS解析的域名. 例如: mydomain.com
- DDOG_ETCD_ENDPOINT: Etcd地址,CoreDNS用来持久化数据. 例如: etcd.com:2379
- DDOG_UP_STREAM: 上游DNS地址, 支持三种配置方式:
   * 配置为/etc/resolv.conf
   * 配置为IP,例如: 10.0.0.1;10.0.0.2
   * 配置为IP+Port,例如: 10.0.0.1:54;10.0.0.2
- DDOG_MONGO_DB: mongo数据库名称
   
   
以下变量若为空会使用默认值
- DDOG_CONF_PATH: Corefile路径,默认为 /

以下变量为可选项
- DDOG_MONGO_NAME: mongo用户名 
- DDOG_MONGO_PASSWD: mongo口令