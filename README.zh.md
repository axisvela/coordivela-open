> [English](README.md) | 中文

# COORDIVELA Open

**COORDIVELA 开发协调控制平面的公开集成面**

COORDIVELA 是**开发协调控制平面（Development Coordination Control Plane,
DCCP）**。它协调人、编码智能体和开发工作，通常面向单个仓库：决定下一步做
什么、分配工作、跟踪进度、发现停滞、发起干预，并记录协调证据。

COORDIVELA 不执行代码、不验证正确性，也不拥有产品运行时。这些职责仍由各自
的产品承担。

## COORDIVELA 解决什么问题

开发工作正变得越来越并发：多个人和智能体可能同时操作同一个仓库、依赖、计划
和交付门禁。没有协调层时，团队会遇到重复认领、更新丢失、所有者过期、依赖发
现过晚、停滞会话不可见，以及协调决策缺少审计记录等问题。

COORDIVELA 把这些协调过程显式化并持久化。工作项有身份、状态、所有者、租约
和版本；依赖和路径冲突会在工作推进前检查；计划在重启后仍然有效；运行时适配
器把协调循环接到执行系统，但不接管执行；配置授权适配器后，干预会被授权；协
调事件可以被观察和重放。

## 它带来什么价值

- **统一的协调事实源**：工作状态、所有权、进度、依赖、计划和干预集中管理。
- **并发安全**：通过租约、fencing token、乐观版本、依赖检查和路径重叠冲突
  警告，减少丢失更新和重复执行。
- **持久计划**：支持 cron、间隔和一次性语义，支持事件触发、时区/DST 处理和
  misfire 策略。
- **运行时协调**：通过 PRAXOVELA、进程监督和 OpenCode 观察/控制适配器连接
  执行系统。
- **跨产品边界**：通过公开合同接入 TEKMOVELA 验证门禁、Rheovela 流程门禁、
  ONTOVELA 现实门禁、Mnemovela/SYMBIVELA 记录、Harmovela 事件和
  GNOSIVELA 语义支撑。
- **可治理**：可选的 AEGIVELA 授权、决策审计、幂等键和可重放的状态转换。
- **可观测**：租户/工作区范围的 SSE 协调事件、有界历史、只读投影和
  Prometheus 指标。
- **社区规模部署**：单二进制、内存或 PostgreSQL 存储，参考规模为少于五个并
  发智能体。

## COORDIVELA 不拥有什么

| 关注点 | 归属产品 |
| --- | --- |
| 代码执行和智能体循环 | PRAXOVELA |
| 验证、保证和发布门禁 | TEKMOVELA |
| 持久流程引擎 | Rheovela |
| 记忆存储与检索 | Mnemovela |
| 人类协作 UI | SYMBIVELA |
| 事件和异步协调协议 | Harmovela |
| 语义与知识底座 | GNOSIVELA |
| 身份、授权和审批 | AEGIVELA / Archon Guard |

COORDIVELA 是协调者，不是控制权威。会改变执行的干预由拥有执行的产品完成。

## COORDIVELA Open

COORDIVELA Open 是 COORDIVELA 的**公开分发和集成面**，包含 Apache-2.0 的合
同、SDK、示例，以及经批准的社区运行时、适配器和探针二进制。

核心运行时采用 AGPL-3.0。此处发布的二进制保留其来源许可证；本仓库中的合同、
SDK、示例和文档采用 Apache-2.0。企业版专用二进制不会在此发布。

## 仓库模型

| 仓库 | 许可证 | 角色 |
| --- | --- | --- |
| `coordivela` | AGPL-3.0 | 面向单人或小团队的社区协调运行时。参考规模：少于五个并发智能体 |
| `coordivela-open` | Apache-2.0 | 本仓库：公开合同、SDK、示例，以及经批准的核心/适配器/探针二进制 |
| `coordivela-ee` | 商业许可 | 企业治理、多仓库协调、服务身份、多租户规模、HA、联邦、合规和托管运维 |

