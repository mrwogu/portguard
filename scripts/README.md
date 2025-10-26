# Release Scripts

Helper scripts for managing releases and testing the release process.

## Scripts

### `create-release.sh`

Creates a new release tag with validation and pushes it to trigger the automated release workflow.

**Usage:**
```bash
./scripts/create-release.sh v1.2.0
```

**What it does:**
1. ✅ Validates version format (semantic versioning)
2. ✅ Checks you're on the main branch (warns if not)
3. ✅ Ensures no uncommitted changes
4. ✅ Pulls latest changes
5. ✅ Runs tests
6. ✅ Runs linter
7. ✅ Checks CHANGELOG.md for version entry (warns if missing)
8. ✅ Creates annotated tag
9. ✅ Pushes tag to origin (triggers GitHub Actions)

**Version formats:**
- `v1.2.3` - Stable release
- `v1.2.3-rc.1` - Release candidate
- `v1.2.3-beta.1` - Beta release
- `v1.2.3-alpha.1` - Alpha release

### `test-release.sh`

Simulates the GitHub Actions release workflow locally. Useful for testing before pushing a tag.

**Usage:**
```bash
# Test with specific version
./scripts/test-release.sh v1.2.0

# Test with current git version
./scripts/test-release.sh
```

**What it does:**
1. ✅ Runs tests
2. ✅ Runs linter
3. ✅ Builds for all platforms
4. ✅ Packages binaries (tar.gz/zip)
5. ✅ Generates SHA256 checksums
6. ✅ Verifies checksums
7. ✅ Tests binary for current platform

**Output:**
All artifacts are created in `dist/` directory, same as the actual release.

## Workflow

### Typical Release Process

```bash
# 1. Update CHANGELOG.md with new version
vim CHANGELOG.md

# 2. Commit changes
git add CHANGELOG.md
git commit -m "chore: prepare release v1.2.0"
git push

# 3. Test release build locally
./scripts/test-release.sh v1.2.0

# 4. Create and push release tag
./scripts/create-release.sh v1.2.0

# 5. Monitor GitHub Actions
# Visit: https://github.com/mrwogu/portguard/actions
```

### Emergency Rollback

If something goes wrong:

```bash
# Delete remote tag
git push --delete origin v1.2.0

# Delete local tag
git tag -d v1.2.0

# Delete GitHub release (via web UI or GitHub CLI)
gh release delete v1.2.0
```

## Requirements

Both scripts require:
- Git
- Go 1.23+
- Make
- Standard Unix tools (tar, zip, shasum)

The `create-release.sh` script also requires:
- Clean working directory (no uncommitted changes)
- Tests and linter passing
- Access to push tags to origin

## Troubleshooting

### "Tests failed"
```bash
# Run tests to see details
make test

# Fix issues and try again
```

### "Linter failed"
```bash
# Run linter to see details
make lint

# Fix issues and try again
```

### "Tag already exists"
```bash
# Delete local tag
git tag -d v1.2.0

# Delete remote tag if pushed
git push --delete origin v1.2.0

# Try again
```

### "Version not found in CHANGELOG.md"
The script will warn but allow you to continue. It's recommended to add a CHANGELOG entry before releasing.

## See Also

- [Release Process Documentation](../docs/RELEASING.md)
- [GitHub Workflows README](../.github/workflows/README.md)
- [Makefile](../Makefile)
