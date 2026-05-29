# build-file-sorter

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

## Install

```sh
go install github.com/ahans/build-file-sorter@latest
```

## Usage

```sh
# Print sorted output to stdout
build-file-sorter path/to/BUILD

# Sort in place
build-file-sorter -w path/to/BUILD

# Read from stdin
build-file-sorter - < path/to/BUILD
```

## Build from source

```sh
git clone https://github.com/ahans/build-file-sorter
cd build-file-sorter
go build .
```
