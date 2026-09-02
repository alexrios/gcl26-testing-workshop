# Exercício 4 — contratos executáveis

Testes de exemplo confirmam casos escolhidos por você. Um teste de contrato percorre regras de um protocolo e pode ser reaplicado a qualquer implementação. Aqui você usará contratos prontos da biblioteca padrão em vez de reescrevê-los parcialmente.

## 1. Diagnostique um `fs.FS` inválido, 10 minutos

Sem abrir `contracts.go`, implemente `TestFixtureFSContract` com `fstest.TestFS`:

1. remova o `t.Skip`;
2. passe `FixtureFS()` e o arquivo esperado;
3. reporte o erro retornado pelo verificador com `t.Fatal`.

```bash
go test -run=TestFixtureFSContract -v ./exercises/04-contracts
```

A primeira execução deve falhar. Leia o diagnóstico, identifique qual caminho válido o adaptador rejeitou e corrija `fixtureFS.Open` sem criar exceções para o arquivo esperado. Depois execute o contrato novamente.

O teste não verifica apenas se um arquivo abre. Ele exercita as regras de `fs.FS`, `fs.File` e dos métodos opcionais que a implementação oferece. Se travar, use o diagnóstico do contrato para localizar a regra violada em `contracts.go`.

`fstest.TestFS` evita uma bateria manual incompleta: percorre a árvore, verifica
arquivos e diretórios e confirma que os caminhos esperados existem. É útil para
`embed.FS`, filesystems em memória e adapters que adicionam prefixos, restringem
caminhos ou envolvem outra implementação de `fs.FS`. O contrato verifica o
filesystem; ele não valida o significado do conteúdo para o produto.

## 2. O contrato de `io.Reader`, 10 minutos

Implemente `TestChunkReaderContract`:

- crie um payload conhecido;
- construa `NewChunkReader(payload, 3)`;
- passe o reader e o conteúdo esperado para `iotest.TestReader`;
- reporte o erro do verificador com `t.Fatal`.

O verificador varia tamanhos de buffer, lê byte a byte, procura além do fim e
confirma o comportamento de `io.EOF`. Uma única chamada a `io.ReadAll` não
cobriria esse contrato.

`iotest.TestReader` é útil para readers em chunks, wrappers que limitam ou
transformam bytes, buffers e adapters de formatos ou protocolos. Se o valor
também implementar `io.ReaderAt` ou `io.Seeker`, o verificador exercita esses
contratos. Ele confirma o comportamento de leitura, não as regras de negócio do
conteúdo.

```bash
go test -run=TestChunkReaderContract -v ./exercises/04-contracts
```

O teste deve falhar primeiro porque o starter reporta uma contagem maior que o número de bytes copiados. Corrija `ChunkReader.Read` e execute novamente.

## Perguntas, 5 minutos

- O que `fstest.TestFS` verifica que um teste de exemplo provavelmente esqueceria?
- O contrato pertence ao produtor da interface ou a cada implementação?
- Quando um teste de contrato fica grande demais para diagnosticar uma falha?

## Concluído quando

- a implementação de `fs.FS` passa pelo verificador padrão;
- o contrato expõe a contagem inválida antes da correção;
- o reader que entrega dados em pequenos chunks passa pelo contrato de `io.Reader`;
- você consegue distinguir contrato reutilizável de uma coleção de casos de exemplo.

Rota mínima: priorize o diagnóstico e a correção do `ChunkReader`, que exercita mais variações de leitura. Se o tempo apertar, o instrutor conduz o diagnóstico de `TestFS` e a turma aplica a correção do adaptador.
