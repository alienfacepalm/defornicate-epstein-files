# Project Structure

This document describes the project structure and organization.

## Directory Structure

```
.
├── main.go                 # Application entry point
├── go.mod                  # Go module definition
├── go.sum                  # Go module checksums
├── Makefile                # Build automation
├── README.md               # Project documentation (index)
├── docs/                   # Documentation directory
│   ├── CONTRIBUTING.md     # Contribution guidelines
│   ├── CODE_OF_CONDUCT.md  # Code of conduct
│   ├── CHANGELOG.md        # Version history
│   ├── RELEASE.md          # Release guide
│   └── PROJECT_STRUCTURE.md # This file
├── .github/
│   └── workflows/
│       └── release.yml     # GitHub Actions release workflow
├── internal/               # Internal packages (not importable)
│   ├── config/             # Configuration management
│   ├── downloader/         # Document downloading with checksum verification
│   ├── extractor/          # Document text extraction
│   ├── pattern/            # Sequential pattern expansion
│   └── pathutil/           # Path resolution utilities
├── documents/              # Document storage (gitignored)
│   ├── pdf/                # PDF files organized by filename
│   ├── docx/               # DOCX files (future)
│   ├── txt/                # TXT files (future)
│   └── ...                 # Other file types
├── bin/                    # Build output (gitignored)
└── releases/               # Release binaries (gitignored, used by GitHub Actions)
```

## Package Organization

### `internal/config`

Handles loading and parsing of configuration files (epstein-files-urls.json).

**Key Functions:**

- `Load(configPath string) (*Config, error)` - Load configuration from file
- `Config.GetInputs() []string` - Get inputs based on config priority

### `internal/downloader`

Manages document downloads from URLs with checksum verification. Supports multiple file types (PDF, DOC, DOCX, RTF, TXT, etc.).

**Key Functions:**

- `New(documentsDir string) *Downloader` - Create new downloader instance
- `Download(url string) (string, error)` - Download document with checksum check
- `GetFileType(filename string) string` - Determine file type from extension
- `GetDocumentsDir(fileType string) string` - Get directory path for file type

**Features:**

- HTTP client with timeout
- User-Agent header
- SHA256 checksum verification
- Automatic duplicate detection
- File type detection and organization

### `internal/extractor`

Handles document text extraction and saving. Currently supports PDF files, with plans to support other formats.

**Key Functions:**

- `New() *Extractor` - Create new extractor instance
- `ExtractText(filePath string) (string, error)` - Extract text from document
- `ExtractTextStructured(filePath string) ([]PageText, string, int, error)` - Extract with page information
- `SaveExtractedText(filePath, text string) (string, error)` - Save extracted text

### `internal/pattern`

Expands sequential patterns into lists of URLs/filenames.

**Key Functions:**

- `ExpandPattern(pattern string) ([]string, error)` - Expand pattern range

**Pattern Format:**

- `{start-end}` or `{start:end}` - Range expansion
- Example: `EFTA{00010724-00010730}.pdf` → 7 files

### `internal/pathutil`

Resolves document file paths, checking the documents directory for filenames. Supports multiple file types.

**Key Functions:**

- `ResolveDocumentPath(input string) string` - Resolve document path (generic)
- `ResolvePDFPath(input string) string` - Resolve PDF path (legacy alias)
- `GetFileType(filename string) string` - Determine file type from extension

**Path Resolution:**

- Checks `documents/{type}/{basename}/{filename}` first (new structure)
- Falls back to `documents/{type}/{filename}`
- Falls back to legacy `pdfs/{basename}/{filename}` for backward compatibility (PDFs only)
- Falls back to legacy `pdfs/{filename}` for backward compatibility (PDFs only)
- Falls back to relative path from current directory

## Directory Structure

Documents are organized by file type under the `documents/` parent directory:

- `documents/pdf/` - PDF files
- `documents/docx/` - DOCX files (when supported)
- `documents/doc/` - DOC files (when supported)
- `documents/txt/` - TXT files (when supported)
- `documents/rtf/` - RTF files (when supported)
- `documents/other/` - Other file types

Each document is stored in its own subdirectory:

- `documents/{type}/{filename}/{filename}.{ext}`
- `documents/{type}/{filename}/{filename}.extracted.json` (or .md/.txt)

## Design Principles

1. **Separation of Concerns**: Each package has a single, well-defined responsibility
2. **Internal Packages**: All packages are in `internal/` to prevent external imports
3. **Error Handling**: All functions return errors with context
4. **Testability**: Functions are designed to be easily testable
5. **Documentation**: All public functions have godoc comments
6. **Extensibility**: Directory structure supports multiple file types

## Testing

Tests are located alongside source files with `_test.go` suffix.

Run tests:

```bash
make test
# or
go test ./...
```

## Building

```bash
make build
# or
go build -o bin/epstein-files-defornicator main.go
```

## Adding New Features

1. Identify the appropriate package (or create new one in `internal/`)
2. Add functions with proper error handling
3. Write tests
4. Update documentation
5. Update [CHANGELOG.md](CHANGELOG.md) in the docs directory
