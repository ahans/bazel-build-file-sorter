package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/bazelbuild/buildtools/build"
)

func main() {
	updateInPlace := flag.Bool("i", false, "update file in place")
	flag.Parse()

	if flag.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "usage: build-file-sorter [-i] <file>\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	filename := flag.Arg(0)

	run(os.Stdout, filename, *updateInPlace)
}

func run(w io.Writer, filename string, updateInPlace bool) {
	var data []byte
	var err error
	data, err = os.ReadFile(filename)
	if err != nil {
		log.Fatalf("reading input file failed: %v", err)
	}

	out, err := sortBuildFile(filename, data)
	if err != nil {
		log.Fatalf("sortBuildFile failed: %v", err)
	}
	if updateInPlace {
		if err := os.WriteFile(filename, []byte(out), 0644); err != nil {
			log.Fatalf("write failed: %v", err)
		}
	} else {
		fmt.Fprint(w, out)
	}
}

func sortBuildFile(filename string, data []byte) (string, error) {
	f, err := build.ParseBuild(filename, data)
	if err != nil {
		return "", err
	}
	sortRules(f)
	return build.FormatString(f), nil
}

// isPinned reports whether expr has a "# nosort" comment before it or on the same line.
func isPinned(expr build.Expr) bool {
	c := expr.Comment()
	return hasNosortComment(c.Before) || hasNosortComment(c.Suffix)
}

func hasNosortComment(comments []build.Comment) bool {
	for _, c := range comments {
		text := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(c.Token), "#"))
		if text == "nosort" {
			return true
		}
	}
	return false
}

// ruleName returns the value of the "name" attribute of a CallExpr, or "" if not found.
func ruleName(call *build.CallExpr) string {
	for _, arg := range call.List {
		assign, ok := arg.(*build.AssignExpr)
		if !ok {
			continue
		}
		lhs, ok := assign.LHS.(*build.Ident)
		if !ok || lhs.Name != "name" {
			continue
		}
		rhs, ok := assign.RHS.(*build.StringExpr)
		if !ok {
			continue
		}
		return rhs.Value
	}
	return ""
}

// sortRules sorts unpinned rules alphabetically by name. Pinned statements
// (marked with # nosort) stay at their original indices; unpinned non-rules
// move to the top; unpinned rules fill the remaining slots sorted.
func sortRules(f *build.File) {
	pinned := map[int]build.Expr{}
	var unpinnedNonRules []build.Expr
	var unpinnedRules []*build.CallExpr

	for i, stmt := range f.Stmt {
		if isPinned(stmt) {
			pinned[i] = stmt
			continue
		}
		call, ok := stmt.(*build.CallExpr)
		if ok && ruleName(call) != "" {
			unpinnedRules = append(unpinnedRules, call)
		} else {
			unpinnedNonRules = append(unpinnedNonRules, stmt)
		}
	}

	sort.Slice(unpinnedRules, func(i, j int) bool {
		return strings.ToLower(ruleName(unpinnedRules[i])) < strings.ToLower(ruleName(unpinnedRules[j]))
	})

	fill := make([]build.Expr, 0, len(unpinnedNonRules)+len(unpinnedRules))
	fill = append(fill, unpinnedNonRules...)
	for _, r := range unpinnedRules {
		fill = append(fill, r)
	}

	result := make([]build.Expr, len(f.Stmt))
	fillIdx := 0
	for i := range f.Stmt {
		if expr, ok := pinned[i]; ok {
			result[i] = expr
		} else {
			result[i] = fill[fillIdx]
			fillIdx++
		}
	}
	f.Stmt = result
}
