# Syllabus

## Objetivo

Ao final, cada participante deverá conseguir aplicar as sete técnicas do
workshop. Para cada uma, deverá explicar a pergunta respondida, a evidência
produzida, uma conclusão que ela não autoriza e quando uma abordagem mais
simples seria suficiente.

## Agenda

| Horário | Duração | Bloco |
|---|---:|---|
| 00:00–00:40 | 40 min | Benchmarks modernos |
| 00:40–01:20 | 40 min | Fuzzing nativo |
| 01:20–01:45 | 25 min | Lifecycle e evidência |
| 01:45–01:55 | 10 min | Intervalo |
| 01:55–02:20 | 25 min | Testes de contrato |
| 02:20–02:35 | 15 min | Failure injection |
| 02:35–03:15 | 40 min | `testing/synctest` |
| 03:15–03:55 | 40 min | Detecção de goroutine leaks |
| 03:55–04:00 | 5 min | Fechamento |

Os sete assuntos somam 225 minutos. O intervalo e o fechamento ocupam os 15
minutos restantes. A preparação do ambiente acontece antes da abertura.

## Resultados esperados

| Técnica | O participante consegue | Evidência de conclusão | Limite ou alternativa mais simples |
|---|---|---|---|
| Benchmark | Medir um workload explícito com `B.Loop` e comparar amostras com `benchstat`. | Quatro sub-benchmarks seriais sem setup ou população no loop, com interpretação de variância e significância. | Não prova superioridade universal nem contenção não medida. Se a pergunta for apenas correção para casos conhecidos, testes comuns bastam. |
| Fuzzing | Converter uma propriedade em fuzz target e uma falha em regressão permanente. | O panic é reproduzível sem busca na versão defeituosa e passa após a correção. | Não prova ausência de falhas para todas as entradas. Poucos casos conhecidos podem ser cobertos por testes comuns. |
| Lifecycle e evidência | Atrelar recursos e evidências ao lifecycle de `testing.T`. | O worker termina com `T.Context` e `T.Cleanup`; atributo, saída e artefato são observáveis. | Não prova regras de negócio. Sem trabalho assíncrono ou evidência persistente, o lifecycle padrão pode bastar. |
| Testes de contrato | Aplicar `fstest.TestFS` e `iotest.TestReader` sem reimplementar parcialmente seus verificadores. | Os contratos expõem as violações antes da correção e passam depois dela. | Não prova regras específicas do produto. Um caso direto basta quando não há um conjunto reutilizável de regras. |
| Failure injection | Injetar uma falha determinística de `Sync` e proteger a visibilidade do estado em memória. | O erro permanece na cadeia e a chave não é publicada. | Não prova rollback, replay ou durabilidade. Use o mecanismo direto da dependência quando ele já produzir o erro necessário. |
| `testing/synctest` | Controlar tempo e quiescência e distinguir o prazo de expiração da reação da goroutine. | O tempo de vida (TTL) e o janitor passam em milissegundos; a suíte também passa com `-race`. | Não prova ausência de data races nem cobre dependências externas à bubble. Uma pós-condição direta basta quando tempo e concorrência não determinam o resultado. |
| `goroutineleak` | Detectar goroutines bloqueadas e inalcançáveis e localizar o worker pelo stack. | A versão defeituosa aumenta o perfil; a corrigida zera o delta, e outro teste prova que `Stop` espera o worker. | Não prova ausência de todos os leaks. Para um worker conhecido, esperar seu canal `done` é uma pós-condição mais direta. |

## Lições bônus fora da agenda

A [lição 08](exercises/08-mutation-testing/README.md) percorre em 20 minutos o
ciclo suíte verde, mutante sobrevivente, hipótese, teste de fronteira e mutante
morto. Ela só pode ser usada em sala quando os sete blocos obrigatórios
terminarem com pelo menos 20 minutos de antecedência.

A [lição 09](exercises/09-property-tests/README.md) aprofunda testes baseados em
propriedades com invariantes, idempotência, comparação com um modelo e relações
metamórficas, além de geração, minimização e análise por mutantes. São 100
minutos opcionais de estudo independente.

Nenhuma lição bônus altera os sete resultados obrigatórios, os 225 minutos de
laboratórios, o intervalo ou o fechamento.

## Fechamento

Cada pessoa registra uma técnica que aplicará em uma suíte real, a pergunta que
ela responderá, uma limitação da evidência e quando uma abordagem mais simples
seria suficiente.

## Fora de escopo

Table tests, mocking frameworks, testcontainers, cobertura como meta, golden files, gRPC, linters e testes de integração não cabem com profundidade junto dos sete laboratórios. Os tópicos podem aparecer como respostas curtas, sem substituir o tempo de prática.
