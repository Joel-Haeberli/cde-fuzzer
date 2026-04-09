package core

import (
	"bytes"
	"fmt"
	"math/rand"
	"strings"
	"text/template"
	"time"

	"github.com/Joel-Haeberli/cde-fuzzer/internal/result"
)

// EnhancedReportGenerator creates synthetic reports with realistic variability
type EnhancedReportGenerator struct {
	baseTemplate *template.Template
	rand         *rand.Rand
	// Data sources for realistic content
	dataSources []map[string]string
}

// NewEnhancedReportGenerator creates a generator with synthetic report capabilities
func NewEnhancedReportGenerator() (*EnhancedReportGenerator, error) {
	// Seed random number generator
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	// Base template with comprehensive structure
	baseTemplate := template.Must(template.New("enhanced_report").Parse(`{{.Header}}

{{if .PatientInfo}}PATIENT INFORMATION:
{{.PatientInfo}}
{{end}}
{{if .ClinicalInfo}}CLINICAL INFORMATION:
{{.ClinicalInfo}}
{{end}}

PROCEDURES PERFORMED:
{{.Procedures}}

FINDINGS:
{{.Findings}}

IMPRESSION:
{{.Impression}}

{{if .Recommendations}}RECOMMENDATIONS:
{{.Recommendations}}
{{end}}
{{if .AdditionalInfo}}
ADDITIONAL INFORMATION:
{{.AdditionalInfo}}
{{end}}
{{.Footer}}
`))

	generator := &EnhancedReportGenerator{
		baseTemplate: baseTemplate,
		rand:         rng,
		dataSources:  make([]map[string]string, 0),
	}

	return generator, nil
}

// AddDataSource adds real data to mix with synthetic data
func (g *EnhancedReportGenerator) AddDataSource(data map[string]string) {
	g.dataSources = append(g.dataSources, data)
}

// GenerateSyntheticReport creates a synthetic report with realistic medical content
func (g *EnhancedReportGenerator) GenerateSyntheticReport(variabilityFactor float64) (string, error) {
	// Generate synthetic data with variability
	syntheticData := g.generateSyntheticData(variabilityFactor)

	// Create empty rule traces (since this is synthetic)
	emptyTraces := make(map[string][]result.RuleTrace)

	return g.ExecuteString(syntheticData, emptyTraces)
}

// GenerateSyntheticReports creates multiple varied synthetic reports
func (g *EnhancedReportGenerator) GenerateSyntheticReports(count int, variabilityFactor float64) ([]string, error) {
	var reports []string

	for i := 0; i < count; i++ {
		report, err := g.GenerateSyntheticReport(variabilityFactor)
		if err != nil {
			return nil, fmt.Errorf("failed to generate report %d: %v", i+1, err)
		}
		reports = append(reports, report)
	}

	return reports, nil
}

// GenerateSyntheticData creates realistic synthetic medical data
func (g *EnhancedReportGenerator) generateSyntheticData(variabilityFactor float64) map[string]string {
	data := make(map[string]string)

	// Patient information
	data["patient_age"] = g.generatePatientAge()
	data["patient_gender"] = g.randomChoice([]string{"female", "male"})

	// Clinical information
	data["clinical_scenario"] = g.generateClinicalScenario()
	data["exam_indication"] = g.generateExamIndication()

	// Procedure information
	data["procedure"] = g.generateProcedure()
	data["body_part"] = g.generateBodyPart()
	data["contrast_used"] = g.generateContrastUsage()

	// Findings
	data["findings"] = g.generateFindings(variabilityFactor)
	data["impression"] = g.generateImpression()
	data["recommendations"] = g.generateRecommendations()

	// Additional realistic details
	if g.rand.Float32() < 0.7 {
		data["radiation_dose"] = g.generateRadiationDose()
	}
	if g.rand.Float32() < 0.5 {
		data["comparison"] = g.generateComparison()
	}

	return data
}

