package timerange

import (
	"testing"
	"time"
)

func mustNew(t *testing.T, from, to string) Range {
	t.Helper()
	r, err := New(from, to)
	if err != nil {
		t.Fatalf("New(%q, %q) unexpected error: %v", from, to, err)
	}
	return r
}

func TestNewValidRange(t *testing.T) {
	r := mustNew(t, "2024-01-01T00:00:00Z", "2024-01-02T00:00:00Z")
	if r.From.IsZero() || r.To.IsZero() {
		t.Fatal("expected non-zero bounds")
	}
}

func TestNewOpenEnded(t *testing.T) {
	r := mustNew(t, "2024-01-01T00:00:00Z", "")
	if r.From.IsZero() {
		t.Fatal("expected non-zero From")
	}
	if !r.To.IsZero() {
		t.Fatal("expected zero To")
	}
}

func TestNewInvalidFrom(t *testing.T) {
	_, err := New("not-a-date", "")
	if err == nil {
		t.Fatal("expected error for invalid from timestamp")
	}
}

func TestNewToBeforeFrom(t *testing.T) {
	_, err := New("2024-06-01T00:00:00Z", "2024-01-01T00:00:00Z")
	if err == nil {
		t.Fatal("expected error when to is before from")
	}
}

func TestContains(t *testing.T) {
	r := mustNew(t, "2024-03-01T00:00:00Z", "2024-03-31T23:59:59Z")

	cases := []struct {
		input string
		want  bool
	}{
		{"2024-03-15T12:00:00Z", true},
		{"2024-02-28T23:59:59Z", false},
		{"2024-04-01T00:00:00Z", false},
		{"2024-03-01T00:00:00Z", true},
		{"2024-03-31T23:59:59Z", true},
	}

	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			ts, err := time.Parse(time.RFC3339, c.input)
			if err != nil {
				t.Fatalf("parse: %v", err)
			}
			if got := r.Contains(ts); got != c.want {
				t.Errorf("Contains(%s) = %v, want %v", c.input, got, c.want)
			}
		})
	}
}

func TestIsZero(t *testing.T) {
	r := mustNew(t, "", "")
	if !r.IsZero() {
		t.Fatal("expected IsZero for empty range")
	}
}

func TestParseLineTime(t *testing.T) {
	cases := []struct {
		line string
		want bool
	}{
		{"2024-03-15T08:30:00Z INFO server started", true},
		{"2024-03-15T08:30:00+02:00 WARN disk low", true},
		{"2024-03-15 08:30:00 ERROR crash", true},
		{"2024/03/15 08:30:00 DEBUG ok", true},
		{"no timestamp here", false},
		{"", false},
	}

	for _, c := range cases {
		t.Run(c.line, func(t *testing.T) {
			_, ok := ParseLineTime(c.line)
			if ok != c.want {
				t.Errorf("ParseLineTime(%q) ok=%v, want %v", c.line, ok, c.want)
			}
		})
	}
}
