package propertytests

import "testing"

func TestNormalizeKeyExamples(t *testing.T) {
	tests := map[string]string{
		"GopherCon LATAM 2026": "gophercon-latam-2026",
		"already-canonical":    "already-canonical",
		"  Go!! Testing  ":     "go-testing",
		"":                     "",
	}

	for input, want := range tests {
		t.Run(input, func(t *testing.T) {
			if got := NormalizeKey(input); got != want {
				t.Fatalf("NormalizeKey(%q) = %q; want %q", input, got, want)
			}
		})
	}
}
