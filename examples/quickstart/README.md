# COORDIVELA Quickstart

1. Start the core stack:

```bash
cd ../../../COORDIVELA
docker compose -f deploy/compose.yaml up --build
```

2. Create and list a work item:

```bash
curl -s -X POST http://localhost:1976/v1/work-items \
  -H 'X-Axis-Tenant: local' \
  -H 'X-Axis-Workspace: default' \
  -H 'Idempotency-Key: quickstart-wi-1' \
  -H 'Content-Type: application/json' \
  -d '{"id":"wi-1","title":"Implement durable scheduler persistence"}'

curl -s http://localhost:1976/v1/work-items \
  -H 'X-Axis-Tenant: local' \
  -H 'X-Axis-Workspace: default'
```

3. Use the Go SDK from `sdk/go/client`.
