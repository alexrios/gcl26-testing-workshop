# Lição bônus 09 — projetando property tests

Esta lição não faz parte dos 240 minutos do workshop e não bloqueia
`mise run progress`.

Pré-requisito: conclua o [exercício 2](../02-fuzzing/README.md). O mecanismo de
busca continua sendo o fuzzing nativo do Go. Nesta lição, o foco é a qualidade
da propriedade: como derivá-la, como construir entradas úteis e como descobrir
quais defeitos ainda sobrevivem.

`NormalizeKey` converte uma entrada arbitrária para uma chave canônica. O
contrato público é:

- letras `ASCII` (`A-Z` e `a-z`) viram minúsculas;
- dígitos são preservados;
- cada sequência de outros bytes vira um único `-` entre termos;
- o resultado não começa nem termina com `-`.

Os testes de exemplo passam. Existe um defeito intencional não coberto por eles.
Não abra `normalize.go` antes do primeiro contraexemplo.

## O mapa usado na lição

Para cada property test, registre cinco decisões:

1. **Domínio:** para quais entradas a afirmação deveria valer?
2. **Construção:** como os argumentos do fuzzer representam esse domínio?
3. **Oráculo:** qual relação observável deve permanecer verdadeira?
4. **Independência:** que código ou hipótese a produção e o oráculo compartilham?
5. **Poder:** qual implementação errada ainda passaria?

Uma execução verde significa somente que aquela busca não encontrou um
contraexemplo. Uma propriedade pode estar verde e ainda ser tautológica, fraca
demais ou baseada no mesmo defeito da produção.

## 1. Derive antes de receber fórmulas, 15 minutos

Leia apenas o contrato acima e liste pelo menos cinco afirmações candidatas. Não
escreva exemplos concretos como `NormalizeKey("Go Test") == "go-test"`. Escreva
relações que façam sentido para conjuntos de entradas.

Para cada candidata, anote:

- o domínio em que ela vale;
- se precisa saber a saída exata;
- um defeito que ela poderia encontrar;
- um defeito que provavelmente sobreviveria.

Escolha uma candidata e implemente `FuzzDerivedProperty`. O scaffold não fornece
a fórmula deliberadamente. Use somente o contrato público.

```bash
go test -run=FuzzDerivedProperty ./exercises/09-property-tests
go test -run='^$' -fuzz=FuzzDerivedProperty -fuzztime=10s \
  ./exercises/09-property-tests
```

Se aparecer um contraexemplo, não conclua imediatamente que a produção está
errada. Primeiro decida se a entrada refuta o contrato, a propriedade ou sua
implementação do oráculo.

Checkpoint: existe uma propriedade escrita em linguagem natural, seu domínio
está explícito e o target contém uma mensagem que explica a relação quebrada.

## 2. Clínica de propriedades ruins, 15 minutos

Execute as quatro candidatas fornecidas em `property_clinic_test.go`:

```bash
go test -run=TestPropertyClinic -v ./exercises/09-property-tests
```

Classifique cada uma antes de seguir:

| Candidata | Válida pelo contrato? | Útil contra algum defeito? | Oráculo independente? |
|---|---|---|---|
| `same-call` | ? | ? | ? |
| `never-panics` | ? | ? | ? |
| `preserves-length` | ? | ? | ? |
| `production-alphabet` | ? | ? | ? |

Use estes diagnósticos:

- **tautológica:** é verdadeira por construção, mesmo com produção errada;
- **fraca:** observa menos que o contrato disponível;
- **forte demais:** rejeita comportamento permitido pelo contrato;
- **acoplada:** reutiliza decisões ou helpers que podem compartilhar o defeito;
- **útil:** elimina pelo menos uma implementação errada plausível sem rejeitar
  comportamento permitido.

O resultado booleano não fornece a classificação sozinho. `preserves-length`,
por exemplo, ficar falso para uma entrada válida pode significar que a
propriedade é forte demais, não que a produção esteja errada.

