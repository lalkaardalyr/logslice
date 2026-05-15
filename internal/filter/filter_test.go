package filter

import (
	"testing"
	"time"
)

const tsLayout = "2006-01-02 15:04:05"

func mustFilter(t *testing.T, opts Options) *Filter {
	t.Helper()
	f, err := New(opts)
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	return f
}

func TestSubstringMatch(t *testing.T) {
	f := mustFilter(t, Options{Pattern: "ERROR"})
	if !f.Match("2024-01-01 00:00:00 ERROR something broke") {
		t.Error("expected match for ERROR line")
	}
	if f.Match("2024-01-01 00:00:00 INFO all good") {
		t.Error("expected no match for INFO line")
	}
}

func TestRegexMatch(t *testing.T) {
	f := mustFilter(t, Options{Pattern: `ERROR|WARN`, UseRegex: true})
	if !f.Match("WARN low disk") {
		t.Error("expected match for WARN")
	}
	if f.Match("INFO startup") {
		t.Error("expected no match for INFO")
	}
}

func TestInvalidRegex(t *testing.T) {
	_, err := New(Options{Pattern: `[invalid`, UseRegex: true})
	if err == nil {
		t.Error("expected error for invalid regex")
	}
}

func TestEmptyPatternMatchesAll(t *testing.T) {
	f := mustFilter(t, Options{})
	if !f.Match("anything at all") {
		t.Error("empty pattern should match every line")
	}
}

func TestTimeRangeFilter(t *testing.T) {
	since := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	until := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	f := mustFilter(t, Options{
		TimestampLayout: tsLayout,
		Since:           since,
		Until:           until,
	})

	cases := []struct {
		line  string
		want  bool
	}{
		{"2024-01-01 09:59:59 before range", false},
		{"2024-01-01 10:00:00 at since boundary", true},
		{"2024-01-01 11:00:00 inside range", true},
		{"2024-01-01 12:00:00 at until boundary", true},
		{"2024-01-01 12:00:01 after range", false},
		{"no timestamp here", true}, // unparseable → let through
	}

	for _, tc := range cases {
		got := f.Match(tc.line)
		if got != tc.want {
			t.Errorf("Match(%q) = %v, want %v", tc.line, got, tc.want)
		}
	}
}

func TestCombinedPatternAndTime(t *testing.T) {
	since := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	f := mustFilter(t, Options{
		Pattern:         "ERROR",
		TimestampLayout: tsLayout,
		Since:           since,
	})

	if f.Match("2024-01-01 09:00:00 ERROR too early") {
		t.Error("should not match: before Since")
	}
	if f.Match("2024-01-01 11:00:00 INFO not an error") {
		t.Error("should not match: wrong pattern")
	}
	if !f.Match("2024-01-01 11:00:00 ERROR in range") {
		t.Error("should match: pattern + time both pass")
	}
}
