package lifecycle

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestWorkerStopsWithTest(t *testing.T) {
	done := StartWorker(t.Context())
	t.Cleanup(func() {
		<-done
	})
}

func TestReportArtifact(t *testing.T) {
	t.Attr("owner", "platform")
	report, err := RenderReport(Report{Conference: "GopherCon LATAM", Tests: 7})
	if err != nil {
		t.Fatal(err)
	}

	path := filepath.Join(t.ArtifactDir(), "report.json")
	if err := os.WriteFile(path, report, 0o644); err != nil {
		t.Fatal(err)
	}
	fmt.Fprintf(t.Output(), "report: %s\n", path)
}
