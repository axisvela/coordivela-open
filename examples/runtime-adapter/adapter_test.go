package example

import (
	"testing"
	"time"

	"github.com/axisvela/coordivela-open/sdk/go/v3/adapter"
	"github.com/axisvela/coordivela-open/sdk/go/v3/adapter/conformance"
)

func TestExampleAdapterConformance(t *testing.T) {
	conformance.Run(t, New(), conformance.Options{
		Start: adapter.StartRequest{
			WorkItemID: "wi-1",
			Workspace:  "default",
			Prompt:     "do work",
		},
		Timeout:         time.Second,
		WaitForTerminal: true,
	})
}