Core 和 Open 共享版本线；EE 独立发布，并声明其支持的 core 版本范围。智能体数
量示例是范围目标，不是许可证强制计数器。

## 本仓库内容

- `contracts/openapi/coordivela.v1.yaml` - 公开 HTTP 合同。
- `sdk/go/` - Go SDK，覆盖工作项、依赖、生命周期、计划、投影和 SSE。
  SDK module 以 `sdk/go/v3.0.0` 发布。
- `examples/quickstart/` - 端到端 API 示例。
- [`docs/OPERATIONS.zh.md`](docs/OPERATIONS.zh.md) - 安装、配置、PostgreSQL、
  备份、升级、监控和故障排查操作手册。
- `releases/` - 经批准的二进制、清单、校验和、许可证声明和 CycloneDX SBOM。

## 当前版本

`1.7.1` 是当前社区版本，与 core 版本线一致。经批准的二进制、校验和及 SBOM
发布在 [`releases/v1.7.1/`](releases/v1.7.1/)，覆盖 windows/amd64、
linux/amd64 和 darwin/arm64。

Open 是分发面，继承 Core 的社区 1.x 成熟度；它不代表 EE 已生产就绪。Core 的
Go module tag `backend/v3.1.1` 位于 `4e8508a`，覆盖 1.7.x 运行时和更新后的
Harmovela module path；Open 二进制由各发布清单记录的 `v1.7.1` 源码提交构建。

2026-10-24 源级门禁：Go SDK 测试、`go vet` 和 `go build` 通过，1.7.1 签名发布
产物齐全。

已知限制：

- PRAXOVELA 的持久执行计时器证据仍由该产品负责。
- 已发布的 `v1.7.1` 二进制由 Go `1.26.9` 构建，包含 GO-2026-6809；需要用 Go
  `1.26.10` 重新构建并发布补丁版本。
- `RuntimeAdapter` v1 自 `1.2.0` 起包含，`1.7.1` 继续提供。

OpenAPI 合同为 revision `1.7.1`，覆盖已发布的 `/v1` 路由，包括协调计划、
claim、finding、directive、按工作项的 SYMBIVELA 投影和 OpenCode 会话附加。

post-1.0 的排期见 [core 路线图](https://github.com/axisvela/coordivela/blob/main/ROADMAP.md)。

## 快速开始

先启动 core 栈：

```bash
cd ../COORDIVELA
docker compose -f deploy/compose.yaml up --build
```

创建工作项并列出：

```bash
curl -s -X POST http://localhost:1976/v1/work-items \
  -H 'X-Axis-Tenant: local' \
  -H 'X-Axis-Workspace: default' \
  -H 'Idempotency-Key: quickstart-wi-1' \
  -H 'Content-Type: application/json' \
  -d '{"id":"wi-1","title":"Implement durable scheduler persistence"}'

curl -s http://localhost:1976/v1/work-items \
  -H 'X-Axis-Tenant: local' \
  -H 'X-Axis-Workspace: default'
```

## 二进制分发

| 二进制 | 来源 | 用途 |
| --- | --- | --- |
| `coordivela` | core | 协调服务器 |
| `coordivela-migrate` | core | PostgreSQL schema 迁移工具 |
| `coordivela-opencode-probe` | core | 只读 OpenCode 观察探针 |
| 适配器二进制 | core/adapters | PRAXOVELA、OpenCode 和进程适配器在单独打包时发布 |

每个发布清单记录二进制版本、来源仓库和提交、许可证、校验和与 SBOM。索引见
[`releases/`](releases/)。

## 合同

OpenAPI 合同是公开 HTTP 面的唯一事实源。SDK 必须按合同生成或校验。破坏性变更
需要提升合同主版本。当前合同 revision 为 `1.7.1`，与运行时/二进制版本线独立
演进。

## 许可证

本仓库的合同、SDK、示例和文档采用 Apache-2.0。发布的二进制采用 AGPL-3.0。
参见 [LICENSE](LICENSE) 以及每个发布目录中的许可证声明。
