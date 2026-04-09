package data_parser

import (
	"fmt"
	"os"
)

// PDFParser parses PDF files (simplified version without external deps)
type PDFParser struct{}

// Parse implements DocumentParser interface for PDF files
func (p *PDFParser) Parse(filePath string) (DocumentContent, error) {
	// For now, treat PDF as binary file and extract basic info
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return DocumentContent{}, fmt.Errorf("failed to stat PDF file: %v", err)
	}

	docContent := DocumentContent{
		FilePath:   filePath,
		FileType:   "pdf",
		RawContent: fmt.Sprintf("[PDF file: %s, size: %d bytes]", fileInfo.Name(), fileInfo.Size()),
		Metadata:   make(map[string]string),
	}

	// Extract basic metadata
	docContent.Metadata["file_size"] = fmt.Sprintf("%d", fileInfo.Size())
	docContent.Metadata["file_type"] = "pdf"

	return docContent, nil
}

// ExtractMetadata extracts metadata from PDF content
func (p *PDFParser) ExtractMetadata(content DocumentContent) map[string]string {
	metadata := make(map[string]string)
	
	// Copy existing metadata
	for k, v := range content.Metadata {
		metadata[k] = v
	}

	return metadata
}