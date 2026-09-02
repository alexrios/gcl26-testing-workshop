package contracts

import (
	"testing"
	"testing/fstest"
	"testing/iotest"
)

func TestFixtureFSContract(t *testing.T) {
	t.Skip("TODO: valide FixtureFS com fstest.TestFS")

	// TODO: execute o contrato para fixtures/default.json.
	_ = fstest.TestFS
}

func TestChunkReaderContract(t *testing.T) {
	t.Skip("TODO: valide ChunkReader com iotest.TestReader")

	// TODO: use um payload conhecido e chunks de três bytes.
	_ = iotest.TestReader
}
