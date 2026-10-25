> [English](OPERATIONS.md) | 中文

# COORDIVELA 操作手册

本手册覆盖社区版 COORDIVELA Core `1.7.2` 版本线：安装、配置、PostgreSQL
准备、启动、健康检查、备份、升级、监控和故障排查。

本手册不覆盖企业身份与策略、多仓库协调、HA leader election、Harmovela 事件
保留、PRAXOVELA/Rheovela 计时器或 SYMBIVELA 人类流。这些能力属于各自产品
或 EE。

## 支持边界

COORDIVELA Core 面向单人或小团队，通常是一个仓库，参考规模为少于五个并发
智能体。租户、工作区和主体身份由受信任网关提供。Core 不签发或验证服务身份。

## 前置条件

- PostgreSQL 18.x。`1.7.2` 使用 PostgreSQL 18.1 验证。
- 以下平台之一的发布二进制：
  - `windows/amd64`
  - `linux/amd64`
  - `darwin/arm64`
- HTTP API 端口 `1976`。COORDIVELA 保留 `1976-1985` 端口段。
- 备份和恢复演练需要 `pg_dump` 与 `pg_restore`。
- 冒烟测试需要 `curl` 或其他 HTTP 客户端。
- 仅从源码构建时需要 Go 1.26.10。
- 由受信任网关或反向代理设置 `X-Axis-Tenant`、`X-Axis-Workspace`，以及可选
  的 `X-Axis-Principal`。不要直接信任不可信客户端传入的这些头。

## 从发布二进制安装

从 [`releases/v1.7.2/`](../releases/v1.7.2/) 或 GitHub Release 下载对应平台
压缩包：

<https://github.com/axisvela/coordivela-open/releases/tag/v1.7.2>

解压前先校验。

PowerShell：

```powershell
Get-Content .\SHA256SUMS
Get-FileHash .\coordivela.exe -Algorithm SHA256
```

Linux 或 macOS：

```sh
sha256sum -c SHA256SUMS
```

每个平台目录包含：

- `coordivela` 或 `coordivela.exe` - HTTP 协调服务器。
- `coordivela-migrate` 或 `coordivela-migrate.exe` - schema 迁移工具。
- `coordivela-opencode-probe` 或 `.exe` - 只读 OpenCode 观察探针。
- `NOTICE` - 源码提交、构建平台和链接模块来源。
- `LICENSE-COORDIVELA-CORE.txt` - AGPL-3.0 许可证。
- `SHA256SUMS` - 目录内容的校验和。

`1.7.2` 二进制压缩包不内嵌 SQL 迁移。请从 core tag 获取匹配的
`backend/migrations` 目录：

```sh
git clone --branch v1.7.2 https://github.com/axisvela/coordivela.git
```

将 `coordivela/backend/migrations` 交给 `coordivela-migrate` 使用。

## 从源码构建

```sh
git clone --branch v1.7.2 https://github.com/axisvela/coordivela.git
cd coordivela/backend
GOWORK=off go build -trimpath -o ../coordivela ./cmd/coordivela
GOWORK=off go build -trimpath -o ../coordivela-migrate ./cmd/coordivela-migrate
GOWORK=off go build -trimpath -o ../coordivela-opencode-probe ./cmd/coordivela-opencode-probe
```

Core 模块在开发期通过本地 `replace` 指向 sibling 产品。CI 会在构建前检出这些
仓库。若要在隔离环境从源码构建，请使用与发布一致的 sibling 检出，或直接使用
已发布二进制。发布 `NOTICE` 记录了批准二进制所用的 sibling 提交。

## 配置

命令行标志优先于环境变量。完整列表运行 `coordivela -h`。

Core 运行时设置：

