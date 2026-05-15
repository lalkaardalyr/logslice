package output

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func newTempWriter(t *testing.T, name string) *Writer {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	w, err := New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return w
}

func TestCounterZero(t *testing.T) {
	w := newTempWriter(t, "zero.log")
	c := NewCounter(w)
	defer c.Close()
	if c.Count() != 0 {
		t.Errorf("expected 0, got %d", c.Count())
	}
}

func TestCounterWriteLine(t *testing.T) {
	w := newTempWriter(t, "wl.log")
	c := NewCounter(w)
	defer c.Close()

	for i := 0; i < 5; i++ {
		if err := c.WriteLine("line"); err != nil {
			t.Fatalf("WriteLine: %v", err)
		}
	}
	if c.Count() != 5 {
		t.Errorf("expected 5, got %d", c.Count())
	}
}

func TestCounterWriteLines(t *testing.T) {
	w := newTempWriter(t, "wls.log")
	c := NewCounter(w)
	defer c.Close()

	ch := make(chan string, 4)
	ch <- "a"
	ch <- "b"
	ch <- "c"
	ch <- "d"
	close(ch)

	if err := c.WriteLines(ch); err != nil {
		t.Fatalf("WriteLines: %v", err)
	}
	if c.Count() != 4 {
		t.Errorf("expected 4, got %d", c.Count())
	}
}

func TestCounterPrintSummary(t *testing.T) {
	w := newTempWriter(t, "sum.log")
	c := NewCounter(w)
	defer c.Close()

	_ = c.WriteLine("hello")
	_ = c.WriteLine("world")

	var buf bytes.Buffer
	c.PrintSummary(&buf)

	expected := "lines written: 2\n"
	if buf.String() != expected {
		t.Errorf("summary: want %q, got %q", expected, buf.String())
	}
}

func TestCounterStdout(t *testing.T) {
	w, err := New("-")
	if err != nil {
		t.Fatalf("New stdout: %v", err)
	}
	// Redirect stdout temporarily
	old := os.Stdout
	r, wr, _ := os.Pipe()
	os.Stdout = wr

	c := NewCounter(w)
	_ = c.WriteLine("test")

	wr.Close()
	os.Stdout = old
	r.Close()

	if c.Count() != 1 {
		t.Errorf("expected 1, got %d", c.Count())
	}
}
