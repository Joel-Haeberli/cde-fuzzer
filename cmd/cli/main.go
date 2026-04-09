package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"

	"github.com/Joel-Haeberli/cde-fuzzer/internal/cde"
	"github.com/Joel-Haeberli/cde-fuzzer/internal/core"
)

func main() {
	filePath := flag.String("file", "", "path to text file to extract from")
	rulesDir := flag.String("rules", "", "path to directory containing rule YAML files")
	useLLM := flag.Bool("llm", false, "enable LLM rules (requires LLM_API_URL and LLM_API_KEY env vars)")
	flag.Parse()

	if *filePath == "" {
		fmt.Fprintln(os.Stderr, "Usage: cde-cli -file <path> -rules <rules-dir>")
		flag.PrintDefaults()
		os.Exit(1)
	}

	data, err := os.ReadFile(*filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading file: %v\n", err)
		os.Exit(1)
	}
	text := string(data)

	// Load rules from directory
	var rules []core.Rule
	if *rulesDir != "" {
		rules, err = core.LoadRulesFromDirectory(*rulesDir, *useLLM)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error loading rules: %v\n", err)
			os.Exit(1)
		}
	} else {
		// Fallback to demo rule if no rules directory is provided
		fmt.Fprintln(os.Stderr, "No rules directory provided, using demo rule")
		rules = append(rules, core.NewRegexRule(
			"age-regex",
			regexp.MustCompile(`\b\d{1,3}\s*(?:years?\s*old|yo)\b`),
			0.85,
		))
	}

	// Example: a simple demo extraction process.
	demoCDE := cde.CDE{
		ID:       "demo-1",
		Question: "What is the patient's age?",
		Answers:  []string{},
	}

	chain := core.NewRuleChain("extraction-chain", rules...)
	estimator := &core.DefaultAccuracyEstimator{}
	process := core.NewExtractionProcess(demoCDE, chain, estimator)

	result, err := process.Run(text)
	if err != nil {
		fmt.Fprintf(os.Stderr, "extraction error: %v\n", err)
		os.Exit(1)
	}

	// New output format: show only rule matches
	if len(result.Traces) == 0 {
		fmt.Println("No matches found.")
	} else {
		for _, t := range result.Traces {
			fmt.Printf("Rule: %s\n", t.RuleName)
			fmt.Printf("  Match: %q\n", t.Match.Value)
			fmt.Printf("  Position: %d-%d\n", t.Match.Start, t.Match.End)
			fmt.Printf("  Accuracy: %.2f\n", t.Accuracy)
			fmt.Println()
		}
	}
}
