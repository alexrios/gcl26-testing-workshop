# Exercício 6 — concorrência determinística no Go 1.27

O cache já implementa tempo de vida (TTL), um janitor periódico e uma API HTTP.
Primeiro observe o custo e o limite de um teste convencional. Depois introduza
`testing/synctest` em etapas, sem polling ou portas do sistema operacional.

## 1. Antes de `synctest`

O starter testa um TTL de 20 milissegundos com `time.Sleep`. Execute dez vezes:

```bash
go test -run=TestTTLExpires -count=10 -v ./exercises/06-synctest
```

O teste passa, mas cada repetição paga a espera no relógio da máquina. Reduzir
uma hora de produção para poucos milissegundos também muda o cenário e aproxima
a asserção do limite temporal.

Antes de editar, formule a hipótese: como testar uma hora sem esperar uma hora
nem trocar o TTL exercitado?

## 2. Uma hora no relógio virtual

Reescreva `TestTTLExpires`:

- envolva o cenário em `synctest.Test`;
- use o TTL real de uma hora;
- mantenha primeiro `time.Sleep(time.Hour)` dentro da bubble;
- confirme que a entrada existe antes do prazo (`deadline`) e expira depois dele.

```bash
go test -run=TestTTLExpires -count=1 -v ./exercises/06-synctest
```

O teste agora deve levar milissegundos de parede. O `time.Sleep` usa o relógio
virtual porque nasceu dentro da bubble.

Mini glossário para explicar o que acabou de acontecer:

- **bubble**: conjunto rastreado de goroutines e objetos de sincronização criado por `synctest.Test`;
- **durable blocking**: espera que a bubble reconhece como estável até outra goroutine da própria bubble agir;
- **quiescência**: momento em que todas as goroutines da bubble estão duravelmente bloqueadas;
- **`synctest.Wait()`**: espera a quiescência sem avançar o relógio;
- **`synctest.Sleep(d)`**: avança o relógio virtual por `d` e depois espera a quiescência, como `time.Sleep(d)` seguido de `synctest.Wait()` dentro da bubble.

## 3. `Deadline` não é reação

Implemente `TestJanitorRemovesExpiredEntries`:

- inicie o janitor com intervalo de um minuto;
- grave entradas que vencem antes e depois do primeiro ciclo;
- mantenha primeiro o `time.Sleep(time.Minute)` entregue no starter;
- verifique `Size`, não `Get`, para provar que a remoção física ocorreu;
- cancele o contexto e espere o canal `done` fechar.

Remova o skip. Antes de adicionar outra sincronização, execute o teste
repetidamente:

```bash
go test -run=TestJanitorRemovesExpiredEntries -count=100 \
  ./exercises/06-synctest
```

O teste pode passar ou falhar conforme a ordem entre o teste e o janitor ao
acordarem no mesmo instante virtual. Um passe isolado não valida a sincronização;
essa incerteza é o resultado do experimento.

Agora adicione `synctest.Wait()` logo depois de `time.Sleep(time.Minute)` e
repita o comando. O primeiro alcança o `deadline`; o segundo espera a bubble
voltar à quiescência depois da reação do janitor.

Por fim, substitua o par por `synctest.Sleep(time.Minute)`. Aplique a mesma forma
composta ao teste de TTL. A abreviação entra somente depois que as duas operações
separadas ficaram observáveis.

## 4. HTTP sem socket real

Em `TestHTTPEntryExpires`, use `httptest.NewTestServer(t, NewHandler(cache))`, introduzido no Go 1.27.

1. envie `PUT /cache/session?ttl=1h`;
2. confirme `GET` com status 200;
3. avance e aquiete a bubble com `synctest.Sleep(time.Hour)`;
4. confirme `GET` com status 404.

Não use `httptest.NewServer`: ele abre loopback real, e I/O externo não é durable blocking.

## 5. Complemente com o race detector

```bash
go test -race -v ./exercises/06-synctest
```

Explique em uma frase o que `-race` observa que `synctest` não substitui.

## Se terminar cedo

- teste que cancelar o janitor antes do primeiro tick preserva a entrada física;
- adicione dois subtests com bubbles independentes. Confirme que ambos começam
  em 2000-01-01 no Tempo Universal Coordenado (`UTC`).

## Concluído quando

- nenhuma espera usa relógio real;
- o teste distingue `deadline` de reação da goroutine;
- o worker sempre termina;
- o teste HTTP não abre porta;
- a suíte passa também com `-race`.

O teste de TTL e a diferença entre `deadline` e reação são o núcleo do
laboratório. Se o tempo apertar, acompanhe a implementação HTTP projetada pelo
instrutor. Execute-a localmente. Reserve os cinco minutos finais para perguntas.
