# Exercício 5 — failure injection sem acaso

O `Store` registra cada escrita em um journal antes de tornar o valor visível em memória. O caminho feliz passa; resta provar o que acontece quando `Sync` falha.

Neste laboratório, a interface `Journal` é o ponto de injeção. Em código sem uma dependência substituível, a mesma ideia pode aparecer como um failpoint nomeado, que precisa ser desativado com `t.Cleanup` e isolado de chamadas concorrentes. Falhas de rede e chaos engineering têm outro escopo.

## Laboratório, 11 minutos

Implemente `TestPutDoesNotCommitWhenSyncFails`:

1. remova o `t.Skip`;
2. configure `journalStub.syncErr` com um erro sentinela;
3. chame `Put`;
4. confirme com `errors.Is` que a falha foi preservada;
5. confirme que a chave não ficou visível.

```bash
go test -run=TestPutDoesNotCommitWhenSyncFails -v \
  ./exercises/05-failure-injection
```

O teste deve falhar primeiro. Antes de abrir `store.go`, use a mensagem para localizar qual invariante foi quebrado. Se travar, a pista é revisar a ordem entre `Sync` e a publicação no mapa.

O `Write` no journal já pode ter acontecido quando `Sync` falha. Portanto, este teste prova apenas a visibilidade em memória nesta execução: ele não prova rollback, atomicidade durável nem o resultado de replay após reiniciar o processo.

## Perguntas, 4 minutos

- Por que uma probabilidade aleatória de falha pioraria este teste?
- Qual é o ponto exato de falha que o stub torna controlável?
- Que falhas reais este teste ainda não representa?
- Um restart poderia reaplicar o registro já escrito no journal? Que protocolo seria necessário testar?
- Quando um failpoint nomeado seria melhor que substituir uma dependência?

## Concluído quando

- o teste falha pelo estado parcialmente confirmado, não por timing;
- o erro injetado permanece identificável;
- uma falha de `Sync` não altera o estado visível.

Rota mínima: todo o laboratório é núcleo. Se faltarem minutos, o instrutor fornece o stub e preserva para o aluno as duas asserções e a decisão de ordem no código de produção.
