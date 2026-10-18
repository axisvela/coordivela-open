# COORDIVELA Open

**Public integration surface for the Development Coordination Control Plane**

Apache-2.0 contracts, SDKs, examples, and approved binaries for COORDIVELA.

> COORDIVELA coordinates development work across humans and agents. It does not
> execute code or verify correctness; those remain with PRAXOVELA and TEKMOVELA.

## Repository Model

| Repository | License | Role |
| --- | --- | --- |
| `coordivela-open` | Apache-2.0 | Public contracts, SDKs, examples, approved binaries |
| `coordivela` | AGPL-3.0 | Open-core coordination runtime |
| `coordivela-ee` | Enterprise | Enterprise features and internal product governance |

## This Repository

- `contracts/openapi/coordivela.v1.yaml` - HTTP contract for work items,
  assignments, progress, interventions, and schedules.
- `sdk/go/` - Go SDK.
- `examples/quickstart/` - end-to-end API example.
- `releases/` - approved binaries and manifests (published by the core release
  process).

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
