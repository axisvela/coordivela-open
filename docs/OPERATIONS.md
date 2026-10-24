> English | [中文](OPERATIONS.zh.md)

# COORDIVELA Operations Manual

This manual covers the community COORDIVELA Core runtime in the `1.0.0` line:
installation, configuration, PostgreSQL setup, startup, health checks, backups,
upgrades, monitoring, and troubleshooting.

It does not cover enterprise identity and policy, multi-repository coordination,
HA leader election, Harmovela event retention, PRAXOVELA/Rheovela timers, or
SYMBIVELA human streams. Those remain with their owning products or with EE.

## Support Boundary

COORDIVELA Core supports a single human or small team, typically in one
repository, with a reference profile of fewer than five concurrent agents.
Tenant, workspace, and principal identity are supplied by a trusted gateway.
Core does not issue or verify service identities.

## Prerequisites

- PostgreSQL 18.x. The 1.7.1 release is tested against PostgreSQL 18.1.
- A released binary for one of:
  - `windows/amd64`
  - `linux/amd64`
  - `darwin/arm64`
- Port `1976` for the HTTP API. COORDIVELA reserves the `1976-1985` block.
- `pg_dump` and `pg_restore` for backups and restore drills.
- `curl` or another HTTP client for smoke tests.
- Go 1.26.9 only when building from source.
- A trusted gateway or reverse proxy that sets `X-Axis-Tenant`,
  `X-Axis-Workspace`, and optionally `X-Axis-Principal`. Do not accept these
  headers directly from untrusted clients.

## Install From Release Binaries

Download the archive for your platform from
[`releases/v1.7.1/`](../releases/v1.7.1/), or from the GitHub release:

<https://github.com/axisvela/coordivela-open/releases/tag/v1.7.1>

Verify the archive before extracting it.

PowerShell:

```powershell
Get-Content .\SHA256SUMS
Get-FileHash .\coordivela.exe -Algorithm SHA256
```

Linux or macOS:

```sh
sha256sum -c SHA256SUMS
```

Each platform directory contains:

- `coordivela` or `coordivela.exe` - HTTP coordination server.
- `coordivela-migrate` or `coordivela-migrate.exe` - schema migration helper.
- `coordivela-opencode-probe` or `.exe` - read-only OpenCode observation probe.
- `NOTICE` - source commit, build platform, and linked-module provenance.
- `LICENSE-COORDIVELA-CORE.txt` - AGPL-3.0 license.
- `SHA256SUMS` - checksums for the directory contents.

The `1.7.1` binary archives do not embed SQL migrations. Obtain the matching
`backend/migrations` directory from the core tag:

```sh
git clone --branch v1.7.1 https://github.com/axisvela/coordivela.git
```

Use `coordivela/backend/migrations` with `coordivela-migrate`.

## Build From Source

```sh
git clone --branch v1.7.1 https://github.com/axisvela/coordivela.git
cd coordivela/backend
GOWORK=off go build -trimpath -o ../coordivela ./cmd/coordivela
GOWORK=off go build -trimpath -o ../coordivela-migrate ./cmd/coordivela-migrate
GOWORK=off go build -trimpath -o ../coordivela-opencode-probe ./cmd/coordivela-opencode-probe
```

The core module uses local `replace` directives for sibling products during
development. CI checks out those repositories before building. For an isolated
source build, use the same sibling checkouts as the release, or use the
published binaries. The release `NOTICE` records the sibling commits used for
the approved binaries.

## Configure

Command-line flags override environment variables. Run `coordivela -h` for the
complete list.

Core runtime settings:

| Environment variable | Flag | Default | Purpose |
| --- | --- | --- | --- |
| `COORDIVELA_ADDR` | `-addr` | `:1976` | HTTP listen address |
| `COORDIVELA_DSN` | `-dsn` | empty | PostgreSQL DSN; empty uses in-memory stores |
| `COORDIVELA_SUPERVISE` | `-supervise` | `false` | Start the autonomous intervention supervisor |
| `COORDIVELA_STALL_SCAN` | `-stall-scan` | `false` | Start the stall reconciliation scan |
| `COORDIVELA_STALL_TENANT` | `-stall-scan-tenant` | empty | Tenant scope for the scan |
| `COORDIVELA_STALL_WORKSPACE` | `-stall-scan-workspace` | empty | Optional workspace scope for the scan |
| `COORDIVELA_STALL_SCAN_INTERVAL` | `-stall-scan-interval` | `1m` | Scan interval |
| `COORDIVELA_REQUIRE_AUTHORIZATION` | `-require-authorization` | `false` | Fail startup unless AEGIVELA authorization is configured |

