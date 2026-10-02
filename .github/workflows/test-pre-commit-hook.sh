#!/usr/bin/env bash
# Installs the pre-commit hook from this checkout (wheel bundling the binary of
# the release pinned in setup.cfg) and runs it on copies of the testdata cases.
set -euo pipefail

cd "$(git rev-parse --show-toplevel)"

pre-commit validate-manifest .pre-commit-hooks.yaml

tmp=$(mktemp -d ./hook-test.XXXXXX)
trap 'rm -rf "$tmp"' EXIT

# Each case as <case>/BUILD; basic also under the other names the hook matches,
# plus one name it must ignore.
files=()
for dir in testdata/*/; do
    case=$(basename "$dir")
    mkdir -p "$tmp/$case"
    cp "${dir}input" "$tmp/$case/BUILD"
    files+=("$tmp/$case/BUILD")
done
for name in BUILD.bazel foo.BUILD foo.BUILD.bazel not_a_build_file; do
    cp testdata/basic/input "$tmp/basic/$name"
    files+=("$tmp/basic/$name")
done

# The copies are untracked, so pre-commit can't see that the hook modified them
# and only fails if the binary does. The diffs below check the actual result.
pre-commit try-repo . bazel-build-file-sorter --files "${files[@]}"

failed=0
for dir in testdata/*/; do
    case=$(basename "$dir")
    diff -u "${dir}expected" "$tmp/$case/BUILD" || failed=1
done
for name in BUILD.bazel foo.BUILD foo.BUILD.bazel; do
    diff -u testdata/basic/expected "$tmp/basic/$name" || failed=1
done
diff -u testdata/basic/input "$tmp/basic/not_a_build_file" || failed=1
exit $failed
