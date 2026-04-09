package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Joel-Haeberli/cde-fuzzer/internal/core"
)

func main() {
	count := flag.Int("count", 3, "number of synthetic reports to generate")
	variability := flag.Float64("variability", 0.8, "variability factor (0.0-1.0) for report diversity")
	outputDir := flag.String("output", "./synthetic_reports", "directory to save synthetic reports")
	flag.Parse()

	fmt.Println("🔮 CDE Extractor - Synthetic Report Generator")
	fmt.Printf("📝 Generating %d synthetic reports with variability %.1f\n", *count, *variability)

	// Create report generator
	generator, err := core.NewEnhancedReportGenerator()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error creating report generator: %v\n", err)
		os.Exit(1)
	}

	// Generate synthetic reports
	reports, err := generator.GenerateSyntheticReports(*count, *variability)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error generating reports: %v\n", err)
		os.Exit(1)
	}

	// Create output directory if it doesn't exist
	err = os.MkdirAll(*outputDir, 0755)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error creating output directory: %v\n", err)
		os.Exit(1)
	}

	// Save reports to files
	for i, report := range reports {
		filePath := fmt.Sprintf("%s/synthetic_report_%d.txt", *outputDir, i+1)
		err = os.WriteFile(filePath, []byte(report), 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Error writing report %d: %v\n", i+1, err)
			continue
		}
		fmt.Printf("✅ Generated: %s\n", filePath)
	}

	fmt.Printf("🎉 Successfully generated %d synthetic reports\n", len(reports))
	fmt.Printf("📁 Reports saved to: %s/\n", *outputDir)
}
