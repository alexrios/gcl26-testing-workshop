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

		want := Entry{Key: key, Value: value}
		data, err := Encode(want)
		if err != nil {
			t.Fatalf("Encode(%#v) inside its domain: %v", want, err)
		}

		got, err := Decode(data)
		if err != nil {
			t.Fatalf("Decode(Encode(%#v)): %v", want, err)
		}
		if got != want {
			t.Fatalf("Decode(Encode(%#v)) = %#v", want, got)
		}
	})
}

func FuzzDecodeNeverPanics(f *testing.F) {
	valid, err := Encode(Entry{Key: "key", Value: "value"})
	if err != nil {
		f.Fatal(err)
	}

	f.Add(valid)
	f.Add([]byte{})
	f.Add([]byte{formatVersion})
	f.Add([]byte{formatVersion, 0})

	f.Fuzz(func(t *testing.T, data []byte) {
		_, _ = Decode(data)
	})
}
