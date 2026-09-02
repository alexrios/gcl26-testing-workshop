# Go Testing Masterclass - GopherCon LATAM 2026

Material dos participantes para o workshop presencial de quatro horas sobre
recursos e técnicas recentes ou menos dominadas de teste em Go.

## Para quem é

O workshop pressupõe que você já:

- escreve testes de unidade em Go;
- conhece subtests e table-driven tests;
- entende goroutines, channels, interfaces e `context.Context`;
- consegue ler um benchmark simples com `b.N`.

Você pode participar sem conhecer fuzzing, `B.Loop`, `synctest` ou o perfil
`goroutineleak`.

## Preparação obrigatória

Use Go 1.27 diretamente ou instale-o com [mise](https://mise.jdx.dev/):

```bash
mise trust
mise install
mise exec -- go version
mise exec -- go mod download
mise run preflight
```

Sem mise:

```bash
go version
go mod download
bash scripts/verify.sh
```

A primeira linha de `go version` deve começar com `go version go1.27` e o
preflight deve terminar em `OK`. Depois de `go mod download`, todo o material
funciona offline. Docker, banco de dados, credenciais e serviços externos não
são necessários.

O workshop usa o race detector. Em Linux e outros sistemas não Darwin, ele
requer cgo, um compilador C instalado e uma plataforma suportada.

## Conteúdo

```text
exercises/
├── 01-benchmarks/         # B.Loop, workloads e benchstat
├── 02-fuzzing/            # propriedade, falha e regressão
├── 03-lifecycle/          # contexto, cleanup, saída e artefatos
├── 04-contracts/          # fstest e iotest
├── 05-failure-injection/  # erro controlado e estado parcial
├── 06-synctest/           # tempo virtual e HTTP em memória
├── 07-goroutine-leaks/    # perfil goroutineleak do Go 1.27
├── 08-mutation-testing/   # bônus: qualidade da suíte por mutantes
└── 09-property-tests/     # bônus: propriedades além do round-trip

support/                   # corpora de fallback e mutantes usados nos exercícios
```

Os sete laboratórios obrigatórios ocupam 225 minutos. O intervalo e o
fechamento completam as quatro horas. Consulte [SYLLABUS.md](SYLLABUS.md) para a
agenda e [REFERENCES.md](REFERENCES.md) para a documentação primária.

As lições 08 e 09 são opcionais. Elas não bloqueiam `mise run progress` nem
alteram a agenda obrigatória.

## Fluxo de trabalho

Cada exercício contém contexto, tarefas, comandos e critério de conclusão. Os
starters compilam e marcam trabalho obrigatório com `TODO:` e `t.Skip`. Remova
cada marcador somente depois de implementar e validar a etapa correspondente.

Use os comandos documentados no `README.md` de cada exercício. Para acompanhar
o conjunto obrigatório:

```bash
mise run progress
```

No checkout inicial, esse comando deve falhar listando os `TODO:` ainda
pendentes. Depois que os sete exercícios obrigatórios estiverem concluídos, ele
deve terminar em `PROGRESS OK`.

## Comandos úteis

```bash
mise run test              # executa testes e corpora já preservados
mise run bench             # exercício 01
mise run fuzz              # exercício 02
mise run lifecycle         # exercício 03
mise run contracts         # exercício 04
mise run failure           # exercício 05
mise run synctest          # exercício 06, também com -race
mise run leaks             # exercício 07, também com -race
mise run mutation          # exercício bônus 08
mise run property          # exercício bônus 09
mise run property-mutants  # força das propriedades da lição 09
mise run preflight         # ambiente e integridade do kit
mise run progress          # progresso nos sete exercícios obrigatórios
```

`preflight` verde significa que o kit está pronto para o workshop. Não significa
que os exercícios foram resolvidos.
