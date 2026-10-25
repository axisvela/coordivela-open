# Production Backup, WAL Retention And PITR Feedback

## Context And Boundary

Core/Open 1.7.2 and EE 1.0.0 do not establish production recovery evidence for
every topology. This public issue requests reproducible operator experience.
Community backup/runbook improvements belong in Core/Open. Managed recovery,
enterprise topology and HA operations remain EE/deployment-owned.

## Evidence Still Needed

- Record PostgreSQL version, dataset scale, backup consistency and integrity,
  storage/encryption/access controls, retention, and tested restore inputs.
- Verify base backups and continuous WAL archival, missing/corrupt WAL alarms,
  retention boundaries and the latest recoverable point.
- Restore to another host and a selected recovery time. Prove pre-target data
  is retained and post-target effects excluded using deterministic markers.
- Measure RTO from failure detection through validated writable-service and
  application-data recovery; measure RPO against committed work, not backup age.
- Document rollback, recovery promotion, integrity checks and cleanup.

## Community Feedback Requested

Which backup tooling, storage and PostgreSQL topology do you operate? Share
sanitized drill procedures, data scale, timings and negative cases. Local toy
restore timings must not be presented as production RTO/RPO.

## Closure Criteria

- Repeatable cross-host drills meet operator-agreed RTO/RPO at representative
  scale, with per-case inputs/results and WAL integrity/retention evidence.
- Operators approve restoration and retention runbooks and accepted risks.
- No customer dumps, secrets, full DSNs or database contents are posted publicly.
