package main

import (
	"fmt"
	"os"
	"strings"

	"defornicate-epstein-files/internal/config"
	"defornicate-epstein-files/internal/downloader"
	"defornicate-epstein-files/internal/extractor"
	"defornicate-epstein-files/internal/pathutil"
	"defornicate-epstein-files/internal/pattern"
)

const (
	configFile = "epstein-files-urls.json"
)

func main() {
	os.Exit(run())
}

// run is the main application logic, separated for testing
func run() int {
	var inputs []string

	// Try to load config file first
	cfg, err := config.Load(configFile)
	if err == nil {
		// Check for pattern first (expands to multiple URLs)
		patternStr := cfg.Pattern
		if patternStr == "" && cfg.PDFPattern != "" {
			patternStr = cfg.PDFPattern // Legacy support
		}
		if patternStr != "" {
			expanded, err := pattern.ExpandPattern(patternStr)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error expanding pattern: %v\n", err)
				return 1
			}
			inputs = expanded
			fmt.Fprintf(os.Stderr, "Using document pattern from epstein-files-urls.json, expanded to %d URL(s)\n", len(inputs))
		} else {
			inputs = cfg.GetInputs()
			if len(inputs) > 0 {
				fmt.Fprintf(os.Stderr, "Using %d document URL(s) from epstein-files-urls.json\n", len(inputs))
			}
		}
	}

	// Fall back to command-line arguments if no config URLs
	if len(inputs) == 0 {
		if len(os.Args) >= 2 {
			inputs = os.Args[1:]
		} else {
			printUsage(err)
			return 1
		}
	}

	// Initialize components
	dl := downloader.New(downloader.DefaultDocumentsDir)
	ext := extractor.New()

	// Process each input
	var hasErrors bool
	for i, input := range inputs {
		if len(inputs) > 1 {
			fmt.Fprintf(os.Stderr, "\n--- Processing %d of %d ---\n", i+1, len(inputs))
		}

		var filePath string
		var err error

		// Check if input is a URL
		if strings.HasPrefix(input, "http://") || strings.HasPrefix(input, "https://") {
			fmt.Fprintf(os.Stderr, "Downloading document from URL: %s\n", input)
			filePath, err = dl.Download(input)
			if err != nil {
				if err == downloader.ErrFileExists {
					fmt.Fprintf(os.Stderr, "Document already exists with same checksum, skipping download: %s\n", filePath)
				} else {
					fmt.Fprintf(os.Stderr, "Error downloading document: %v\n", err)
					if len(inputs) == 1 {
						return 1
					}
					hasErrors = true
					continue
				}
			} else {
				fmt.Fprintf(os.Stderr, "Document saved to: %s\n", filePath)
			}
		} else {
			// Resolve local file path (handles filenames in documents directory)
			filePath = pathutil.ResolveDocumentPath(input)
		}

		// Extract text from document
		text, err := ext.ExtractText(filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error extracting text: %v\n", err)
			if len(inputs) == 1 {
				return 1
			}
			hasErrors = true
			continue
		}

		// Save extracted text to file next to the document
		extractedFilePath, err := ext.SaveExtractedText(filePath, text)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error saving extracted text: %v\n", err)
			if len(inputs) == 1 {
				return 1
			}
			hasErrors = true
			continue
		}
		fmt.Fprintf(os.Stderr, "Extracted text saved to: %s\n", extractedFilePath)

		// Also output the extracted text to stdout
		if len(inputs) > 1 {
			fmt.Fprintf(os.Stderr, "--- Text from %s ---\n", filePath)
		}
		fmt.Print(text)
		if len(inputs) > 1 && i < len(inputs)-1 {
			fmt.Print("\n\n")
		}
	}

	if hasErrors {
		return 1
	}
	return 0
}

func printUsage(configErr error) {
	fmt.Fprintf(os.Stderr, "Usage: %s [document-file-path-or-url ...]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  If no argument is provided, will use urls (or url) from epstein-files-urls.json\n")
	fmt.Fprintf(os.Stderr, "  If epstein-files-urls.json doesn't exist or has no URLs, argument(s) are required\n")
	fmt.Fprintf(os.Stderr, "\nExample: %s document.pdf\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Example: %s https://example.com/document.pdf\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Example: %s doc1.pdf doc2.docx file.txt\n", os.Args[0])
	if configErr != nil {
		fmt.Fprintf(os.Stderr, "\nConfig file error: %v\n", configErr)
	}
}
