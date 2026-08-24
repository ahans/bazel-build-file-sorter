#!/usr/bin/env bash
# Cross-compiles bazel-build-file-sorter-<goos>-<goarch>[.exe] for all release platforms.
# An optional version (e.g. v0.3.0) is embedded instead of `git describe`.
set -euo pipefail

make_args=()
if [[ $# -ge 1 ]]; then
    make_args+=("VERSION=$1")
fi

for goos in linux darwin windows; do
    ext=''
    if [ "$goos" == "windows" ]; then
        ext='.exe'
    fi
    for goarch in amd64 arm64; do
        rm -f "bazel-build-file-sorter${ext}"
        GOOS=$goos GOARCH=$goarch make build "${make_args[@]+"${make_args[@]}"}"
        mv "bazel-build-file-sorter${ext}" "bazel-build-file-sorter-${goos}-${goarch}${ext}"
    done
done
