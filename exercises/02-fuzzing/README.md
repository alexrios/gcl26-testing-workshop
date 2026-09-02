# Exercício 2 — de propriedade a regressão

O codec usa um formato pequeno: um byte de versão, dois bytes com o tamanho da chave, a chave e o valor. Os testes de exemplo passam, mas o parser foi escrito supondo que a entrada sempre veio do encoder.

Antes da primeira execução, não abra `codec.go`: formule a propriedade e deixe
a falha produzir a primeira evidência. Depois de observar a falha, leia o parser
como parte do diagnóstico.

Mini mapa da API:

- `f.Add(...)` registra seeds executados também por `go test` comum;
- `f.Fuzz(func(t *testing.T, ...))` define a propriedade e os tipos de entrada;
- durante a busca, o motor deriva novos inputs e mantém um corpus interno;
- um caso gravado em `testdata/fuzz/<Target>` vira regressão versionável.

## 1. Escreva uma propriedade de round-trip

Implemente `FuzzEntryRoundTrip` em `codec_fuzz_test.go`:

```text
Decode(Encode(Entry{Key: key, Value: value})) == Entry{Key: key, Value: value}
```

Adicione seeds vazios, texto `ASCII`, Unicode e strings com bytes arbitrários.
Rode apenas o corpus inicial:

```bash
go test -run=FuzzEntryRoundTrip ./exercises/02-fuzzing
```

O domínio da propriedade é o conjunto aceito por `Encode`: chaves com no máximo `math.MaxUint16` bytes. Se `Encode` rejeitar uma chave maior, trate o caso como fora do domínio com `t.Skip`; não o reporte como defeito do round-trip.

## 2. Faça o decoder enfrentar bytes arbitrários

Implemente `FuzzDecodeNeverPanics`. O scaffold contém três seeds válidos; acrescente outros se ajudarem sua hipótese. `Decode` pode retornar erro para dados inválidos, mas não pode entrar em panic.

```bash
go test -run='^$' -fuzz=FuzzDecodeNeverPanics -fuzztime=20s \
  ./exercises/02-fuzzing
```

Leia o caminho do input que o Go gravou e reproduza sem busca:

```bash
go test -run=FuzzDecodeNeverPanics ./exercises/02-fuzzing
```

## 3. Corrija sem apagar a evidência

Valide o tamanho antes de indexar ou fatiar `data`. Mantenha o arquivo criado em `testdata/fuzz`: ele agora é uma regressão executada por `go test` comum.

## Se o fuzzer não achar em dois minutos

Use o corpus preparado no kit:

```bash
mkdir -p exercises/02-fuzzing/testdata/fuzz/FuzzDecodeNeverPanics
cp support/fuzz-crasher \
  exercises/02-fuzzing/testdata/fuzz/FuzzDecodeNeverPanics/short-input
go test ./exercises/02-fuzzing
```

## Stretch goals

- rejeite bytes restantes depois do valor se o protocolo exigir representação canônica;
- limite o tamanho de entradas antes de alocar;
- adicione uma propriedade `Encode(Decode(data))` apenas para entradas válidas e explique por que ela não é simétrica à primeira.

## Concluído quando

- os dois fuzz targets têm pelo menos três seeds úteis;
- a versão original do parser entra em panic com um input reproduzível;
- o parser corrigido retorna erro, não panic;
- `go test ./exercises/02-fuzzing` executa o caso encontrado como regressão.

Reserve os cinco minutos finais do bloco para discutir limites de entrada e a diferença entre corpus de busca e corpus de regressão.

Rota mínima: implemente os dois targets, reproduza o panic sem busca e corrija o parser. Os stretch goals ficam para depois do workshop.
