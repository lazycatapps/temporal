# Temporal on LazyCat

Temporal 工作流引擎在 LazyCat 平台上的部署项目，包含完整的服务配置和示例应用。

## 什么是 Temporal？

Temporal 是一个开源的微服务编排平台，用于构建可靠、可扩展的分布式应用程序。它通过工作流（Workflow）和活动（Activity）的概念，帮助开发者：

- ✅ **自动处理重试** - 活动失败时自动重试，无需手动编写重试逻辑
- ✅ **保证执行完整性** - 即使服务重启，工作流也会从中断处继续执行
- ✅ **实现长时间运行的流程** - 支持几秒到几年的工作流执行时间
- ✅ **简化分布式事务** - 提供补偿事务（Saga）模式支持
- ✅ **提供可观测性** - 内置 Web UI 查看工作流执行历史

### Temporal 核心概念

**Workflow（工作流）**
- 业务流程的定义和编排逻辑
- 确定性执行，可以被重放（replay）
- 长时间运行，状态持久化

**Activity（活动）**
- 具体的业务操作实现
- 可以有副作用（调用 API、写数据库等）
- 支持自动重试和超时控制

**Worker（工作者）**
- 执行工作流和活动的进程
- 从任务队列拉取任务并执行
- 可以水平扩展

**Task Queue（任务队列）**
- Worker 和 Client 之间的通信桥梁
- 负载均衡和任务分发

## 项目结构

```
temporal/
├── order-demo/              # 订单处理系统示例
│   ├── workflows/           # 工作流和活动实现
│   ├── worker/              # Worker 程序
│   ├── client/              # 客户端程序
│   └── README.md            # 详细说明文档
└── README.md                # 本文件
```

## 部署架构

本项目的 Temporal 服务已部署在 LazyCat 平台上：

- **应用包名**：`cloud.lazycat.app.liu.temporal`
- **子域名**：`temporal`
- **Temporal gRPC 服务**：`temporal.{BOXNAME}.heiyu.space:7233`
- **Web UI**：`https://temporal.{BOXNAME}.heiyu.space`

### LazyCat 地址规则

LazyCat 平台使用以下地址格式：
- **HTTP/HTTPS**：`https://{subdomain}.{boxname}.heiyu.space`
- **TCP（如 gRPC）**：`{subdomain}.{boxname}.heiyu.space:{port}`

其中：
- `{subdomain}` 是应用配置的子域名（本应用为 `temporal`）
- `{boxname}` 是你的 Box 名称（通过 `lzc-cli box default` 获取）
- `{port}` 是对外暴露的端口（本应用 gRPC 端口为 `7233`）

## 快速开始

### 1. 获取你的 Box 名称

```bash
BOXNAME=$(lzc-cli box default)
echo $BOXNAME
# 输出示例: demo
```

### 2. 访问 Temporal Web UI

在浏览器中打开：
```
https://temporal.${BOXNAME}.heiyu.space
```

例如，如果你的 BOXNAME 是 `demo`，则访问：
```
https://temporal.demo.heiyu.space
```

### 3. 运行订单处理 Demo

进入示例项目目录：
```bash
cd order-demo
```

启动 Worker（终端 1）：
```bash
go run worker/main.go -temporal-address temporal.${BOXNAME}.heiyu.space:7233
```

运行客户端触发订单（终端 2）：
```bash
go run client/main.go -temporal-address temporal.${BOXNAME}.heiyu.space:7233
```

查看详细说明：
```bash
cat order-demo/README.md
```

## 示例项目：订单处理系统

order-demo 是一个完整的订单处理工作流示例，演示了：

- **5 步订单流程**：验证 → 预留库存 → 支付 → 发货 → 通知
- **自动重试**：活动失败时自动重试最多 3 次
- **补偿事务**：支付失败时自动释放库存（Saga 模式）
- **容错处理**：通知失败不影响订单完成
- **随机失败模拟**：10% 库存不足，5% 支付失败

详细文档请查看：[order-demo/README.md](order-demo/README.md)

## 学习资源

