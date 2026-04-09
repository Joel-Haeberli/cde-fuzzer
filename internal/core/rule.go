package core

import "github.com/Joel-Haeberli/cde-fuzzer/internal/result"

// Rule is the fundamental extraction unit. It matches against raw text
// and produces extraction results with accuracy estimates.
type Rule interface {
	// Name returns a human-readable identifier for traceability.
	Name() string

	// Match reports whether this rule finds a relevant pattern in the text.
	Match(text string) bool

	// Apply runs the rule against the text and returns all matches.
	Apply(text string) ([]result.Match, error)

	// Accuracy returns a confidence score in [0, 1] for this rule's extractions.
	Accuracy() float64
}
