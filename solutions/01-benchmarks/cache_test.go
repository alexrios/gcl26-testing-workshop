package benchmarks

import "testing"

func TestImplementationsShareTheContract(t *testing.T) {
	factories := map[string]func() Cache{
		"RWMutex":  func() Cache { return NewRWCache() },
		"sync.Map": func() Cache { return NewSyncMapCache() },
	}

	for name, newCache := range factories {
		t.Run(name, func(t *testing.T) {
			c := newCache()
			c.Set("conference", "GopherCon LATAM")
			got, ok := c.Get("conference")
			if !ok || got != "GopherCon LATAM" {
				t.Fatalf("Get(conference) = (%q, %v)", got, ok)
			}
		})
	}
}