### 官方文档
- [Temporal 官方网站](https://temporal.io/)
- [Go SDK 文档](https://docs.temporal.io/dev-guide/go)
- [核心概念](https://docs.temporal.io/concepts)

### 参考项目
- [Temporal 核心项目](https://github.com/temporalio/temporal) - Temporal Server 源代码
- [Temporal Docker Compose](https://github.com/temporalio/docker-compose) - 官方 Docker Compose 部署配置
- [Go SDK 示例代码库](https://github.com/temporalio/samples-go) - 官方 Go 示例集合

## 常用命令

### 运行 Worker 和 Client

```bash
# 获取 Box 名称
BOXNAME=$(lzc-cli box default)

# 设置 Temporal 地址变量
TEMPORAL_ADDR=temporal.${BOXNAME}.heiyu.space:7233

# 启动 Worker
cd order-demo
go run worker/main.go -temporal-address $TEMPORAL_ADDR

# 运行客户端
cd order-demo
go run client/main.go -temporal-address $TEMPORAL_ADDR
```

### 访问服务

```bash
# 获取 Box 名称
BOXNAME=$(lzc-cli box default)

# Web UI 地址
echo "Temporal Web UI: https://temporal.${BOXNAME}.heiyu.space"

# gRPC 服务地址
echo "Temporal gRPC: temporal.${BOXNAME}.heiyu.space:7233"
```

### Temporal CLI（如果安装了）

连接到 LazyCat 上的 Temporal 服务：

```bash
# 设置环境变量
export TEMPORAL_ADDRESS=temporal.${BOXNAME}.heiyu.space:7233

# 查看工作流列表
temporal workflow list --address $TEMPORAL_ADDRESS

# 查看工作流详情
temporal workflow describe --workflow-id <workflow-id> --address $TEMPORAL_ADDRESS

# 终止工作流
temporal workflow terminate --workflow-id <workflow-id> --address $TEMPORAL_ADDRESS
```

## 故障排除

### Worker 无法连接

```bash
# 1. 检查 Box 名称是否正确
BOXNAME=$(lzc-cli box default)
echo "Your box name: $BOXNAME"

# 2. 检查 Temporal 服务是否运行
# 访问 Web UI，应该能看到 Temporal 管理界面
echo "Open in browser: https://temporal.${BOXNAME}.heiyu.space"

# 3. 检查网络连接（测试 TCP 端口）
nc -zv temporal.${BOXNAME}.heiyu.space 7233

# 4. 使用正确的地址参数运行 Worker
go run worker/main.go -temporal-address temporal.${BOXNAME}.heiyu.space:7233
```

### 地址配置错误

确保使用正确的 LazyCat 地址格式：

❌ **错误示例**：
- `localhost:7233`（本地地址）
- `temporal:7233`（容器内部地址）
- `temporal.heiyu.space:7233`（缺少 boxname）

✅ **正确示例**：
- `temporal.demo.heiyu.space:7233`（假设 boxname 是 demo）
- `temporal.${BOXNAME}.heiyu.space:7233`（使用变量）

### 工作流执行失败

- 在 Temporal Web UI 查看错误详情
- 检查 Worker 终端的日志输出
- 确认活动函数已正确注册
- 确认 Worker 已成功连接到 Temporal 服务

## 使用指南

### 基于此项目开发你的工作流

1. **参考示例项目**
   - 查看 order-demo 的实现方式
   - 复制项目结构创建新的工作流应用
   - 连接到 LazyCat 上的 Temporal 服务

2. **部署你的 Worker**
   - Worker 可以运行在任何可以访问 LazyCat 的环境
   - 使用 `-temporal-address` 参数指定 Temporal 服务地址
   - Worker 支持水平扩展以提高吞吐量

3. **集成到现有系统**
   - 在你的应用中引入 Temporal Go SDK
   - 调用 Temporal 服务执行工作流
   - 通过 Web UI 监控工作流执行状态

### 扩展建议

- **添加更多工作流**：实现其他业务流程（如数据处理、任务调度等）
- **配置持久化**：确保 PostgreSQL 数据定期备份
- **监控告警**：集成监控系统跟踪工作流执行情况
- **多环境部署**：在不同的 Box 上部署开发/测试/生产环境

## 项目背景

本项目将 Temporal 工作流引擎迁移到 LazyCat 平台，使团队能够：

- 快速部署和使用 Temporal 服务
- 无需维护复杂的 Kubernetes 集群
- 通过 LazyCat 平台统一管理服务
- 利用 LazyCat 的域名和网络能力简化访问

## 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件。

Temporal 本身也遵循 MIT 许可证。
