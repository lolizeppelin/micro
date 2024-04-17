## 微服务框架,缝合来源
### - [go-micro](https://github.com/go-micro/go-micro)
### - [pitaya](https://github.com/topfreegames/pitaya)


# 文件夹说明

- breaker 熔断器
- broker 异步推送(默认kafka)
- breaker grpc客户端
- codec 序列化工具
- config 全局配置
- errors 通用错误
- log 日志
- registry 注册服务
- resolver 
- selector 选择器
- tracing 链路追踪
- transport grpc传输层
- utils 通用工具



## IDEA 设置
- Ctrl+Shift+A 搜索Registry并点击
```text
取消勾选下列值
go.run.processes.with.pty
terminal.use.conpty.on.windows
```