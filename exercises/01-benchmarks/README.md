# Exercício 1 — benchmarks que respondem uma pergunta

Temos duas implementações corretas do mesmo contrato: uma protegida por `sync.RWMutex` e outra baseada em `sync.Map`. O objetivo não é coroar uma vencedora universal. Você deverá medir workloads específicos sem incluir trabalho acidental.

## 1. Veja um benchmark enganoso

```bash
go test -run='^$' -bench=Misleading -benchmem ./exercises/01-benchmarks
```

Antes de editar, liste o que ele mede além de `Get`.

## 2. Implemente os benchmarks seriais

Em `cache_bench_test.go`:

1. remova os `b.Skip`;
2. crie e popule o cache antes de `for b.Loop()`;
3. use um conjunto de chaves pré-computado;
4. escreva sub-benchmarks `RWMutex` e `sync.Map` para `Get` e sobrescrita via `Set`;
5. ative `b.ReportAllocs()`.

Execute:

```bash
go test -run='^$' -bench='Benchmark(Get|SetOverwrite)$' -benchmem \
  ./exercises/01-benchmarks
```

## 3. Compare distribuições

Colete duas amostras independentes do mesmo código. A hipótese é que não há mudança real: qualquer diferença aparente deve ser tratada como ruído até que os dados sustentem outra conclusão.

```bash
go test -run='^$' -bench='Benchmark(Get|SetOverwrite)$' -benchmem -count=10 \
  ./exercises/01-benchmarks > bench-a.txt

go test -run='^$' -bench='Benchmark(Get|SetOverwrite)$' -benchmem -count=10 \
  ./exercises/01-benchmarks > bench-b.txt

go tool benchstat bench-a.txt bench-b.txt
```

Procure o intervalo de cada amostra e a indicação de significância. Um `~` no delta informa que o `benchstat` não encontrou diferença estatisticamente significativa; não significa que as medições foram idênticas. Se houver uma diferença inesperada entre amostras do mesmo código, investigue ruído, throttling e competição pela máquina antes de atribuí-la à implementação.

## 4. Se terminar cedo

Implemente `BenchmarkGetParallel` com `b.RunParallel`. Dentro dele, use `pb.Next()`, não `b.Loop()`.

Depois do workshop, uma extensão útil é comparar duas variantes com nomes compatíveis: chave fixa versus rotação pelo keyspace. Declare antes qual custo adicional você espera incluir.

## Concluído quando

- setup e população estão fora de `B.Loop`;
- os quatro benchmarks seriais reportam tempo e alocações;
- você consegue explicar por que os resultados seriais não predizem contenção;
- o benchmark paralelo, se implementado, não mistura `B.Loop` e `RunParallel`.

Reserve os cinco minutos finais do bloco para comparar hipóteses e discutir o que o workload de sobrescrita não representa.

Rota mínima: conclua os quatro sub-benchmarks seriais e interprete a comparação A/B. `RunParallel` é stretch goal e não bloqueia `mise run progress`.
