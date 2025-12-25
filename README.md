# PDF Text Extractor

A command-line tool written in Go that extracts all text from PDF files. It supports both local files and remote URLs.

## Features

- Extract text from local PDF files
- Download and extract text from PDFs via URL
- **Checksum verification** - Skips re-downloading identical files
- **Sequential pattern support** - Download multiple PDFs using pattern ranges
- **Automatic text file saving** - Saves extracted text to `.extracted.txt` files
- Simple command-line interface
- Handles multi-page documents
- Support for multiple PDFs via config file

## Installation

### Download Pre-built Binaries

Download the latest release from the [Releases page](https://github.com/alienfacepalm/defornicate-epstein-files/releases). Binaries are available for:
- Windows (amd64, 386)
- Linux (amd64, 386, arm64)
- macOS (amd64, arm64)

### Build from Source

1. Make sure you have Go 1.24 or later installed
2. Clone or download this project
3. Install dependencies:
   ```bash
   go mod download
   ```
4. Build:
   ```bash
   go build -o pdf-extractor main.go
   ```

## Usage

### Using config.json

Create a `config.json` file with PDF URLs or patterns:

```json
{
  "pdf_pattern": "https://www.justice.gov/epstein/files/DataSet%208/EFTA{00010724-00010730}.pdf"
}
```

Or use multiple URLs:

```json
{
  "pdf_urls": [
    "https://example.com/file1.pdf",
    "file2.pdf"
  ]
}
```

Then run:
```bash
./pdf-extractor
```

### Command-line Usage

#### Extract text from a local PDF file:

```bash
./pdf-extractor document.pdf
```

#### Extract text from a PDF URL:

```bash
./pdf-extractor https://www.justice.gov/epstein/files/DataSet%208/EFTA00010724.pdf
```

#### Extract from multiple files:

```bash
./pdf-extractor file1.pdf file2.pdf https://example.com/file3.pdf
```

### Sequential Patterns

Use pattern ranges in `config.json` to download multiple sequential PDFs:

```json
{
  "pdf_pattern": "https://example.com/EFTA{00010724-00010730}.pdf"
}
```

This will download EFTA00010724.pdf through EFTA00010730.pdf (7 files total).

### Output

- Extracted text is saved to `[filename].extracted.txt` next to each PDF
- Text is also output to stdout for piping/redirection

## Example

To extract text from the specified PDF:

```bash
go run main.go https://www.justice.gov/epstein/files/DataSet%208/EFTA00010724.pdf > extracted_text.txt
```

## Dependencies

- `github.com/ledongthuc/pdf` - PDF text extraction library

## Notes

- **Checksum verification**: If a PDF already exists, the tool checks if it's identical before re-downloading
- **Automatic text saving**: Extracted text is automatically saved to `.extracted.txt` files in the same directory as the PDFs
- **PDF storage**: Downloaded PDFs are saved to the `pdfs/` directory
- Multi-page documents include page separators in the output
- The tool outputs all extracted text to stdout in addition to saving to files
