package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestSortBuildFile(t *testing.T) {
	entries, err := os.ReadDir("testdata")
	if err != nil {
		t.Fatalf("read testdata: %v", err)
	}

	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		name := e.Name()
		t.Run(name, func(t *testing.T) {
			inputPath := filepath.Join("testdata", name, "input")
			expectedPath := filepath.Join("testdata", name, "expected")

			input, err := os.ReadFile(inputPath)
			if err != nil {
				t.Fatalf("read input: %v", err)
			}
			expected, err := os.ReadFile(expectedPath)
			if err != nil {
				t.Fatalf("read expected: %v", err)
			}

			got, err := sortBuildFile(inputPath, input)
			if err != nil {
				t.Fatalf("sortBuildFile: %v", err)
			}

			if got != string(expected) {
				t.Errorf("output mismatch\ngot:\n%s\nwant:\n%s", got, expected)
			}
		})
	}
}

func TestRun(t *testing.T) {
	inputPath := filepath.Join("testdata", "basic", "input")
	expected, err := os.ReadFile(filepath.Join("testdata", "basic", "expected"))
	if err != nil {
		t.Fatalf("read expected: %v", err)
	}

	for _, inPlace := range []bool{true, false} {
		t.Run(fmt.Sprintf("inPlace=%v", inPlace), func(t *testing.T) {
			buildFilePath := copyToTemp(t, inputPath)
			var buf bytes.Buffer
			run(&buf, buildFilePath, inPlace)
			got, err := os.ReadFile(buildFilePath)
			if err != nil {
				t.Fatalf("read build file: %v", err)
			}
			if inPlace {
				if string(got) != string(expected) {
					t.Errorf("output mismatch\ngot:\n%s\nwant:\n%s", got, expected)
				}

				if buf.Len() != 0 {
					t.Errorf("expected no output, got: %q", buf.String())
				}
			} else {
				originalInput, err := os.ReadFile(inputPath)
				if err != nil {
					t.Fatalf("read input file: %v", err)
				}
				if string(got) != string(originalInput) {
					t.Errorf("input build file was modified")
				}
				if buf.String() != string(expected) {
					t.Errorf("output mismatch\ngot:\n%s\nwant:\n%s", buf.String(), expected)
				}
			}
		})
	}
}

func copyToTemp(t *testing.T, src string) string {
	t.Helper()
	data, err := os.ReadFile(src)
	if err != nil {
		t.Fatalf("read src: %v", err)
	}
	dst := filepath.Join(t.TempDir(), filepath.Base(src))
	err = os.WriteFile(dst, data, 0644)
	if err != nil {
		t.Fatalf("write dst: %v", err)
	}
	return dst
}
