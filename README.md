# bazel-build-file-sorter

[![CI](https://github.com/ahans/bazel-build-file-sorter/actions/workflows/main.yml/badge.svg)](https://github.com/ahans/bazel-build-file-sorter/actions/workflows/main.yml)

Sorts targets in Bazel BUILD files alphabetically by name, using the same parser as [buildozer](https://github.com/bazelbuild/buildtools/tree/main/buildozer).

`load()` statements and variable assignments are left at the top in their original order. Only named rules (targets with a `name` attribute) are sorted.

## Pinning statements

Add a `# nosort` comment before any statement to keep it at its current position:

```python
# nosort
VARIABLE = [...]   # stays in place; rules sort around it

# nosort
cc_library(        # stays in place; other rules sort around it
    name = "...",
)
```

Without `# nosort`, variable assignments and other non-rule statements move to the top of the file ahead of all rules.

## pre-commit hook

Sorting BUILD files is also possible via a [pre-commit](https://pre-commit.com/) hook.
Add this to your `.pre-commit-config.yaml`:

```yaml
repos:
  - repo: https://github.com/ahans/bazel-build-file-sorter
    rev: v0.2.1
    hooks:
      - id: bazel-build-file-sorter
```

## Install

```sh
go install github.com/ahans/bazel-build-file-sorter@latest
```

## Usage

```sh
# Print sorted output to stdout
bazel-build-file-sorter path/to/BUILD

# Sort one or more files in place
bazel-build-file-sorter -i path/to/BUILD other/BUILD.bazel

# Print version
bazel-build-file-sorter -v
```
