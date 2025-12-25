package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/ledongthuc/pdf"
)

type Config struct {
	PDFURL  string   `json:"pdf_url"`  // Single URL (for backward compatibility)
	PDFURLs []string `json:"pdf_urls"` // Multiple URLs
	// Sequential pattern: base URL with {start-end} or {start:end} for range
	// Example: "https://example.com/EFTA{10724-10730}.pdf" or "EFTA{10724:10730}.pdf"
	PDFPattern string `json:"pdf_pattern,omitempty"` // Pattern with {start-end} or {start:end}
}

func main() {
	var inputs []string

	// Try to load config file first
	config, err := loadConfig("config.json")
	if err == nil {
		// Check for pattern first (expands to multiple URLs)
		if config.PDFPattern != "" {
			inputs = expandPattern(config.PDFPattern)
			fmt.Fprintf(os.Stderr, "Using PDF pattern from config.json, expanded to %d URL(s)\n", len(inputs))
		} else if len(config.PDFURLs) > 0 {
			// Check for multiple URLs
			inputs = config.PDFURLs
			fmt.Fprintf(os.Stderr, "Using %d PDF URL(s) from config.json\n", len(inputs))
		} else if config.PDFURL != "" {
			// Fall back to single URL for backward compatibility
			inputs = []string{config.PDFURL}
			fmt.Fprintf(os.Stderr, "Using PDF URL from config.json: %s\n", config.PDFURL)
		}
	}

	// Fall back to command-line arguments if no config URLs
	if len(inputs) == 0 {
		if len(os.Args) >= 2 {
			inputs = os.Args[1:]
		} else {
			fmt.Fprintf(os.Stderr, "Usage: %s [pdf-file-path-or-url ...]\n", os.Args[0])
			fmt.Fprintf(os.Stderr, "  If no argument is provided, will use pdf_urls (or pdf_url) from config.json\n")
			fmt.Fprintf(os.Stderr, "  If config.json doesn't exist or has no URLs, argument(s) are required\n")
			fmt.Fprintf(os.Stderr, "\nExample: %s document.pdf\n", os.Args[0])
			fmt.Fprintf(os.Stderr, "Example: %s https://example.com/document.pdf\n", os.Args[0])
			fmt.Fprintf(os.Stderr, "Example: %s doc1.pdf doc2.pdf\n", os.Args[0])
			if err != nil {
				fmt.Fprintf(os.Stderr, "\nConfig file error: %v\n", err)
			}
			os.Exit(1)
		}
	}

	// Process each input
	for i, input := range inputs {
		if len(inputs) > 1 {
			fmt.Fprintf(os.Stderr, "\n--- Processing %d of %d ---\n", i+1, len(inputs))
		}

		var pdfPath string

		// Check if input is a URL
		if strings.HasPrefix(input, "http://") || strings.HasPrefix(input, "https://") {
			fmt.Fprintf(os.Stderr, "Downloading PDF from URL: %s\n", input)
			pdfPath, err = downloadPDF(input)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error downloading PDF: %v\n", err)
				if len(inputs) == 1 {
					os.Exit(1)
				}
				continue // Skip to next URL if multiple
			}
			fmt.Fprintf(os.Stderr, "PDF saved to: %s\n", pdfPath)
		} else {
			// Resolve local file path (handles filenames in pdfs directory)
			pdfPath = resolvePDFPath(input)
		}

		// Extract text from PDF
		text, err := extractTextFromPDF(pdfPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error extracting text: %v\n", err)
			if len(inputs) == 1 {
				os.Exit(1)
			}
			continue // Skip to next URL if multiple
		}

		// Save extracted text to file next to the PDF
		extractedFilePath, err := saveExtractedText(pdfPath, text)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error saving extracted text: %v\n", err)
			if len(inputs) == 1 {
				os.Exit(1)
			}
			continue // Skip to next URL if multiple
		}
		fmt.Fprintf(os.Stderr, "Extracted text saved to: %s\n", extractedFilePath)

		// Also output the extracted text to stdout
		if len(inputs) > 1 {
			fmt.Fprintf(os.Stderr, "--- Text from %s ---\n", pdfPath)
		}
		fmt.Print(text)
		if len(inputs) > 1 && i < len(inputs)-1 {
			fmt.Print("\n\n")
		}
	}
}

func downloadPDF(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status: %s", resp.Status)
	}

	// Create pdfs directory if it doesn't exist
	pdfsDir := "pdfs"
	if err := os.MkdirAll(pdfsDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create pdfs directory: %w", err)
	}

	// Extract filename from URL or generate one
	filename := extractFilenameFromURL(url)
	if filename == "" {
		filename = "downloaded.pdf"
	}

	// Ensure filename ends with .pdf
	if !strings.HasSuffix(strings.ToLower(filename), ".pdf") {
		filename += ".pdf"
	}

	filePath := filepath.Join(pdfsDir, filename)

	// Read the response body into memory to compute checksum
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// Compute checksum of the downloaded content
	downloadedHash := sha256.Sum256(bodyBytes)

	// Check if file already exists
	if _, err := os.Stat(filePath); err == nil {
		// File exists, compute its checksum
		existingHash, err := computeFileChecksum(filePath)
		if err != nil {
			// If we can't read the existing file, replace it
			fmt.Fprintf(os.Stderr, "Warning: failed to compute checksum of existing file, will replace: %v\n", err)
		} else {
			// Compare checksums
			if downloadedHash == existingHash {
				fmt.Fprintf(os.Stderr, "PDF already exists with same checksum, skipping download: %s\n", filePath)
				return filePath, nil
			}
			// Checksums don't match, will replace the file
			fmt.Fprintf(os.Stderr, "PDF exists but checksum differs, replacing: %s\n", filePath)
		}
	}

	// Create or replace the file
	outFile, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer outFile.Close()

	// Write the downloaded content to file
	_, err = outFile.Write(bodyBytes)
	if err != nil {
		os.Remove(filePath)
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	return filePath, nil
}

