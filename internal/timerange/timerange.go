package timerange

import (
	"fmt"
	"time"
)

// Range represents an inclusive time window [From, To].
type Range struct {
	From time.Time
	To   time.Time
}

// New creates a Range from two RFC3339 timestamp strings.
// Either value may be empty to indicate an open-ended range.
func New(from, to string) (Range, error) {
	var r Range
	var err error

	if from != "" {
		r.From, err = time.Parse(time.RFC3339, from)
		if err != nil {
			return Range{}, fmt.Errorf("invalid 'from' timestamp %q: %w", from, err)
		}
	}

	if to != "" {
		r.To, err = time.Parse(time.RFC3339, to)
		if err != nil {
			return Range{}, fmt.Errorf("invalid 'to' timestamp %q: %w", to, err)
		}
	}

	if !r.From.IsZero() && !r.To.IsZero() && r.To.Before(r.From) {
		return Range{}, fmt.Errorf("'to' (%s) must not be before 'from' (%s)", to, from)
	}

	return r, nil
}

// Contains reports whether t falls within the range.
// A zero From means no lower bound; a zero To means no upper bound.
func (r Range) Contains(t time.Time) bool {
	if !r.From.IsZero() && t.Before(r.From) {
		return false
	}
	if !r.To.IsZero() && t.After(r.To) {
		return false
	}
	return true
}

// IsZero reports whether the range has no bounds set.
func (r Range) IsZero() bool {
	return r.From.IsZero() && r.To.IsZero()
}

// ParseLineTime attempts to extract a timestamp from the beginning of a log
// line using common log timestamp layouts.
func ParseLineTime(line string) (time.Time, bool) {
	layouts := []struct {
		layout string
		width   int
	}{
		{time.RFC3339, 25},
		{"2006-01-02T15:04:05", 19},
		{"2006-01-02 15:04:05", 19},
		{"2006/01/02 15:04:05", 19},
	}

	for _, l := range layouts {
		if len(line) >= l.width {
			t, err := time.Parse(l.layout, line[:l.width])
			if err == nil {
				return t, true
			}
		}
	}
	return time.Time{}, false
}
