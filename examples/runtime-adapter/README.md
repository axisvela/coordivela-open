# RuntimeAdapter v1 Example

This example shows a minimal adapter built only from the public
`RuntimeAdapter` v1 contract. It is deterministic and intended for contract
tests, not production execution.

```bash
cd examples/runtime-adapter
GOWORK=off go test ./...
```

The test runs the public conformance kit:

```go
conformance.Run(t, New(), conformance.Options{
    Start: adapter.StartRequest{WorkItemID: "wi-1", Prompt: "do work"},
    Timeout: time.Second,
    WaitForTerminal: true,
})
```

See `contracts/runtime-adapter/v1/README.md` for lifecycle and event semantics.
