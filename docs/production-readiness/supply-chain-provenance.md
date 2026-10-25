# Release SBOM, Provenance And Verification Feedback

## Context And Boundary

Open publishes approved community binaries with their source licenses. Public
1.7.2 manifests/signatures are not a claim of EE 1.0.0 artifact coverage.
Enterprise binaries and private materials must not be published in Open.

## Evidence Still Needed

- Tie every release artifact to the exact source revision, build toolchain,
  platform, dependency graph and build parameters used for that artifact.
- Distinguish the installed binary's SBOM from a source/module inventory.
  Resolve or explicitly record local dependency replacements.
- Verify signatures and hashes from a clean download/checkout, with documented
  trust-root distribution and signing-key rotation/revocation.
- Validate SBOM completeness, license attribution, vulnerability review and
  accepted-risk handling; cover container provenance when containers ship.
- Publish only sanitized, appropriately licensed public evidence. A signed
  checksum document alone does not establish build provenance or completeness.

## Community Feedback Requested

Which SBOM/provenance formats, verification tools and distribution methods do
you require? Please report public artifact verification problems with version,
platform, command and expected/actual result, without uploading private assets.

## Closure Criteria

- A per-release platform matrix binds source, binary hash, SBOM and provenance.
- Clean verification and policy/license/CVE review pass with documented limits.
- Historical artifacts are not silently relabeled or overwritten.
