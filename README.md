# COORDIVELA Open

**Public integration surface for the Development Coordination Control Plane**

Apache-2.0 public contracts, SDKs, examples, and the approved binary
distribution surface for COORDIVELA.

> COORDIVELA coordinates development work across humans and agents. It does not
> execute code or verify correctness; those remain with PRAXOVELA and TEKMOVELA.

## Repository Model

| Repository | License | Role |
| --- | --- | --- |
| `coordivela` | AGPL-3.0 | Community coordination runtime for a single human or small team. Reference profile: fewer than five concurrent agents |
| `coordivela-open` | Apache-2.0 | Public contracts, SDKs, examples, and approved core/adapter/probe binaries |
| `coordivela-ee` | Commercial | Enterprise governance, multi-tenant scale, HA, federation, compliance, and managed operations |

Core and open share a version line. EE releases independently and states the
core range it supports. The agent-count example is a scope target, not a
license-enforced counter.

## This Repository

- `contracts/openapi/coordivela.v1.yaml` - HTTP contract for work items,
  assignments, progress, interventions, and schedules.
- `sdk/go/` - Go SDK.
- `examples/quickstart/` - end-to-end API example.
- `releases/` - approved binaries, manifests, checksums, and SBOMs published by
  the core release process.

## Binary Distribution

`coordivela-open` is the distribution home for approved COORDIVELA binaries:

| Binary | Source | Purpose |
| --- | --- | --- |
| `coordivela` | core | Coordination server |
| `coordivela-migrate` | core | Schema migration helper |
| `coordivela-opencode-probe` | core | Read-only OpenCode observation and probe |
| Adapter binaries | core/adapters | PRAXOVELA, OpenCode, and process adapters when packaged as separate binaries |

Each release manifest records the binary version, source repository, license,
checksums, and SBOM. Binaries retain the license of their source; the contracts,
SDK, and examples in this repository are Apache-2.0. Enterprise-only binaries
are never published here.

The current release is `1.0.0`, matching the core version line. Approved
binaries are published under [`releases/v1.0.0/`](releases/v1.0.0/) with
per-platform checksums, license notices, and CycloneDX SBOMs. The OpenAPI
contract in this repository remains revision `0.1.0`.

## Quickstart

```bash
# Start the core stack first
cd ../COORDIVELA
docker compose -f deploy/compose.yaml up --build

# Create a work item
curl -s http://localhost:1976/v1/work-items \
  -H 'X-Axis-Tenant: local' \
  -H 'X-Axis-Workspace: default' \
  -H 'Content-Type: application/json' \
  -d '{"id":"wi-1","title":"Implement durable scheduler persistence"}'

# List work items
curl -s http://localhost:1976/v1/work-items \
  -H 'X-Axis-Tenant: local' \
  -H 'X-Axis-Workspace: default'
```

## Contracts

The OpenAPI contract is the source of truth for the public HTTP surface. SDKs
must be generated or checked against it. Breaking changes require a major
contract version.

## License

Apache-2.0. See [LICENSE](LICENSE).
