# Lição bônus 08 — mutation testing

O contrato de `EligibleForFreeShipping` é inclusivo: um pedido recebe frete
grátis quando seu total é **igual ou superior** ao mínimo. A produção já está
correta. O trabalho deste exercício é encontrar uma fronteira que a suíte verde
executa de forma insuficiente e melhorar o teste.

## 1. Confirme a suíte verde, 3 minutos

Execute da raiz do repositório:

```bash
go test -v ./exercises/08-mutation-testing
```

Os casos abaixo e acima do mínimo devem passar. O caso exatamente igual ao
mínimo aparece como `SKIP`; ele é bônus e não contém `TODO:` obrigatório.

## 2. Produza um mutante sobrevivente, 4 minutos

O Gremlins `v0.6.0` já está fixado como tool do módulo. Execute o gate local:

```bash
bash scripts/mutation-check.sh ./exercises/08-mutation-testing
```

O script resolve o Gremlins pinado pelo módulo e copia somente os arquivos Go
da lição para um módulo temporário. Isso evita que a ferramenta replique o deck
e outras dependências locais em seu workdir. A execução isolada mantém
`--workers=1`, `--timeout-coefficient=30` e `--threshold-efficacy=100`.

O gate deve terminar com status diferente de zero. Procure esta evidência:

```text
LIVED CONDITIONALS_BOUNDARY
```

O mutante troca `>=` por `>`. Os testes acima e abaixo continuam verdes, então
a suíte não distingue o limite inclusivo do exclusivo.

## 3. Formule a hipótese, 3 minutos

Antes de editar, responda:

- qual entrada produz resultados diferentes nas duas versões?
- esse resultado faz parte do contrato declarado?
- qual menor teste diferenciaria a produção do mutante?

`LIVED` é uma hipótese para investigação. Pode revelar uma lacuna real, um
mutante equivalente ou um comportamento fora do contrato. Neste exemplo, o
limite inclusivo torna a lacuna observável.

## 4. Teste a fronteira, 7 minutos

Abra `exercises/08-mutation-testing/shipping_test.go`. Remova somente o
`t.Skip` do subtest `igual ao mínimo`. Depois execute:

```bash
go test -v ./exercises/08-mutation-testing
bash scripts/mutation-check.sh ./exercises/08-mutation-testing
```

Agora todos os mutantes devem aparecer como `KILLED`, a eficácia deve ser
`100.00%` e o threshold deve passar.

## 5. Encerre com evidência, 3 minutos

A execução demonstrou que esta suíte diferencia as mutações geradas pelo
Gremlins `v0.6.0` para esta função. Ela não prova que o contrato inteiro está
correto, que todo defeito possível seria detectado ou que cada mutante
sobrevivente em outro código representa uma falha real.

## Concluído quando

- a suíte comum permanece verde;
- o starter produz exatamente um `LIVED CONDITIONALS_BOUNDARY`;
- você explica por que o valor igual ao mínimo é a entrada discriminante;
- o novo teste mata o mutante sem alterar a produção;
- a execução final alcança `100.00%` de eficácia.

O gate local existe porque o Gremlins `v0.6.0` não converte corretamente o flag
numérico `--threshold-efficacy` recebido pela CLI. Ele não recalcula a métrica:
apenas exige da saída do próprio Gremlins `Lived: 0` e eficácia `100.00%`.
