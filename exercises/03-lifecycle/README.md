# Exercício 3 — lifecycle e evidência de testes

O objetivo é fazer o próprio teste controlar o tempo de vida dos recursos e preservar evidência útil quando algo falha.

## 1. Contexto e cleanup, 10 minutos

Implemente `TestWorkerStopsWithTest`:

1. remova o `t.Skip`;
2. inicie o worker com `t.Context()`;
3. registre um `t.Cleanup` que espera o canal `done` fechar;
4. não crie um `context.WithCancel` manual.

O teste termina porque `t.Context()` é cancelado antes dos cleanups.

```bash
go test -run=TestWorkerStopsWithTest -v ./exercises/03-lifecycle
```

## 2. Metadados, saída e artefatos, 10 minutos

Implemente `TestReportArtifact`:

- associe `owner=platform` com `t.Attr`;
- grave `report.json` em `t.ArtifactDir()`;
- escreva o caminho usando `fmt.Fprintf(t.Output(), ...)`.

`t.Attr` associa metadados curtos à execução, em vez de exigir que uma
ferramenta extraia campos de uma mensagem de log. Casos comuns incluem
`owner`, `component`, `kind`, `implementation` e `scenario`. Use atributos para
classificar resultados; mantenha diagnóstico humano em `t.Log` ou `t.Output` e
evidências maiores em `t.ArtifactDir`.

Com `-v`, o atributo aparece no log:

```text
=== ATTR  TestReportArtifact owner platform
```

Com `go test -json`, a mesma emissão também produz um evento estruturado:

```json
{"Action":"attr","Test":"TestReportArtifact","Key":"owner","Value":"platform"}
```

Um coletor pode associar o evento ao pacote e ao teste, agrupar ou filtrar
resultados e encaminhar uma falha ao owner sem interpretar texto livre. O Go
não define o significado das chaves: a ferramenta ou a organização estabelece
esse vocabulário. O runner precisa consumir a ação `attr` ou traduzi-la para seu
formato; integrações que conhecem apenas eventos antigos podem ignorá-la. A
ordem dos atributos não é significativa; a chave não pode conter espaços e o
valor não pode conter quebras de linha.

Compare:

```bash
go test -run=TestReportArtifact -v ./exercises/03-lifecycle
mkdir -p artifacts
go test -count=1 -run=TestReportArtifact -artifacts \
  -outputdir="$(pwd)/artifacts" \
  ./exercises/03-lifecycle
```

Sem `-artifacts`, a evidência é temporária. Com a flag, ela permanece abaixo de `./artifacts/_artifacts`. O `-count=1` impede que o cache reutilize uma execução que não produziu o arquivo neste comando.

## Perguntas, 5 minutos

- Por que `t.Context` é cancelado antes dos cleanups?
- Quando usar `t.Output` em vez de um buffer?
- Que atributo estável ajudaria uma ferramenta a classificar este teste?
- Qual evidência ajudaria a diagnosticar localmente uma falha real?

## Concluído quando

- o worker termina sem cancelamento manual;
- o teste emite um atributo e uma saída associada ao próprio teste;
- `report.json` só permanece quando `-artifacts` é informado em uma execução real.

Rota mínima: conclua o worker com `T.Context` e `T.Cleanup`. Se a turma atrasar, o instrutor conduz `T.Attr`, `T.Output` e `T.ArtifactDir`, e todos executam a solução para observar a evidência.
