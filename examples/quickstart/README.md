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

4. Read projections:

```go
package main

import (
    "context"
    "fmt"

    "github.com/axisvela/coordivela-open/sdk/go/v3/client"
)

func main() {
    c := client.New("http://localhost:1976", "local", "default")
    c.Principal = "alice"
    ctx := context.Background()

    mine, err := c.MyWork(ctx)
    if err != nil { panic(err) }
    attention, err := c.Attention(ctx)
    if err != nil { panic(err) }
    fmt.Println("my work:", len(mine), "attention:", len(attention))
}
```

`Program(ctx, programID)` returns work items plus state counts, and
`WorkItemProjection(ctx, id)` returns one item with its outgoing dependencies.

5. Stream coordination events (using the same `c` and `ctx`, plus `bufio` and
   `strings` imports):

```go
resp, err := c.StreamEvents(ctx, 0)
if err != nil { panic(err) }
defer resp.Body.Close()

scanner := bufio.NewScanner(resp.Body)
for scanner.Scan() {
    line := scanner.Text()
    if strings.HasPrefix(line, "data: ") {
        fmt.Println(strings.TrimPrefix(line, "data: "))
    }
}
```

Pass the last SSE `id` as the second argument to `StreamEvents` to request
bounded replay. Replay is an in-process convenience; Harmovela owns durable
retention when kernel publishing is enabled.
