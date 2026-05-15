package reader_test

import (
	"compress/gzip"
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/logslice/internal/reader"
)

func writePlain(t *testing.T, lines []string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "plain-*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	for _, l := range lines {
		f.WriteString(l + "\n")
	}
	return f.Name()
}

func writeGzip(t *testing.T, lines []string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "compressed.log.gz")
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	gw := gzip.NewWriter(f)
	for _, l := range lines {
		gw.Write([]byte(l + "\n"))
	}
	gw.Close()
	return path
}

func collectLines(t *testing.T, path string) []string {
	t.Helper()
	lr, err := reader.New(path)
	if err != nil {
		t.Fatalf("reader.New: %v", err)
	}
	defer lr.Close()
	var got []string
	for line := range lr.Lines() {
		got = append(got, line)
	}
	return got
}

func TestPlainFile(t *testing.T) {
	want := []string{"hello world", "second line", "third line"}
	path := writePlain(t, want)
	got := collectLines(t, path)
	if len(got) != len(want) {
		t.Fatalf("expected %d lines, got %d", len(want), len(got))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("line %d: want %q, got %q", i, want[i], got[i])
		}
	}
}

func TestGzipFile(t *testing.T) {
	want := []string{"compressed line 1", "compressed line 2"}
	path := writeGzip(t, want)
	got := collectLines(t, path)
	if len(got) != len(want) {
		t.Fatalf("expected %d lines, got %d", len(want), len(got))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("line %d: want %q, got %q", i, want[i], got[i])
		}
	}
}

func TestMissingFile(t *testing.T) {
	_, err := reader.New("/nonexistent/path/to/file.log")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
