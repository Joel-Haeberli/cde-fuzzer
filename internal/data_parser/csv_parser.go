package data_parser

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

// CSVParser parses CSV files
type CSVParser struct{}

// Parse implements DocumentParser interface for CSV files
func (p *CSVParser) Parse(filePath string) (DocumentContent, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return DocumentContent{}, fmt.Errorf("failed to open CSV file: %v", err)
	}
	defer file.Close()

	// Read CSV data
	reader := csv.NewReader(file)
	// Allow flexible number of fields per row
	reader.FieldsPerRecord = -1

	// Read all records
	records, err := reader.ReadAll()
	if err != nil {
		return DocumentContent{}, fmt.Errorf("failed to read CSV data: %v", err)
	}

	if len(records) == 0 {
		return DocumentContent{}, fmt.Errorf("CSV file is empty: %s", filePath)
	}

	// Parse CSV structure
	headers, rows, metadata := parseCSVStructure(records)

	// Build raw content representation
	var content strings.Builder
	content.WriteString("CSV File: ")
	content.WriteString(filePath)
	content.WriteString("\n\n")
	content.WriteString("Headers: ")
	content.WriteString(strings.Join(headers, ", "))
	content.WriteString("\n\n")
	content.WriteString("Row Count: ")
	content.WriteString(fmt.Sprintf("%d", len(rows)))
	content.WriteString("\n\n")
	content.WriteString("=== CSV Data ===\n")

	// Add table representation
	for i, row := range rows {
		content.WriteString(fmt.Sprintf("Row %d: %v\n", i+1, row))
	}

	docContent := DocumentContent{
		FilePath:   filePath,
		FileType:   "csv",
		RawContent: content.String(),
		Metadata:   metadata,
		Tables:     []Table{{Headers: headers, Rows: rows}},
	}

	return docContent, nil
}

// ExtractMetadata extracts metadata from CSV content
func (p *CSVParser) ExtractMetadata(content DocumentContent) map[string]string {
	metadata := make(map[string]string)
	
	// Copy existing metadata
	for k, v := range content.Metadata {
		metadata[k] = v
	}

	// Add CSV-specific metadata
	if len(content.Tables) > 0 {
		table := content.Tables[0]
		metadata["header_count"] = fmt.Sprintf("%d", len(table.Headers))
		metadata["row_count"] = fmt.Sprintf("%d", len(table.Rows))
		if len(table.Headers) > 0 {
			metadata["first_header"] = table.Headers[0]
		}
		if len(table.Rows) > 0 && len(table.Rows[0]) > 0 {
			metadata["first_cell"] = table.Rows[0][0]
		}
	}

	return metadata
}

func parseCSVStructure(records [][]string) ([]string, [][]string, map[string]string) {
	var headers []string
	var rows [][]string
	metadata := make(map[string]string)

	if len(records) == 0 {
		return headers, rows, metadata
	}

	// First row is typically headers
	headers = records[0]
	metadata["header_count"] = fmt.Sprintf("%d", len(headers))

	// Remaining rows are data
	if len(records) > 1 {
		rows = records[1:]
		metadata["row_count"] = fmt.Sprintf("%d", len(rows))
	}

	// Count non-empty cells
	nonEmptyCount := 0
	for _, row := range records {
		for _, cell := range row {
			if strings.TrimSpace(cell) != "" {
				nonEmptyCount++
			}
		}
	}
	metadata["non_empty_cells"] = fmt.Sprintf("%d", nonEmptyCount)

	// Calculate density
	totalCells := 0
	for _, row := range records {
		totalCells += len(row)
	}
	if totalCells > 0 {
		density := float64(nonEmptyCount) / float64(totalCells) * 100
		metadata["data_density"] = fmt.Sprintf("%.1f%%", density)
	}

	return headers, rows, metadata
}

// IdentifyCSVDataElements analyzes CSV content to identify potential CDEs
func IdentifyCSVDataElements(content DocumentContent) []PotentialCDE {
	var potentialCDEs []PotentialCDE

	if len(content.Tables) == 0 {
		return potentialCDEs
	}

	table := content.Tables[0]
	if len(table.Headers) == 0 {
		return potentialCDEs
	}

	// Check each header for common CDE patterns
	for _, header := range table.Headers {
		cleanHeader := strings.ToLower(strings.TrimSpace(header))
		
		// Common CDE patterns in medical data
		if strings.Contains(cleanHeader, "age") || 
		   strings.Contains(cleanHeader, "alter") ||
		   strings.Contains(cleanHeader, "patient") && strings.Contains(cleanHeader, "id") {
			potentialCDEs = append(potentialCDEs, PotentialCDE{
				Name:    "patient_identifier",
				Pattern:  header,
				Source:   content.FilePath,
				Confidence: 0.9,
			})
		} else if strings.Contains(cleanHeader, "date") || 
		          strings.Contains(cleanHeader, "datum") ||
		          strings.Contains(cleanHeader, "time") {
			potentialCDEs = append(potentialCDEs, PotentialCDE{
				Name:    "date_time",
				Pattern:  header,
				Source:   content.FilePath,
				Confidence: 0.85,
			})
		} else if strings.Contains(cleanHeader, "diagnosis") || 
		          strings.Contains(cleanHeader, "befund") ||
		          strings.Contains(cleanHeader, "finding") {
			potentialCDEs = append(potentialCDEs, PotentialCDE{
				Name:    "clinical_finding",
				Pattern:  header,
				Source:   content.FilePath,
				Confidence: 0.9,
			})
		} else if strings.Contains(cleanHeader, "procedure") || 
		          strings.Contains(cleanHeader, "exam") ||
		          strings.Contains(cleanHeader, "test") {
			potentialCDEs = append(potentialCDEs, PotentialCDE{
				Name:    "medical_procedure",
				Pattern:  header,
				Source:   content.FilePath,
				Confidence: 0.8,
			})
		} else if strings.Contains(cleanHeader, "result") || 
		          strings.Contains(cleanHeader, "value") ||
		          strings.Contains(cleanHeader, "score") {
			potentialCDEs = append(potentialCDEs, PotentialCDE{
				Name:    "test_result",
				Pattern:  header,
				Source:   content.FilePath,
				Confidence: 0.75,
			})
		}
	}

	return potentialCDEs
}