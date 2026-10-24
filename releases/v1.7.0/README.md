# COORDIVELA v1.7.0

- Source: [COORDIVELA core](https://github.com/axisvela/coordivela) at commit
  `47ed9ab3e84e0614725f534240cd42e70605decc` (`v1.7.0` plus the darwin
  cross-compile fix).
- Released: 2026-10-24.
- License: binaries are AGPL-3.0. See `LICENSE-COORDIVELA-CORE.txt`. The
  contracts, SDK, examples, and documentation in this repository remain
  Apache-2.0.
- Build: Go 1.26.9, CGO disabled, `-trimpath`, `-buildvcs=false`.
- Contract: OpenAPI revision `1.7.0`.
- Manifest: [`manifest.json`](manifest.json).

## Binaries

| Platform | Directory | Binaries |
| --- | --- | --- |
| Windows / amd64 | [`windows-amd64/`](windows-amd64/) | `coordivela.exe`, `coordivela-migrate.exe`, `coordivela-opencode-probe.exe` |
| Linux / amd64 | [`linux-amd64/`](linux-amd64/) | `coordivela`, `coordivela-migrate`, `coordivela-opencode-probe` |
| macOS / arm64 | [`darwin-arm64/`](darwin-arm64/) | `coordivela`, `coordivela-migrate`, `coordivela-opencode-probe` |

`coordivela` is the coordination server, `coordivela-migrate` applies the
PostgreSQL schema, and `coordivela-opencode-probe` observes OpenCode sessions.

## Verification

GitHub release assets include per-platform zip archives, a top-level
`SHA256SUMS` with its signature, the manifest, SBOMs, and the public key. Each
platform directory also contains `NOTICE`, `SHA256SUMS`, and `SHA256SUMS.asc`.
Import the public key and verify:

```bash
gpg --import COORDIVELA-RELEASE-SIGNING-KEY.asc
gpg --verify manifest.json.asc manifest.json
gpg --verify sbom/coordivela-backend.cdx.json.asc sbom/coordivela-backend.cdx.json
gpg --verify sbom/coordivela-frontend.cdx.json.asc sbom/coordivela-frontend.cdx.json
gpg --verify windows-amd64/SHA256SUMS.asc windows-amd64/SHA256SUMS
```

Then verify the binaries:

```sh
cd windows-amd64
sha256sum -c SHA256SUMS
```

```powershell
Get-FileHash .\coordivela.exe -Algorithm SHA256
Get-Content .\SHA256SUMS
```

## SBOM

- [`sbom/coordivela-backend.cdx.json`](sbom/coordivela-backend.cdx.json)
- [`sbom/coordivela-frontend.cdx.json`](sbom/coordivela-frontend.cdx.json)

## Known Limitations

- PRAXOVELA durable execution-timer evidence remains outstanding in that
  product. Core ships stall reconciliation and timer-signal ingestion but does
  not claim PRAXOVELA timer durability.
- Multi-instance leader election remains EE/HA-owned; core uses PostgreSQL
  advisory locks for conflict-checked assignment and per-schedule leases.
- Durable event retention, replay, and dead-letter handling are Harmovela-owned.
