package data_parser

import (
	"fmt"
	"os"
)

// DOCXParser parses Word documents (simplified version without external deps)
type DOCXParser struct{}

// Parse implements DocumentParser interface for DOCX files
func (p *DOCXParser) Parse(filePath string) (DocumentContent, error) {
	// For now, treat DOCX as binary file and extract basic info
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return DocumentContent{}, fmt.Errorf("failed to stat DOCX file: %v", err)
	}

	docContent := DocumentContent{
		FilePath:   filePath,
		FileType:   "docx",
		RawContent: fmt.Sprintf("[DOCX file: %s, size: %d bytes]", fileInfo.Name(), fileInfo.Size()),
		Metadata:   make(map[string]string),
	}

	// Extract basic metadata
	docContent.Metadata["file_size"] = fmt.Sprintf("%d", fileInfo.Size())
	docContent.Metadata["file_type"] = "docx"

	return docContent, nil
}

// ExtractMetadata extracts metadata from DOCX content
func (p *DOCXParser) ExtractMetadata(content DocumentContent) map[string]string {
	metadata := make(map[string]string)
	
	// Copy existing metadata
	for k, v := range content.Metadata {
		metadata[k] = v
	}

	return metadata
}