Optional runtime and integration settings:

| Environment variable | Flag | Purpose |
| --- | --- | --- |
| `PRAXOVELA_URL` | `-praxovela-url` | PRAXOVELA base URL; empty disables the adapter |
| `PRAXOVELA_TOKEN` | `-praxovela-token` | PRAXOVELA bearer token |
| `COORDIVELA_OPENCODE_URL` or `OPENCODE_SERVER_URL` | `-opencode-url` | OpenCode server URL; empty disables the control adapter |
| `HARMOVELA_WS_URL` | `-harmovela-ws` | Harmovela WebSocket URL; empty disables observation consumption |
| `HARMOVELA_TOPIC` | `-harmovela-topic` | Observation topic; default `opencode.observation` |
| `HARMOVELA_SIGNAL_TOPIC` | `-harmovela-signal-topic` | Optional timer-owner signal topic |
| `HARMOVELA_TENANT_ID` | `-harmovela-tenant` | Optional tenant filter |
| `HARMOVELA_WORKSPACE_ID` | `-harmovela-workspace` | Optional workspace filter |
| `TEKMOVELA_URL` | `-tekmovela-url` | Verification gate base URL |
| `TEKMOVELA_TOKEN` | `-tekmovela-token` | Optional bearer token |
| `COORDIVELA_REQUIRE_VERIFICATION` | `-require-verification` | Fail startup unless the verification gate is configured |
| `RHEOVELA_URL` | `-rheovela-url` | Process gate base URL |
| `RHEOVELA_TOKEN` | `-rheovela-token` | Optional bearer token |
| `COORDIVELA_REQUIRE_PROCESS_GATE` | `-require-process-gate` | Fail startup unless the process gate is configured |
| `ONTOVELA_URL` | `-ontovela-url` | Reality gate base URL |
| `ONTOVELA_TOKEN` | `-ontovela-token` | Optional bearer token |
| `ONTOVELA_PURPOSE` | `-ontovela-purpose` | Reality View purpose; default `coordivela.complete` |
| `ONTOVELA_REQUIRED_STATE` | `-ontovela-required-state` | Comma-separated `property:maxAgeSeconds` requirements |
| `COORDIVELA_REQUIRE_REALITY_GATE` | `-require-reality-gate` | Fail startup unless the reality gate is configured |
| `SYMBIVELA_URL` | `-symbivela-url` | Decision projection base URL |
| `SYMBIVELA_ACTOR` | `-symbivela-actor` | Actor for projected decisions; default `coordivela` |
| `MNEMOVELA_URL` | `-mnemovela-url` | Decision capture base URL |
| `MNEMOVELA_BRANCH` | `-mnemovela-branch` | Optional branch; defaults to the work-item workspace |
| `GNOSIVELA_URL` | `-gnosivela-url` | Grounding base URL |
| `GNOSIVELA_PURPOSE` | `-gnosivela-purpose` | Grounding purpose; default `coordivela.complete` |
| `GNOSIVELA_BUDGET` | `-gnosivela-budget` | Grounding context budget; default `800` |

