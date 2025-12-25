// Package pathutil provides path resolution utilities for document files.
package pathutil

import (
	"os"
	"path/filepath"
	"strings"
)

const (
	// DefaultDocumentsDir is the default parent directory for documents
	DefaultDocumentsDir = "documents"
)

// GetFileType determines the file type from a filename
// Uses the same logic as downloader.GetFileType for consistency
func GetFileType(filename string) string {
	// Import downloader's GetFileType to avoid duplication
	// For now, we'll keep a local implementation but make it consistent
	ext := strings.ToLower(filepath.Ext(filename))
	
	// Map extensions to file types
	extMap := map[string]string{
		".pdf":  "pdf",
		".doc":  "doc",
		".docx": "docx",
		".rtf":  "rtf",
		".txt":  "txt",
		".odt":  "odt",
	}
	
	if fileType, ok := extMap[ext]; ok {
		return fileType
	}
	
	// Default to "pdf" for backward compatibility if no extension
	return "pdf"
}

// ResolveDocumentPath resolves a document file path, checking the documents directory if it's just a filename
func ResolveDocumentPath(input string) string {
	// If it's already an absolute path or contains directory separators, use as-is
	if filepath.IsAbs(input) || strings.Contains(input, string(filepath.Separator)) {
		return input
	}

	// Determine file type from extension
	fileType := GetFileType(input)
	
	// Get base name without extension for subdirectory lookup
	ext := filepath.Ext(input)
	baseName := strings.TrimSuffix(input, ext)
	baseName = strings.TrimSuffix(baseName, strings.ToUpper(ext))
	
	// Check in documents/{type}/{basename}/{filename} first (new structure)
	typeDir := filepath.Join(DefaultDocumentsDir, fileType)
	docSubDir := filepath.Join(typeDir, baseName)
	docPath := filepath.Join(docSubDir, input)
	if _, err := os.Stat(docPath); err == nil {
		return docPath
	}

	// Fall back to old structure: documents/{type}/{filename}
	docPath = filepath.Join(typeDir, input)
	if _, err := os.Stat(docPath); err == nil {
		return docPath
	}

	// Fall back to legacy structure: pdfs/{basename}/{filename} (for backward compatibility with PDFs)
	if fileType == "pdf" {
		legacySubDir := filepath.Join("pdfs", baseName)
		docPath = filepath.Join(legacySubDir, input)
		if _, err := os.Stat(docPath); err == nil {
			return docPath
		}

		// Fall back to legacy structure: pdfs/{filename}
		docPath = filepath.Join("pdfs", input)
		if _, err := os.Stat(docPath); err == nil {
			return docPath
		}
	}

	// Fall back to treating as relative path from current directory
	return input
}

// ResolvePDFPath is a legacy alias for ResolveDocumentPath (for backward compatibility)
func ResolvePDFPath(input string) string {
	return ResolveDocumentPath(input)
}

