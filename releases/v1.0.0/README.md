# COORDIVELA v1.0.0

- Source: [COORDIVELA core](https://github.com/axisvela/coordivela) at commit
  `7dc57f0a5a6f8b5ffec63be8ac7f8656fdda1b8f` (`v1.0.0`).
- Released: 2026-10-20.
- License: binaries are AGPL-3.0. See `LICENSE-COORDIVELA-CORE.txt`. The
  contracts, SDK, examples, and documentation in this repository remain
  Apache-2.0.
- Build: Go 1.26.9, CGO disabled, `-trimpath`.
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

Each platform directory contains `NOTICE` and `SHA256SUMS`. Verify a download
with:

```powershell
Get-FileHash .\coordivela.exe -Algorithm SHA256
Get-Content .\SHA256SUMS
```

```sh
sha256sum -c SHA256SUMS
```

## SBOM

- [`sbom/coordivela-backend.cdx.json`](sbom/coordivela-backend.cdx.json)
- [`sbom/coordivela-frontend.cdx.json`](sbom/coordivela-frontend.cdx.json)

## Known Limitations

- PRAXOVELA durable execution-timer evidence remains outstanding in that
  product. Core ships stall reconciliation and timer-signal ingestion but does
  not claim PRAXOVELA timer durability.
- `RuntimeAdapter` v1 is not part of this release.
- The OpenAPI contract in this repository remains revision `0.1.0`; the `/v1`
  runtime surface is version `1.0.0`.
