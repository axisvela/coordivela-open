# Changelog

All notable changes to COORDIVELA Open are documented here.

## [Unreleased]

### Changed

- Expanded the English and Chinese READMEs with COORDIVELA positioning, the
  problem it solves, its value, and the Core/Open/EE boundaries.
- Added English and Chinese operations manuals covering installation,
  configuration, PostgreSQL setup, startup, health checks, backups, upgrades,
  monitoring, and troubleshooting.

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
