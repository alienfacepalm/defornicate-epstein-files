// Package pattern provides sequential pattern expansion for batch processing.
// Patterns support ranges like {start-end} or {start:end} to generate multiple URLs/filenames.
package pattern

import (
	"fmt"
	"regexp"
	"strconv"
)

// ExpandPattern expands a sequential pattern into a list of URLs/filenames
// Pattern format: {start-end} or {start:end}
// Example: "EFTA{00010724-00010730}.pdf" expands to EFTA00010724.pdf through EFTA00010730.pdf
func ExpandPattern(pattern string) ([]string, error) {
	// Pattern: {start-end} or {start:end}
	re := regexp.MustCompile(`\{(\d+)[-:](\d+)\}`)
	matches := re.FindStringSubmatch(pattern)

	if len(matches) != 3 {
		// No pattern found, return as single item
		return []string{pattern}, nil
	}

	startStr := matches[1]
	endStr := matches[2]
	start, err1 := strconv.Atoi(startStr)
	end, err2 := strconv.Atoi(endStr)
	if err1 != nil || err2 != nil {
		return nil, fmt.Errorf("invalid pattern numbers: %w", err1)
	}
	if start > end {
		return nil, fmt.Errorf("start number (%d) must be <= end number (%d)", start, end)
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

	return results, nil
}

