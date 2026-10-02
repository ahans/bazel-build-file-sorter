# bazel-build-file-sorter

## What this is

A CLI tool that sorts named rules (targets) in Bazel BUILD files alphabetically by name, using `github.com/bazelbuild/buildtools/build` — the same parser as buildozer and buildifier.

## Key design decisions

- **Non-rules (load, variables) move to the top** by default. This is intentional: it's always semantically safe (variables are defined before any rule that might use them). The alternative — keeping non-rules in place — risks sorting a rule above a variable it depends on.
- **`# nosort` pins a statement in place.** Pinned statements stay at their original index; everything else sorts around them. Applies to both rules and non-rules. Place the comment on the line before the statement.
- **Comments travel with their statement.** The buildtools AST attaches comments as `Before`/`Suffix` fields on nodes, so they move with whatever statement they're attached to — even when separated by a blank line.
- **Case-insensitive sort** (`strings.ToLower`).

## Future: package-named targets

Targets whose `name` matches the Bazel package name (e.g., `cc_library(name = "foo")` in `//foo/BUILD`) may warrant special sort treatment (e.g., always sorted first or excluded from sorting). Not yet implemented.

To support this, `build.File.WorkspaceRoot` needs to be populated — the parser leaves it empty. Auto-detect by walking up the directory tree from the BUILD file to find `WORKSPACE`, `WORKSPACE.bazel`, or `MODULE.bazel`. The package name is then `"//" + filepath.Dir(f.Path)` relative to `WorkspaceRoot`.

## pre-commit hook

The hook lives in `.pre-commit-hooks.yaml` with `language: python`. pre-commit installs `language: python` hooks by running `pip install .` in a fresh virtualenv, so this repo is packaged as a wheel that bundles the Go binary — following the same pattern as [shellcheck-py](https://github.com/shellcheck-py/shellcheck-py).

The wheel is built with [setuptools-download](https://github.com/asottile/setuptools-download): `setup.cfg` lists one `download_scripts` entry per platform (URL to the GitHub release asset + sha256, selected by a PEP 508 marker), and setuptools-download fetches the matching one at build time and installs it as a wheel *script*, so `pip install` drops `bazel-build-file-sorter` into the venv's `bin/`. `setup.py` only overrides `bdist_wheel` to mark the wheel non-pure and platform-tagged. Nothing is downloaded at runtime. The hook's `entry` is the binary itself (`bazel-build-file-sorter`), with `args: [-i]` so it sorts files in place. Marker values are Python's (`sys_platform` `linux`/`darwin`/`win32`, `platform_machine` `x86_64`/`aarch64`/`arm64`/`AMD64`/`ARM64`), not Go's.

**Releasing** happens locally, not in CI, because the sha256 hashes in `setup.cfg` must be committed under the release tag: run `./release.sh v0.x.y` on an up-to-date, clean `main`. It cross-compiles all binaries via `./build-all.sh <tag>` (which passes `VERSION=<tag>` to `make`, so the embedded version is right before the tag exists), regenerates the `download_scripts` section and `version` in `setup.cfg`, and bumps `rev:` in the README. `main` is protected (changes only via PR, required checks), so it commits on a `release/<tag>` branch, tags that commit, asks for confirmation, pushes branch and tag, uploads the binaries with `gh release create`, and opens a PR into `main`. The release is live as soon as the tag and assets are up; the PR's CI then installs the hook from the published binaries. Merge release PRs with a merge commit so the tag stays reachable from `main` (`git describe` relies on it). Do not reintroduce a tag-triggered release workflow — rebuilt binaries would not match the committed hashes.

## Test structure

Tests live in `testdata/<case-name>/input` and `testdata/<case-name>/expected`. The test runner in `main_test.go` discovers cases by reading the directory. Add a new case by creating a new subdirectory with both files.
