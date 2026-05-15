package output

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"strings"
)

// Format represents the output format for log lines.
type Format int

const (
	FormatPlain Format = iota
	FormatGzip
)

// Writer wraps an io.WriteCloser and writes log lines to it.
type Writer struct {
	w      io.WriteCloser
	inner  io.WriteCloser // non-nil when gzip wraps a file
	format Format
}

// New creates a new Writer. If path is "-", output goes to stdout.
// Format is inferred from the file extension (.gz => gzip).
func New(path string) (*Writer, error) {
	if path == "-" {
		return &Writer{w: os.Stdout, format: FormatPlain}, nil
	}

	f, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("output: create %q: %w", path, err)
	}

	if strings.HasSuffix(path, ".gz") {
		gw := gzip.NewWriter(f)
		return &Writer{w: gw, inner: f, format: FormatGzip}, nil
	}

	return &Writer{w: f, format: FormatPlain}, nil
}

// WriteLine writes a single log line followed by a newline.
func (w *Writer) WriteLine(line string) error {
	_, err := fmt.Fprintln(w.w, line)
	return err
}

// WriteLines writes all lines from the provided channel.
func (w *Writer) WriteLines(lines <-chan string) error {
	for line := range lines {
		if err := w.WriteLine(line); err != nil {
			return err
		}
	}
	return nil
}

// Close flushes and closes the underlying writer(s).
func (w *Writer) Close() error {
	if err := w.w.Close(); err != nil {
		return err
	}
	if w.inner != nil {
		return w.inner.Close()
	}
	return nil
}
