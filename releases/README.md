# Release Manifests

Published binary releases include checksums, the originating core version,
license notices, CycloneDX SBOMs, and GPG signatures.

## Signing Key

The COORDIVELA release signing key is published as
`v1.7.1/COORDIVELA-RELEASE-SIGNING-KEY.asc`. Import it and verify a release:

```bash
gpg --import releases/v1.7.1/COORDIVELA-RELEASE-SIGNING-KEY.asc
gpg --verify releases/v1.7.1/manifest.json.asc releases/v1.7.1/manifest.json
gpg --verify releases/v1.7.1/sbom/coordivela-backend.cdx.json.asc \
  releases/v1.7.1/sbom/coordivela-backend.cdx.json
gpg --verify releases/v1.7.1/windows-amd64/SHA256SUMS.asc \
  releases/v1.7.1/windows-amd64/SHA256SUMS
```

## v1.7.1

- Manifest: [`v1.7.1/`](v1.7.1/)
- Binaries: `coordivela`, `coordivela-migrate`, and
  `coordivela-opencode-probe` for windows/amd64, linux/amd64, and darwin/arm64.
- SBOMs: [`v1.7.1/sbom/`](v1.7.1/sbom/)
- Signatures: `manifest.json.asc`, each `SHA256SUMS.asc`, and each SBOM `.asc`.
- Source: [COORDIVELA core](https://github.com/axisvela/coordivela) at commit
  `49ba950ab4c9753a1e220f5c42ff68f40a10952b` (`v1.7.1`). Binaries are AGPL-3.0.

## v1.3.0

- Manifest: [`v1.3.0/`](v1.3.0/)
- Binaries: `coordivela`, `coordivela-migrate`, and
  `coordivela-opencode-probe` for windows/amd64, linux/amd64, and darwin/arm64.
- SBOMs: [`v1.3.0/sbom/`](v1.3.0/sbom/)
- Signatures: `manifest.json.asc`, each `SHA256SUMS.asc`, and each SBOM `.asc`.
- Source: [COORDIVELA core](https://github.com/axisvela/coordivela) at commit
  `9db1aa040a2b8b35a6796374d4ad284565b24f12` (`v1.3.0`). Binaries are AGPL-3.0.

## v1.0.0

- Manifest: [`v1.0.0/`](v1.0.0/)
- Binaries: `coordivela`, `coordivela-migrate`, and
  `coordivela-opencode-probe` for windows/amd64, linux/amd64, and darwin/arm64.
- SBOMs:
  [`v1.0.0/sbom/`](v1.0.0/sbom/)
- Source: [COORDIVELA core](https://github.com/axisvela/coordivela) at commit
  `7dc57f0a5a6f8b5ffec63be8ac7f8656fdda1b8f` (`v1.0.0`). Binaries are AGPL-3.0.
