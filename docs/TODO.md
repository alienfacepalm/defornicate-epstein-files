# TODO - Implementation Roadmap

This document tracks features and improvements to be implemented.

## High Priority

### Auto-Discovery of Documents

- [ ] **Remove config JSON requirement for manual URL specification**

  - Currently requires `epstein-files-urls.json` with manual URLs/patterns
  - Should automatically discover documents from source sites

- [ ] **Implement auto-discovery on justice.gov**

  - Crawl/scrape the Epstein files section: `https://www.justice.gov/epstein/files/`
  - Automatically detect and download all available documents
  - Handle different dataset sections (DataSet 1, DataSet 2, etc.)
  - Parse HTML to extract document links
  - Handle pagination if documents are spread across multiple pages

- [ ] **Implement auto-discovery on jmail.world**

  - Crawl/scrape the jmail.world site for Epstein-related documents
  - Automatically detect document links
  - Handle site structure and navigation

- [ ] **Smart discovery mode**
  - Command-line flag: `--auto-discover` or `--discover`
  - Option to specify source: `--source=justice.gov` or `--source=jmail.world`
  - Option to discover all sources: `--discover-all`

## Medium Priority

### Additional File Format Support

- [ ] **DOC/DOCX support**

  - Implement text extraction for Microsoft Word documents
  - Use appropriate library (e.g., `github.com/unidoc/unioffice` or similar)

- [ ] **RTF support**

  - Implement Rich Text Format text extraction

- [ ] **TXT support**

  - Simple text file reading (already straightforward)

- [ ] **ODT support**
  - OpenDocument Text format extraction

### Enhanced Discovery Features

- [ ] **Incremental discovery**

  - Track which documents have already been downloaded
  - Only download new/updated documents
  - Maintain a local index/database of discovered documents

- [ ] **Discovery filters**

  - Filter by date range
  - Filter by document type
  - Filter by keywords in filenames
  - Filter by file size

- [ ] **Parallel downloads**
  - Download multiple documents concurrently
  - Configurable concurrency limit
  - Rate limiting to avoid overwhelming servers

### Configuration Improvements

- [ ] **Discovery configuration**

  - Config file for discovery settings (sources, filters, etc.)
  - Default discovery sources
  - Discovery schedule/automation

- [ ] **Output format selection**
  - Command-line flag to choose output format: `--format=json|markdown|plain`
  - Per-document format selection
  - Default format configuration

## Low Priority

### User Experience

- [ ] **Progress indicators**

  - Better progress bars for downloads
  - ETA for large batches
  - Download speed indicators

- [ ] **Verbose/debug mode**

  - `--verbose` flag for detailed logging
  - `--debug` flag for developer debugging
  - Log file output option

- [ ] **Resume interrupted downloads**
  - Save partial downloads
  - Resume from where it left off
  - Handle network interruptions gracefully

### Data Management

- [ ] **Database/index for documents**

  - SQLite database to track downloaded documents
  - Metadata storage (download date, checksum, source, etc.)
  - Query/search capabilities
  - Deduplication across sources

- [ ] **Document metadata extraction**

  - Extract PDF metadata (author, creation date, etc.)
  - Extract email metadata (from, to, date, subject)
  - Store in structured format

- [ ] **Full-text search**
  - Index extracted text for fast searching
  - Search across all documents
  - Search by keywords, phrases, dates

### Advanced Features

- [ ] **Email parsing**

  - Parse email content from documents
  - Extract attachments
  - Thread email conversations
  - Export emails in standard formats (mbox, eml)

- [ ] **Content analysis**

  - Named entity recognition (people, organizations, dates)
  - Topic modeling
  - Sentiment analysis
  - Relationship graph building

- [ ] **Export formats**
  - Export to CSV
  - Export to database
  - Export to Elasticsearch
  - Custom export formats

### Infrastructure

- [ ] **Docker support**

  - Dockerfile for containerized deployment
  - Docker Compose for full stack

- [ ] **API/Server mode**

  - REST API for programmatic access
  - Web interface for browsing documents
  - GraphQL API option

- [ ] **Scheduled discovery**
  - Cron-like scheduling
  - Automatic periodic discovery
  - Notifications for new documents

## Notes

- Auto-discovery should respect robots.txt
- Implement rate limiting to be a good citizen
- Cache discovery results to avoid repeated requests
- Handle site structure changes gracefully
- Provide fallback mechanisms if auto-discovery fails
