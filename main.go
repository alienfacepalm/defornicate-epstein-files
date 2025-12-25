package main

import (
	"fmt"
	"os"
	"path/filepath"
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

// findConfigFile searches for the config file in multiple locations:
// 1. Current working directory
// 2. Directory where the executable is located
// 3. Parent directory of the executable (for bin/ structure)
func findConfigFile() string {
	// Try current working directory first
	if _, err := os.Stat(configFile); err == nil {
		return configFile
	}

	// Get the executable's directory
	execPath, err := os.Executable()
	if err == nil {
		execDir := filepath.Dir(execPath)
		
		// Try in executable's directory
		configPath := filepath.Join(execDir, configFile)
		if _, err := os.Stat(configPath); err == nil {
			return configPath
		}

		// Try parent directory (for bin/ structure)
		parentDir := filepath.Dir(execDir)
		configPath = filepath.Join(parentDir, configFile)
		if _, err := os.Stat(configPath); err == nil {
			return configPath
		}
	}

	// Fall back to current directory (will fail gracefully if not found)
	return configFile
}

func main() {
	os.Exit(run())
}

// run is the main application logic, separated for testing
func run() int {
	var inputs []string

	// Try to load config file first
	configPath := findConfigFile()
	cfg, err := config.Load(configPath)
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
	var successCount, errorCount int
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
					errorCount++
					continue
				}
			} else {
				fmt.Fprintf(os.Stderr, "Document saved to: %s\n", filePath)
				successCount++
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
			errorCount++
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
			errorCount++
			continue
		}
		fmt.Fprintf(os.Stderr, "Extracted text saved to: %s\n", extractedFilePath)
		successCount++

		// Also output the extracted text to stdout
		if len(inputs) > 1 {
			fmt.Fprintf(os.Stderr, "--- Text from %s ---\n", filePath)
		}
		fmt.Print(text)
		if len(inputs) > 1 && i < len(inputs)-1 {
			fmt.Print("\n\n")
		}
	}

	// Print summary if processing multiple files
	if len(inputs) > 1 {
		fmt.Fprintf(os.Stderr, "\n--- Summary ---\n")
		fmt.Fprintf(os.Stderr, "Total processed: %d\n", len(inputs))
		fmt.Fprintf(os.Stderr, "Successful: %d\n", successCount)
		if errorCount > 0 {
			fmt.Fprintf(os.Stderr, "Errors: %d\n", errorCount)
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
