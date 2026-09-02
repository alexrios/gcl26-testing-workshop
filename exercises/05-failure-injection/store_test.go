package failureinjection

import (
	"errors"
	"testing"
)

type journalStub struct {
	writeErr error
	syncErr  error
}

func (j *journalStub) Write([]byte) error { return j.writeErr }
func (j *journalStub) Sync() error        { return j.syncErr }

func TestPutHappyPath(t *testing.T) {
	store := NewStore(&journalStub{})
	if err := store.Put("region", "latam"); err != nil {
		t.Fatal(err)
	}
	if value, ok := store.Get("region"); !ok || value != "latam" {
		t.Fatalf("Get(region) = %q, %v; want latam, true", value, ok)
	}
}

func TestPutDoesNotCommitWhenSyncFails(t *testing.T) {
	t.Skip("TODO: injete uma falha de Sync e verifique erro e estado em memória")

	// TODO: use um erro sentinela em journalStub.syncErr.
	_ = errors.Is
}
