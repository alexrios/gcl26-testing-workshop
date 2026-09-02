package propertytests

import "testing"

func TestPropertyClinic(t *testing.T) {
	candidates := []struct {
		name  string
		check func(string) bool
	}{
		{name: "same-call", check: sameCall},
		{name: "never-panics", check: neverPanics},
		{name: "preserves-length", check: preservesLength},
		{name: "production-alphabet", check: usesProductionAlphabet},
	}
	inputs := []string{"a_b", "Go Testing", "!!!", "already-canonical"}

	for _, candidate := range candidates {
		t.Run(candidate.name, func(t *testing.T) {
			for _, input := range inputs {
				t.Logf("input=%q resultado=%t", input, candidate.check(input))
			}
		})
	}
}

func sameCall(input string) bool {
	return NormalizeKey(input) == NormalizeKey(input)
}

func neverPanics(input string) bool {
	_ = NormalizeKey(input)
	return true
}

func preservesLength(input string) bool {
	return len(NormalizeKey(input)) == len(input)
}

func usesProductionAlphabet(input string) bool {
	result := NormalizeKey(input)
	for i := 0; i < len(result); i++ {
		current := result[i]
		if !isASCIIAlphaNumeric(current) && current != '-' {
			return false
		}
	}
	return true
}
