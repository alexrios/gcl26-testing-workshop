package testing

import stdtesting "testing"

func TestEligibleForFreeShipping(t *stdtesting.T) {
	minimum := 10_000

	t.Run("abaixo do mínimo", func(t *stdtesting.T) {
		if EligibleForFreeShipping(9_999, minimum) {
			t.Fatal("total abaixo do mínimo foi considerado elegível")
		}
	})

	t.Run("acima do mínimo", func(t *stdtesting.T) {
		if !EligibleForFreeShipping(10_001, minimum) {
			t.Fatal("total acima do mínimo não foi considerado elegível")
		}
	})

	t.Run("igual ao mínimo", func(t *stdtesting.T) {
		t.Skip("BÔNUS: remova este skip para testar o limite inclusivo")
		if !EligibleForFreeShipping(minimum, minimum) {
			t.Fatal("total igual ao mínimo não foi considerado elegível")
		}
	})
}
