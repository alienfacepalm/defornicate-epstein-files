# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed
- Refactored codebase to be file-type agnostic (removed PDF-specific naming)
- Updated configuration to use generic `url`, `urls`, and `pattern` fields (legacy `pdf_url`, `pdf_urls`, `pdf_pattern` still supported)
- Document storage organized by file type: `documents/{type}/{filename}/`
- Updated all documentation to reflect multi-format support

### Added
- File type detection and organization system
- Support for multiple document formats (structure ready for DOC, DOCX, RTF, TXT, etc.)
- Generic path resolution utilities for all file types

## [0.0.1] - 2025-12-24

### Added
- Initial release
- PDF text extraction from local files and URLs
- Checksum verification to avoid duplicate downloads
- Sequential pattern support for batch processing
- Automatic extraction text file saving (`.extracted.txt`, `.extracted.json`, `.extracted.md`)
- Support for multiple documents via epstein-files-urls.json
- Multi-platform binary releases (Windows, Linux, macOS)
- GitHub Actions workflow for automated releases

### Features
- Download documents from URLs with automatic checksum checking
- Extract text from PDFs (local or remote)
- Save extracted text to files next to documents in structured formats
- Support for sequential filename patterns (e.g., `EFTA{00010724-00010730}.pdf`)
- Config file support for batch processing
- Command-line interface for single or multiple files

[Unreleased]: https://github.com/alienfacepalm/defornicate-epstein-files/compare/v0.0.1...HEAD
[0.0.1]: https://github.com/alienfacepalm/defornicate-epstein-files/releases/tag/v0.0.1

