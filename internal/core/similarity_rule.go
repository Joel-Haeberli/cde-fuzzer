package core

import (
	"strings"

	"github.com/Joel-Haeberli/cde-fuzzer/internal/result"
)

// SimilarityFunc computes a similarity score in [0, 1] between two strings.
type SimilarityFunc func(a, b string) float64

// SimilarityRule implements Rule using a string similarity algorithm.
// It slides a window over the text looking for spans similar to the target.
type SimilarityRule struct {
	name      string
	target    string
	threshold float64
	simFunc   SimilarityFunc
}

// NewSimilarityRule creates a rule that matches when the similarity between
// a text span and the target exceeds the given threshold.
func NewSimilarityRule(name, target string, threshold float64, simFunc SimilarityFunc) *SimilarityRule {
	return &SimilarityRule{
		name:      name,
		target:    target,
		threshold: threshold,
		simFunc:   simFunc,
	}
}

func (r *SimilarityRule) Name() string { return r.name }

func (r *SimilarityRule) Match(text string) bool {
	words := strings.Fields(text)
	targetWords := len(strings.Fields(r.target))
	for i := 0; i <= len(words)-targetWords; i++ {
		span := strings.Join(words[i:i+targetWords], " ")
		if r.simFunc(span, r.target) >= r.threshold {
			return true
		}
	}
	return false
}

func (r *SimilarityRule) Apply(text string) ([]result.Match, error) {
	var matches []result.Match
	words := strings.Fields(text)
	targetWords := len(strings.Fields(r.target))
	if targetWords == 0 || len(words) < targetWords {
		return matches, nil
	}
	for i := 0; i <= len(words)-targetWords; i++ {
		span := strings.Join(words[i:i+targetWords], " ")
		if r.simFunc(span, r.target) >= r.threshold {
			// Find the byte offsets of this span in the original text
			start := strings.Index(text, span)
			if start >= 0 {
				matches = append(matches, result.Match{
					Value: span,
					Start: start,
					End:   start + len(span),
				})
			}
		}
	}
	return matches, nil
}

func (r *SimilarityRule) Accuracy() float64 { return r.threshold }

// Levenshtein computes a normalized Levenshtein similarity in [0, 1].
func Levenshtein(a, b string) float64 {
	if a == b {
		return 1.0
	}
	la, lb := len(a), len(b)
	if la == 0 || lb == 0 {
		return 0.0
	}

	// Compute edit distance using two-row approach.
	prev := make([]int, lb+1)
	curr := make([]int, lb+1)
	for j := range prev {
		prev[j] = j
	}
	for i := 1; i <= la; i++ {
		curr[0] = i
		for j := 1; j <= lb; j++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}
			curr[j] = min(curr[j-1]+1, min(prev[j]+1, prev[j-1]+cost))
		}
		prev, curr = curr, prev
	}

	maxLen := max(la, lb)
	return 1.0 - float64(prev[lb])/float64(maxLen)
}
