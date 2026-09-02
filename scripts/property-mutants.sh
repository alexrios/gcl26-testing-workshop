#!/usr/bin/env bash
set -euo pipefail

repo_dir=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
cd "$repo_dir"

property_tests_dir=${1:-exercises/09-property-tests}
if [[ ! -d "$property_tests_dir" ]]; then
	echo "Diretorio de property tests inexistente: $property_tests_dir" >&2
	exit 1
fi
observe_derived=true

mutation_dir=$(mktemp -d "${TMPDIR:-/tmp}/testing-workshop-property-mutants.XXXXXX")
if [[ -z "$mutation_dir" || ! -d "$mutation_dir" ]]; then
	echo "Nao foi possivel criar o diretorio temporario dos mutantes." >&2
	exit 1
fi
trap 'rm -rf "$mutation_dir"' EXIT

(
	cd "$mutation_dir"
	go mod init workshop-property-mutants >/dev/null 2>&1
	go mod edit -go=1.27
)

targets=(
	FuzzNormalizeCanonical
	FuzzNormalizeIdempotent
	FuzzNormalizeMatchesModel
	FuzzEquivalentSeparators
	FuzzDerivedProperty
)

expected_status() {
	case "$1:$2" in
		drops-underscore:FuzzNormalizeMatchesModel | drops-underscore:FuzzEquivalentSeparators)
			echo KILL
			;;
		double-separator:FuzzNormalizeCanonical | double-separator:FuzzNormalizeMatchesModel)
			echo KILL
			;;
		canonical-cycle:FuzzNormalizeIdempotent | canonical-cycle:FuzzNormalizeMatchesModel)
			echo KILL
			;;
		*)
			echo PASS
			;;
	esac
}

prepare_package() {
	local implementation=$1
	local package_dir=$2

	mkdir -p "$package_dir/testdata/fuzz"
	cp "$implementation" "$package_dir/normalize.go"
	cp "$property_tests_dir"/*_test.go "$package_dir/"
	for target in "${targets[@]}"; do
		cp -R "support/property-corpus/$target" \
			"$package_dir/testdata/fuzz/$target"
	done
}

printf '%-20s %-28s %-6s\n' MUTANTE PROPRIEDADE STATUS

for mutant_dir in support/property-mutants/*; do
	mutant=$(basename "$mutant_dir")
	implementation="$mutant_dir/normalize.go"
	package_dir="$mutation_dir/$mutant"
	prepare_package "$implementation" "$package_dir"

	for target in "${targets[@]}"; do
		output_file="$mutation_dir/$mutant-$target.log"
		if (cd "$package_dir" && go test -run="^$target$" .) >"$output_file" 2>&1; then
			actual=PASS
		else
			actual=KILL
		fi
		expected=$(expected_status "$mutant" "$target")
		printf '%-20s %-28s %-6s\n' "$mutant" "$target" "$actual"
		if [[ "$target" == FuzzDerivedProperty && "$observe_derived" == true ]]; then
			continue
		fi
		if [[ "$actual" != "$expected" ]]; then
			echo "Esperado $expected para $mutant com $target." >&2
			cat "$output_file" >&2
			exit 1
		fi
	done
done

echo "PROPERTY MUTANTS OK"
