#!/bin/bash
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

MODULE_URL="https://github.com/SnapdragonPartners/maestro-demo.git"
GO_VERSION="1.21"

echo -e "${GREEN}Go Module Initialization and Maintenance${NC}"
echo "================================================"

# Check if go.mod exists
if [ ! -f "go.mod" ]; then
    echo -e "${YELLOW}go.mod not found. Initializing module...${NC}"
    go mod init "$MODULE_URL"
    echo -e "${GREEN}✓ Module initialized with URL: $MODULE_URL${NC}"
else
    echo -e "${GREEN}✓ go.mod file exists${NC}"
    
    # Check if module name matches expected URL
    CURRENT_MODULE=$(head -1 go.mod | awk '{print $2}')
    if [ "$CURRENT_MODULE" != "$MODULE_URL" ]; then
        echo -e "${YELLOW}Module name mismatch. Current: $CURRENT_MODULE, Expected: $MODULE_URL${NC}"
        echo -e "${YELLOW}Updating module name...${NC}"
        
        # Create backup
        cp go.mod go.mod.bak
        
        # Update module name
        sed -i "1s|module .*|module $MODULE_URL|" go.mod
        echo -e "${GREEN}✓ Module name updated to: $MODULE_URL${NC}"
    else
        echo -e "${GREEN}✓ Module name is correct: $CURRENT_MODULE${NC}"
    fi
fi

# Check Go version in go.mod
if grep -q "go $GO_VERSION" go.mod; then
    echo -e "${GREEN}✓ Go version $GO_VERSION is specified${NC}"
else
    echo -e "${YELLOW}Updating Go version to $GO_VERSION...${NC}"
    
    # Add or update go version
    if grep -q "^go " go.mod; then
        sed -i "s/^go .*/go $GO_VERSION/" go.mod
    else
        echo "" >> go.mod
        echo "go $GO_VERSION" >> go.mod
    fi
    echo -e "${GREEN}✓ Go version updated to $GO_VERSION${NC}"
fi

# Run go mod tidy
echo -e "${YELLOW}Running go mod tidy...${NC}"
if go mod tidy; then
    echo -e "${GREEN}✓ go mod tidy completed successfully${NC}"
else
    echo -e "${RED}✗ go mod tidy failed${NC}"
    exit 1
fi

# Verify module dependencies
echo -e "${YELLOW}Verifying module dependencies...${NC}"
if go mod verify; then
    echo -e "${GREEN}✓ Module dependencies verified${NC}"
else
    echo -e "${RED}✗ Module verification failed${NC}"
    exit 1
fi

echo ""
echo -e "${GREEN}Go module initialization and maintenance completed successfully!${NC}"
echo ""
echo "Final go.mod content:"
cat go.mod
