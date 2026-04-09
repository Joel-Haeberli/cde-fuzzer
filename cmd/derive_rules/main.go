package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Joel-Haeberli/cde-fuzzer/internal/data_parser"
	"github.com/Joel-Haeberli/cde-fuzzer/internal/rule_derivation"
)

func main() {
	dataDir := flag.String("data", "./data", "directory containing data sources")
	outputDir := flag.String("output", "./derived_rules", "directory to save derived rules")
	recursive := flag.Bool("recursive", false, "parse directories recursively")
	flag.Parse()

	fmt.Println("🔍 CDE Extractor - Rule Derivation Tool")
	fmt.Printf("📁 Scanning directory: %s\n", *dataDir)

	// Parse all documents
	var documents []data_parser.DocumentContent
	var err error

	if *recursive {
		documents, err = parseDirectoryRecursive(*dataDir)
	} else {
		documents, err = data_parser.ParseDirectory(*dataDir)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error parsing documents: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("📄 Found %d documents\n", len(documents))

	// Derive rules from documents
	rules, err := rule_derivation.DeriveRulesFromDocuments(documents)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error deriving rules: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("🔧 Derived %d rules\n", len(rules))

	// Save rules as YAML files
	err = rule_derivation.SaveRulesAsYAML(rules, *outputDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error saving rules: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Successfully derived %d rules from %d documents\n", len(rules), len(documents))
	fmt.Printf("📁 Rules saved to: %s/\n", *outputDir)
}

// parseDirectoryRecursive parses documents in a directory and its subdirectories
func parseDirectoryRecursive(dirPath string) ([]data_parser.DocumentContent, error) {
	var allDocuments []data_parser.DocumentContent

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			// Try to parse the file
			doc, err := data_parser.ParseDocument(path)
			if err != nil {
				fmt.Printf("⚠️  Warning: failed to parse %s: %v\n", path, err)
				return nil // Continue with other files
			}
			allDocuments = append(allDocuments, doc)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk directory: %v", err)
	}

	return allDocuments, nil
}
