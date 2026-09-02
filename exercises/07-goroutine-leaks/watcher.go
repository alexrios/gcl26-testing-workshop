package goroutineleaks

import "sync"

type Watcher struct {
	stop     chan struct{}
	done     chan struct{}
	stopOnce sync.Once
}

func NewWatcher() *Watcher {
	watcher := &Watcher{
		stop: make(chan struct{}),
		done: make(chan struct{}),
	}
	go func() {
		defer close(watcher.done)
		<-watcher.stop
	}()
	return watcher
}

func (w *Watcher) Stop() {
	// O starter retorna sem sinalizar nem esperar pelo worker.
}
