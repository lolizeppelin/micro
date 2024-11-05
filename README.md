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

## 组件方法命名

#### 1. 标准HTTP restful方法(不可以网关转发)

```text
Get
List
Create
Update
Patch
Delete
```

#### 2. 非标准HTTP方法,使用下述前缀,内部提取前缀后的字符串小写为方法名(不可以网关转发)

```text
GET_
POST_
PUT_
PATCH_
DELETE_
```

#### 3. 内部RPC,使用下述前缀,内部提取前缀后的字符串小写为方法名(不可以网关转发)

```text
RPC_
```

#### 4. 其他公开方法为外部RPC方法,统一POST请求(可以网关转发)

