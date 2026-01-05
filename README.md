# Hook & Event Demo

这是一个展示如何在 Go 项目中结合使用 **Hook（钩子）** 和 **Event（事件）** 模式来解耦业务逻辑的示例项目。

详细的内容介绍全在微信公众号中。干货持续更新，敬请关注「代码扳手」微信公众号：

<img width="430" height="430" alt="image" src="wx.jpg" />

## 📖 项目介绍

在复杂的业务系统中，核心业务逻辑往往容易被各种辅助逻辑（如风控、通知、积分、数据分析等）纠缠在一起，导致代码难以维护。本项目演示了如何通过两种设计模式来解决这个问题：

1.  **Hook 模式 (钩子)**：
    *   **用途**：在核心流程的特定生命周期节点（如 `Before` 执行前、`After` 执行后）插入逻辑。
    *   **特点**：可以感知上下文，支持同步阻断（如风控拦截）或异步执行。
2.  **Event 模式 (事件)**：
    *   **用途**：核心流程完成后，向外部广播消息，触发副作用。
    *   **特点**：完全解耦，通常用于“发后即忘”的场景（如发短信、加积分）。

## 📂 目录结构

```text
hook-event-demo/
├── event/          # 事件总线实现
│   ├── bus.go      # 简单的内存事件总线
│   ├── consumer.go # 事件消费者注册 (SMS, Point, BI)
│   └── event.go    # 事件定义
├── example/        # 业务逻辑示例
│   └── order.go    # 订单服务，集成 Hook 和 Event
├── hook/           # Hook 引擎实现
│   ├── engine.go   # Hook 执行引擎
│   ├── hook.go     # Hook 定义 (Sync/Async, Priority)
│   └── context.go  # Hook 上下文数据
├── main.go         # 程序入口
└── go.mod
```

## 🛠️ 核心设计

### 1. Hook 引擎 (`hook/`)
Hook 引擎允许我们在业务方法的执行前后注入逻辑。

*   **类型 (`Type`)**：支持 `Before` (前置) 和 `After` (后置)。
*   **模式 (`Mode`)**：
    *   `Sync`: 同步执行。如果设置了 `MustSucceed: true`，且 Hook 返回错误，则会阻断主流程。
    *   `Async`: 异步执行，不阻塞主流程。
*   **优先级 (`Priority`)**：数值越大优先级越高，按顺序执行。

### 2. Event 总线 (`event/`)
一个轻量级的内存事件总线，用于处理业务完成后的后续动作。

*   **Publish/Subscribe**：发布订阅模式。
*   **异步处理**：事件处理函数在独立的 Goroutine 中执行，不影响主链路响应时间。

## 🚀 运行示例

直接运行 `main.go`：

```bash
go run main.go
```

### 运行逻辑分析

`main.go` 中模拟了两个场景：

1.  **创建普通订单 (金额 500)**
    *   **Hook (Before)**: `RiskCheck` (风控) -> 通过。
    *   **Core**: 保存订单。
    *   **Hook (After)**: `PublishOrderCreatedEvent` -> 发布事件。
    *   **Event**: 触发 SMS、积分、BI 系统的监听器。

2.  **创建高风险订单 (金额 20,000)**
    *   **Hook (Before)**: `RiskCheck` (风控) -> **拦截** (金额 > 10,000)。
    *   **Result**: 返回错误 "risk rejected"，后续流程（保存订单、发事件）均不会执行。

### 预期输出

```text
=== create normal order ===
[CORE] order order_1001 saved
[SMS] order=order_1001
[POINT] user=user_1
[BI] order=order_1001 amount=500
=== create risk order ===
order failed: risk rejected
```
*(注：由于异步执行，SMS/POINT/BI 的日志顺序可能不同)*

## 💡 最佳实践总结

*   **使用 Hook 做控制**：当你需要决定“是否允许继续执行”或者“必须在核心逻辑前后紧密执行”时（如参数校验、风控、事务管理），使用 Hook。
*   **使用 Event 做解耦**：当你只需要“通知别人这件事发生了”，且不关心别人处理的结果或耗时（如发送通知、更新报表），使用 Event。
