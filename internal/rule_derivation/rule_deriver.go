package rule_derivation

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"

	"github.com/Joel-Haeberli/cde-fuzzer/internal/core"
	"github.com/Joel-Haeberli/cde-fuzzer/internal/data_parser"
	"gopkg.in/yaml.v3"
)

// RuleConfig represents a rule configuration for derivation
type RuleConfig struct {
	Name      string  `yaml:"name"`
	Type      string  `yaml:"type"`
	Pattern   string  `yaml:"pattern,omitempty"`
	Accuracy  float64 `yaml:"accuracy,omitempty"`
	Target    string  `yaml:"target,omitempty"`
	Threshold float64 `yaml:"threshold,omitempty"`
	Prompt    string  `yaml:"prompt,omitempty"`
}

// DeriveRulesFromDocuments analyzes documents and derives extraction rules
func DeriveRulesFromDocuments(documents []data_parser.DocumentContent) ([]RuleConfig, error) {
	var allRules []RuleConfig

	// Analyze each document for potential rules
	for _, doc := range documents {
		fmt.Printf("Analyzing document: %s\n", doc.FilePath)

		// Identify potential CDEs
		var potentialCDEs []data_parser.PotentialCDE
		if doc.FileType == "csv" {
			// Use CSV-specific CDE identification for better results
			potentialCDEs = data_parser.IdentifyCSVDataElements(doc)
		} else {
			potentialCDEs = data_parser.IdentifyDataElements(doc)
		}
		fmt.Printf("Found %d potential CDEs\n", len(potentialCDEs))

		// Generate rules for each potential CDE
		for _, cde := range potentialCDEs {
			rules, err := generateRulesForCDE(cde, doc)
			if err != nil {
				fmt.Printf("Warning: failed to generate rules for %s: %v\n", cde.Name, err)
				continue
			}
			allRules = append(allRules, rules...)
		}
	}

	// Validate and filter rules
	validatedRules := validateRules(allRules)

	return validatedRules, nil
}

func generateRulesForCDE(cde data_parser.PotentialCDE, doc data_parser.DocumentContent) ([]RuleConfig, error) {
	var rules []RuleConfig

	// Generate regex rule
	regexRule, err := generateRegexRule(cde, doc)
	if err == nil {
		rules = append(rules, regexRule)
	}

	// Generate LLM rule for complex patterns
	llmRule := generateLLMRule(cde)
	rules = append(rules, llmRule)

	return rules, nil
}

func generateRegexRule(cde data_parser.PotentialCDE, doc data_parser.DocumentContent) (RuleConfig, error) {
	// Simple regex generation based on CDE name
	var pattern string

	switch cde.Name {
	case "age":
		pattern = "\\b(\\d{1,3})\\s*(?:Jahre?|Jahr|yo|years?|months?|mos?|days?|dys?|Monate?|Tage?|Jahren?)\\b"
	case "date":
		pattern = "\\b(\\d{2}\\.\\d{2}\\.\\d{4}|\\d{4}-\\d{2}-\\d{2}|\\d{1,2}\\s*(?:Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec)[a-z]*\\s*\\d{4})\\b"
	case "patient_id":
		pattern = "\\b(?:Patienten-ID|Patient ID|Fallnummer|Case ID|ID):?\\s*([A-Z0-9\\-]+)\\b"
	case "diagnosis":
		pattern = "\\b(?:Diagnose|Befund|Findings|Diagnosis):?\\s*([^.?!]*[.?!])"
	default:
		// Generic pattern for the CDE name
		pattern = fmt.Sprintf("\\b%s\\b", regexp.QuoteMeta(cde.Pattern))
	}

	return RuleConfig{
		Name:     fmt.Sprintf("%s_extractor", cde.Name),
		Type:     "regex",
		Pattern:  pattern,
		Accuracy: 0.85,
	}, nil
}

