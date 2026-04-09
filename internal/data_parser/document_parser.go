package data_parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// DocumentContent represents parsed document content with metadata
type DocumentContent struct {
	FilePath    string
	FileType    string
	RawContent  string
	Metadata    map[string]string
	Tables      []Table
	Structure   DocumentStructure
}

// Table represents a parsed table from documents
type Table struct {
	Headers []string
	Rows    [][]string
}

// DocumentStructure represents document organization
type DocumentStructure struct {
	Sections []Section
}

// Section represents a document section
type Section struct {
	Title     string
	Content   string
	Level     int
}

// DocumentParser interface for document parsing
type DocumentParser interface {
	Parse(filePath string) (DocumentContent, error)
	ExtractMetadata(content DocumentContent) map[string]string
}

// ParseDocument parses a document based on its file extension
func ParseDocument(filePath string) (DocumentContent, error) {
	ext := strings.ToLower(filepath.Ext(filePath))

	var parser DocumentParser
	switch ext {
	case ".csv":
		parser = &CSVParser{}
	case ".docx":
		parser = &DOCXParser{}
	case ".pdf":
		parser = &PDFParser{}
	case ".xlsx", ".xls":
		parser = &XLSXParser{}
	case ".txt":
		parser = &TXTParser{}
	default:
		return DocumentContent{}, fmt.Errorf("unsupported file type: %s", ext)
	}

	return parser.Parse(filePath)
}

// ParseDirectory parses all documents in a directory
func ParseDirectory(dirPath string) ([]DocumentContent, error) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %v", err)
	}

	var documents []DocumentContent
	for _, file := range files {
		if file.IsDir() {
			// Skip directories for now, could add recursive parsing later
			continue
		}

		filePath := filepath.Join(dirPath, file.Name())
		doc, err := ParseDocument(filePath)
		if err != nil {
			fmt.Printf("Warning: failed to parse %s: %v\n", filePath, err)
			continue
		}

		documents = append(documents, doc)
	}

	return documents, nil
}

// IdentifyDataElements analyzes document content to identify potential CDEs
func IdentifyDataElements(content DocumentContent) []PotentialCDE {
	// This will be implemented with more sophisticated analysis
	var potentialCDEs []PotentialCDE
	
	// Basic pattern matching for common medical data elements
	commonPatterns := []struct {
		name     string
		patterns []string
	}{
		{"age", []string{"alter", "age", "jahre", "years"}},
		{"date", []string{"datum", "date", "\\d{2}\\.\\d{2}\\.\\d{4}"}},
		{"patient_id", []string{"patienten-id", "patient id", "fallnummer"}},
		{"diagnosis", []string{"diagnose", "befund", "findings"}},
	}

	for _, patternGroup := range commonPatterns {
		for _, pattern := range patternGroup.patterns {
			if strings.Contains(strings.ToLower(content.RawContent), pattern) {
				potentialCDEs = append(potentialCDEs, PotentialCDE{
					Name:    patternGroup.name,
					Pattern:  pattern,
					Source:   content.FilePath,
				})
				break // Only add each CDE type once per document
			}
		}
	}

	return potentialCDEs
}

// IdentifyCSVDataElements is a specialized function for CSV files
// This is defined in csv_parser.go and provides better CDE identification for structured data

// PotentialCDE represents a potential common data element identified in documents
type PotentialCDE struct {
	Name    string
	Pattern  string
	Source   string
	Confidence float64
}