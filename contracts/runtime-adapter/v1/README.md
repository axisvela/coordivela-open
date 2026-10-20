# RuntimeAdapter v1

**Contract version:** `1.0.0`

The RuntimeAdapter contract is the public boundary between COORDIVELA and a
runtime that executes coding work (PRAXOVELA, OpenCode, a process adapter, or an
external implementation). COORDIVELA owns coordination and observation. The
runtime owns execution, tools, credentials, and its own durable state.

## Go surface

The Go contract and conformance kit are published in
`sdk/go/adapter` and `sdk/go/adapter/conformance`.

```go
type Adapter interface {
    Name() string
    Capabilities() Capabilities
    Start(ctx context.Context, req StartRequest) (string, error)
    SendInput(ctx context.Context, runID, text string) error
    RespondApproval(ctx context.Context, runID, requestID string, decision ApprovalDecision) error
    Cancel(ctx context.Context, runID string) error
    Snapshot(ctx context.Context, runID string) (Snapshot, error)
    Subscribe(ctx context.Context, runID string) (<-chan Event, error)
    Resume(ctx context.Context, runID string) error
}
```

`Capabilities` must be truthful. If `interactive` is false, callers must not
depend on `SendInput`. If `approvalHook` is false, callers must not depend on
`RespondApproval`. If `cancel` is false, `Cancel` must return
`ErrUnsupported`. If `sessionResume` is false, `Resume` must return
`ErrUnsupported`.

## Lifecycle

1. `Start` returns a non-empty run id. The same id must be accepted by
   `Snapshot`, `Subscribe`, `Cancel`, and the optional control methods.
2. `Snapshot.State` is non-empty and normalized by the adapter.
3. `Subscribe` emits normalized events until a terminal event (`exited`,
   `error`, or `completed`) and then closes the channel.
4. Events carry the run id when the adapter can supply it. `at` is required.
5. `waiting_for_approval` events must carry `requestId`.
6. `Cancel` is best-effort and must not block indefinitely.
7. Adapters must not expose unredacted credentials or secrets in events,
   snapshots, or errors.

## Event vocabulary

| Event | Meaning |
| --- | --- |
| `started` | Runtime accepted the run |
| `output_activity` | New output or progress |
| `waiting_for_input` | Blocked on human text input |
| `waiting_for_approval` | Blocked on approve/deny |
| `idle` | Alive with no active work |
| `heartbeat` | Liveness signal |
| `stalled` | No progress beyond the adapter threshold |
| `exited` | Runtime exited |
| `error` | Adapter/runtime error |
| `completed` | Run completed successfully |

## Errors

Use the sentinel errors in the `adapter` package where possible:
`ErrUnsupported`, `ErrUnavailable`, `ErrUnauthenticated`, `ErrInvalidRequest`,
and `ErrConflict`. Wrap them with context; callers use `errors.Is`.

## Conformance

Run the public kit from a test:

```go
import (
    "testing"

    "github.com/axisvela/coordivela-open/sdk/go/v3/adapter"
    "github.com/axisvela/coordivela-open/sdk/go/v3/adapter/conformance"
)

func TestRuntimeAdapterConformance(t *testing.T) {
    conformance.Run(t, NewMyAdapter(), conformance.Options{
        Start: adapter.StartRequest{WorkItemID: "wi-1", Prompt: "go"},
        Timeout: 5 * time.Second,
        WaitForTerminal: true,
    })
}
```

The kit checks the lifecycle, snapshot shape, event vocabulary, event run-id
consistency, timestamps, approval request ids, and terminal stream closure.
Capability checks (`CheckInput`, `CheckCancel`, `CheckResume`) are opt-in
because they may change runtime state.

## Versioning

`1.x` is backward compatible. Additive event kinds or optional fields are minor
changes. Removing or renaming fields, changing event semantics, or changing the
Go interface requires a new major contract version. The contract version is
independent of the COORDIVELA product version.
