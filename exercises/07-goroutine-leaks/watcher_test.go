package goroutineleaks

import (
	"fmt"
	"io"
	"runtime"
	"runtime/pprof"
	"strings"
	"testing"
	"testing/synctest"
)

func TestWatcherStopDoesNotLeak(t *testing.T) {
	t.Skip("TODO: compare o perfil goroutineleak antes e depois do cenário")

	profile := pprof.Lookup("goroutineleak")
	if err := profile.WriteTo(io.Discard, 0); err != nil {
		t.Fatal(err)
	}
	before := profile.Count()

	createAndStopWatcher()
	for range 1_000 {
		runtime.Gosched()
	}

	var stacks strings.Builder
	// TODO: colete novamente, compare Count e imprima stacks antes de falhar.
	_, _, _ = fmt.Fprint, before, stacks
}

func TestWatcherStopWaitsForWorker(t *testing.T) {
	t.Skip("TODO: prove que Stop só retorna depois que o worker termina")

	assertStopWaitsForWorker(t)
}

func createAndStopWatcher() {
	watcher := NewWatcher()
	watcher.Stop()
	watcher.Stop()
}

func assertStopWaitsForWorker(t *testing.T) {
	t.Helper()
	synctest.Test(t, func(t *testing.T) {
		watcher := &Watcher{
			stop: make(chan struct{}),
			done: make(chan struct{}),
		}
		returned := make(chan struct{})

		go func() {
			watcher.Stop()
			close(returned)
		}()
		synctest.Wait()

		select {
		case <-returned:
			close(watcher.done)
			t.Fatal("Stop returned before the worker finished")
		default:
		}
		select {
		case <-watcher.stop:
		default:
			close(watcher.done)
			t.Fatal("Stop did not signal the worker")
		}

		close(watcher.done)
		synctest.Wait()
		select {
		case <-returned:
		default:
			t.Fatal("Stop did not return after the worker finished")
		}
	})
}
