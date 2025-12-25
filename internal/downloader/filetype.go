package downloader

import (
	"path/filepath"
	"strings"
)

// GetFileType determines the file type from a filename or URL
func GetFileType(filename string) string {
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
	
	// Default to "other" for unknown types
	return "other"
}

// GetDocumentsDir returns the directory path for a specific file type
func GetDocumentsDir(fileType string) string {
	return filepath.Join(DefaultDocumentsDir, fileType)
}

