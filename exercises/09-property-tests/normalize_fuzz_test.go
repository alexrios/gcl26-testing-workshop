package propertytests

import "testing"

func TestModelExamples(t *testing.T) {
	t.Skip("BONUS: valide o modelo independentemente da producao")
}

func FuzzNormalizeCanonical(f *testing.F) {
	f.Add("GopherCon LATAM 2026")
	f.Add("already-canonical")
	f.Add("\x00Go\xffTesting")

	f.Fuzz(func(t *testing.T, input string) {
		t.Skip("BONUS: verifique apenas a forma canonica da saida")

		_ = input
	})
}

func FuzzNormalizeIdempotent(f *testing.F) {
	f.Add("GopherCon LATAM 2026")
	f.Add("already-canonical")
	f.Add("a/b/c")

	f.Fuzz(func(t *testing.T, input string) {
		t.Skip("BONUS: verifique NormalizeKey(NormalizeKey(x)) == NormalizeKey(x)")

		_ = input
	})
}

func FuzzNormalizeMatchesModel(f *testing.F) {
	f.Add("GopherCon LATAM")
	f.Add("already-canonical")
	f.Add("")

	f.Fuzz(func(t *testing.T, input string) {
		t.Skip("BONUS: compare NormalizeKey com um modelo independente")

		_ = input
	})
}

func FuzzEquivalentSeparators(f *testing.F) {
	f.Add("property", "tests", uint8(0))
	f.Add("", "boundary", uint8(2))
	f.Add("left", "", uint8(3))

	f.Fuzz(func(t *testing.T, left, right string, selector uint8) {
		t.Skip("BONUS: compare execucoes relacionadas por separadores equivalentes")

		_, _, _ = left, right, selector
	})
}

// FuzzDerivedProperty fica deliberadamente sem uma formula. Derive uma
// propriedade nova apenas do contrato publico antes de consultar a solucao.
func FuzzDerivedProperty(f *testing.F) {
	f.Add("derive uma propriedade")
	f.Add("\x00\xff")
	f.Add("")

	f.Fuzz(func(t *testing.T, input string) {
		t.Skip("BONUS: formule, nomeie e implemente uma propriedade nova")

		_ = input
	})
}

// BONUS: implemente sem chamar NormalizeKey nem helpers da producao.
func isCanonical(value string) bool {
	return false
}

// BONUS: implemente com operacoes diferentes das usadas em NormalizeKey.
func modelNormalize(input string) string {
	return ""
}
