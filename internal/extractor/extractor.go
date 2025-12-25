// Package extractor provides document text extraction functionality.
// Currently supports PDF files, with plans to support doc, docx, rtf, txt, and other formats.
package extractor

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ledongthuc/pdf"
)

// Extractor handles document text extraction
type Extractor struct {
	outputFormat string // "json", "markdown", or "plain"
}

// New creates a new Extractor instance with default plain text format
func New() *Extractor {
	return &Extractor{
		outputFormat: "json", // Default to JSON for structured output
	}
}

// NewWithFormat creates a new Extractor instance with specified format
func NewWithFormat(format string) *Extractor {
	validFormats := map[string]bool{"json": true, "markdown": true, "plain": true}
	if !validFormats[format] {
		format = "json" // Default to JSON if invalid
	}
	return &Extractor{
		outputFormat: format,
	}
}

// ExtractText extracts all text from a document file and returns plain text
func (e *Extractor) ExtractText(filePath string) (string, error) {
	_, fullText, _, err := e.ExtractTextStructured(filePath)
	return fullText, err
}

// ExtractTextStructured extracts all text from a document file with page information
func (e *Extractor) ExtractTextStructured(filePath string) ([]PageText, string, int, error) {
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, "", 0, fmt.Errorf("file does not exist: %s", filePath)
	}

	// Determine file type and extract accordingly
	ext := strings.ToLower(filepath.Ext(filePath))
	
	// Currently only PDF is supported, but structure is ready for other formats
	if ext == ".pdf" {
		return e.extractFromPDF(filePath)
	}
	
	// For other file types, return error (to be implemented)
	return nil, "", 0, fmt.Errorf("file type %s not yet supported (currently only PDF is supported)", ext)
}

// extractFromPDF extracts text from a PDF file
func (e *Extractor) extractFromPDF(filePath string) ([]PageText, string, int, error) {
	// Read PDF file
	file, reader, err := pdf.Open(filePath)
	if err != nil {
		return nil, "", 0, fmt.Errorf("failed to open document: %w (document may be encrypted or in an unsupported format)", err)
	}
	defer file.Close()

	var textBuilder strings.Builder
	var pages []PageText
	totalPages := reader.NumPage()

	if totalPages == 0 {
		return nil, "", 0, fmt.Errorf("document has no pages")
	}

	// Extract text from each page
	for i := 1; i <= totalPages; i++ {
		page := reader.Page(i)
		if page.V.IsNull() {
			// Skip null pages silently
			continue
		}

		text, err := page.GetPlainText(nil)
		if err != nil {
			// Try to continue with other pages
			continue
		}

		if text != "" {
			// Add page separator for multi-page documents in plain text
			if i > 1 {
				textBuilder.WriteString(fmt.Sprintf("\n\n--- Page %d ---\n\n", i))
			}
			textBuilder.WriteString(text)
			
			// Store page text
			pages = append(pages, PageText{
				PageNumber: i,
				Text:       text,
			})
		}
	}

	fullText := textBuilder.String()
	if fullText == "" {
		return nil, "", 0, fmt.Errorf("no text could be extracted from the document (document may be encrypted, image-based, or in an unsupported format)")
	}

	return pages, fullText, totalPages, nil
}

// SaveExtractedText saves extracted text to a file next to the document
func (e *Extractor) SaveExtractedText(filePath string, text string) (string, error) {
	// Get structured data for formatting
	pages, fullText, _, err := e.ExtractTextStructured(filePath)
	if err != nil {
		// Fall back to plain text if structured extraction fails
		return e.savePlainText(filePath, text)
	}

	// Get the directory and base filename of the document
	dir := filepath.Dir(filePath)
	baseName := filepath.Base(filePath)
	
	// Remove extension (any extension)
	ext := filepath.Ext(baseName)
	baseNameNoExt := strings.TrimSuffix(baseName, ext)
	baseNameNoExt = strings.TrimSuffix(baseNameNoExt, strings.ToUpper(ext))
	
	var extractedPath string
	var content []byte
	
	// Format based on output format
	switch e.outputFormat {
	case "json":
		extractedPath = filepath.Join(dir, baseNameNoExt+".extracted.json")
		content, err = FormatAsJSON(filePath, pages, fullText)
		if err != nil {
			return "", fmt.Errorf("failed to format as JSON: %w", err)
		}
	case "markdown":
		extractedPath = filepath.Join(dir, baseNameNoExt+".extracted.md")
		content, err = FormatAsMarkdown(filePath, pages, fullText)
		if err != nil {
			return "", fmt.Errorf("failed to format as Markdown: %w", err)
		}
	default: // plain
		extractedPath = filepath.Join(dir, baseNameNoExt+".extracted.txt")
		content = []byte(fullText)
	}
	
	// Write the formatted content to the file
	err = os.WriteFile(extractedPath, content, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write extracted text file: %w", err)
	}
	
	return extractedPath, nil
}

// savePlainText saves plain text as fallback
func (e *Extractor) savePlainText(filePath string, text string) (string, error) {
	dir := filepath.Dir(filePath)
	baseName := filepath.Base(filePath)
	
	// Remove any extension
	ext := filepath.Ext(baseName)
	baseNameNoExt := strings.TrimSuffix(baseName, ext)
	baseNameNoExt = strings.TrimSuffix(baseNameNoExt, strings.ToUpper(ext))
	extractedPath := filepath.Join(dir, baseNameNoExt+".extracted.txt")
	
	err := os.WriteFile(extractedPath, []byte(text), 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write extracted text file: %w", err)
	}
	
	return extractedPath, nil
}

