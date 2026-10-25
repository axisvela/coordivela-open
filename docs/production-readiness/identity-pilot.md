# Identity And Authorization Deployment Feedback

## Context And Boundary

Core/Open are on 1.7.2. EE 1.0.0 is a source-level baseline, not a production
acceptance claim. This public issue solicits integration experience after that
release; it does not move enterprise identity features into Core/Open.
Identity and authentication remain with the IdP/EIDOVELA; authorization remains
with AEGIVELA. EE owns enterprise composition and policy integration.

## Evidence Still Needed

- Exercise actual human login and service/workload authentication as separate
  cases, with the correct authorization mode and trusted identity source.
- Validate issuer, audience, expiry, signature, tenant/workspace scope, and
  rotation/revocation using supported provider deployments.
- Test allow, deny, provider outage, unknown key, unsupported obligations, and
  approval action/resource/version binding on real effect paths.
- Establish that denied, expired, revoked, or unavailable authority produces no
  unauthorized effect; retain redacted request/result and persistence evidence.

## Community Feedback Requested

Which IdP/AEGIVELA versions, login flows and service-identity mechanisms do you
use? Which mappings or error classifications prevent adoption? Please share
sanitized reproduction steps or a minimal test fixture, not credentials/tokens.

## Closure Criteria

- A reviewed deployment matrix records exact versions, safe configuration,
  per-case expected/actual outcomes and failure recovery.
- Human and service cases are not conflated; mode binding is explicitly tested.
- Deployment/security owners approve remaining risks. Comments or source tags
  alone do not close production acceptance.

Public sharing rule: never post bearer tokens, private keys, full DSNs, customer
identities, or private policy payloads. Chinese or English feedback is welcome.
