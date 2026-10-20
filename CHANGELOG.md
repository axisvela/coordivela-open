# Changelog

All notable changes to COORDIVELA Open are documented here.

## [1.3.0] - 2026-10-20

### Added

- Work-item conflict override routes and SDK methods
  (`POST/DELETE /work-items/{id}/conflict-override`).
- Resource-scoped collaborator lease routes and SDK methods
  (`GET/POST /work-items/{id}/leases`, `DELETE /work-items/{id}/leases/{leaseId}`),
  including optional `canComplete` completion grants.

## [1.2.0] - 2026-10-20

### Added

- RuntimeAdapter v1 contract under `contracts/runtime-adapter/v1/`, including
  lifecycle semantics, event vocabulary, error taxonomy, and JSON Schema.
- Go RuntimeAdapter v1 interface in `sdk/go/adapter`.
- Public conformance kit in `sdk/go/adapter/conformance`.
- Core now aliases the published RuntimeAdapter types and runs the public
  conformance kit in its adapter tests.
- External RuntimeAdapter v1 example with a passing conformance test under
  `examples/runtime-adapter/`.
- Projection and SSE usage examples in the quickstart.
- `GET /projections/symbivela-attention` and a matching Go SDK method for
  read-only SYMBIVELA attention.

### Changed

- Expanded the English and Chinese READMEs with COORDIVELA positioning, the
  problem it solves, its value, and the Core/Open/EE boundaries.
- Added English and Chinese operations manuals covering installation,
  configuration, PostgreSQL setup, startup, health checks, backups, upgrades,
  monitoring, and troubleshooting.
- Completed the OpenAPI contract for all published `/v1` routes and set the
  contract revision to `1.0.0`.
- Expanded the Go SDK to cover work-item lifecycle, dependencies, assignments,
  schedules, projections, and SSE.

### Fixed

- The Go SDK now sends `Idempotency-Key` on work-item creation and exposes
  `CreateWorkItemWithKey` for caller-controlled retry keys.
- Quickstart requests now include the required `Idempotency-Key` header.

## [1.0.0] - 2026-10-20

### Added

- Approved core binaries for windows/amd64, linux/amd64, and darwin/arm64:
  `coordivela`, `coordivela-migrate`, and `coordivela-opencode-probe`.
- Release manifest, per-platform checksums, AGPL license notices, and
  CycloneDX SBOMs under `releases/v1.0.0/`.

### Changed

- Synced the open repository version to the core `1.0.0` line.

## [0.1.0] - 2026-10-18

### Added

- Initial public OpenAPI contract for work items, assignments, progress,
  interventions, and schedules.
- Go SDK skeleton for the work-item surface.
- Quickstart example.