| 环境变量 | 标志 | 默认值 | 用途 |
| --- | --- | --- | --- |
| `COORDIVELA_ADDR` | `-addr` | `:1976` | HTTP 监听地址 |
| `COORDIVELA_DSN` | `-dsn` | 空 | PostgreSQL DSN；为空时使用内存存储 |
| `COORDIVELA_SUPERVISE` | `-supervise` | `false` | 启动自主干预 supervisor |
| `COORDIVELA_STALL_SCAN` | `-stall-scan` | `false` | 启动停滞协调扫描 |
| `COORDIVELA_STALL_TENANT` | `-stall-scan-tenant` | 空 | 扫描的租户范围 |
| `COORDIVELA_STALL_WORKSPACE` | `-stall-scan-workspace` | 空 | 扫描的可选工作区范围 |
| `COORDIVELA_STALL_SCAN_INTERVAL` | `-stall-scan-interval` | `1m` | 扫描间隔 |
| `COORDIVELA_REQUIRE_AUTHORIZATION` | `-require-authorization` | `false` | 未配置 AEGIVELA 授权时启动失败 |

可选的运行时和集成设置：

| 环境变量 | 标志 | 用途 |
| --- | --- | --- |
| `PRAXOVELA_URL` | `-praxovela-url` | PRAXOVELA 基础 URL；为空时禁用适配器 |
| `PRAXOVELA_TOKEN` | `-praxovela-token` | PRAXOVELA bearer token |
| `COORDIVELA_OPENCODE_URL` 或 `OPENCODE_SERVER_URL` | `-opencode-url` | OpenCode 服务器 URL；为空时禁用控制适配器 |
| `HARMOVELA_WS_URL` | `-harmovela-ws` | Harmovela WebSocket URL；为空时禁用观察消费 |
| `HARMOVELA_TOPIC` | `-harmovela-topic` | 观察主题；默认 `opencode.observation` |
| `HARMOVELA_SIGNAL_TOPIC` | `-harmovela-signal-topic` | 可选的计时器所有者信号主题 |
| `HARMOVELA_TENANT_ID` | `-harmovela-tenant` | 可选租户过滤 |
| `HARMOVELA_WORKSPACE_ID` | `-harmovela-workspace` | 可选工作区过滤 |
| `TEKMOVELA_URL` | `-tekmovela-url` | 验证门禁基础 URL |
| `TEKMOVELA_TOKEN` | `-tekmovela-token` | 可选 bearer token |
| `COORDIVELA_REQUIRE_VERIFICATION` | `-require-verification` | 未配置验证门禁时启动失败 |
| `RHEOVELA_URL` | `-rheovela-url` | 流程门禁基础 URL |
| `RHEOVELA_TOKEN` | `-rheovela-token` | 可选 bearer token |
| `COORDIVELA_REQUIRE_PROCESS_GATE` | `-require-process-gate` | 未配置流程门禁时启动失败 |
| `ONTOVELA_URL` | `-ontovela-url` | 现实门禁基础 URL |
| `ONTOVELA_TOKEN` | `-ontovela-token` | 可选 bearer token |
| `ONTOVELA_PURPOSE` | `-ontovela-purpose` | Reality View 用途；默认 `coordivela.complete` |
| `ONTOVELA_REQUIRED_STATE` | `-ontovela-required-state` | 逗号分隔的 `property:maxAgeSeconds` 要求 |
| `COORDIVELA_REQUIRE_REALITY_GATE` | `-require-reality-gate` | 未配置现实门禁时启动失败 |
| `SYMBIVELA_URL` | `-symbivela-url` | 决策投影基础 URL |
| `SYMBIVELA_ACTOR` | `-symbivela-actor` | 投影决策的 actor；默认 `coordivela` |
| `MNEMOVELA_URL` | `-mnemovela-url` | 决策捕获基础 URL |
| `MNEMOVELA_BRANCH` | `-mnemovela-branch` | 可选分支；默认使用工作项工作区 |
| `GNOSIVELA_URL` | `-gnosivela-url` | 语义支撑基础 URL |
| `GNOSIVELA_PURPOSE` | `-gnosivela-purpose` | 语义支撑用途；默认 `coordivela.complete` |
| `GNOSIVELA_BUDGET` | `-gnosivela-budget` | 语义支撑上下文预算；默认 `800` |

