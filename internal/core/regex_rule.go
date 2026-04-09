package core

import (
	"regexp"

	"github.com/Joel-Haeberli/cde-fuzzer/internal/result"
)

// RegexRule implements Rule using a compiled regular expression.
type RegexRule struct {
	name    string
	pattern *regexp.Regexp
	acc     float64
}

// NewRegexRule creates a rule that matches text against the given regex pattern.
// The accuracy value should be in [0, 1].
func NewRegexRule(name string, pattern *regexp.Regexp, accuracy float64) *RegexRule {
	return &RegexRule{
		name:    name,
		pattern: pattern,
		acc:     accuracy,
	}
}

func (r *RegexRule) Name() string { return r.name }

func (r *RegexRule) Match(text string) bool {
	return r.pattern.MatchString(text)
}

func (r *RegexRule) Apply(text string) ([]result.Match, error) {
	indices := r.pattern.FindAllStringIndex(text, -1)
	matches := make([]result.Match, 0, len(indices))
	for _, idx := range indices {
		matches = append(matches, result.Match{
			Value: text[idx[0]:idx[1]],
			Start: idx[0],
			End:   idx[1],
		})
	}
	return matches, nil
}

func (r *RegexRule) Accuracy() float64 { return r.acc }
