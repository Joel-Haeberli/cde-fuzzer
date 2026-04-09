package core

import (
	"bytes"
	"math/rand"
	"text/template"
	"time"

	"github.com/Joel-Haeberli/cde-fuzzer/internal/result"
)

// DiverseReportGenerator creates more human-like, varied reports
type DiverseReportGenerator struct {
	baseTemplate *template.Template
	rand         *rand.Rand
	// Template variations for different report sections
	scenarioTemplates   []string
	procedureTemplates  []string
	findingsTemplates   []string
	impressionTemplates []string
}

// NewDiverseReportGenerator creates a generator with built-in variability
func NewDiverseReportGenerator() (*DiverseReportGenerator, error) {
	// Seed random number generator
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	// Base template with placeholders
	baseTemplate := template.Must(template.New("diverse_report").Parse(`{{.Header}}

CLINICAL INFORMATION:
{{.ClinicalInfo}}

PROCEDURES PERFORMED:
{{.Procedures}}

FINDINGS:
{{.Findings}}

IMPRESSION:
{{.Impression}}

{{if .Recommendations}}RECOMMENDATIONS:
{{.Recommendations}}
{{end}}{{.AdditionalInfo}}
`))

	generator := &DiverseReportGenerator{
		baseTemplate: baseTemplate,
		rand:         rng,
	}

	// Initialize template variations
	generator.initTemplateVariations()

	return generator, nil
}

func (g *DiverseReportGenerator) initTemplateVariations() {
	// Different header formats
	g.scenarioTemplates = []string{
		"Patient presents with {{.Scenario}}.",
		"Clinical indication: {{.Scenario}}.",
		"Reason for exam: {{.Scenario}}.",
		"History: {{.Scenario}}.",
	}

	// Different procedure description formats
	g.procedureTemplates = []string{
		"- {{.Procedure}}",
		"• {{.Procedure}} was performed.",
		"Exam: {{.Procedure}}",
		"Technique: {{.Procedure}}",
	}

	// Different findings formats
	g.findingsTemplates = []string{
		"The examination demonstrates {{.Findings}}.",
		"Findings include: {{.Findings}}",
		"Imaging reveals {{.Findings}}.",
		"Notable findings: {{.Findings}}",
	}

	// Different impression formats
	g.impressionTemplates = []string{
		"1. {{.Impression}}",
		"Impression: {{.Impression}}",
		"Conclusion: {{.Impression}}",
		"Summary: {{.Impression}}",
	}
}

