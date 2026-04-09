package core

import (
	"github.com/Joel-Haeberli/cde-fuzzer/internal/cde"
	"github.com/Joel-Haeberli/cde-fuzzer/internal/result"
)

// ExtractionProcess binds a CDE definition to a rule chain and accuracy estimator,
// forming a complete extraction pipeline.
type ExtractionProcess struct {
	CDE       cde.CDE
	Chain     *RuleChain
	Estimator AccuracyEstimator
}

// NewExtractionProcess creates a new extraction process.
func NewExtractionProcess(c cde.CDE, chain *RuleChain, estimator AccuracyEstimator) *ExtractionProcess {
	return &ExtractionProcess{
		CDE:       c,
		Chain:     chain,
		Estimator: estimator,
	}
}

// Run executes the extraction against the given text and returns the result.
func (ep *ExtractionProcess) Run(text string) (*result.ExtractionResult, error) {
	traces, err := ep.Chain.Run(text)
	if err != nil {
		return nil, err
	}
	if len(traces) == 0 {
		return &result.ExtractionResult{
			CDEID: ep.CDE.ID,
		}, nil
	}

	// Pick the best match by accuracy.
	best := traces[0]
	for _, t := range traces[1:] {
		if t.Accuracy > best.Accuracy {
			best = t
		}
	}

	// Compute overall accuracy as the average across all traces.
	var totalAcc float64
	for _, t := range traces {
		totalAcc += t.Accuracy
	}

	return &result.ExtractionResult{
		CDEID:    ep.CDE.ID,
		Answer:   best.Match.Value,
		Traces:   traces,
		Accuracy: totalAcc / float64(len(traces)),
	}, nil
}
