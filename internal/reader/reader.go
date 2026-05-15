package reader

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"strings"
)

// LogReader wraps a scanner over a potentially gzip-compressed log file.
type LogReader struct {
	file    *os.File
	gzReader *gzip.Reader
	Scanner *bufio.Scanner
}

// New opens a log file (plain text or .gz) and returns a LogReader.
func New(path string) (*LogReader, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening file %q: %w", path, err)
	}

	lr := &LogReader{file: f}

	var src io.Reader = f
	if strings.HasSuffix(path, ".gz") {
		gr, err := gzip.NewReader(f)
		if err != nil {
			f.Close()
			return nil, fmt.Errorf("creating gzip reader for %q: %w", path, err)
		}
		lr.gzReader = gr
		src = gr
	}

	lr.Scanner = bufio.NewScanner(src)
	return lr, nil
}

// Close releases all resources held by the LogReader.
func (lr *LogReader) Close() error {
	if lr.gzReader != nil {
		if err := lr.gzReader.Close(); err != nil {
			return err
		}
	}
	return lr.file.Close()
}

// Lines returns a channel that emits each line in the log file.
// The channel is closed when the file is exhausted or an error occurs.
func (lr *LogReader) Lines() <-chan string {
	ch := make(chan string, 64)
	go func() {
		defer close(ch)
		for lr.Scanner.Scan() {
			ch <- lr.Scanner.Text()
		}
	}()
	return ch
}
