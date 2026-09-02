package fuzzing

import (
	"math"
	"testing"
)

func FuzzEntryRoundTrip(f *testing.F) {
	f.Add("conference", "GopherCon LATAM")
	f.Add("", "")
	f.Add("chave-🔥", "valor\x00com\nbytes")

	f.Fuzz(func(t *testing.T, key, value string) {
		if len(key) > math.MaxUint16 {
			t.Skip("key is outside Encode's domain")
		}
		t.Skip("TODO: implemente a propriedade no domínio aceito por Encode")

		_, _ = key, value
	})
}

func FuzzDecodeNeverPanics(f *testing.F) {
	f.Add([]byte{formatVersion, 0, 0})
	f.Add([]byte{formatVersion, 0, 1, 'k'})
	f.Add([]byte{formatVersion, 0, 0, 'v'})

	f.Fuzz(func(t *testing.T, data []byte) {
		t.Skip("TODO: passe bytes arbitrários para Decode; erros são permitidos")

		_ = data
	})
}
