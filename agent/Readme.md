# DDog Agent

## Agents列表

* DestoryAgent
    > 销毁Agent，从DDog Main中接受销毁任务。然后调用集群API销毁对应的服务实例

* MonitorAgent
    > 监控Agent, 接受所有处理失败的任务数据。然后择机重启任务，同时统计失败次数。 如果失败次数过高，则发送告警信息

* RetriAgent
    > 状态监测Agent，主动检索命名空间中的服务状态。发现有失败服务之后，告之MonitorAgent,同时将此服务标记为失败

* SpiderAgent
    > 服务探针Agent, 伴随服务同步启动，用来侦测服务状态是否正常

## Change Log

* v0.6.9
 - 默认负载模式由集群内访问改为VPC内网访问

* v0.6.7
 - 增加SidecarAgent

* v0.6.6
 - 增加状态检查Agent RetirAgent

* v0.6.5
 - MonitorAgent 在保存数据时会判断当前是否存在。如果存在则合并，否则插入.
 - MonitorAgent 修复状态收集时不增加事件次数的问题