Checkpoint: você consegue explicar por que `same-call` e `never-panics` podem
ficar verdes para muitas implementações determinísticas incorretas, e por que
`production-alphabet` não é um oráculo independente.

Depois de registrar suas respostas, compare os diagnósticos com outra pessoa e
explique qual implementação incorreta cada candidata ainda aceitaria.

## 3. Separe forma canônica de idempotência, 15 minutos

Agora implemente dois targets independentes:

```text
isCanonical(NormalizeKey(input))

NormalizeKey(NormalizeKey(input)) == NormalizeKey(input)
```

Implemente `isCanonical` diretamente a partir do contrato. Não chame
`NormalizeKey`, `isASCIIAlphaNumeric` nem outro helper da produção. Separar os
targets preserva o diagnóstico: uma implementação pode produzir apenas strings
canônicas e ainda mudar uma chave canônica na segunda chamada.

```bash
go test -run='FuzzNormalize(Canonical|Idempotent)' \
  ./exercises/09-property-tests
go test -run='^$' -fuzz=FuzzNormalizeCanonical -fuzztime=10s \
  ./exercises/09-property-tests
go test -run='^$' -fuzz=FuzzNormalizeIdempotent -fuzztime=10s \
  ./exercises/09-property-tests
```

As duas buscas devem permanecer verdes mesmo com o defeito intencional. Isso é
evidência de que a saída tem a forma esperada e é um ponto fixo, não de que cada
byte da entrada recebeu a interpretação correta.

Checkpoint: registre uma implementação errada que passaria pelas duas
propriedades.

## 4. Use um modelo sem copiar a produção, 15 minutos

Implemente `modelNormalize` com operações estruturalmente diferentes. Uma rota
possível é:

1. substituir `[^A-Za-z0-9]+` por `-`;
2. aplicar `strings.ToLower` somente depois da substituição;
3. remover `-` das extremidades.

A ordem importa. Aplicar Unicode lowercase antes de excluir caracteres não
`ASCII` mudaria o domínio descrito pelo contrato e poderia criar um falso
contraexemplo.

Antes de comparar produção e modelo, remova o skip de `TestModelExamples`. Em
seguida, adicione casos diretos para `ASCII`, separadores repetidos, Unicode e
UTF-8 inválido:

```bash
go test -run=TestModelExamples -v ./exercises/09-property-tests
```

Implemente:

```text
NormalizeKey(input) == modelNormalize(input)
```

```bash
go test -run='^$' -fuzz=FuzzNormalizeMatchesModel -fuzztime=20s \
  ./exercises/09-property-tests
```

Leia o input minimizado e reproduza sem busca:

```bash
go test -run=FuzzNormalizeMatchesModel ./exercises/09-property-tests
```

O modelo cobre mais do contrato que os invariantes anteriores, mas também cria
um segundo artefato que pode estar errado. Mantenha exemplos diretos para o
modelo e revise divergências antes de corrigir a produção.

Checkpoint: o contraexemplo reproduz sem `-fuzz` e você consegue explicar por
que forma canônica e idempotência não o detectaram.

## 5. Controle a construção e escreva uma relação metamórfica, 15 minutos

Uma propriedade metamórfica compara execuções relacionadas sem calcular a saída
exata. Pelo contrato, trocar um byte separador por outro deve preservar o
resultado:

```text
NormalizeKey(left + separator + right) ==
NormalizeKey(left + "/" + right)
```

O fuzzer gera `selector uint8`. Transforme-o por módulo em um índice de:

```text
' ', '_', '-', '/', '.', 0x00, 0xff
```

Essa representação constrói somente separadores relevantes. Usar um byte
arbitrário e chamar `t.Skip` para letras e dígitos faria parte do orçamento de
busca ser desperdiçada em entradas rejeitadas.

```bash
go test -run='^$' -fuzz=FuzzEquivalentSeparators -fuzztime=20s \
  ./exercises/09-property-tests
```

