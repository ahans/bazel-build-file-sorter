package main

import (
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
