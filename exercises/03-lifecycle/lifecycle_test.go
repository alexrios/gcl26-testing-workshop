package lifecycle

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestWorkerStopsWithTest(t *testing.T) {
	t.Skip("TODO: use t.Context e espere o worker em t.Cleanup")

	// TODO: inicie o worker e registre o cleanup.
}

func TestReportArtifact(t *testing.T) {
	t.Skip("TODO: use t.Attr, t.Output e t.ArtifactDir")

	report, err := RenderReport(Report{Conference: "GopherCon LATAM", Tests: 7})
	if err != nil {
		t.Fatal(err)
	}

	// TODO: associe o owner, grave report.json e informe o caminho.
	_, _, _ = fmt.Fprintf, os.WriteFile, filepath.Join
	_ = report
}
