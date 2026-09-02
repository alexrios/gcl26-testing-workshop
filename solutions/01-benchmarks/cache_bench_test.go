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

func BenchmarkGet(b *testing.B) {
	for _, implementation := range implementations {
		b.Run(implementation.name, func(b *testing.B) {
			c := implementation.new()
			populate(c)
			b.ReportAllocs()

			i := 0
			for b.Loop() {
				c.Get(benchKeys[i%keySpace])
				i++
			}
		})
	}
}

func BenchmarkSetOverwrite(b *testing.B) {
	for _, implementation := range implementations {
		b.Run(implementation.name, func(b *testing.B) {
			c := implementation.new()
			populate(c)
			b.ReportAllocs()

			i := 0
			for b.Loop() {
				c.Set(benchKeys[i%keySpace], "value")
				i++
			}
		})
	}
}

func BenchmarkGetParallel(b *testing.B) {
	for _, implementation := range implementations {
		b.Run(implementation.name, func(b *testing.B) {
			c := implementation.new()
			populate(c)
			b.ReportAllocs()

			b.RunParallel(func(pb *testing.PB) {
				i := 0
				for pb.Next() {
					c.Get(benchKeys[i%keySpace])
					i++
				}
			})
		})
	}
}
