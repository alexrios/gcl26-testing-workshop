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
	wantErr := errors.New("disk unavailable")
	store := NewStore(&journalStub{syncErr: wantErr})

	err := store.Put("region", "latam")
	if !errors.Is(err, wantErr) {
		t.Fatalf("Put error = %v, want an error wrapping %v", err, wantErr)
	}
	if value, ok := store.Get("region"); ok {
		t.Fatalf("Get(region) = %q, true; want key absent after Sync failure", value)
	}
}
