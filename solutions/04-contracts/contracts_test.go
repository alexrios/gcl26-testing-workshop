package contracts

import (
	"testing"
	"testing/fstest"
	"testing/iotest"
)

func TestFixtureFSContract(t *testing.T) {
	if err := fstest.TestFS(FixtureFS(), "fixtures/default.json"); err != nil {
		t.Fatal(err)
	}
}

func TestChunkReaderContract(t *testing.T) {
	want := []byte("GopherCon LATAM")
	if err := iotest.TestReader(NewChunkReader(want, 3), want); err != nil {
		t.Fatal(err)
	}
}
