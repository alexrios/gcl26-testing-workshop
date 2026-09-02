#!/usr/bin/env bash
set -euo pipefail

repo_dir=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
cd "$repo_dir"

bash -n scripts/*.sh

run_quiet() {
	local name=$1
	shift

	local output
	if ! output=$("$@" 2>&1); then
		echo "PREFLIGHT FAIL [$name]" >&2
		printf '%s\n' "$output" >&2
		exit 1
	fi
}

unformatted=$(gofmt -l exercises support/property-mutants)
if [[ -n "$unformatted" ]]; then
	echo "Arquivos sem gofmt:"
	echo "$unformatted"
	exit 1
fi

run_quiet go-mod-tidy go mod tidy -diff

if ! GOPROXY=off go tool benchstat -h >/dev/null 2>&1; then
	echo "benchstat indisponível offline; execute go mod download com acesso à rede"
	exit 1
fi
if ! GOPROXY=off go tool -n gremlins >/dev/null 2>&1; then
	echo "gremlins indisponível offline; execute go mod download com acesso à rede"
	exit 1
fi

run_quiet tests go test ./...
run_quiet benchmark-compile go test -run='^$' -bench=. -benchtime=1x \
	./exercises/01-benchmarks
run_quiet race go test -race \
	./exercises/06-synctest ./exercises/07-goroutine-leaks
run_quiet vet go vet ./...

echo "OK"