// GenerateReport creates a varied, human-like report
func (g *DiverseReportGenerator) GenerateReport(data map[string]string, ruleMatches map[string][]result.RuleTrace) (string, error) {
	// Extract data with fallbacks
	scenario := g.getValueOrFallback(data, "scenario_description_extractor", "clinical evaluation")
	procedure := g.getValueOrFallback(data, "procedure_extractor", "imaging examination")
	bodyPart := g.getValueOrFallback(data, "body_part_extractor", "the area of interest")
	appropriateness := g.getValueOrFallback(data, "appropriateness_extractor", "appropriate for the clinical scenario")
	radiationDose := g.getValueOrFallback(data, "radiation_dose_extractor", "within standard limits")

	// Extract additional custom fields
	patientAge := g.getValueOrFallback(data, "patient_age_extractor", "")
	contrastUsed := g.getValueOrFallback(data, "contrast_used_extractor", "")
	examIndication := g.getValueOrFallback(data, "exam_indication_extractor", "")

	// Generate varied content for each section
	header := g.generateHeader()
	clinicalInfo := g.generateClinicalInfo(scenario, patientAge, examIndication)
	procedures := g.generateProcedures(procedure, contrastUsed)
	findings := g.generateFindings(bodyPart)
	impression := g.generateImpression(bodyPart, appropriateness)
	recommendations := g.generateRecommendations(appropriateness, radiationDose)
	additionalInfo := g.generateAdditionalInfo(data)

	templateData := struct {
		Header          string
		ClinicalInfo    string
		Procedures      string
		Findings        string
		Impression      string
		Recommendations string
		AdditionalInfo  string
	}{
		Header:          header,
		ClinicalInfo:    clinicalInfo,
		Procedures:      procedures,
		Findings:        findings,
		Impression:      impression,
		Recommendations: recommendations,
		AdditionalInfo:  additionalInfo,
	}

	var buf bytes.Buffer
	err := g.baseTemplate.Execute(&buf, templateData)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (g *DiverseReportGenerator) getValueOrFallback(data map[string]string, key, fallback string) string {
	if val, ok := data[key]; ok && val != "" {
		return val
	}
	return fallback
}

func (g *DiverseReportGenerator) generateHeader() string {
	headers := []string{
		"RADIOLOGY REPORT",
		"IMAGING REPORT",
		"RADIOLOGY CONSULTATION",
		"DIAGNOSTIC IMAGING REPORT",
	}
	return g.randomChoice(headers)
}

func (g *DiverseReportGenerator) generateClinicalInfo(scenario, patientAge, examIndication string) string {
	// Build clinical info with available data
	clinicalParts := []string{}

	if patientAge != "" {
		clinicalParts = append(clinicalParts, patientAge+" year old patient")
	}

	if examIndication != "" {
		clinicalParts = append(clinicalParts, "presenting for "+examIndication+" examination")
	}

	if scenario != "" {
		clinicalParts = append(clinicalParts, "with "+scenario)
	}

	clinicalInfo := ""
	if len(clinicalParts) > 0 {
		clinicalInfo = g.randomChoice([]string{
			"Patient is a " + g.joinWithCommas(clinicalParts) + ".",
			"Clinical history: " + g.joinWithCommas(clinicalParts) + ".",
			"Presentation: " + g.joinWithCommas(clinicalParts) + ".",
		})
	} else {
		clinicalInfo = g.randomChoice(g.scenarioTemplates)
		data := struct{ Scenario string }{Scenario: scenario}
		var result bytes.Buffer
		tmpl := template.Must(template.New("").Parse(clinicalInfo))
		_ = tmpl.Execute(&result, data)
		clinicalInfo = result.String()
	}

	return clinicalInfo
}

func (g *DiverseReportGenerator) generateProcedures(procedure, contrastUsed string) string {
	procedureText := procedure
	if contrastUsed != "" {
		// Remove redundant contrast info if already in procedure
		if !g.stringContains(procedure, "contrast") && !g.stringContains(procedure, "Contrast") {
			procedureText = procedure + " " + contrastUsed
		}
	}

	tmplText := g.randomChoice(g.procedureTemplates)
	data := struct{ Procedure string }{Procedure: procedureText}
	var result bytes.Buffer
	tmpl := template.Must(template.New("").Parse(tmplText))
	_ = tmpl.Execute(&result, data)
	return result.String()
}

func (g *DiverseReportGenerator) generateFindings(bodyPart string) string {
	// Add some random variability to findings
	findings := []string{
		"normal appearance of " + bodyPart,
		"unremarkable " + bodyPart,
		"no acute abnormalities in " + bodyPart,
		"expected postoperative changes in " + bodyPart,
		"age-appropriate findings in " + bodyPart,
	}

	baseFinding := g.randomChoice(findings)

	// Sometimes add additional details
	if g.rand.Float32() < 0.5 {
		baseFinding += ". " + g.randomChoice([]string{
			"No suspicious lesions identified.",
			"No evidence of acute pathology.",
			"Findings are stable compared to prior examination.",
			"No significant interval change.",
		})
	}

	tmplText := g.randomChoice(g.findingsTemplates)
	data := struct{ Findings string }{Findings: baseFinding}
	var result bytes.Buffer
	tmpl := template.Must(template.New("").Parse(tmplText))
	_ = tmpl.Execute(&result, data)
	return result.String()
}

func (g *DiverseReportGenerator) generateImpression(bodyPart, appropriateness string) string {
	impressions := []string{
		"The " + bodyPart + " appears within normal limits.",
		"No significant abnormality detected in " + bodyPart + ".",
		"Findings are consistent with " + appropriateness + " clinical scenario.",
		"Imaging correlates with the clinical presentation.",
	}

	baseImpression := g.randomChoice(impressions)

	tmplText := g.randomChoice(g.impressionTemplates)
	data := struct{ Impression string }{Impression: baseImpression}
	var result bytes.Buffer
	tmpl := template.Must(template.New("").Parse(tmplText))
	_ = tmpl.Execute(&result, data)
	return result.String()
}

func (g *DiverseReportGenerator) generateRecommendations(appropriateness, radiationDose string) string {
	if g.rand.Float32() < 0.3 {
		// Sometimes omit recommendations
		return ""
	}

	recommendations := []string{
		"Correlation with clinical findings is recommended.",
		"Follow-up imaging may be considered based on clinical evolution.",
		"This examination is " + appropriateness + " for the given clinical scenario.",
		"Radiation dose was " + radiationDose + " and is considered appropriate.",
	}

	return g.randomChoice(recommendations)
}

func (g *DiverseReportGenerator) randomChoice(options []string) string {
	if len(options) == 0 {
		return ""
	}
	return options[g.rand.Intn(len(options))]
}

func (g *DiverseReportGenerator) joinWithCommas(parts []string) string {
	if len(parts) == 0 {
		return ""
	}
	if len(parts) == 1 {
		return parts[0]
	}
	if len(parts) == 2 {
		return parts[0] + " and " + parts[1]
	}
	// For 3+ items
	result := ""
	for i, part := range parts {
		if i > 0 {
			if i == len(parts)-1 {
				result += ", and "
			} else {
				result += ", "
			}
		}
		result += part
	}
	return result
}

func (g *DiverseReportGenerator) stringContains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || g.containsSubstring(s, substr)))
}

