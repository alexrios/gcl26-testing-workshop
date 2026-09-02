package lifecycle

import (
	"context"
	"encoding/json"
)

func StartWorker(ctx context.Context) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)
		<-ctx.Done()
	}()
	return done
}

type Report struct {
	Conference string `json:"conference"`
	Tests      int    `json:"tests"`
}

func RenderReport(report Report) ([]byte, error) {
	return json.MarshalIndent(report, "", "  ")
}
