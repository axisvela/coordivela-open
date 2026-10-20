> English | [中文](README.zh.md)

# COORDIVELA Open

**Public integration surface for the COORDIVELA Development Coordination Control Plane**

COORDIVELA is the **Development Coordination Control Plane (DCCP)**. It
coordinates humans, coding agents, and development work, typically within one
repository. It decides what should happen next, assigns work, tracks progress,
detects stalls, requests intervention, and records coordination evidence.

COORDIVELA does not execute code, verify correctness, or own product runtimes.
Those responsibilities remain with the products that own them.

## What COORDIVELA Solves

Development work is increasingly concurrent: several people and agents may
touch the same repository, dependencies, schedules, and delivery gates at the
same time. Without a coordination layer, teams see duplicate claims, lost
updates, stale ownership, late dependency discovery, invisible stalled
sessions, and coordination decisions with no audit trail.

COORDIVELA makes that coordination explicit and durable. Work items have
identity, state, owners, leases, and versions. Dependencies and path conflicts
are checked before work proceeds. Schedules survive restarts. Runtime adapters
connect the coordination loop to execution systems without taking over
execution. Interventions are authorized when an authorization adapter is
configured, and coordination events can be observed and replayed.

## What Value It Delivers

- **One coordination source of truth** for work state, ownership, progress,
  dependencies, schedules, and interventions.
- **Safe concurrency** through leases, fencing tokens, optimistic versioning,
  dependency checks, and path-overlap conflict warnings.
- **Durable schedules** with cron, interval, and one-shot semantics,
  event-triggered dispatch, timezone/DST handling, and misfire policy.
- **Runtime coordination** through adapters for PRAXOVELA, process supervision,
  and OpenCode observation and control.
- **Cross-product boundaries** through public contracts: TEKMOVELA verification
  gates, Rheovela process gates, ONTOVELA reality gates, Mnemovela and
  SYMBIVELA recording, Harmovela events, and GNOSIVELA grounding.
- **Governable operations** with optional AEGIVELA authorization, decision
  audit, idempotency keys, and replay-safe state transitions.
- **Observability** through scoped SSE coordination events, bounded history,
  read projections, and Prometheus metrics.
- **Community-scale deployment** with a single binary, memory or PostgreSQL
  stores, and a reference profile of fewer than five concurrent agents.

## What COORDIVELA Does Not Own

| Concern | Owning product |
| --- | --- |
| Code execution and agent loop | PRAXOVELA |
| Verification, assurance, and release gates | TEKMOVELA |
| Durable process engine | Rheovela |
| Memory storage and retrieval | Mnemovela |
| Human collaboration UI | SYMBIVELA |
| Events and asynchronous coordination protocol | Harmovela |
| Semantic and knowledge fabric | GNOSIVELA |
| Identity, authorization, and approvals | AEGIVELA / Archon Guard |

COORDIVELA is a coordinator, not a control authority. Interventions that change
execution are performed by the owning runtime.

## COORDIVELA Open

COORDIVELA Open is the **public distribution and integration surface** for
COORDIVELA. It contains Apache-2.0 contracts, SDKs, examples, and approved
binaries for the community runtime, adapters, and probes.

The core runtime is AGPL-3.0. Binaries published here retain the license of
their source; the contracts, SDKs, examples, and documentation in this
repository are Apache-2.0. Enterprise-only binaries are never published here.

## Repository Model

| Repository | License | Role |
| --- | --- | --- |
| `coordivela` | AGPL-3.0 | Community coordination runtime for a single human or small team. Reference profile: fewer than five concurrent agents |
| `coordivela-open` | Apache-2.0 | This repository: public contracts, SDKs, examples, and approved core/adapter/probe binaries |
| `coordivela-ee` | Commercial | Enterprise governance, multi-repository coordination, service identity, multi-tenant scale, HA, federation, compliance, and managed operations |

Core and Open share a version line; EE releases independently and states the
core range it supports. The agent-count example is a scope target, not a
license-enforced counter.

## What Is In This Repository

- `contracts/openapi/coordivela.v1.yaml` - public HTTP contract.
- `contracts/runtime-adapter/v1/` - RuntimeAdapter v1 semantics and JSON Schema.
- `sdk/go/` - Go SDK covering work items, dependencies, lifecycle, schedules,
  projections, SSE, and the RuntimeAdapter v1 contract/conformance kit.
- `examples/quickstart/` - end-to-end API example.
- `examples/runtime-adapter/` - external RuntimeAdapter v1 example that passes
  the public conformance kit.
- [`docs/OPERATIONS.md`](docs/OPERATIONS.md) - installation, configuration,
  PostgreSQL, backup, upgrade, monitoring, and troubleshooting.
- `releases/` - approved binaries, manifests, checksums, license notices, and
  CycloneDX SBOMs.

## Current Release

`1.3.0` is the current release, matching the core version line. Signed
binaries, checksums, and SBOMs are published under
[`releases/v1.3.0/`](releases/v1.3.0/) for windows/amd64, linux/amd64, and
darwin/arm64.

Known limitations:

- PRAXOVELA durable execution-timer evidence remains outstanding in that
  product.

The OpenAPI contract is revision `1.3.0` and covers the published `/v1` routes.

The RuntimeAdapter v1 contract, Go interface, JSON Schema, and conformance kit
are included in the `v1.1.0` through `v1.3.0` source releases under
`contracts/runtime-adapter/v1/`, `sdk/go/adapter/`, and
`sdk/go/adapter/conformance/`.

See the [core roadmap](https://github.com/axisvela/coordivela/blob/main/ROADMAP.md)
for post-1.0 sequencing.

## Quickstart

Start the core stack first:

```bash
cd ../COORDIVELA
docker compose -f deploy/compose.yaml up --build
```

Create and list a work item:

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

## Binary Distribution

| Binary | Source | Purpose |
| --- | --- | --- |
| `coordivela` | core | Coordination server |
| `coordivela-migrate` | core | PostgreSQL schema migration helper |
| `coordivela-opencode-probe` | core | Read-only OpenCode observation and probe |
| Adapter binaries | core/adapters | PRAXOVELA, OpenCode, and process adapters when packaged separately |

Each release manifest records the binary version, source repository and commit,
license, checksums, and SBOM. See [`releases/`](releases/) for the index.

## Contracts

The OpenAPI contract is the source of truth for the public HTTP surface. SDKs
must be generated or checked against it. Breaking changes require a major
contract version. The current contract revision is `1.0.0`; it is versioned
independently from the runtime and binary release line.

## License

Apache-2.0 for this repository's contracts, SDKs, examples, and documentation.
Published binaries are AGPL-3.0. See [LICENSE](LICENSE) and the license notice
inside each release directory.
