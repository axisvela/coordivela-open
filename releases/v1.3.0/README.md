# COORDIVELA v1.3.0

- Source: [COORDIVELA core](https://github.com/axisvela/coordivela) at commit
  `9db1aa040a2b8b35a6796374d4ad284565b24f12` (`v1.3.0`).
- Released: 2026-10-20.
- License: binaries are AGPL-3.0. See `LICENSE-COORDIVELA-CORE.txt`. The
  contracts, SDK, examples, and documentation in this repository remain
  Apache-2.0.
- Build: Go 1.26.9, CGO disabled, `-trimpath`, `-buildvcs=false`.
- Contract: OpenAPI revision `1.3.0`.
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

Each platform directory contains `NOTICE`, `SHA256SUMS`, and
`SHA256SUMS.asc`. The manifest and SBOMs are also signed. Import the public key
and verify:

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