func generateLLMRule(cde data_parser.PotentialCDE) RuleConfig {
	// Generate LLM prompt based on CDE type
	var prompt string

	switch cde.Name {
	case "age":
		prompt = "Extract the patient's age from the text. Look for numbers followed by age-related terms like 'years', 'Jahre', 'yo', etc. If no age information is present, respond with 'no data'."
	case "date":
		prompt = "Extract any dates mentioned in the text. Look for date formats like DD.MM.YYYY, YYYY-MM-DD, or month names. Return dates in ISO format (YYYY-MM-DD) if possible."
	case "patient_id":
		prompt = "Extract the patient ID or case number from the text. Look for terms like 'Patienten-ID', 'Patient ID', 'Fallnummer', or 'Case ID' followed by alphanumeric identifiers."
	case "diagnosis":
		prompt = "Extract the diagnosis or main findings from the text. Look for sections headed by 'Diagnose', 'Befund', 'Findings', or 'Diagnosis' and extract the relevant clinical information."
	default:
		prompt = fmt.Sprintf("Extract information related to %s from the text. Be precise and return only the relevant data.", cde.Name)
	}

	return RuleConfig{
		Name:     fmt.Sprintf("llm_%s_extractor", cde.Name),
		Type:     "llm",
		Prompt:   prompt,
		Accuracy: 0.90,
	}
}

func validateRules(rules []RuleConfig) []RuleConfig {
	var validated []RuleConfig

	for _, rule := range rules {
		// Skip empty rules
		if rule.Name == "" || rule.Type == "" {
			continue
		}

		// Validate regex rules
		if rule.Type == "regex" && rule.Pattern != "" {
			_, err := regexp.Compile(rule.Pattern)
			if err != nil {
				fmt.Printf("Warning: invalid regex pattern in rule %s: %v\n", rule.Name, err)
				continue
			}
		}

		// Validate LLM rules
		if rule.Type == "llm" && rule.Prompt == "" {
			fmt.Printf("Warning: LLM rule %s has empty prompt\n", rule.Name)
			continue
		}

		// Set default accuracy if not specified
		if rule.Accuracy == 0 {
			if rule.Type == "regex" {
				rule.Accuracy = 0.85
			} else if rule.Type == "llm" {
				rule.Accuracy = 0.90
			} else {
				rule.Accuracy = 0.80
			}
		}

		validated = append(validated, rule)
	}

	return validated
}

// SaveRulesAsYAML saves rule configurations as YAML files
func SaveRulesAsYAML(rules []RuleConfig, outputDir string) error {
	// Create output directory if it doesn't exist
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	// First, load all existing rules to check for duplicates
	existingRules, err := loadExistingRules(outputDir)
	if err != nil {
		fmt.Printf("Warning: failed to load existing rules for duplicate check: %v\n", err)
		// Continue without duplicate checking if we can't load existing rules
		existingRules = []RuleConfig{}
	}

	// Save each rule as a separate YAML file
	for _, rule := range rules {
		// Check if this rule is a duplicate of an existing rule
		if isDuplicateRule(rule, existingRules) {
			fmt.Printf("⚠️  Skipping duplicate rule: %s\n", rule.Name)
			continue
		}

		yamlContent, err := convertRuleToYAML(rule)
		if err != nil {
			fmt.Printf("Warning: failed to convert rule %s to YAML: %v\n", rule.Name, err)
			continue
		}

		// Find a unique filename that doesn't already exist
		baseFileName := rule.Name + ".yaml"
		filePath := filepath.Join(outputDir, baseFileName)

		// Check if file already exists
		if _, err := os.Stat(filePath); err == nil {
			// File exists, find next available number
			nextNumber := findNextAvailableNumber(outputDir, rule.Name, ".yaml")
			if nextNumber > 0 {
				filePath = filepath.Join(outputDir, fmt.Sprintf("%s_%d.yaml", rule.Name, nextNumber))
			}
		} else if !os.IsNotExist(err) {
			// Other error occurred
			fmt.Printf("Warning: error checking file %s: %v\n", filePath, err)
			continue
		}

		err = os.WriteFile(filePath, []byte(yamlContent), 0644)
		if err != nil {
			fmt.Printf("Warning: failed to write rule file %s: %v\n", filePath, err)
			continue
		}

		fmt.Printf("Saved rule: %s\n", filePath)
	}

	return nil
}

// loadExistingRules loads all existing rule configurations from YAML files in a directory
func loadExistingRules(dirPath string) ([]RuleConfig, error) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %v", err)
	}

	var rules []RuleConfig
	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".yaml" {
			continue
		}

		filePath := filepath.Join(dirPath, file.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Printf("Warning: failed to read rule file %s: %v\n", filePath, err)
			continue
		}

		var ruleConfig RuleConfig
		err = yaml.Unmarshal(data, &ruleConfig)
		if err != nil {
			fmt.Printf("Warning: failed to parse rule file %s: %v\n", filePath, err)
			continue
		}

		rules = append(rules, ruleConfig)
	}

	return rules, nil
}

