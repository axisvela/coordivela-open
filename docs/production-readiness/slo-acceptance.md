# Production SLIs, SLOs And Support Acceptance Feedback

## Context And Boundary

Source releases and local drill timings do not establish an operational SLO.
This public issue requests community measurement requirements after the EE
1.0.0 source baseline. Numeric production commitments require deployment and
product owner sign-off; enterprise support remains EE-owned.

## Evidence Still Needed

- Define service boundaries, eligible requests/events, success/failure,
  denominator, latency buckets, measurement windows and maintenance exclusions.
- Keep readiness probing separate from request availability; retain dependency
  failures in end-to-end SLIs unless an explicit agreement defines exclusions.
- Use request/event budgets for request/event SLIs, not time budgets by analogy.
- Measure authorization unavailability separately from expected policy denial;
  never permit known-revoked authority or unauthorized effects as error budget.
- Base RPO alerts on the latest recoverable point, including archived WAL, not
  merely backup file age. Validate failover through application operations.
- Validate telemetry coverage, alerts, routing, on-call response and incident
  review using representative load and failure scenarios.

## Community Feedback Requested

Which coordination latencies, availability boundaries and recovery targets
matter for your use case? What workload size and client behavior should a
representative test include? Share sanitized metrics, not customer traces.

## Closure Criteria

- Explicit SLIs are implemented and validated; targets and error budgets are
  based on measured representative operation and agreed service boundaries.
- Alert/runbook drills pass and responsible owners approve targets/support.
- Draft targets remain clearly labeled until that approval exists.