func (g *DiverseReportGenerator) containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func (g *DiverseReportGenerator) generateAdditionalInfo(data map[string]string) string {
	// Collect all fields that aren't part of our standard sections
	additionalFields := []string{}

	standardFields := map[string]bool{
		"scenario_description_extractor": true,
		"procedure_extractor":            true,
		"radiation_dose_extractor":       true,
		"body_part_extractor":            true,
		"appropriateness_extractor":      true,
		"patient_age_extractor":          true,
		"contrast_used_extractor":        true,
		"exam_indication_extractor":      true,
	}

	for fieldName, value := range data {
		if !standardFields[fieldName] && value != "" {
			// Make field name more readable
			readableName := g.makeReadableFieldName(fieldName)
			additionalFields = append(additionalFields, readableName+": "+value)
		}
	}

	if len(additionalFields) == 0 {
		return ""
	}

	return "\nADDITIONAL INFORMATION:\n" + g.joinAdditionalFields(additionalFields)
}

func (g *DiverseReportGenerator) makeReadableFieldName(fieldName string) string {
	// Remove _extractor suffix and replace underscores with spaces
	name := fieldName
	if len(name) > 10 && name[len(name)-10:] == "_extractor" {
		name = name[:len(name)-10]
	}
	name = g.replaceUnderscores(name)
	return g.titleCase(name)
}

func (g *DiverseReportGenerator) replaceUnderscores(s string) string {
	result := ""
	for i, char := range s {
		if char == '_' {
			result += " "
		} else if i == 0 {
			result += string(char)
		} else {
			result += string(char)
		}
	}
	return result
}

func (g *DiverseReportGenerator) titleCase(s string) string {
	if s == "" {
		return ""
	}
	words := g.splitWords(s)
	for i, word := range words {
		if len(word) > 0 {
			words[i] = string(word[0]-32) + word[1:] // Simple uppercase first letter
		}
	}
	return g.joinWords(words)
}

func (g *DiverseReportGenerator) splitWords(s string) []string {
	var words []string
	currentWord := ""
	for _, char := range s {
		if char == ' ' {
			if currentWord != "" {
				words = append(words, currentWord)
				currentWord = ""
			}
		} else {
			currentWord += string(char)
		}
	}
	if currentWord != "" {
		words = append(words, currentWord)
	}
	return words
}

func (g *DiverseReportGenerator) joinWords(words []string) string {
	result := ""
	for i, word := range words {
		if i > 0 {
			result += " "
		}
		result += word
	}
	return result
}

func (g *DiverseReportGenerator) joinAdditionalFields(fields []string) string {
	result := ""
	for _, field := range fields {
		result += "- " + field + "\n"
	}
	return result
}
