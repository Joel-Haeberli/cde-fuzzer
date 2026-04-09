package core

import "github.com/Joel-Haeberli/cde-fuzzer/internal/result"

// RuleChain composes multiple rules into a sequential extraction pipeline.
// Each rule in the chain is evaluated in order; all matching results are collected.
type RuleChain struct {
	name  string
	rules []Rule
}

// NewRuleChain creates a named chain from the given rules.
func NewRuleChain(name string, rules ...Rule) *RuleChain {
	return &RuleChain{name: name, rules: rules}
}

// Run evaluates all rules in the chain against the text.
// Returns traces for every rule that matched.
func (c *RuleChain) Run(text string) ([]result.RuleTrace, error) {
	var traces []result.RuleTrace
	for _, rule := range c.rules {
		if !rule.Match(text) {
			continue
		}
		matches, err := rule.Apply(text)
		if err != nil {
			return nil, err
		}
		for _, m := range matches {
			traces = append(traces, result.RuleTrace{
				RuleName: rule.Name(),
				Accuracy: rule.Accuracy(),
				Match:    m,
			})
		}
	}
	return traces, nil
}
