# Epstein Files Defornicator

A command-line tool written in Go that extracts all text from document files. It supports multiple file formats (PDF, DOC, DOCX, RTF, TXT, etc.) and both local files and remote URLs.

> **Note**: This README serves as the main documentation index. For detailed documentation, see the [Documentation](#documentation) section below.

## Features

- Extract text from local document files (PDF, DOC, DOCX, RTF, TXT, and more)
- Download and extract text from documents via URL
- **Checksum verification** - Skips re-downloading identical files
- **Sequential pattern support** - Download multiple documents using pattern ranges
- **Automatic text file saving** - Saves extracted text in structured formats (JSON, Markdown, or plain text)
- **Multi-format support** - Organized by file type in `documents/{type}/` directories
- Simple command-line interface
- Handles multi-page documents
- Support for multiple documents via config file

## Installation

### Download Pre-built Binaries

Download the latest release from the [Releases page](https://github.com/alienfacepalm/defornicate-epstein-files/releases). Binaries are available for:

- Windows (amd64, 386)
- Linux (amd64, 386, arm64)
- macOS (amd64, arm64)

### Build from Source

1. Make sure you have Go 1.25 or later installed
2. Clone or download this project
3. Install dependencies:
   ```bash
   go mod download
   ```
4. Build:

   ```bash
   # Using Makefile (recommended)
   make build

   # Or using go directly
   go build -o bin/epstein-files-defornicator main.go
   ```

### Development

See [docs/CONTRIBUTING.md](docs/CONTRIBUTING.md) for development setup and guidelines.

## Documentation

- **[Contributing Guide](docs/CONTRIBUTING.md)** - How to contribute to this project
- **[Code of Conduct](docs/CODE_OF_CONDUCT.md)** - Community guidelines
- **[Project Structure](docs/PROJECT_STRUCTURE.md)** - Code organization and architecture
- **[Release Guide](docs/RELEASE.md)** - How to create releases
- **[Changelog](docs/CHANGELOG.md)** - Version history and changes

## Usage

### Using epstein-files-urls.json

Create an `epstein-files-urls.json` file with document URLs or patterns:

```json
{
  "pattern": "https://www.justice.gov/epstein/files/DataSet%208/EFTA{00010724-00010730}.pdf"
}
```

Or use multiple URLs:

```json
{
  "urls": ["https://example.com/file1.pdf", "file2.docx", "document.txt"]
}
```

Then run:

```bash
./epstein-files-defornicator
```

### Command-line Usage

#### Extract text from a local document file:

```bash
./epstein-files-defornicator document.pdf
./epstein-files-defornicator document.docx
./epstein-files-defornicator document.txt
```

#### Extract text from a document URL:

```bash
./epstein-files-defornicator https://www.justice.gov/epstein/files/DataSet%208/EFTA00010724.pdf
```

#### Extract from multiple files:

```bash
./epstein-files-defornicator file1.pdf file2.docx file3.txt https://example.com/file4.rtf
```

### Sequential Patterns

Use pattern ranges in `epstein-files-urls.json` to download multiple sequential documents:

```json
{
  "pattern": "https://example.com/EFTA{00010724-00010730}.pdf"
}
```

This will download EFTA00010724.pdf through EFTA00010730.pdf (7 files total).

### Output Formats

Extracted text is saved in structured formats next to each document:

- **JSON** (default): `[filename].extracted.json` - Structured format with metadata and page-by-page content
- **Markdown**: `[filename].extracted.md` - Human-readable Markdown format
- **Plain Text**: `[filename].extracted.txt` - Simple text format

The JSON format includes:

- Metadata (filename, extraction date, page count)
- Full text
- Page-by-page breakdown with word counts

Text is also output to stdout for piping/redirection (always in plain text format).

## Example

To extract text from a document:

```bash
go run main.go https://www.justice.gov/epstein/files/DataSet%208/EFTA00010724.pdf > extracted_text.txt
```

## Dependencies

- `github.com/ledongthuc/pdf` - PDF text extraction library (for PDF support)

## Supported File Types

Currently supported:

- **PDF** (.pdf) - Full support

Planned support:

- **Word Documents** (.doc, .docx)
- **Rich Text Format** (.rtf)
- **Text Files** (.txt)
- **OpenDocument Text** (.odt)

## Notes

- **Checksum verification**: If a document already exists, the tool checks if it's identical before re-downloading
- **Automatic text saving**: Extracted text is automatically saved in structured formats (JSON by default) in the same directory as the document
- **Document storage**: Each document is stored in its own subdirectory organized by file type: `documents/{type}/{filename}/{filename}.{ext}` and `documents/{type}/{filename}/{filename}.extracted.json` (or .md/.txt)
- **Structured output**: JSON format includes metadata, full text, and page-by-page breakdown with word counts
- Multi-page documents include page separators in the output
- The tool outputs all extracted text to stdout in addition to saving to files
