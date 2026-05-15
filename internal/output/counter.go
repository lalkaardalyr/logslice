package output

import (
	"fmt"
	"io"
	"sync/atomic"
)

// Counter wraps a Writer and tracks the number of lines written.
type Counter struct {
	w     *Writer
	count int64
}

// NewCounter wraps an existing Writer with line counting.
func NewCounter(w *Writer) *Counter {
	return &Counter{w: w}
}

// WriteLine writes a line and increments the counter.
func (c *Counter) WriteLine(line string) error {
	if err := c.w.WriteLine(line); err != nil {
		return err
	}
	atomic.AddInt64(&c.count, 1)
	return nil
}

// WriteLines drains the channel, writing each line and counting it.
func (c *Counter) WriteLines(lines <-chan string) error {
	for line := range lines {
		if err := c.WriteLine(line); err != nil {
			return err
		}
	}
	return nil
}

// Count returns the total number of lines written so far.
func (c *Counter) Count() int64 {
	return atomic.LoadInt64(&c.count)
}

// PrintSummary writes a human-readable summary to w.
func (c *Counter) PrintSummary(w io.Writer) {
	fmt.Fprintf(w, "lines written: %d\n", c.Count())
}

// Close closes the underlying Writer.
func (c *Counter) Close() error {
	return c.w.Close()
}
