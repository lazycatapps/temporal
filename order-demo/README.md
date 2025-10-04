# Temporal 订单处理系统 Demo

这是一个基于 Temporal 的订单处理系统示例，演示了如何使用 Temporal 构建可靠的分布式工作流。

## 功能特性

本 Demo 实现了一个完整的订单处理流程，包含 5 个活动（Activity）：

### 核心流程

1. **订单验证（ValidateOrderActivity）**
   - 验证订单数量是否大于 0
   - 验证订单金额是否大于 0
   - 模拟耗时：1 秒

2. **库存预留（ReserveInventoryActivity）**
   - 为订单预留商品库存
   - **模拟失败**：10% 概率库存不足
   - 模拟耗时：2 秒

3. **支付处理（ProcessPaymentActivity）**
   - 处理订单支付
   - **模拟失败**：5% 概率支付网关超时
   - 失败时触发补偿事务（释放库存）
   - 模拟耗时：2 秒

4. **订单发货（ShipOrderActivity）**
   - 安排订单发货并生成物流单号
   - 物流单号格式：`TRACK-{订单ID}-{时间戳}`
   - 模拟耗时：3 秒

5. **通知发送（SendNotificationActivity）**
   - 向客户发送订单完成通知
   - **容错处理**：通知失败不影响订单完成
   - 模拟耗时：1 秒

### 补偿事务

**库存释放（ReleaseInventoryActivity）**
- 在支付失败时自动调用
- 释放已预留的库存
- 实现 Saga 模式的补偿逻辑

## Temporal 特性演示

✅ **自动重试机制**
   - 所有活动失败时自动重试，最多 3 次
   - 配置在工作流的 `RetryPolicy` 中

✅ **补偿事务（Saga 模式）**
   - 支付失败时自动释放已预留的库存
   - 保证数据一致性

✅ **容错处理**
   - 通知发送失败只记录警告，不影响订单完成
   - 工作流继续执行并返回成功状态

✅ **可靠性保证**
   - 活动超时时间：10 秒（StartToCloseTimeout）
   - Worker 中断后可恢复执行
   - 所有状态持久化到 Temporal 服务

## 项目结构

```
order-demo/
├── workflows/           # 工作流和活动实现包
│   ├── workflow.go      # 订单工作流定义（OrderWorkflow）
│   └── activities.go    # 6 个活动函数实现
├── worker/              # Worker 程序（执行工作流和活动）
│   └── main.go          # 注册并启动 Worker
├── client/              # 客户端程序（触发工作流）
│   └── main.go          # 创建订单并启动工作流
├── go.mod               # Go 模块定义
├── go.sum               # 依赖锁定文件
└── README.md
```

### 核心组件说明

**workflows/workflow.go**
- 定义 `Order` 数据结构
- 实现 `OrderWorkflow` 工作流
- 按顺序执行 5 个步骤
- 实现补偿逻辑（支付失败时释放库存）

**workflows/activities.go**
- 实现 6 个活动函数
- 包含随机失败模拟
- 打印执行日志便于观察

**worker/main.go**
- 连接到 Temporal 服务器
- 注册工作流和所有活动
- 监听 `order-queue` 任务队列
- 支持 `-temporal-address` 参数（默认 localhost:7233）

**client/main.go**
- 创建示例订单数据
- 启动工作流执行
- 等待并显示执行结果
- 提供 Web UI 访问链接
- 支持 `-temporal-address` 参数（默认 localhost:7233）

## 运行步骤

### 1. 确保 Temporal 服务已启动

```bash
cd ..
docker compose ps
```

确认以下服务正在运行：
- temporal (端口 7233)
- temporal-ui (端口 8080)
- temporal-postgresql (端口 5432)

### 2. 启动 Worker

在一个终端窗口运行：

```bash
go run worker/main.go
```

可选参数：
```bash
go run worker/main.go -temporal-address=localhost:7233
```

你应该看到：
```
Starting Order Processing Worker (connected to localhost:7233)...
Worker is listening on task queue: order-queue
```

### 3. 运行客户端触发订单

在另一个终端窗口运行：

```bash
go run client/main.go
```

可选参数：
```bash
go run client/main.go -temporal-address=localhost:7233
```

客户端会：
- 创建一个包含时间戳的订单 ID
- 显示订单详情（客户 ID、产品名称、数量、金额）
- 启动工作流并显示 Workflow ID 和 RunID
- 提供 Web UI 链接用于查看执行详情
- 等待并显示最终执行结果

### 4. 查看执行过程

**客户端终端输出示例：**
```
=================================================
Temporal Server: localhost:7233
Starting Order Workflow for Order: ORDER-1728012345
Customer: CUST-12345
Product: Laptop (Quantity: 2)
Total Amount: $2999.98
=================================================

Workflow started with ID: order-workflow-ORDER-1728012345
RunID: a1b2c3d4-e5f6-7890-abcd-ef1234567890

You can view the workflow in Temporal UI:
http://localhost:8080/namespaces/default/workflows/order-workflow-ORDER-1728012345

Waiting for workflow to complete...
```

