# Production Load Balancing And Failover Feedback

## Context And Boundary

EE 1.0.0 is a source baseline, not a production HA claim. This public issue
requests deployment feedback without publishing an EE implementation or making
multi-instance leader election a Core/Open requirement.

## Evidence Still Needed

- Record supported proxy/LB versions, TLS termination, trusted headers,
  readiness/health checks, routing, leader selection and timeout policies.
- Exercise process crashes, rolling restart, dependency outages, partial
  network failure, connection draining and restoration of a standby.
- Prove authentication/scope is not bypassed through an alternate backend.
- Validate real reads/writes across failover, including long-lived streams and
  in-flight mutation retries; reject stale fencing and avoid duplicate effects.
- Measure user-visible disruption and accepted/lost/uncertain commands, not
  just an HA status endpoint becoming available.

## Community Feedback Requested

Which LB, Kubernetes/VM topology and streaming-client patterns should the
deployment matrix cover? Sanitized failed requests and reproducible fault
scenarios are particularly useful. A custom loopback proxy is not equivalent
to production load-balancer acceptance.

## Closure Criteria

- Representative multi-host tests meet agreed failover/draining targets.
- There is persistence evidence for no unauthorized/duplicate mutation and
  explicit handling of uncertain outcomes.
- Deployment owners approve runbooks and residual risk. Do not share private
  topology identifiers, keys, credentials or customer payloads.
