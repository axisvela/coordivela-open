# Production Readiness And Community Feedback

Baseline checked: 2026-10-25. Core/Open: **1.7.2**. EE: **1.0.0 source
baseline**. Core Go module: `backend/v3.1.2`; Open Go SDK: `sdk/go/v3.0.0`.
Product, module and contract versions are separate identifiers.

EE 1.0.0 is a maintainer-approved source release, not production acceptance.
Production gates remain open after 1.0.0. Community Core/Open scope is unchanged;
EE releases independently and enterprise-only features/binaries remain private.

## Open Public Issues

These are requests for sanitized operational experience and public integration
requirements. Issue discussion is not a production approval or a commitment to
implement enterprise features in Core/Open.

| Gate | Issue | Scope and acceptance criteria |
| --- | --- | --- |
| Identity and authorization | [#1](https://github.com/axisvela/coordivela-open/issues/1) | [Deployment pilot](identity-pilot.md) |
| Backup and recovery | [#5](https://github.com/axisvela/coordivela-open/issues/5) | [WAL retention, PITR, RTO/RPO](recovery-pitr.md) |
| Availability and failover | [#3](https://github.com/axisvela/coordivela-open/issues/3) | [LB, draining and multi-host failure](load-balancer-failover.md) |
| Artifact supply chain | [#2](https://github.com/axisvela/coordivela-open/issues/2) | [SBOM, provenance and verification](supply-chain-provenance.md) |
| Compliance | [#4](https://github.com/axisvela/coordivela-open/issues/4) | [Evidence mapping and owner approval](compliance-evidence.md) |
| Operational acceptance | [#6](https://github.com/axisvela/coordivela-open/issues/6) | [SLIs, SLOs, alerting and support](slo-acceptance.md) |

## Feedback Status

At the initial check on 2026-10-25, the repository had no existing issues or pull
requests and Discussions was disabled. The six issues above were then created
to invite feedback. Check their live comments for subsequent community input;
creation is not evidence of adoption, external review, or approval.

## How To Contribute

Include product/module versions, operating system, dependency versions,
sanitized configuration, reproduction commands, per-case expected/actual
results, failure handling, cleanup and limitations. Separate helper tests,
local drills, deployed pilots and production evidence.

Never publish credentials, bearer tokens, private keys, full DSNs, customer
data, private audit reports, or EE-only implementation artifacts here.
Public Core/Open fixes remain public; enterprise-specific implementation work
is tracked in its owning repository. Chinese or English feedback is welcome.
