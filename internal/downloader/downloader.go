// Package downloader provides document downloading functionality with checksum verification
// to avoid duplicate downloads.
package downloader

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	// DefaultTimeout is the default HTTP client timeout
	DefaultTimeout = 30 * time.Second
	// DefaultUserAgent is the user agent string for HTTP requests
	DefaultUserAgent = "epstein-files-defornicator/1.0"
	// DefaultDocumentsDir is the default parent directory for storing documents
	DefaultDocumentsDir = "documents"
	// DefaultFilePerm is the default file permission (0644)
	DefaultFilePerm = 0644
	// DefaultDirPerm is the default directory permission (0755)
	DefaultDirPerm = 0755
)

// Downloader handles document downloads with checksum verification
type Downloader struct {
	client    *http.Client
	documentsDir string
	userAgent string
}

// New creates a new Downloader instance
func New(documentsDir string) *Downloader {
	return &Downloader{
		client: &http.Client{
			Timeout: DefaultTimeout,
		},
		documentsDir: documentsDir,
		userAgent:    DefaultUserAgent,
	}
}

// Download downloads a document from a URL, checking checksums to avoid duplicates
func (d *Downloader) Download(url string) (string, error) {
	// Create request with user agent
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", d.userAgent)

	// Make request
	resp, err := d.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status: %s", resp.Status)
	}

	// Extract filename from URL or generate one
	filename := extractFilenameFromURL(url)
	if filename == "" {
		filename = "downloaded"
	}

	// Determine file type from filename
	fileType := GetFileType(filename)
	if filename == "downloaded" {
		// If no extension detected, default to pdf for now
		fileType = "pdf"
		filename = "downloaded.pdf"
	}

	// Get the documents directory for this file type
	typeDir := GetDocumentsDir(fileType)
	
	// Create documents directory structure if it doesn't exist
	if err := os.MkdirAll(typeDir, DefaultDirPerm); err != nil {
		return "", fmt.Errorf("failed to create documents directory: %w", err)
	}

	// Get base name without extension for subdirectory
	ext := filepath.Ext(filename)
	baseName := strings.TrimSuffix(filename, ext)
	baseName = strings.TrimSuffix(baseName, strings.ToUpper(ext))
	
	// Create subdirectory for this document
	docSubDir := filepath.Join(typeDir, baseName)
	if err := os.MkdirAll(docSubDir, DefaultDirPerm); err != nil {
		return "", fmt.Errorf("failed to create document subdirectory: %w", err)
	}

	// Store document in its own subdirectory
	filePath := filepath.Join(docSubDir, filename)

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
			return filePath, fmt.Errorf("warning: failed to compute checksum of existing file, will replace: %w", err)
		}

		// Compare checksums
		if downloadedHash == existingHash {
			return filePath, ErrFileExists
		}
		// Checksums don't match, will replace the file
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

// computeFileChecksum calculates the SHA256 checksum of a file
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

// extractFilenameFromURL extracts a safe filename from a URL
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

// ErrFileExists is returned when a file with the same checksum already exists
var ErrFileExists = fmt.Errorf("file already exists with same checksum")

