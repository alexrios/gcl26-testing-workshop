# Referências

O workshop privilegia documentação primária. Links de bibliotecas externas aparecem apenas quando a ferramenta é usada diretamente.

## Go 1.27

- [Go 1.27 Release Notes](https://go.dev/doc/go1.27)
- [`testing/synctest`](https://pkg.go.dev/testing/synctest)
- [`synctest.Sleep`](https://pkg.go.dev/testing/synctest#Sleep)
- [`httptest.NewTestServer`](https://pkg.go.dev/net/http/httptest#NewTestServer)
- [Perfis de `runtime/pprof`](https://pkg.go.dev/runtime/pprof#Profile)

## Benchmarks

- [`testing.B`](https://pkg.go.dev/testing#B)
- [`testing.B.Loop`](https://pkg.go.dev/testing#B.Loop)
- [`testing.B.RunParallel`](https://pkg.go.dev/testing#B.RunParallel)
- [`benchstat`](https://pkg.go.dev/golang.org/x/perf/cmd/benchstat)
- [Go 1.26: correção de inlining em `B.Loop`](https://go.dev/doc/go1.26#testingpkgtesting)

## Fuzzing

- [Tutorial oficial de fuzzing](https://go.dev/doc/tutorial/fuzz)
- [Fuzzing no pacote `testing`](https://pkg.go.dev/testing#hdr-Fuzzing)
- [Go fuzzing design draft](https://go.dev/s/draft-fuzzing-design)
- [QuickCheck: a lightweight tool for random testing of Haskell programs](https://www.cs.tufts.edu/~nr/cs257/archive/john-hughes/quick.pdf)
- [Metamorphic Testing: A New Approach for Generating Next Test Cases](https://www.cse.ust.hk/faculty/scc/publ/CS98-01-metamorphictesting.pdf)

## Lifecycle e artefatos

- [`T.Context`](https://pkg.go.dev/testing#T.Context)
- [`T.Cleanup`](https://pkg.go.dev/testing#T.Cleanup)
- [`T.Output`](https://pkg.go.dev/testing#T.Output)
- [`T.Attr`](https://pkg.go.dev/testing#T.Attr)
- [`T.ArtifactDir`](https://pkg.go.dev/testing#T.ArtifactDir)

## Contratos e failure injection

- [`testing/fstest.TestFS`](https://pkg.go.dev/testing/fstest#TestFS)
- [`testing/iotest`](https://pkg.go.dev/testing/iotest)
- [`errors.Is`](https://pkg.go.dev/errors#Is)

## Detecção de goroutine leaks

- [Perfil `goroutineleak` nas release notes do Go 1.27](https://go.dev/doc/go1.27#runtime)
- [`runtime/pprof.Lookup`](https://pkg.go.dev/runtime/pprof#Lookup)
- [Série *Goroutine Leak Detection*](https://alexrios.me/series/Goroutine%20Leak%20Detection/) — escrita sobre a API experimental do Go 1.26; use como contexto e confirme os comandos na documentação do Go 1.27
- [`go.uber.org/goleak`](https://pkg.go.dev/go.uber.org/goleak)
- [Race detector](https://go.dev/doc/articles/race_detector)

## Mutation testing — lição bônus 08

- [Instalação do Gremlins `v0.6.0`](https://gremlins.dev/latest/install/)
- [Comando `unleash`, flags e cálculo dos thresholds](https://gremlins.dev/latest/usage/commands/unleash/)
- [Statuses `KILLED`, `LIVED`, `NOT COVERED`, `TIMED OUT` e `NOT VIABLE`](https://gremlins.dev/latest/quick-start/)
- [Mutação `CONDITIONALS_BOUNDARY`](https://gremlins.dev/latest/usage/mutations/conditionals_boundary/)

## Para continuar depois do workshop

- [`testing/cryptotest`, Go 1.26](https://pkg.go.dev/testing/cryptotest)
- [`testing/slogtest`](https://pkg.go.dev/testing/slogtest)
