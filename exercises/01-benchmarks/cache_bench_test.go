package benchmarks

import (
	"strconv"
	"testing"
)

const keySpace = 128

var benchKeys = makeKeys(keySpace)

var implementations = []struct {
	name string
	new  func() Cache
}{
	{name: "RWMutex", new: func() Cache { return NewRWCache() }},
	{name: "sync.Map", new: func() Cache { return NewSyncMapCache() }},
}

func makeKeys(n int) []string {
	keys := make([]string, n)
	for i := range n {
		keys[i] = "key-" + strconv.Itoa(i)
	}
	return keys
}

func populate(c Cache) {
	for _, key := range benchKeys {
		c.Set(key, "value")
	}
}

// BenchmarkGetMisleading é o diagnóstico inicial do laboratório.
func BenchmarkGetMisleading(b *testing.B) {
	for b.Loop() {
		c := NewRWCache()
		populate(c)
		c.Get(benchKeys[0])
	}
}

func BenchmarkGet(b *testing.B) {
	for _, implementation := range implementations {
		b.Run(implementation.name, func(b *testing.B) {
			b.Skip("TODO: prepare o cache e meça Get dentro de B.Loop")

			_ = implementation.new
		})
	}
}

func BenchmarkSetOverwrite(b *testing.B) {
	for _, implementation := range implementations {
		b.Run(implementation.name, func(b *testing.B) {
			b.Skip("TODO: prepare o cache e meça Set dentro de B.Loop")

			_ = implementation.new
		})
	}
}

func BenchmarkGetParallel(b *testing.B) {
	b.Skip("stretch goal: use b.RunParallel e pb.Next")
}
