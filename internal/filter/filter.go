package filter

import (
	"regexp"
	"strings"
	"time"
)

// Options holds filtering configuration.
type Options struct {
	// Pattern is a substring or regex to match against log lines.
	Pattern string
	// UseRegex treats Pattern as a regular expression.
	UseRegex bool
	// Since filters out lines with timestamps before this time (zero = no lower bound).
	Since time.Time
	// Until filters out lines with timestamps after this time (zero = no upper bound).
	Until time.Time
	// TimestampLayout is the Go time layout used to parse timestamps in log lines.
	TimestampLayout string
}

// Filter applies the configured options to a stream of log lines.
type Filter struct {
	opts    Options
	regex   *regexp.Regexp
	hasTime bool
}

// New creates a new Filter from the given Options.
// Returns an error if UseRegex is true and Pattern is not a valid regexp.
func New(opts Options) (*Filter, error) {
	f := &Filter{
		opts:    opts,
		hasTime: opts.TimestampLayout != "" && (!opts.Since.IsZero() || !opts.Until.IsZero()),
	}
	if opts.UseRegex && opts.Pattern != "" {
		re, err := regexp.Compile(opts.Pattern)
		if err != nil {
			return nil, err
		}
		f.regex = re
	}
	return f, nil
}

// Match returns true when the line passes all active filter conditions.
func (f *Filter) Match(line string) bool {
	if !f.matchPattern(line) {
		return false
	}
	if f.hasTime && !f.matchTime(line) {
		return false
	}
	return true
}

func (f *Filter) matchPattern(line string) bool {
	if f.opts.Pattern == "" {
		return true
	}
	if f.regex != nil {
		return f.regex.MatchString(line)
	}
	return strings.Contains(line, f.opts.Pattern)
}

func (f *Filter) matchTime(line string) bool {
	if len(line) < len(f.opts.TimestampLayout) {
		return true // cannot parse — let it through
	}
	t, err := time.Parse(f.opts.TimestampLayout, line[:len(f.opts.TimestampLayout)])
	if err != nil {
		return true // unparseable timestamp — let it through
	}
	if !f.opts.Since.IsZero() && t.Before(f.opts.Since) {
		return false
	}
	if !f.opts.Until.IsZero() && t.After(f.opts.Until) {
		return false
	}
	return true
}
