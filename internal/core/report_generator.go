package core

import (
	"bytes"
	"text/template"

	"github.com/Joel-Haeberli/cde-fuzzer/internal/result"
)

// ReportGenerator creates structured reports from extracted data
type ReportGenerator interface {
	// GenerateReport creates a report from extracted data and rule matches
	GenerateReport(data map[string]string, ruleMatches map[string][]result.RuleTrace) (string, error)
}

// TemplateReportGenerator uses Go templates to format reports
type TemplateReportGenerator struct {
	template *template.Template
}

// NewTemplateReportGenerator creates a new template-based report generator
func NewTemplateReportGenerator(templateContent string) (*TemplateReportGenerator, error) {
	tmpl, err := template.New("report").Parse(templateContent)
	if err != nil {
		return nil, err
	}
	return &TemplateReportGenerator{template: tmpl}, nil
}

// GenerateReport implements ReportGenerator interface
func (g *TemplateReportGenerator) GenerateReport(data map[string]string, ruleMatches map[string][]result.RuleTrace) (string, error) {
	// Prepare template data
	templateData := struct {
		Data        map[string]string
		RuleMatches map[string][]result.RuleTrace
	}{
		Data:        data,
		RuleMatches: ruleMatches,
	}

	var buf bytes.Buffer
	err := g.template.Execute(&buf, templateData)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

// RuleReverser maps extracted data back to report sections
type RuleReverser struct {
	rules []Rule
}

// NewRuleReverser creates a new rule reverser
func NewRuleReverser(rules []Rule) *RuleReverser {
	return &RuleReverser{rules: rules}
}

// ReverseMap creates a mapping from rule names to their extracted values
func (r *RuleReverser) ReverseMap(traces []result.RuleTrace) map[string]string {
	result := make(map[string]string)
	for _, trace := range traces {
		result[trace.RuleName] = trace.Match.Value
	}
	return result
}
