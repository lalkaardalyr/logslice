package filter

import (
	"bufio"
	"strings"
	"testing"
)

func TestPipelineSubstring(t *testing.T) {
	input := strings.Join([]string{
		"INFO server started",
		"ERROR disk full",
		"WARN low memory",
		"ERROR connection refused",
		"INFO shutdown",
	}, "\n")

	f := mustFilter(t, Options{Pattern: "ERROR"})
	scanner := bufio.NewScanner(strings.NewReader(input))

	var results []string
	n, err := Pipeline(scanner, f, func(line string) error {
		results = append(results, line)
		return nil
	})

	if err != nil {
		t.Fatalf("Pipeline error: %v", err)
	}
	if n != 2 {
		t.Errorf("matched %d lines, want 2", n)
	}
	if results[0] != "ERROR disk full" {
		t.Errorf("unexpected first match: %q", results[0])
	}
	if results[1] != "ERROR connection refused" {
		t.Errorf("unexpected second match: %q", results[1])
	}
}

func TestPipelineNoMatch(t *testing.T) {
	input := "INFO all good\nINFO still good\n"
	f := mustFilter(t, Options{Pattern: "ERROR"})
	scanner := bufio.NewScanner(strings.NewReader(input))

	n, err := Pipeline(scanner, f, func(string) error { return nil })
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 0 {
		t.Errorf("expected 0 matches, got %d", n)
	}
}

func TestPipelineEmptyInput(t *testing.T) {
	f := mustFilter(t, Options{})
	scanner := bufio.NewScanner(strings.NewReader(""))

	n, err := Pipeline(scanner, f, func(string) error { return nil })
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 0 {
		t.Errorf("expected 0 matches on empty input, got %d", n)
	}
}
