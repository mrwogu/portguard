# Release Process

This document describes how to create a new release of PortGuard.

## Overview

Releases are fully automated through GitHub Actions. When you push a version tag, the CI/CD pipeline:

1. Runs all tests and linters
2. Builds binaries for all platforms
3. Packages them with checksums
4. Creates a GitHub Release with all artifacts
5. Builds and publishes Docker images

## Prerequisites

Before creating a release:

- [ ] All tests pass locally: `make test`
- [ ] Code is linted: `make lint`
- [ ] Documentation is up to date
- [ ] CHANGELOG.md is updated with release notes
- [ ] Version follows [Semantic Versioning](https://semver.org/)

## Release Steps

### 1. Update CHANGELOG.md

Add a new version section at the top:

```markdown
## [1.2.0] - 2024-10-26

### Added
- New feature X
- New feature Y

### Changed
- Improved Z

### Fixed
- Bug fix for issue #123
```

### 2. Test Release Build Locally

```bash
# Test the complete release build process
./scripts/test-release.sh v1.2.0

# Or test with current version
./scripts/test-release.sh
```

This simulates the GitHub Actions workflow and creates artifacts in `dist/`.

### 3. Commit Changes

```bash
git add CHANGELOG.md
git commit -m "chore: prepare release v1.2.0"
git push origin main
```

### 4. Create and Push Tag

```bash
# Create annotated tag
git tag -a v1.2.0 -m "Release v1.2.0"

# Push tag to trigger release workflow
git push origin v1.2.0
```

### 5. Monitor Release

1. Go to [GitHub Actions](https://github.com/mrwogu/portguard/actions)
2. Watch the "Release" workflow
3. Check for any errors

### 6. Verify Release

Once complete:

1. Visit [Releases page](https://github.com/mrwogu/portguard/releases)
2. Verify all artifacts are present:
   - `portguard-linux-amd64.tar.gz` + `.sha256`
   - `portguard-linux-arm64.tar.gz` + `.sha256`
   - `portguard-darwin-amd64.tar.gz` + `.sha256`
   - `portguard-darwin-arm64.tar.gz` + `.sha256`
   - `portguard-windows-amd64.zip` + `.sha256`
3. Check [Docker images](https://github.com/mrwogu/portguard/pkgs/container/portguard)
4. Test download and run:

```bash
# Download and test
wget https://github.com/mrwogu/portguard/releases/download/v1.2.0/portguard-linux-amd64.tar.gz
tar -xzf portguard-linux-amd64.tar.gz
./portguard-linux-amd64 --version
```

## Version Conventions

### Stable Releases

```bash
v1.0.0    # Major version
v1.1.0    # Minor version (new features)
v1.1.1    # Patch version (bug fixes)
```

### Pre-releases

```bash
v1.0.0-rc.1     # Release candidate
v1.0.0-beta.1   # Beta version
v1.0.0-alpha.1  # Alpha version
```

Pre-release versions are automatically marked as "Pre-release" on GitHub.

## Rollback

If you need to remove a bad release:

```bash
# Delete remote tag
git push --delete origin v1.2.0

# Delete local tag
git tag -d v1.2.0

# Delete GitHub release (manually via web UI or gh CLI)
gh release delete v1.2.0
```

Then fix the issue and create a new tag (e.g., v1.2.1).

## Hotfix Releases

For urgent bug fixes:

1. Create a hotfix branch from the release tag:
   ```bash
   git checkout -b hotfix/1.2.1 v1.2.0
   ```

2. Make the fix and commit:
   ```bash
   git commit -am "fix: critical bug"
   ```

3. Merge back to main:
   ```bash
   git checkout main
   git merge --no-ff hotfix/1.2.1
   ```

4. Tag and push:
   ```bash
   git tag -a v1.2.1 -m "Release v1.2.1 - Hotfix"
   git push origin main v1.2.1
   ```

## Troubleshooting

### Release workflow fails on tests

- Fix the failing tests
- Delete the tag: `git push --delete origin v1.2.0`
- Recreate and push the tag after fixes

### Release workflow fails on build

- Check the build works locally: `make build-all`
- Ensure Go version matches (1.23+)
- Check for platform-specific issues

### Artifacts are missing

- Check GitHub Actions logs for upload errors
- Ensure all platforms built successfully
- Verify packaging step completed

### Docker build fails

- Test locally: `docker build -t test -f docker/Dockerfile .`
- Check Dockerfile exists and is correct
- Verify multi-platform build support

## Manual Release (Emergency)

If GitHub Actions is unavailable:

```bash
# Build locally
make build-all

# Package
cd dist
tar -czf portguard-linux-amd64.tar.gz portguard-linux-amd64
# ... repeat for other platforms

# Create checksums
shasum -a 256 *.tar.gz *.zip > checksums.txt

# Create release via GitHub CLI
gh release create v1.2.0 \
  --title "v1.2.0" \
  --notes "See CHANGELOG.md" \
  dist/*.tar.gz dist/*.zip dist/*.sha256
```

## Post-Release

After a successful release:

1. Announce on relevant channels
2. Update documentation sites if any
3. Create a new section in CHANGELOG.md for next version:
   ```markdown
   ## [Unreleased]
   
   ### Added
   ### Changed
   ### Fixed
   ```

## Getting Help

- Check [GitHub Actions documentation](https://docs.github.com/en/actions)
- Review workflow files in `.github/workflows/`
- Open an issue for automation problems
