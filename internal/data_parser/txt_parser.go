package data_parser

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// TXTParser parses plain text files
type TXTParser struct{}

// Parse implements DocumentParser interface for TXT files
func (p *TXTParser) Parse(filePath string) (DocumentContent, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return DocumentContent{}, err
	}

	doc := DocumentContent{
		FilePath:   filePath,
		FileType:   "txt",
		RawContent: string(content),
		Metadata:   make(map[string]string),
	}

	// Extract basic metadata
	doc.Metadata["file_size"] = fmt.Sprintf("%d", len(content))
	doc.Metadata["line_count"] = fmt.Sprintf("%d", strings.Count(string(content), "\n"))

	// Parse structure
	doc.Structure = parseTXTStructure(string(content))

	return doc, nil
}

// ExtractMetadata extracts metadata from TXT content
func (p *TXTParser) ExtractMetadata(content DocumentContent) map[string]string {
	metadata := make(map[string]string)
	
	// Copy existing metadata
	for k, v := range content.Metadata {
		metadata[k] = v
	}

	// Add content-based metadata
	lines := strings.Split(content.RawContent, "\n")
	if len(lines) > 0 {
		metadata["first_line"] = lines[0]
	}
	if len(lines) > 1 {
		metadata["second_line"] = lines[1]
	}

	return metadata
}

func parseTXTStructure(content string) DocumentStructure {
	var structure DocumentStructure
	var currentSection *Section

	scanner := bufio.NewScanner(strings.NewReader(content))
	lineNum := 0

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			lineNum++
			continue
		}

		// Detect section headers (lines with special formatting)
		if isSectionHeader(line) {
			// Save previous section if exists
			if currentSection != nil {
				structure.Sections = append(structure.Sections, *currentSection)
			}

			// Start new section
			currentSection = &Section{
				Title:   cleanSectionTitle(line),
				Content: "",
				Level:   detectSectionLevel(line),
			}
		} else if currentSection != nil {
			// Add to current section content
			if currentSection.Content != "" {
				currentSection.Content += "\n"
			}
			currentSection.Content += line
		}

		lineNum++
	}

	// Don't forget the last section
	if currentSection != nil {
		structure.Sections = append(structure.Sections, *currentSection)
	}

	return structure
}

func isSectionHeader(line string) bool {
	// Section headers often have special formatting
	if len(line) == 0 {
		return false
	}

	// Check for common header patterns
	uppercaseChars := 0
	for _, char := range line {
		if char >= 'A' && char <= 'Z' {
			uppercaseChars++
		}
	}

	// If most characters are uppercase, likely a header
	if float64(uppercaseChars)/float64(len(line)) > 0.7 {
		return true
	}

	// Check for common header endings
	if strings.HasSuffix(line, ":") || strings.HasSuffix(line, "-") {
		return true
	}

	// Check for underlines (next line would be dashes/equals)
	return false
}

func cleanSectionTitle(line string) string {
	// Remove trailing punctuation
	line = strings.TrimRight(line, ":-.")
	return strings.TrimSpace(line)
}

func detectSectionLevel(line string) int {
	// Simple level detection based on indentation or formatting
	if strings.HasPrefix(line, "  ") {
		return 2
	}
	if len(line) < 20 {
		return 1
	}
	return 1
}