**Worker 终端输出示例：**
```
[Activity] Validating order ORDER-1728012345...
[Activity] Reserving 2 units of Laptop for order ORDER-1728012345...
[Activity] Processing payment of $2999.98 for order ORDER-1728012345...
[Activity] Shipping order ORDER-1728012345 to customer CUST-12345...
[Activity] Sending notification to customer CUST-12345 for order ORDER-1728012345...
```

**成功完成时客户端显示：**
```
=================================================
Workflow Result: Order ORDER-1728012345 completed successfully!
=================================================
```

### 5. 在 Web UI 中查看详细信息

打开浏览器访问：http://localhost:8080

在 Temporal Web UI 中可以看到：
- **工作流列表**：所有执行的订单工作流
- **执行历史**：每个活动的开始、完成时间
- **事件时间线**：工作流的完整执行路径
- **输入输出**：每个活动的参数和返回值
- **重试记录**：失败活动的重试次数和错误信息
- **补偿事务**：支付失败时的库存释放记录

## 测试场景

### 场景 1：正常流程（约 85% 概率）

所有步骤成功完成，订单状态从创建到发货通知：
1. ✅ 订单验证通过
2. ✅ 库存预留成功
3. ✅ 支付处理成功
4. ✅ 订单发货成功（生成物流单号）
5. ✅ 通知发送成功

**预期耗时**：约 9 秒（1+2+2+3+1）

### 场景 2：库存不足（约 10% 概率）

在库存预留步骤失败：
1. ✅ 订单验证通过
2. ❌ 库存预留失败（insufficient inventory）
3. 🔄 自动重试最多 3 次
4. ❌ 最终工作流失败

**观察点**：
- Worker 会显示重试日志
- Web UI 可以看到 3 次重试记录
- 客户端会收到错误信息

### 场景 3：支付失败（约 5% 概率）

在支付处理步骤失败，触发补偿事务：
1. ✅ 订单验证通过
2. ✅ 库存预留成功
3. ❌ 支付处理失败（payment gateway timeout）
4. 🔄 自动重试最多 3 次
5. ⚠️ 触发补偿：释放已预留的库存
6. ❌ 最终工作流失败

**观察点**：
- Worker 会显示 `[Activity] Releasing X units...`
- Web UI 可以看到补偿活动执行记录
- 实现了 Saga 模式的数据一致性

### 多次测试

运行多次客户端观察不同场景：
```bash
# 运行 10 次观察各种情况
for i in {1..10}; do
  echo "Test run $i"
  go run client/main.go
  sleep 2
done
```

## 代码要点

### 工作流定义（workflows/workflow.go）

```go
// 配置活动选项
ao := workflow.ActivityOptions{
    StartToCloseTimeout: 10 * time.Second,  // 活动超时时间
    RetryPolicy: &temporal.RetryPolicy{
        MaximumAttempts: 3,                  // 最多重试 3 次
    },
}

// 补偿事务示例
err = workflow.ExecuteActivity(ctx, ProcessPaymentActivity, order).Get(ctx, &paymentResult)
if err != nil {
    logger.Error("Payment failed, rolling back inventory", "error", err)
    // 释放已预留的库存
    _ = workflow.ExecuteActivity(ctx, ReleaseInventoryActivity, order).Get(ctx, nil)
    return "", err
}
```

### 订单数据结构

```go
type Order struct {
    OrderID       string    // 格式：ORDER-{时间戳}
    CustomerID    string    // 示例：CUST-12345
    ProductName   string    // 示例：Laptop
    Quantity      int       // 示例：2
    TotalAmount   float64   // 示例：2999.98
    Status        string    // 示例：pending
}
```

### 任务队列

- **队列名称**：`order-queue`
- **Worker** 监听该队列等待任务
- **Client** 向该队列提交工作流执行请求

## 扩展建议

基于当前实现，可以尝试以下扩展：

1. **添加 Signal 支持**
   - 实现订单取消功能
   - 人工审批流程（等待审批信号）

2. **添加 Query 支持**
   - 查询订单当前状态
   - 查询已完成的步骤

3. **增强错误处理**
   - 不同错误类型的重试策略
   - 发货失败时的退款逻辑

4. **添加定时功能**
   - 订单超时自动取消（使用 Timer）
   - 延迟发货（使用 Sleep）

5. **集成真实服务**
   - 连接真实数据库存储订单
   - 集成第三方支付网关
   - 对接物流系统 API

6. **监控和可观测性**
   - 添加 Metrics 指标
   - 集成 OpenTelemetry
   - 设置告警规则
