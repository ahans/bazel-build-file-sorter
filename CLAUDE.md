# build-file-sorter

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

## Test structure

Tests live in `testdata/<case-name>/input` and `testdata/<case-name>/expected`. The test runner in `main_test.go` discovers cases by reading the directory. Add a new case by creating a new subdirectory with both files.
