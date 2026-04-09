package data_parser

import (
	"fmt"
	"os"
)

// XLSXParser parses Excel files (simplified version without external deps)
type XLSXParser struct{}

// Parse implements DocumentParser interface for XLSX files
func (p *XLSXParser) Parse(filePath string) (DocumentContent, error) {
	// For now, treat XLSX as binary file and extract basic info
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return DocumentContent{}, fmt.Errorf("failed to stat XLSX file: %v", err)
	}

	docContent := DocumentContent{
		FilePath:   filePath,
		FileType:   "xlsx",
		RawContent: fmt.Sprintf("[XLSX file: %s, size: %d bytes]", fileInfo.Name(), fileInfo.Size()),
		Metadata:   make(map[string]string),
	}

	// Extract basic metadata
	docContent.Metadata["file_size"] = fmt.Sprintf("%d", fileInfo.Size())
	docContent.Metadata["file_type"] = "xlsx"

	return docContent, nil
}

// ExtractMetadata extracts metadata from XLSX content
func (p *XLSXParser) ExtractMetadata(content DocumentContent) map[string]string {
	metadata := make(map[string]string)
	
	// Copy existing metadata
	for k, v := range content.Metadata {
		metadata[k] = v
	}

	return metadata
}