// ExecuteString executes the template with the given data
func (g *EnhancedReportGenerator) ExecuteString(data map[string]string, ruleMatches map[string][]result.RuleTrace) (string, error) {
	templateData := struct {
		Header          string
		PatientInfo     string
		ClinicalInfo    string
		Procedures      string
		Findings        string
		Impression      string
		Recommendations string
		AdditionalInfo  string
		Footer          string
	}{
		Header:          g.generateHeader(),
		PatientInfo:     g.generatePatientInfoSection(data),
		ClinicalInfo:    g.generateClinicalInfoSection(data),
		Procedures:      g.generateProceduresSection(data),
		Findings:        g.generateFindingsSection(data),
		Impression:      g.generateImpressionSection(data),
		Recommendations: g.generateRecommendationsSection(data),
		AdditionalInfo:  g.generateAdditionalInfoSection(data),
		Footer:          g.generateFooter(),
	}

	var buf bytes.Buffer
	err := g.baseTemplate.Execute(&buf, templateData)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

// Section generation functions
func (g *EnhancedReportGenerator) generateHeader() string {
	headers := []string{
		"RADIOLOGY REPORT",
		"IMAGING REPORT",
		"DIAGNOSTIC RADIOLOGY REPORT",
		"MEDICAL IMAGING REPORT",
	}
	return g.randomChoice(headers)
}

func (g *EnhancedReportGenerator) generatePatientInfoSection(data map[string]string) string {
	if data["patient_age"] == "" && data["patient_gender"] == "" {
		return ""
	}

	infoParts := []string{}
	if data["patient_age"] != "" {
		infoParts = append(infoParts, fmt.Sprintf("%s-year-old", data["patient_age"]))
	}
	if data["patient_gender"] != "" {
		infoParts = append(infoParts, data["patient_gender"])
	}

	return fmt.Sprintf("Patient: %s", strings.Join(infoParts, ", "))
}

func (g *EnhancedReportGenerator) generateClinicalInfoSection(data map[string]string) string {
	parts := []string{}

	if data["clinical_scenario"] != "" {
		parts = append(parts, data["clinical_scenario"])
	}
	if data["exam_indication"] != "" {
		parts = append(parts, fmt.Sprintf("Indication: %s", data["exam_indication"]))
	}

	if len(parts) == 0 {
		return ""
	}
	return strings.Join(parts, ". ") + "."
}

func (g *EnhancedReportGenerator) generateProceduresSection(data map[string]string) string {
	procedure := data["procedure"]
	if procedure == "" {
		return "Standard imaging examination"
	}

	if data["contrast_used"] != "" {
		return fmt.Sprintf("%s with %s", procedure, data["contrast_used"])
	}
	return procedure
}

func (g *EnhancedReportGenerator) generateFindingsSection(data map[string]string) string {
	findings := data["findings"]
	if findings == "" {
		return "No significant abnormalities detected."
	}

	if data["comparison"] != "" {
		return fmt.Sprintf("%s %s", findings, data["comparison"])
	}
	return findings
}

func (g *EnhancedReportGenerator) generateImpressionSection(data map[string]string) string {
	impression := data["impression"]
	if impression == "" {
		return "Normal examination with no acute abnormalities."
	}
	return impression
}

func (g *EnhancedReportGenerator) generateRecommendationsSection(data map[string]string) string {
	recommendations := data["recommendations"]
	if recommendations == "" {
		return ""
	}
	return recommendations
}

func (g *EnhancedReportGenerator) generateAdditionalInfoSection(data map[string]string) string {
	var infoParts []string

	if data["radiation_dose"] != "" {
		infoParts = append(infoParts, fmt.Sprintf("Radiation dose: %s", data["radiation_dose"]))
	}

	if len(infoParts) == 0 {
		return ""
	}
	return strings.Join(infoParts, "\n")
}

func (g *EnhancedReportGenerator) generateFooter() string {
	footers := []string{
		"End of report",
		"Report completed",
		"---",
		"",
	}
	return g.randomChoice(footers)
}

// Data generation functions
func (g *EnhancedReportGenerator) generatePatientAge() string {
	ages := []string{"25", "32", "45", "58", "67", "72", "81", "38", "42", "55", "62", "78"}
	return g.randomChoice(ages)
}

func (g *EnhancedReportGenerator) generateClinicalScenario() string {
	scenarios := []string{
		"Patient presents with right breast lump",
		"Follow-up examination for known breast cancer",
		"Screening mammography for high-risk patient",
		"Evaluation of palpable mass in left breast",
		"Post-treatment surveillance",
		"Assessment of breast pain and nipple discharge",
		"Preoperative staging for newly diagnosed carcinoma",
		"Routine screening examination",
	}
	return g.randomChoice(scenarios)
}

func (g *EnhancedReportGenerator) generateExamIndication() string {
	indications := []string{
		"screening",
		"diagnostic workup",
		"follow-up",
		"staging",
		"treatment response assessment",
		"preoperative planning",
		"surveillance",
	}
	return g.randomChoice(indications)
}

func (g *EnhancedReportGenerator) generateProcedure() string {
	procedures := []string{
		"Bilateral mammography",
		"Digital breast tomosynthesis",
		"Breast MRI with contrast",
		"Targeted ultrasound examination",
		"Stereotactic core biopsy",
		"MRI-guided vacuum-assisted biopsy",
		"Diagnostic mammography with magnification views",
		"Whole breast ultrasound",
	}
	return g.randomChoice(procedures)
}

func (g *EnhancedReportGenerator) generateBodyPart() string {
	bodyParts := []string{
		"right breast",
		"left breast",
		"bilateral breasts",
		"upper outer quadrant",
		"retroareolar region",
		"axillary tail",
		"central breast",
	}
	return g.randomChoice(bodyParts)
}

func (g *EnhancedReportGenerator) generateContrastUsage() string {
	if g.rand.Float32() < 0.6 {
		return ""
	}
	usage := []string{
		"intravenous contrast administration",
		"gadolinium-based contrast agent",
		"contrast-enhanced imaging",
	}
	return g.randomChoice(usage)
}

func (g *EnhancedReportGenerator) generateFindings(variabilityFactor float64) string {
	baseFindings := []string{
		"The breasts demonstrate heterogeneous fibroglandular tissue",
		"Scattered fibroglandular densities are present bilaterally",
		"The breast parenchyma shows age-appropriate involution",
		"No suspicious masses, calcifications, or architectural distortions are identified",
		"Benign-appearing cysts are noted in the upper outer quadrants",
		"A well-circumscribed oval mass is seen in the right breast at 10 o'clock position",
		"Multiple bilateral simple cysts are present",
		"The implants appear intact with no evidence of rupture",
	}

	findings := g.randomChoice(baseFindings)

	// Add variability based on factor
	if variabilityFactor > 0.5 && g.rand.Float64() < variabilityFactor {
		additionalDetails := []string{
			"No suspicious microcalcifications are present.",
			"No axillary adenopathy is identified.",
			"The skin and nipple-areolar complexes appear normal.",
			"No architectural distortion or focal asymmetry is seen.",
			"The findings are stable compared to prior examination.",
		}
		findings += " " + g.randomChoice(additionalDetails)
	}

	return findings
}

func (g *EnhancedReportGenerator) generateImpression() string {
	impressions := []string{
		"1. Normal screening examination. BIRADS category 1.",
		"2. Benign findings. BIRADS category 2.",
		"3. Probably benign finding. BIRADS category 3. Short-term follow-up recommended.",
		"4. Suspicious abnormality. BIRADS category 4. Biopsy should be considered.",
		"5. Highly suggestive of malignancy. BIRADS category 5. Appropriate action should be taken.",
		"No evidence of malignancy or suspicious findings.",
		"Findings consistent with benign breast disease.",
		"Recommend correlation with clinical findings.",
	}
	return g.randomChoice(impressions)
}

func (g *EnhancedReportGenerator) generateRecommendations() string {
	if g.rand.Float32() < 0.4 {
		return ""
	}

	recommendations := []string{
		"Correlation with clinical examination is recommended.",
		"Follow-up imaging in 6 months is suggested for the probably benign finding.",
		"Biopsy is recommended for the suspicious mass.",
		"Continue annual screening mammography.",
		"No additional imaging is necessary at this time.",
		"Clinical correlation and consideration of biopsy is advised.",
		"Recommend multidisciplinary tumor board review.",
	}
	return g.randomChoice(recommendations)
}

func (g *EnhancedReportGenerator) generateRadiationDose() string {
	doses := []string{
		"2.5 mGy",
		"3.1 mGy",
		"1.8 mGy",
		"2.2 mGy",
		"within standard limits",
		"as low as reasonably achievable",
	}
	return g.randomChoice(doses)
}

func (g *EnhancedReportGenerator) generateComparison() string {
	comparisons := []string{
		"These findings are stable compared to the prior examination from 6 months ago.",
		"No significant interval change is noted.",
		"The current findings represent expected postoperative changes.",
		"The mass has decreased in size since the previous study.",
		"New findings are identified since the last examination.",
	}
	return g.randomChoice(comparisons)
}

// Utility functions
func (g *EnhancedReportGenerator) randomChoice(options []string) string {
	if len(options) == 0 {
		return ""
	}
	return options[g.rand.Intn(len(options))]
}

// GenerateReport implements the ReportGenerator interface
func (g *EnhancedReportGenerator) GenerateReport(data map[string]string, ruleMatches map[string][]result.RuleTrace) (string, error) {
	// For synthetic reports, we generate our own data but can mix with provided data
	syntheticData := g.generateSyntheticData(0.8)

	// Mix in any provided real data
	for k, v := range data {
		if v != "" {
			syntheticData[k] = v
		}
	}

	return g.ExecuteString(syntheticData, ruleMatches)
}
