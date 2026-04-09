package core

// AccuracyEstimator computes an accuracy estimate for a rule applied to text.
type AccuracyEstimator interface {
	Estimate(rule Rule, text string) float64
}

// DefaultAccuracyEstimator returns the rule's own accuracy rating.
type DefaultAccuracyEstimator struct{}

func (e *DefaultAccuracyEstimator) Estimate(rule Rule, text string) float64 {
	return rule.Accuracy()
}