// isDuplicateRule checks if a rule is a duplicate of any existing rule
func isDuplicateRule(newRule RuleConfig, existingRules []RuleConfig) bool {
	for _, existingRule := range existingRules {
		// Rules are considered duplicates if they have the same name, type, and pattern/prompt
		if newRule.Name == existingRule.Name &&
			newRule.Type == existingRule.Type {

			// For regex rules, compare patterns
			if newRule.Type == "regex" && newRule.Pattern == existingRule.Pattern {
				return true
			}

			// For LLM rules, compare prompts
			if newRule.Type == "llm" && newRule.Prompt == existingRule.Prompt {
				return true
			}

			// For similarity rules, compare targets and thresholds
			if newRule.Type == "similarity" &&
				newRule.Target == existingRule.Target &&
				newRule.Threshold == existingRule.Threshold {
				return true
			}
		}
	}
	return false
}

// findNextAvailableNumber finds the next available number for a filename pattern
func findNextAvailableNumber(dirPath, baseName, extension string) int {
	// List all files in the directory
	files, err := os.ReadDir(dirPath)
	if err != nil {
		fmt.Printf("Warning: failed to read directory %s: %v\n", dirPath, err)
		return 1 // Default to 1 if we can't read directory
	}

	// Find all files that match the pattern baseName_N.extension
	var existingNumbers []int
	pattern := fmt.Sprintf("^%s_(\\d+)%s$", regexp.QuoteMeta(baseName), regexp.QuoteMeta(extension))
	regex, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Printf("Warning: invalid regex pattern for finding existing numbers: %v\n", err)
		return 1
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fileName := file.Name()
		matches := regex.FindStringSubmatch(fileName)
		if len(matches) == 2 {
			// matches[1] contains the number
			number, err := strconv.Atoi(matches[1])
			if err == nil {
				existingNumbers = append(existingNumbers, number)
			}
		}
	}

	// Find the maximum existing number
	if len(existingNumbers) > 0 {
		sort.Ints(existingNumbers)
		return existingNumbers[len(existingNumbers)-1] + 1
	}
	return 1
}

func convertRuleToYAML(rule RuleConfig) (string, error) {
	yamlData, err := yaml.Marshal(rule)
	if err != nil {
		return "", err
	}
	return string(yamlData), nil
}

// ConvertToCoreRules converts RuleConfig to core.Rule interface
// Note: This function creates rules manually since the core package doesn't export
// the rule creation function. For full functionality, you should save rules as YAML
// and use core.LoadRulesFromDirectory().
func ConvertToCoreRules(ruleConfigs []RuleConfig) ([]core.Rule, error) {
	var rules []core.Rule

	for _, config := range ruleConfigs {
		var rule core.Rule
		var err error

		switch config.Type {
		case "regex":
			if config.Pattern == "" {
				return nil, fmt.Errorf("pattern is required for regex rule %s", config.Name)
			}
			regex, err := regexp.Compile(config.Pattern)
			if err != nil {
				return nil, fmt.Errorf("invalid regex pattern in rule %s: %v", config.Name, err)
			}
			rule = core.NewRegexRule(config.Name, regex, config.Accuracy)

		case "llm":
			if config.Prompt == "" {
				return nil, fmt.Errorf("prompt is required for llm rule %s", config.Name)
			}
			// Use mock LLM client for now
			mockClient := &MockLLMClient{}
			rule = core.NewLLMRule(config.Name, config.Prompt, config.Accuracy, mockClient)

		case "similarity":
			if config.Target == "" {
				return nil, fmt.Errorf("target is required for similarity rule %s", config.Name)
			}
			rule = core.NewSimilarityRule(config.Name, config.Target, config.Threshold, core.Levenshtein)

		default:
			return nil, fmt.Errorf("unknown rule type: %s", config.Type)
		}

		if err != nil {
			return nil, fmt.Errorf("failed to create rule %s: %v", config.Name, err)
		}
		rules = append(rules, rule)
	}

	return rules, nil
}

// MockLLMClient is a simple mock implementation for testing
type MockLLMClient struct{}

func (m *MockLLMClient) Complete(prompt, text string) (string, error) {
	// Simple mock that returns a placeholder response
	return "mock response", nil
}

// DeriveRulesFromDirectory convenience function to derive rules from a directory
func DeriveRulesFromDirectory(dirPath string) ([]RuleConfig, error) {
	// Parse all documents in directory
	documents, err := data_parser.ParseDirectory(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse directory: %v", err)
	}

	// Derive rules from documents
	return DeriveRulesFromDocuments(documents)
}
