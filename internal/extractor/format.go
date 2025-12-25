package extractor

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

// ExtractedText represents the structured format for extracted document text
type ExtractedText struct {
	Metadata Metadata `json:"metadata"`
	Content  Content  `json:"content"`
}

// Metadata contains information about the document and extraction
type Metadata struct {
	Filename       string    `json:"filename"`
	ExtractedAt    time.Time `json:"extracted_at"`
	TotalPages     int       `json:"total_pages"`
	PagesExtracted int       `json:"pages_extracted"`
	FormatVersion  string    `json:"format_version"`
}

// Content contains the extracted text organized by pages
type Content struct {
	FullText string   `json:"full_text"`
	Pages    []Page   `json:"pages"`
}

// Page represents text from a single page
type Page struct {
	PageNumber int    `json:"page_number"`
	Text       string `json:"text"`
	WordCount  int    `json:"word_count"`
}

// FormatVersion is the current format version
const FormatVersion = "1.0"

// FormatAsJSON formats extracted text as structured JSON
func FormatAsJSON(filePath string, pages []PageText, fullText string) ([]byte, error) {
	filename := filepath.Base(filePath)
	pagesExtracted := len(pages)
	
	// Count total pages (may include null pages)
	totalPages := 0
	for _, page := range pages {
		if page.PageNumber > totalPages {
			totalPages = page.PageNumber
		}
	}

	extracted := ExtractedText{
		Metadata: Metadata{
			Filename:       filename,
			ExtractedAt:    time.Now(),
			TotalPages:     totalPages,
			PagesExtracted: pagesExtracted,
			FormatVersion:  FormatVersion,
		},
		Content: Content{
			FullText: fullText,
			Pages:    make([]Page, 0, len(pages)),
		},
	}

	// Convert page text to structured pages
	for _, pageText := range pages {
		wordCount := len(strings.Fields(pageText.Text))
		extracted.Content.Pages = append(extracted.Content.Pages, Page{
			PageNumber: pageText.PageNumber,
			Text:       pageText.Text,
			WordCount:  wordCount,
		})
	}

	return json.MarshalIndent(extracted, "", "  ")
}

// FormatAsMarkdown formats extracted text as Markdown
func FormatAsMarkdown(filePath string, pages []PageText, fullText string) ([]byte, error) {
	filename := filepath.Base(filePath)
	var builder strings.Builder

	// Header
	builder.WriteString(fmt.Sprintf("# Document Text Extraction: %s\n\n", filename))
	builder.WriteString(fmt.Sprintf("**Extracted:** %s\n\n", time.Now().Format(time.RFC3339)))
	builder.WriteString(fmt.Sprintf("**Pages:** %d\n\n", len(pages)))
	builder.WriteString("---\n\n")

	// Full text
	builder.WriteString("## Full Text\n\n")
	builder.WriteString("```\n")
	builder.WriteString(fullText)
	builder.WriteString("\n```\n\n")

	// Pages
	if len(pages) > 1 {
		builder.WriteString("## Pages\n\n")
		for _, page := range pages {
			builder.WriteString(fmt.Sprintf("### Page %d\n\n", page.PageNumber))
			builder.WriteString("```\n")
			builder.WriteString(page.Text)
			builder.WriteString("\n```\n\n")
		}
	}

	return []byte(builder.String()), nil
}

// PageText represents text extracted from a single page
type PageText struct {
	PageNumber int
	Text       string
}

