#!/bin/bash
# Create a new release tag with validation
# Usage: ./scripts/create-release.sh v1.2.0

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

VERSION=$1

if [ -z "$VERSION" ]; then
    echo -e "${RED}Error: Version not provided${NC}"
    echo "Usage: $0 v1.2.0"
    exit 1
fi

# Validate version format
if ! [[ $VERSION =~ ^v[0-9]+\.[0-9]+\.[0-9]+(-[a-z]+\.[0-9]+)?$ ]]; then
    echo -e "${RED}Error: Invalid version format${NC}"
    echo "Version should be: v1.2.3 or v1.2.3-rc.1"
    exit 1
fi

echo -e "${GREEN}🚀 Creating release $VERSION${NC}"
echo ""

# Check if on main branch
CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
if [ "$CURRENT_BRANCH" != "main" ]; then
    echo -e "${YELLOW}Warning: Not on main branch (current: $CURRENT_BRANCH)${NC}"
    read -p "Continue anyway? (y/N) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

# Check for uncommitted changes
if ! git diff-index --quiet HEAD --; then
    echo -e "${RED}Error: Uncommitted changes detected${NC}"
    echo "Commit or stash changes before creating a release"
    exit 1
fi

# Check if tag already exists
if git rev-parse "$VERSION" >/dev/null 2>&1; then
    echo -e "${RED}Error: Tag $VERSION already exists${NC}"
    exit 1
fi

# Pull latest changes
echo "📥 Pulling latest changes..."
git pull --rebase
echo ""

# Run tests
echo "🧪 Running tests..."
if ! make test > /dev/null 2>&1; then
    echo -e "${RED}Error: Tests failed${NC}"
    echo "Fix tests before releasing"
    exit 1
fi
echo -e "${GREEN}✓ Tests passed${NC}"
echo ""

# Run linter
echo "🔍 Running linter..."
if ! make lint > /dev/null 2>&1; then
    echo -e "${RED}Error: Linter failed${NC}"
    echo "Fix linting issues before releasing"
    exit 1
fi
echo -e "${GREEN}✓ Linter passed${NC}"
echo ""

# Check CHANGELOG.md
echo "📝 Checking CHANGELOG.md..."
if ! grep -q "## \[$VERSION" CHANGELOG.md; then
    echo -e "${YELLOW}Warning: Version $VERSION not found in CHANGELOG.md${NC}"
    read -p "Continue without CHANGELOG entry? (y/N) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "Add an entry to CHANGELOG.md and try again"
        exit 1
    fi
else
    echo -e "${GREEN}✓ CHANGELOG.md updated${NC}"
fi
echo ""

# Confirm
echo -e "${YELLOW}Ready to create and push tag $VERSION${NC}"
echo ""
echo "This will:"
echo "  1. Create annotated tag $VERSION"
echo "  2. Push to origin"
echo "  3. Trigger GitHub Actions release workflow"
echo ""
read -p "Proceed? (y/N) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Cancelled"
    exit 0
fi

# Create tag
echo ""
echo "🏷️  Creating tag..."
git tag -a "$VERSION" -m "Release $VERSION"
echo -e "${GREEN}✓ Tag created${NC}"

# Push tag
echo ""
echo "📤 Pushing tag to origin..."
git push origin "$VERSION"
echo -e "${GREEN}✓ Tag pushed${NC}"

echo ""
echo -e "${GREEN}✅ Release $VERSION created successfully!${NC}"
echo ""
echo "📋 Next steps:"
echo "  1. Monitor workflow: https://github.com/mrwogu/portguard/actions"
echo "  2. Check release: https://github.com/mrwogu/portguard/releases/tag/$VERSION"
echo "  3. Verify artifacts and Docker images"
echo ""
echo "To rollback if needed:"
echo "  git push --delete origin $VERSION"
echo "  git tag -d $VERSION"
