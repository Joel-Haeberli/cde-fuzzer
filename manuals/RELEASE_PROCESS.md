# Release Process

## GitHub Release Workflow

This project uses a GitHub Actions workflow to automatically build and release binaries when you create version tags.

## How to Create a Release

### 1. Tag the Main Branch

```bash
# Create a new version tag (e.g., v1.0.0)
git tag v1.0.0

# Push the tag to GitHub
git push origin v1.0.0
```

### 2. GitHub Actions Will Automatically:

1. **Build all binaries** for Linux, macOS, and Windows
2. **Create a GitHub Release** with the tag name
3. **Upload all binaries** as release assets

### 3. Released Binaries

The workflow builds and releases these binaries:

#### Linux (amd64):
- `cde-extractor-linux`
- `cde-extractor-server-linux`
- `generate-report-linux`
- `generate-diverse-report-linux`
- `derive-rules-linux`
- `generate-synthetic-linux`

#### macOS (amd64):
- `cde-extractor-mac`
- `cde-extractor-server-mac`
- `generate-report-mac`
- `generate-diverse-report-mac`
- `derive-rules-mac`
- `generate-synthetic-mac`

#### Windows (amd64):
- `cde-extractor-windows.exe`
- `cde-extractor-server-windows.exe`
- `generate-report-windows.exe`
- `generate-diverse-report-windows.exe`
- `derive-rules-windows.exe`
- `generate-synthetic-windows.exe`

## Version Tagging Convention

Use semantic versioning for tags:
- `v1.0.0` - Major release
- `v1.0.1` - Patch release
- `v1.1.0` - Minor release
- `v2.0.0` - Major breaking changes

## Manual Release Process (Alternative)

If you prefer to create releases manually:

### 1. Build Binaries Locally

```bash
make build
```

This will create all binaries in the `bin/` directory.

### 2. Create GitHub Release Manually

1. Go to GitHub repository → Releases
2. Click "Draft a new release"
3. Enter tag version (e.g., `v1.0.0`)
4. Add release title and description
5. Upload binaries from `bin/` directory
6. Publish release

## Workflow File

The release workflow is defined in: `.github/workflows/release.yml`

### Key Features:
- **Automatic triggering** on version tags (`v*` pattern)
- **Cross-platform builds** for Linux, macOS, Windows
- **Automatic release creation** with proper naming
- **Asset uploading** with correct content types

### Customization

To modify the workflow:
1. Edit `.github/workflows/release.yml`
2. Update binary names or build commands as needed
3. Adjust Go version if required

## Troubleshooting

### Workflow Not Triggering
- Ensure tag follows `v*` pattern (e.g., `v1.0.0`)
- Verify tag was pushed to GitHub: `git push origin v1.0.0`
- Check GitHub Actions tab for workflow status

### Build Failures
- Check Go version compatibility
- Verify all dependencies are in `go.mod`
- Test local build first: `make build`

### Release Not Created
- Check GitHub token permissions
- Verify workflow has write access to repository
- Review GitHub Actions logs for errors

## Best Practices

1. **Test before releasing**: Run `make build` locally first
2. **Update CHANGELOG**: Document changes for each release
3. **Tag from main branch**: Ensure you're on the correct branch
4. **Use semantic versioning**: Follow `MAJOR.MINOR.PATCH` convention
5. **Include release notes**: Describe new features and fixes

## Example Release Process

```bash
# Commit all changes
git add .
git commit -m "Prepare for v1.0.0 release"

# Create annotated tag
git tag -a v1.0.0 -m "Version 1.0.0"

# Push tag to trigger workflow
git push origin v1.0.0

# Monitor workflow in GitHub Actions
# Download released binaries from GitHub Releases page
```

## Automatic Release Benefits

✅ **Consistent builds** - Same environment every time
✅ **Cross-platform** - All platforms built automatically
✅ **Version history** - Clear release tracking
✅ **Asset management** - Easy binary distribution
✅ **Time saving** - No manual uploads required