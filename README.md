# DDog
A tool for auto generate coredns configure file

## API Reference

## How to Run?
推荐使用封装好的镜像:`vikings/ddog`

## ENV Reference
以下变量不允许为空
- EnvDomain: 需要DNS解析的域名. 例如: mydomain.com
- EnvEtcd: Etcd地址,CoreDNS用来持久化数据. 例如: etcd.com:2379
- EnvUpStream: 上游DNS地址, 支持三种配置方式:
   * 配置为/etc/resolv.conf
   * 配置为IP,例如: 10.0.0.1;10.0.0.2
   * 配置为IP+Port,例如: 10.0.0.1:54;10.0.0.2
   
以下变量若为空会使用默认值
- EnvConfPath: Corefile路径,默认为 /