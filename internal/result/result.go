package result

// Match represents a single matched span in the source text.
type Match struct {
	Value string // The extracted value
	Start int    // Start offset in source text
	End   int    // End offset in source text
}

// RuleTrace records which rule produced a match and its confidence.
type RuleTrace struct {
	RuleName string
	Accuracy float64
	Match    Match
}

// ExtractionResult holds the output of an extraction process,
// including traceability metadata linking output back to the raw input.
type ExtractionResult struct {
	CDEID    string
	Answer   string      // The selected/extracted answer
	Traces   []RuleTrace // Ordered trace of rules that contributed
	Accuracy float64     // Overall accuracy estimate
}
