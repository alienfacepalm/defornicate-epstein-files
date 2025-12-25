// Package config provides configuration loading and management.
package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config represents the application configuration.
// It supports multiple ways to specify document sources:
// - URL: Single URL (for backward compatibility)
// - URLs: Multiple URLs
// - Pattern: Sequential pattern with {start-end} or {start:end}
type Config struct {
	URL     string   `json:"url"`      // Single URL (for backward compatibility, also accepts pdf_url)
	URLs    []string `json:"urls"`     // Multiple URLs (also accepts pdf_urls)
	Pattern string   `json:"pattern"`  // Pattern with {start-end} or {start:end} (also accepts pdf_pattern)
	// Legacy fields for backward compatibility
	PDFURL    string   `json:"pdf_url,omitempty"`
	PDFURLs   []string `json:"pdf_urls,omitempty"`
	PDFPattern string  `json:"pdf_pattern,omitempty"`
}

// Load reads and parses the configuration file
func Load(configPath string) (*Config, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	var cfg Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &cfg, nil
}

// GetInputs returns the list of inputs to process based on config priority
func (c *Config) GetInputs() []string {
	// Normalize legacy fields to new fields
	if c.Pattern == "" && c.PDFPattern != "" {
		c.Pattern = c.PDFPattern
	}
	if len(c.URLs) == 0 && len(c.PDFURLs) > 0 {
		c.URLs = c.PDFURLs
	}
	if c.URL == "" && c.PDFURL != "" {
		c.URL = c.PDFURL
	}

	// Priority: Pattern > URLs > URL
	if c.Pattern != "" {
		return []string{} // Will be expanded by pattern expander
	}
	if len(c.URLs) > 0 {
		return c.URLs
	}
	if c.URL != "" {
		return []string{c.URL}
	}
	return []string{}
}

