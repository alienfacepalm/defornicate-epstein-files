# Contributing to Epstein Files Defornicator

Thank you for your interest in contributing! This document provides guidelines and instructions for contributing to this project.

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/your-username/defornicate-epstein-files.git`
3. Create a branch: `git checkout -b feature/your-feature-name`
4. Make your changes
5. Test your changes: `make test`
6. Commit your changes: `git commit -m "Add your feature"`
7. Push to your fork: `git push origin feature/your-feature-name`
8. Open a Pull Request

## Development Setup

### Prerequisites

- Go 1.25 or later
- Make (optional, for using Makefile)

### Building

```bash
# Build the binary
make build

# Or use go directly
go build -o bin/epstein-files-defornicator main.go
```

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage
```

### Code Style

- Follow Go standard formatting: `make fmt`
- Run linter: `make lint` (requires golangci-lint)
- Run go vet: `make vet`

## Code Organization

The project is organized into internal packages:

- `internal/config` - Configuration loading
- `internal/downloader` - Document downloading with checksum verification (supports multiple file types)
- `internal/extractor` - Document text extraction (currently PDF, extensible to other formats)
- `internal/pattern` - Sequential pattern expansion
- `internal/pathutil` - Path resolution utilities (supports multiple file types)

## Commit Messages

- Use clear, descriptive commit messages
- Reference issue numbers if applicable
- Follow conventional commit format when possible

## Pull Request Guidelines

- Keep PRs focused on a single feature or bug fix
- Include tests for new functionality
- Update documentation as needed
- Ensure all tests pass
- Update [CHANGELOG.md](CHANGELOG.md) in the docs directory if applicable

## Reporting Issues

When reporting issues, please include:

- Description of the issue
- Steps to reproduce
- Expected behavior
- Actual behavior
- Environment (OS, Go version, etc.)
- Relevant logs or error messages

## Questions?

Feel free to open an issue for questions or discussions.
