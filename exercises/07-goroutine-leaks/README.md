# Exercício 7 — detectar goroutines abandonadas

Uma suíte pode ficar verde e ainda deixar goroutines bloqueadas. No Go 1.27, o perfil `goroutineleak` identifica goroutines bloqueadas que não são alcançáveis a partir de goroutines executáveis nem de raízes globais.

## 1. Transforme o vazamento em falha, 18 minutos

Implemente `TestWatcherStopDoesNotLeak`:

1. remova o `t.Skip`;
2. obtenha o perfil com `pprof.Lookup("goroutineleak")`, chame `profile.WriteTo(io.Discard, 0)` e só então capture `Count()`;
3. crie e pare um `Watcher` dentro de `createAndStopWatcher`; o helper chama `Stop` duas vezes para também testar idempotência;
4. ceda o processador repetidamente com `runtime.Gosched` para o worker bloquear; um loop de 1.000 iterações mantém o exercício simples;
5. chame `WriteTo` novamente, agora em um `strings.Builder`, e compare a nova contagem.

`Count` informa o resultado da detecção mais recente; ele não inicia a coleta. `WriteTo(..., 0)` dispara a detecção e grava o perfil compactado, descartado na linha de base. `WriteTo(..., 2)` dispara outra detecção e produz stacks legíveis de todas as goroutines, não apenas das classificadas como leak. Use o delta de `Count` para decidir se houve vazamento e os stacks para localizar o worker suspeito; imprima-os em `t.Output()` antes de falhar.

```bash
go test -run=TestWatcherStopDoesNotLeak -count=1 -v \
  ./exercises/07-goroutine-leaks
```

O experimento tem quatro causas separadas:

1. o retorno do helper torna o watcher inalcançável;
2. os `Gosched` dão oportunidade para o worker chegar ao bloqueio;
3. `WriteTo` dispara a nova coleta;
4. `Count` lê o resultado dessa última coleta.

O loop de 1.000 chamadas é uma acomodação de escalonamento deste exemplo, não uma API nem uma garantia geral de que uma goroutine chegou a determinado estado.

## 2. Corrija o lifecycle, 12 minutos

Faça `Watcher.Stop`:

- fechar `stop` uma única vez;
- esperar `done` antes de retornar;
- continuar seguro se chamado mais de uma vez.

Remova também o skip de `TestWatcherStopWaitsForWorker`. Esse teste usa `synctest` para provar uma pós-condição diferente da coleta do perfil: `Stop` não pode retornar enquanto o worker ainda não terminou.

Rode também com o race detector:

```bash
go test -race -count=1 -v ./exercises/07-goroutine-leaks
```

## Limites e perguntas, 10 minutos

- Por que comparar `runtime.NumGoroutine` seria mais ruidoso?
- Um perfil vazio prova ausência de todos os vazamentos?
- Como `synctest`, `goleak`, `goroutineleak` e `-race` respondem perguntas diferentes?

O perfil usa alcançabilidade. Uma goroutine abandonada ainda referenciada por estado global ou por uma goroutine executável pode não aparecer. Já `synctest` encontra workers que impedem uma bubble de terminar; ele exige que o cenário caiba nessa bubble. A biblioteca `go.uber.org/goleak` compara snapshots de stacks e cobre código fora da bubble, mas precisa lidar com goroutines legítimas e transitórias. Neste workshop, o exercício usa o perfil estável do Go 1.27 por ser a alternativa mais nova.

## Concluído quando

- o teste original imprime o stack da goroutine abandonada e falha;
- `Stop` encerra o worker e é idempotente;
- o teste de pós-condição falha se `Stop` apenas sinaliza e retorna;
- a suíte passa com `-race`;
- você consegue explicar pelo menos um falso negativo possível do perfil.

Rota mínima: implemente a coleta, observe o stack, corrija `Stop` e rode os dois testes com `-race`. A comparação entre quatro ferramentas pode ser conduzida pelo instrutor se o bloco alcançar 03:45 antes da turma.