Compare os contraexemplos do modelo e da relação metamórfica. O Go minimiza os
argumentos suportados pelo target, não o conceito de “caso mais simples” do seu
domínio. A representação escolhida determina o que pode ser reduzido. A API
`testing.F` não oferece um gerador ou shrinker customizado; quando a estrutura
for complexa, codifique-a nos tipos suportados e faça a transformação dentro do
target.

Checkpoint: explique por que mapear `selector` constrói o domínio desejado e
qual distribuição essa transformação induz sobre os sete separadores.

## 6. Preserve a evidência e corrija, 10 minutos

Mantenha os arquivos criados em `testdata/fuzz`. Só agora abra `normalize.go`,
explique qual ramo viola o contrato e faça a menor correção que trate todos os
bytes separadores da mesma maneira.

```bash
go test ./exercises/09-property-tests
for target in FuzzNormalizeCanonical FuzzNormalizeIdempotent \
  FuzzNormalizeMatchesModel FuzzEquivalentSeparators FuzzDerivedProperty; do
  go test -run='^$' -fuzz="^${target}$" -fuzztime=3s \
    ./exercises/09-property-tests
done
```

Checkpoint: os contraexemplos preservados passam e a correção não removeu nem
suavizou nenhuma propriedade.

## Se a busca não encontrar a falha

Use o contraexemplo preparado, que é um corpus real do target:

```bash
mkdir -p exercises/09-property-tests/testdata/fuzz/FuzzNormalizeMatchesModel
cp support/property-tests-crasher \
  exercises/09-property-tests/testdata/fuzz/FuzzNormalizeMatchesModel/underscore
go test -run=FuzzNormalizeMatchesModel ./exercises/09-property-tests
```

## 7. Meça força com mutantes, 10 minutos

Antes de executar, preveja `PASS` ou `KILL` para cada célula:

| Mutante | Forma canônica | Idempotência | Modelo | Metamórfica | Sua derivada |
|---|---:|---:|---:|---:|---:|
| descarta `_` | ? | ? | ? | ? | ? |
| preserva separadores repetidos | ? | ? | ? | ? | ? |
| alterna entre `a` e `b` | ? | ? | ? | ? | ? |

Depois compare sua hipótese com o laboratório local:

```bash
mise run property-mutants
```

O comando executa seus targets contra três mutantes preparados. A coluna
`FuzzDerivedProperty` será apenas observada, porque propriedades derivadas
válidas podem eliminar conjuntos diferentes de mutantes.

`KILL` significa que um contraexemplo conhecido faz aquela propriedade rejeitar
o mutante. Não é uma pontuação universal: os três mutantes são uma amostra
pedagógica escolhida. Observe especialmente que o modelo elimina mais mutantes,
enquanto propriedades menores localizam melhor o mecanismo quebrado e carregam
menos lógica de oráculo.

Checkpoint: nenhuma propriedade isolada é apresentada como suficiente e cada
linha da matriz tem uma explicação causal.

## 8. Debrief, 5 minutos

Para sua propriedade derivada e para cada target fornecido, responda:

- Que parte do contrato ele observa?
- Qual mutante ele consegue eliminar sem apoio dos outros targets?
- Qual implementação incorreta ainda passa?
- Que hipótese é compartilhada com a produção?
- A construção concentra a busca no domínio relevante?
- A falha mínima continua semanticamente legível?

Revise então `FuzzDerivedProperty`. Não existe uma propriedade oficial. Se a sua
for válida, independente e eliminar uma classe plausível de defeitos, preserve-a.

## Concluído quando

- você derivou ao menos uma propriedade sem receber a fórmula;
- classificou propriedades tautológicas, fracas, fortes demais e acopladas;
- forma canônica e idempotência têm targets e diagnósticos separados;
- um modelo independente e uma relação metamórfica encontram o defeito inicial;
- você explicou como representação, filtragem e minimização interagem;
- a matriz de mutantes demonstra forças diferentes entre as propriedades;
- os casos encontrados reproduzem sem busca e passam depois da correção.
