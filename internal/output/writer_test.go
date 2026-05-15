package output

import (
	"bufio"
	"compress/gzip"
	"os"
	"path/filepath"
	"testing"
)

func TestWritePlainFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.log")

	w, err := New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	lines := []string{"alpha", "beta", "gamma"}
	for _, l := range lines {
		if err := w.WriteLine(l); err != nil {
			t.Fatalf("WriteLine: %v", err)
		}
	}
	if err := w.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	f, _ := os.Open(path)
	defer f.Close()
	var got []string
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		got = append(got, sc.Text())
	}
	if len(got) != len(lines) {
		t.Fatalf("expected %d lines, got %d", len(lines), len(got))
	}
	for i, l := range lines {
		if got[i] != l {
			t.Errorf("line %d: want %q, got %q", i, l, got[i])
		}
	}
}

func TestWriteGzipFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.log.gz")

	w, err := New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if w.format != FormatGzip {
		t.Errorf("expected FormatGzip")
	}

	if err := w.WriteLine("compressed line"); err != nil {
		t.Fatalf("WriteLine: %v", err)
	}
	if err := w.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	f, _ := os.Open(path)
	defer f.Close()
	gr, err := gzip.NewReader(f)
	if err != nil {
		t.Fatalf("gzip.NewReader: %v", err)
	}
	sc := bufio.NewScanner(gr)
	if !sc.Scan() {
		t.Fatal("expected at least one line")
	}
	if sc.Text() != "compressed line" {
		t.Errorf("got %q", sc.Text())
	}
}

func TestWriteLinesChannel(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "chan.log")

	w, err := New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	ch := make(chan string, 3)
	ch <- "one"
	ch <- "two"
	ch <- "three"
	close(ch)

	if err := w.WriteLines(ch); err != nil {
		t.Fatalf("WriteLines: %v", err)
	}
	w.Close()

	f, _ := os.Open(path)
	defer f.Close()
	sc := bufio.NewScanner(f)
	count := 0
	for sc.Scan() {
		count++
	}
	if count != 3 {
		t.Errorf("expected 3 lines, got %d", count)
	}
}