授权设置（如 `AEGIVELA_URL`、`AEGIVELA_AUDIENCE`、`AEGIVELA_SCOPE`、
`AEGIVELA_CLIENT_ID`、`AEGIVELA_WORKLOAD_REF` 和 supervisor assertion）见
core 的
[安全指南](https://github.com/axisvela/coordivela/blob/main/docs/08-security-and-governance.md)。

密钥应通过密钥管理系统或权限受限的环境文件传入。不要把 token、DSN、workload
assertion 或 prompt 写入受版本控制的文件、事件或日志。

最小生产示例：

```dotenv
COORDIVELA_ADDR=:1976
COORDIVELA_DSN=postgres://coordivela:CHANGE_ME@db:5432/coordivela?sslmode=require
COORDIVELA_STALL_SCAN=true
COORDIVELA_STALL_TENANT=acme
COORDIVELA_STALL_WORKSPACE=default
COORDIVELA_STALL_SCAN_INTERVAL=1m
```

## 准备 PostgreSQL

创建专用角色和数据库：

```sql
CREATE ROLE coordivela LOGIN PASSWORD 'CHANGE_ME';
CREATE DATABASE coordivela OWNER coordivela;
```

使用发布二进制和匹配的迁移目录执行迁移：

```sh
./coordivela-migrate \
  -dsn "postgres://coordivela:CHANGE_ME@db:5432/coordivela?sslmode=require" \
  -dir ./coordivela/backend/migrations
```

或设置 `COORDIVELA_DSN` 后运行：

```sh
./coordivela-migrate -dir ./coordivela/backend/migrations
```

迁移文件只追加、按字典序执行。每个文件与其记录行在同一事务中执行。PostgreSQL
会话级 advisory lock 防止并发迁移进程。已执行文件名记录在
`schema_migrations`，后续运行会跳过。

确认结果：

```sql
SELECT name, applied_at FROM schema_migrations ORDER BY name;
```

## 启动服务

```sh
export COORDIVELA_DSN="postgres://coordivela:CHANGE_ME@db:5432/coordivela?sslmode=require"
./coordivela -addr :1976
```

收到 `SIGINT` 或 `SIGTERM` 时，HTTP 服务停止接收新连接，并在最多十秒内完成
排空。

未设置 `COORDIVELA_DSN` 时，COORDIVELA 使用内存存储。该模式仅适合冒烟测试；
进程退出后工作项、计划、幂等记录和消费者状态都会丢失。

### systemd 示例

```ini
[Unit]
Description=COORDIVELA coordination server
After=network-online.target postgresql.service
Wants=network-online.target

[Service]
User=coordivela
EnvironmentFile=/etc/coordivela/coordivela.env
ExecStart=/opt/coordivela/coordivela -addr :1976
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
```

升级时，在启动或重启服务前使用新二进制运行 `coordivela-migrate`。

### Docker Compose

Core 仓库提供本地 compose profile：

```sh
git clone --branch v1.7.2 https://github.com/axisvela/coordivela.git
cd coordivela
docker compose -f deploy/compose.yaml up --build
```

该 profile 启动 PostgreSQL 和 API。它面向本地/开发，不是生产编排方案。
`1.7.2` 没有发布官方容器镜像；请从源码构建或使用发布二进制。

## 验证安装

存活和版本：

```sh
curl -s http://127.0.0.1:1976/healthz
```

预期：

```json
{"status":"ok","version":"1.7.2"}
```

就绪：

```sh
curl -s http://127.0.0.1:1976/readyz
```

预期：

```json
{"status":"ready"}
```

配置 DSN 时，`/readyz` 会 ping PostgreSQL。未配置 DSN 时没有数据库依赖，因此
直接报告 ready。

指标：

```sh
curl -s http://127.0.0.1:1976/metrics
```

创建工作项并列出：

```sh
curl -s -X POST http://127.0.0.1:1976/v1/work-items \
  -H 'X-Axis-Tenant: local' \
  -H 'X-Axis-Workspace: default' \
  -H 'X-Axis-Principal: operator' \
  -H 'Idempotency-Key: smoke-1' \
  -H 'Content-Type: application/json' \
  -d '{"id":"wi-1","title":"Smoke test"}'

curl -s http://127.0.0.1:1976/v1/work-items \
  -H 'X-Axis-Tenant: local' \
  -H 'X-Axis-Workspace: default'
```

## 请求约定

| 头或字段 | 是否必需 | 用途 |
| --- | --- | --- |
| `X-Axis-Tenant` | `/v1` 路由必需 | 租户范围 |
| `X-Axis-Workspace` | 工作区范围路由必需 | 工作区范围 |
| `X-Axis-Principal` | `projections/my-work` 必需；其他路由可选 actor | 网关提供的主体身份 |
| `X-Axis-Profile` | 否 | 场景 profile 回退 |
| `Idempotency-Key` | 变更方法必需 | 可重放的变更 |
| `If-Match` 或 `expected_version` | 状态转换需要 | 乐观版本检查 |
| `X-Aegivela-Workload-Assertion` | 启用授权时必需 | AEGIVELA `system_api` assertion |

所有请求体上限为 1 MiB。`GET /v1/events` 的 SSE 端点每进程最多允许 64 个
订阅者，并要求同时提供租户和工作区头。

## 备份与恢复

创建 PostgreSQL custom-format 备份：

```sh
pg_dump --format=custom --file coordivela.dump "$COORDIVELA_DSN"
```

只恢复到隔离数据库：

```sh
pg_restore --exit-on-error --no-owner --no-privileges \
  --dbname "$RESTORE_DSN" coordivela.dump
```

Core 仓库也提供 PowerShell 辅助脚本：

```powershell
./deploy/backup-postgres.ps1 -OutputPath backups/coordivela.dump
./deploy/restore-postgres.ps1 -InputPath backups/coordivela.dump -ConfirmRestore
```

恢复后运行迁移、启动服务、确认 `/readyz`，并读取一个租户范围内的工作项。
将备份存入经批准的加密存储，并记录备份校验和、目标环境和演练结果。

## 升级

1. 备份数据库。
2. 停止旧进程或排空流量。
3. 替换二进制。
4. 使用新的 `coordivela-migrate` 和匹配的 `backend/migrations` 目录执行迁移。
5. 使用相同 `COORDIVELA_DSN` 启动新进程。
6. 确认 `/healthz` 报告预期版本，`/readyz` 成功。
7. 恢复流量，检查 `/metrics` 和一次租户范围读取。

迁移文件只追加。不要编辑或重排已应用的迁移。

## 监控

抓取 `GET /metrics`。主要指标：

- `coordivela_events_total{type="..."}` - 按类型统计的协调事件。
- `coordivela_interventions_total{kind="..."}` - 按类型统计的干预。
- `coordivela_intervention_latency_seconds` - 从停滞信号到已应用干预的延迟
  直方图。
- `coordivela_event_stream_subscribers` - 活跃 SSE 订阅者。
- `coordivela_build_info{version="..."}` - 构建版本。

该端点不包含租户或工作项标识。告警、保留和外部 SLO 执行由部署方负责。

## 停滞扫描与 Supervisor

使用 `-stall-scan` 或 `COORDIVELA_STALL_SCAN=true` 启用协调扫描。设置
`COORDIVELA_STALL_TENANT`，需要时设置 `COORDIVELA_STALL_WORKSPACE`。扫描在
启动时运行，之后按 `COORDIVELA_STALL_SCAN_INTERVAL` 运行，每个工作项版本只
发出一次 `work.stalled`。

使用 `-supervise` 或 `COORDIVELA_SUPERVISE=true` 启用自主干预 supervisor。
启用授权时，配置 supervisor workload assertion 或 assertion 文件。supervisor
动作会经 AEGIVELA 评估并审计。

持久计时器恢复仍由 Rheovela 和 PRAXOVELA 负责。Harmovela 仍是持久事件主
干，负责重放和 dead-letter。

## OpenCode 探针

探针为只读，不发送 prompt、审批或中止操作。

```sh
./coordivela-opencode-probe -url http://127.0.0.1:4096 -health
./coordivela-opencode-probe -url http://127.0.0.1:4096 -list
./coordivela-opencode-probe -url http://127.0.0.1:4096 -session SESSION_ID -json
```

使用 `-harmovela ws://127.0.0.1:8787/harmovela` 向 Harmovela 发布观察，使用
`-outbox /var/lib/coordivela/outbox` 启用可选持久 outbox。不要发布凭据或原始
prompt。

## 故障排查

| 症状 | 可能原因 | 处理 |
| --- | --- | --- |
| `400 X-Axis-Tenant required` | 缺少受信任范围头 | 在网关设置该头；不要信任不可信客户端 |
| `400 Idempotency-Key required` | 变更请求缺少 key | 为同一逻辑操作发送稳定 key |
| `401 unauthenticated` | 已启用授权但 assertion 缺失或无效 | 提供 `X-Aegivela-Workload-Assertion` |
| `403 authorization_denied` | AEGIVELA 拒绝操作 | 检查策略、动作、资源和 trace ID |
| 已知工作项返回 `404` | 租户/工作区错误或跨租户访问 | 检查范围头；404 是隔离设计 |
| 状态转换返回 `409` | 版本过期、租约冲突、非法转换或幂等 key 正在处理 | 重新读取工作项、更新 `expected_version`，或等待首个请求完成 |
| `409 event replay unavailable` | `Last-Event-ID` 超出内存窗口 | 去掉 `Last-Event-ID` 重连，或通过 Harmovela 重放 |
| `400 invalid json` | 请求体超过 1 MiB 或 JSON 格式错误 | 发送更小的合法 JSON 请求体或使用外部引用 |
| `429 event stream busy` | 单进程 SSE 订阅者超过 64 | 减少订阅者或谨慎扩容；按 `Retry-After` 重试 |
| `/readyz` 返回 `503` | PostgreSQL ping 失败 | 检查 DSN、网络、凭据和数据库健康 |
| `503 authorization_unavailable` | 授权服务不可达 | 保持 fail-closed；恢复 AEGIVELA |
| 迁移挂起 | 另一个迁移进程持有 advisory lock | 等待其完成；不要在事务中途终止 |
| 重启后数据消失 | 未设置 `COORDIVELA_DSN` | 配置 PostgreSQL 以获得持久状态 |

## 安全检查清单

- 将 COORDIVELA 置于受信任网关之后，由网关设置租户、工作区和主体头。
- 客户端、数据库和集成流量使用 TLS。
- DSN 和 token 存入密钥管理系统。
- 可用时对变更路由启用 AEGIVELA 授权。
- 尽可能将 `/metrics` 和 `/healthz` 限制在监控网络。
- 以非 root 用户运行并限制写权限。
- 备份 PostgreSQL，并演练恢复到隔离数据库。
- 不要把凭据或原始 prompt 写入工作项描述、事件或日志。

## 已知限制

- PRAXOVELA 的持久执行计时器证据仍由该产品负责。
- `RuntimeAdapter` v1 自 `1.2.0` 起包含，`1.7.2` 继续提供。
- 多写者租约、HA leader election、多仓库协调、服务身份、联邦和企业策略不是
  core 功能。