Authorization settings such as `AEGIVELA_URL`, `AEGIVELA_AUDIENCE`,
`AEGIVELA_SCOPE`, `AEGIVELA_CLIENT_ID`, `AEGIVELA_WORKLOAD_REF`, and supervisor
assertions are described in the core
[security guide](https://github.com/axisvela/coordivela/blob/main/docs/08-security-and-governance.md).

Pass secrets through a secret manager or an environment file with restricted
permissions. Do not put tokens, DSNs, workload assertions, or prompts in tracked
files, events, or logs.

Minimal production example:

```dotenv
COORDIVELA_ADDR=:1976
COORDIVELA_DSN=postgres://coordivela:CHANGE_ME@db:5432/coordivela?sslmode=require
COORDIVELA_STALL_SCAN=true
COORDIVELA_STALL_TENANT=acme
COORDIVELA_STALL_WORKSPACE=default
COORDIVELA_STALL_SCAN_INTERVAL=1m
```

## Prepare PostgreSQL

Create a dedicated role and database:

```sql
CREATE ROLE coordivela LOGIN PASSWORD 'CHANGE_ME';
CREATE DATABASE coordivela OWNER coordivela;
```

Apply migrations with the release binary and the matching migration directory:

```sh
./coordivela-migrate \
  -dsn "postgres://coordivela:CHANGE_ME@db:5432/coordivela?sslmode=require" \
  -dir ./coordivela/backend/migrations
```

Or set `COORDIVELA_DSN` and run:

```sh
./coordivela-migrate -dir ./coordivela/backend/migrations
```

Migrations are append-only and ordered lexically. Each file runs in a
transaction together with its bookkeeping row. A PostgreSQL session advisory
lock prevents concurrent migration runners. Applied file names are recorded in
`schema_migrations`; later runs skip them.

Confirm the result:

```sql
SELECT name, applied_at FROM schema_migrations ORDER BY name;
```

## Start The Server

```sh
export COORDIVELA_DSN="postgres://coordivela:CHANGE_ME@db:5432/coordivela?sslmode=require"
./coordivela -addr :1976
```

On `SIGINT` or `SIGTERM`, the HTTP server stops accepting connections and drains
for up to ten seconds.

Without `COORDIVELA_DSN`, COORDIVELA uses in-memory stores. That mode is useful
for a smoke test only; work items, schedules, idempotency records, and consumer
state are lost when the process exits.

### systemd Example

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

Run `coordivela-migrate` with the new binary before starting or restarting the
service during an upgrade.

### Docker Compose

The core repository includes a local compose profile:

```sh
git clone --branch v1.7.1 https://github.com/axisvela/coordivela.git
cd coordivela
docker compose -f deploy/compose.yaml up --build
```

This profile starts PostgreSQL and the API. It is a local/development profile,
not a production orchestrator. No official container image is published for
`1.7.1`; build from source or use the released binaries.

## Verify The Installation

Liveness and version:

```sh
curl -s http://127.0.0.1:1976/healthz
```

Expected:

```json
{"status":"ok","version":"1.7.1"}
```

Readiness:

```sh
curl -s http://127.0.0.1:1976/readyz
```

Expected:

```json
{"status":"ready"}
```

`/readyz` pings PostgreSQL when a DSN is configured. Without a DSN it reports
ready because there is no database dependency.

Metrics:

```sh
curl -s http://127.0.0.1:1976/metrics
```

Create and list a work item:

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

## Request Conventions

| Header or field | Required | Purpose |
| --- | --- | --- |
| `X-Axis-Tenant` | Yes on `/v1` routes | Tenant scope |
| `X-Axis-Workspace` | Yes on workspace-scoped routes | Workspace scope |
| `X-Axis-Principal` | For `projections/my-work`; optional actor elsewhere | Principal identity from the gateway |
| `X-Axis-Profile` | No | Scenario profile fallback |
| `Idempotency-Key` | Yes on mutating methods | Replay-safe mutation |
| `If-Match` or `expected_version` | For state transitions | Optimistic version check |
| `X-Aegivela-Workload-Assertion` | When authorization is enabled | AEGIVELA `system_api` assertion |

All request bodies are capped at 1 MiB. The SSE endpoint at `GET /v1/events`
admits at most 64 subscribers per process and requires both tenant and workspace
headers.

## Backups And Restore

Create a PostgreSQL custom-format backup:

```sh
pg_dump --format=custom --file coordivela.dump "$COORDIVELA_DSN"
```

Restore only into an isolated database:

```sh
pg_restore --exit-on-error --no-owner --no-privileges \
  --dbname "$RESTORE_DSN" coordivela.dump
```

The core repository also provides PowerShell helpers:

```powershell
./deploy/backup-postgres.ps1 -OutputPath backups/coordivela.dump
./deploy/restore-postgres.ps1 -InputPath backups/coordivela.dump -ConfirmRestore
```

After a restore, run migrations, start the server, confirm `/readyz`, and read a
tenant-scoped work item. Store backups in approved encrypted storage and record
the backup checksum, target environment, and drill result.

## Upgrade

1. Back up the database.
2. Stop or drain traffic from the old process.
3. Replace the binary.
4. Apply migrations with the new `coordivela-migrate` and the matching
   `backend/migrations` directory.
5. Start the new process with the same `COORDIVELA_DSN`.
6. Confirm `/healthz` reports the expected version and `/readyz` succeeds.
7. Route traffic back and check `/metrics` plus a tenant-scoped read.

Migration files are append-only. Do not edit or reorder an applied migration.

## Monitoring

Scrape `GET /metrics`. The main metrics are:

- `coordivela_events_total{type="..."}` - coordination events by type.
- `coordivela_interventions_total{kind="..."}` - interventions by kind.
- `coordivela_intervention_latency_seconds` - stall-signal to applied
  intervention latency histogram.
- `coordivela_event_stream_subscribers` - active SSE subscribers.
- `coordivela_build_info{version="..."}` - build version.

The endpoint contains no tenant or work-item identifiers. Alerting, retention,
and external SLO enforcement are deployment responsibilities.

## Stall Scan And Supervisor

Enable the reconciliation scan with `-stall-scan` or
`COORDIVELA_STALL_SCAN=true`. Set `COORDIVELA_STALL_TENANT` and, when required,
`COORDIVELA_STALL_WORKSPACE`. The scan runs at startup and then on
`COORDIVELA_STALL_SCAN_INTERVAL`, emitting `work.stalled` once per item version.

Enable the autonomous intervention supervisor with `-supervise` or
`COORDIVELA_SUPERVISE=true`. When authorization is enabled, configure a
supervisor workload assertion or assertion file. Supervisor actions are
evaluated and audited through AEGIVELA.

Durable timer recovery remains owned by Rheovela and PRAXOVELA. Harmovela
remains the durable event spine and owns replay and dead-letter handling.

## OpenCode Probe

The probe is read-only and does not send prompts, approvals, or aborts.

```sh
./coordivela-opencode-probe -url http://127.0.0.1:4096 -health
./coordivela-opencode-probe -url http://127.0.0.1:4096 -list
./coordivela-opencode-probe -url http://127.0.0.1:4096 -session SESSION_ID -json
```

Publish observations to Harmovela with `-harmovela ws://127.0.0.1:8787/harmovela`
and an optional durable outbox with `-outbox /var/lib/coordivela/outbox`.
Credentials and raw prompts must not be published.

## Troubleshooting

| Symptom | Likely cause | Action |
| --- | --- | --- |
| `400 X-Axis-Tenant required` | Missing trusted scope header | Set the header at the gateway; do not accept it from untrusted clients |
| `400 Idempotency-Key required` | Mutating request without a key | Send a stable key for the logical operation |
| `401 unauthenticated` | Authorization enabled and assertion missing or invalid | Supply `X-Aegivela-Workload-Assertion` |
| `403 authorization_denied` | AEGIVELA denied the action | Check the policy, action, resource, and trace ID |
| `404` on a known item | Wrong tenant/workspace or cross-tenant access | Check scope headers; 404 is intentional for isolation |
| `409` on a transition | Stale version, lease conflict, invalid transition, or in-progress idempotency key | Re-read the item, update `expected_version`, or retry after the first request completes |
| `409 event replay unavailable` | `Last-Event-ID` is outside the in-memory window | Reconnect without `Last-Event-ID` or replay through Harmovela |
| `400 invalid json` | Body exceeds 1 MiB or is malformed | Send a smaller valid JSON body or use an external reference |
| `429 event stream busy` | More than 64 SSE subscribers in one process | Reduce subscribers or scale carefully; retry with `Retry-After` |
| `503` on `/readyz` | PostgreSQL ping failed | Check DSN, network, credentials, and database health |
| `503 authorization_unavailable` | Authorization service unreachable | Fail closed; restore AEGIVELA availability |
| Migration hangs | Another runner holds the advisory lock | Wait for the other runner; do not kill it mid-transaction |
| Data disappears after restart | Running without `COORDIVELA_DSN` | Configure PostgreSQL for durable state |

## Security Checklist

- Put COORDIVELA behind a trusted gateway that sets tenant, workspace, and
  principal headers.
- Use TLS for client, database, and integration traffic.
- Store DSNs and tokens in a secret manager.
- Enable AEGIVELA authorization for mutating routes when available.
- Restrict `/metrics` and `/healthz` to the monitoring network where possible.
- Run as a non-root user and restrict write access.
- Back up PostgreSQL and test restore into an isolated database.
- Do not put credentials or raw prompts in work-item descriptions, events, or
  logs.

## Known Limitations

- PRAXOVELA durable execution-timer evidence remains outstanding in that
  product.
- `RuntimeAdapter` v1 has been part of the release line since `1.1.0` and is
  included in `1.7.1`.
- Multi-writer leases, HA leader election, multi-repository coordination,
  service identity, federation, and enterprise policy are not core features.
