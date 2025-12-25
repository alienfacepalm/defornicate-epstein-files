# Release Guide

## Creating a Release

### Automatic Release (Recommended)

1. Create and push a version tag:
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

2. GitHub Actions will automatically:
   - Build binaries for all platforms (Windows, Linux, macOS) into the `releases/` directory
   - Create a GitHub Release
   - Attach all binaries from `releases/` to the release

### Manual Release

You can also trigger the workflow manually from the GitHub Actions tab:
1. Go to Actions → Release
2. Click "Run workflow"
3. Select the branch and enter a tag name (e.g., `v1.0.0`)

## Supported Platforms

The release workflow builds binaries for:
- **Windows**: amd64, 386
- **Linux**: amd64, 386, arm64
- **macOS**: amd64 (Intel), arm64 (Apple Silicon)

## Versioning

Use semantic versioning (e.g., `v1.0.0`, `v1.1.0`, `v2.0.0`).

## Release Notes

When creating a release, include:
- New features
- Bug fixes
- Breaking changes (if any)
- Platform-specific notes

