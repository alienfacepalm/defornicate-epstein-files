# PDF Text Extractor

A command-line tool written in Go that extracts all text from PDF files. It supports both local files and remote URLs.

## Features

- Extract text from local PDF files
- Download and extract text from PDFs via URL
- Simple command-line interface
- Handles multi-page documents

## Installation

1. Make sure you have Go 1.21 or later installed
2. Clone or download this project
3. Install dependencies:
   ```bash
   go mod download
   ```

## Usage

### Extract text from a local PDF file:

```bash
go run main.go document.pdf
```

### Extract text from a PDF URL:

```bash
go run main.go https://www.justice.gov/epstein/files/DataSet%208/EFTA00010724.pdf
```

### Build and run:

```bash
# Build the executable
go build -o pdf-extractor main.go

# Run with local file
./pdf-extractor document.pdf

# Run with URL
./pdf-extractor https://example.com/document.pdf
```

### Save output to a file:

```bash
go run main.go document.pdf > output.txt
```

## Example

To extract text from the specified PDF:

```bash
go run main.go https://www.justice.gov/epstein/files/DataSet%208/EFTA00010724.pdf > extracted_text.txt
```

## Dependencies

- `github.com/ledongthuc/pdf` - PDF text extraction library

## Notes

- When using a URL, the PDF is temporarily downloaded to extract text, then automatically deleted
- The tool outputs all extracted text to stdout
- Multi-page documents include page separators in the output