func computeFileChecksum(filePath string) ([32]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return [32]byte{}, err
	}
	defer file.Close()

	hasher := sha256.New()
	_, err = io.Copy(hasher, file)
	if err != nil {
		return [32]byte{}, err
	}

	var hash [32]byte
	copy(hash[:], hasher.Sum(nil))
	return hash, nil
}

func extractFilenameFromURL(url string) string {
	// Remove query parameters
	if idx := strings.Index(url, "?"); idx != -1 {
		url = url[:idx]
	}

	// Extract the last part of the URL path
	parts := strings.Split(url, "/")
	if len(parts) > 0 {
		filename := parts[len(parts)-1]
		// URL decode basic characters
		filename = strings.ReplaceAll(filename, "%20", "_")
		filename = strings.ReplaceAll(filename, "%2F", "_")
		filename = strings.ReplaceAll(filename, "%", "_")
		// Remove any invalid characters for filenames
		invalidChars := []string{"<", ">", ":", "\"", "|", "?", "*", "\\"}
		for _, char := range invalidChars {
			filename = strings.ReplaceAll(filename, char, "_")
		}
		return filename
	}
	return ""
}

func extractTextFromPDF(pdfPath string) (string, error) {
	// Check if file exists
	if _, err := os.Stat(pdfPath); os.IsNotExist(err) {
		return "", fmt.Errorf("file does not exist: %s", pdfPath)
	}

	// Read PDF file
	file, reader, err := pdf.Open(pdfPath)
	if err != nil {
		return "", fmt.Errorf("failed to open PDF: %w\nNote: This PDF may be encrypted or in an unsupported format", err)
	}
	defer file.Close()

	var textBuilder strings.Builder
	totalPages := reader.NumPage()

	if totalPages == 0 {
		return "", fmt.Errorf("PDF has no pages")
	}

	// Extract text from each page
	for i := 1; i <= totalPages; i++ {
		page := reader.Page(i)
		if page.V.IsNull() {
			fmt.Fprintf(os.Stderr, "Warning: page %d is null, skipping\n", i)
			continue
		}

		text, err := page.GetPlainText(nil)
		if err != nil {
			// Try to continue with other pages
			fmt.Fprintf(os.Stderr, "Warning: failed to extract text from page %d: %v\n", i, err)
			continue
		}

		// Add page separator for multi-page documents
		if i > 1 {
			textBuilder.WriteString("\n\n--- Page " + fmt.Sprintf("%d", i) + " ---\n\n")
		}

		if text != "" {
			textBuilder.WriteString(text)
		}
	}

	result := textBuilder.String()
	if result == "" {
		return "", fmt.Errorf("no text could be extracted from the PDF. The PDF may be encrypted, image-based, or in an unsupported format")
	}

	return result, nil
}

func expandPattern(pattern string) []string {
	// Pattern: {start-end} or {start:end}
	// Example: "EFTA{00010724-00010730}.pdf" or "https://example.com/EFTA{10724:10730}.pdf"
	// The padding is determined by the format in the pattern itself
	re := regexp.MustCompile(`\{(\d+)[-:](\d+)\}`)
	matches := re.FindStringSubmatch(pattern)
	
	if len(matches) != 3 {
		// No pattern found, return as single item
		return []string{pattern}
	}

	startStr := matches[1]
	endStr := matches[2]
	start, err1 := strconv.Atoi(startStr)
	end, err2 := strconv.Atoi(endStr)
	if err1 != nil || err2 != nil || start > end {
		return []string{pattern}
	}

	// Determine padding length from the longer of the two numbers (to preserve format)
	paddingLen := len(startStr)
	if len(endStr) > paddingLen {
		paddingLen = len(endStr)
	}

	var results []string
	for i := start; i <= end; i++ {
		// Format number with same padding as the pattern
		numStr := fmt.Sprintf("%0*d", paddingLen, i)
		// Replace pattern with the number
		expanded := re.ReplaceAllString(pattern, numStr)
		results = append(results, expanded)
	}

	return results
}

func saveExtractedText(pdfPath string, text string) (string, error) {
	// Get the directory and base filename of the PDF
	dir := filepath.Dir(pdfPath)
	baseName := filepath.Base(pdfPath)
	
	// Remove .pdf extension and add .extracted.txt
	extractedName := strings.TrimSuffix(baseName, ".pdf")
	extractedName = strings.TrimSuffix(extractedName, ".PDF")
	extractedName += ".extracted.txt"
	
	// Create the full path for the extracted text file
	extractedPath := filepath.Join(dir, extractedName)
	
	// Write the text to the file
	err := os.WriteFile(extractedPath, []byte(text), 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write extracted text file: %w", err)
	}
	
	return extractedPath, nil
}

func resolvePDFPath(input string) string {
	// If it's already an absolute path or contains directory separators, use as-is
	if filepath.IsAbs(input) || strings.Contains(input, string(filepath.Separator)) {
		return input
	}

	// If it's just a filename, check in pdfs directory first
	pdfsDir := "pdfs"
	pdfPath := filepath.Join(pdfsDir, input)
	if _, err := os.Stat(pdfPath); err == nil {
		return pdfPath
	}

	// Fall back to treating as relative path from current directory
	return input
}

func loadConfig(configPath string) (*Config, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

