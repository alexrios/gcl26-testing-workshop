#!/usr/bin/env bash
set -euo pipefail

repo_dir=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
cd "$repo_dir"

required_exercises=(
	./exercises/01-benchmarks
	./exercises/02-fuzzing
	./exercises/03-lifecycle
	./exercises/04-contracts
	./exercises/05-failure-injection
	./exercises/06-synctest
	./exercises/07-goroutine-leaks
)

pending=$(
	for exercise_dir in "${required_exercises[@]}"; do
		grep -nH 'TODO:' "$exercise_dir"/*.go || true
	done
)
if [[ -n "$pending" ]]; then
	echo "Exercícios ainda têm marcadores obrigatórios:"
	echo "$pending"
	echo
	echo "Remova cada TODO somente depois de implementar e validar o exercício."
	exit 1
fi

run_quiet() {
	local name=$1
	shift

	local output
	if ! output=$("$@" 2>&1); then
		echo "PROGRESS FAIL [$name]" >&2
		printf '%s\n' "$output" >&2
		exit 1
	fi
}

run_quiet tests go test "${required_exercises[@]}"
run_quiet benchmarks go test -run='^$' \
	-bench='Benchmark(Get|SetOverwrite)$' -benchtime=1x \
	./exercises/01-benchmarks
run_quiet race go test -race \
	./exercises/06-synctest ./exercises/07-goroutine-leaks

echo "PROGRESS OK"
