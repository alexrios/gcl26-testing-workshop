#!/usr/bin/env bash
set -euo pipefail

repo_dir=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
cd "$repo_dir"

target=${1:-./exercises/08-mutation-testing}
target_dir=$(cd "$target" && pwd)
gremlins=$(GOPROXY=off go tool -n gremlins)
mutation_dir=$(mktemp -d "${TMPDIR:-/tmp}/testing-workshop-mutation.XXXXXX")
if [[ -z "$mutation_dir" || ! -d "$mutation_dir" ]]; then
	echo "Não foi possível criar o módulo temporário do mutation testing." >&2
	exit 1
fi
trap 'rm -rf "$mutation_dir"' EXIT

cp "$target_dir"/*.go "$mutation_dir/"
(
	cd "$mutation_dir"
	GOWORK=off go mod init workshop-mutation >/dev/null 2>&1
)

set +e
output=$(
	cd "$mutation_dir"
	GOPROXY=off GOWORK=off "$gremlins" unleash . \
		--workers=1 --timeout-coefficient=30 --threshold-efficacy=100 2>&1
)
status=$?
set -e

printf '%s\n' "$output"

if [[ $status -ne 0 ]]; then
	exit "$status"
fi

# Gremlins v0.6.0 lê flags numéricos do Viper como strings e não aplica o
# threshold pela CLI. Preserve o comando documentado e valide o mesmo contrato
# pela saída até que a versão pinada corrija o comportamento.
if ! grep -Fq "Lived: 0" <<<"$output" ||
	! grep -Fq "Test efficacy: 100.00%" <<<"$output"; then
	echo "MUTATION FAIL: eficácia abaixo do threshold de 100%." >&2
	exit 1
fi

echo "MUTATION OK"
