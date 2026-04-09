package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Joel-Haeberli/cde-fuzzer/internal/core"
)

func main() {
	rulesDir := flag.String("rules", "", "path to directory containing rule YAML files")
	dataDir := flag.String("data", "", "path to directory containing extracted data files")
	output := flag.String("output", "generated_report.txt", "output report file path")
	templateFile := flag.String("template", "templates/report_template.tmpl", "path to report template")
	flag.Parse()

	if *rulesDir == "" || *dataDir == "" {
		fmt.Fprintln(os.Stderr, "Usage: generate-report -rules <rules-dir> -data <data-dir> [-output <file>] [-template <template-file>]")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Load rules
	_, err := core.LoadRulesFromDirectory(*rulesDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading rules: %v\n", err)
		os.Exit(1)
	}

	// Load template
	templateContent, err := os.ReadFile(*templateFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading template: %v\n", err)
		os.Exit(1)
	}

	// Create report generator
	generator, err := core.NewTemplateReportGenerator(string(templateContent))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating generator: %v\n", err)
		os.Exit(1)
	}

	// Process data files
	dataFiles, err := filepath.Glob(filepath.Join(*dataDir, "*.txt"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error finding data files: %v\n", err)
		os.Exit(1)
	}

	if len(dataFiles) == 0 {
		fmt.Fprintf(os.Stderr, "no data files found in %s\n", *dataDir)
		os.Exit(1)
	}

	// For now, process the first file
	_, err = os.ReadFile(dataFiles[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading data file: %v\n", err)
		os.Exit(1)
	}

	// TODO: Parse extracted data and create rule traces
	// This is a simplified version - in a real implementation, we'd parse
	// the extracted data format and create proper RuleTrace objects

	// Create a simple demo report

	// For demo purposes, create some mock data
	demoData := map[string]string{
		"scenario_description_extractor": "Breast cancer staging evaluation",
		"procedure_extractor":            "MRI breast without and with IV contrast",
		"radiation_dose_extractor":       "0 mSvO",
		"body_part_extractor":            "breast",
		"appropriateness_extractor":      "Usually appropriate",
	}

	// Generate report
	report, err := generator.GenerateReport(demoData, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error generating report: %v\n", err)
		os.Exit(1)
	}

	// Write output
	err = os.WriteFile(*output, []byte(report), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error writing report: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated report: %s\n", *output)
}
