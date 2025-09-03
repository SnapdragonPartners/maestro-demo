#!/bin/bash
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}Testing Go Module Initialization${NC}"
echo "=================================="

# Create temporary directory for testing
TEST_DIR=$(mktemp -d)
echo "Test directory: $TEST_DIR"

# Test 1: Initialize module in empty directory
echo -e "\n${YELLOW}Test 1: Initialize module in empty directory${NC}"
cd "$TEST_DIR"
if ../workspace/scripts/init-go-module.sh; then
    echo -e "${GREEN}✓ Test 1 passed: Module initialization in empty directory${NC}"
else
    echo -e "${RED}✗ Test 1 failed${NC}"
    exit 1
fi

# Test 2: Run again on existing module (should update if needed)
echo -e "\n${YELLOW}Test 2: Run on existing module${NC}"
if ../workspace/scripts/init-go-module.sh; then
    echo -e "${GREEN}✓ Test 2 passed: Existing module handling${NC}"
else
    echo -e "${RED}✗ Test 2 failed${NC}"
    exit 1
fi

# Test 3: Verify go.mod content
echo -e "\n${YELLOW}Test 3: Verify go.mod content${NC}"
if grep -q "https://github.com/SnapdragonPartners/maestro-demo.git" go.mod && grep -q "go 1.21" go.mod; then
    echo -e "${GREEN}✓ Test 3 passed: go.mod has correct content${NC}"
else
    echo -e "${RED}✗ Test 3 failed: go.mod content incorrect${NC}"
    cat go.mod
    exit 1
fi

# Test 4: Test with wrong module name
echo -e "\n${YELLOW}Test 4: Test module name correction${NC}"
echo "module wrong-module" > go.mod
echo "go 1.20" >> go.mod
if ../workspace/scripts/init-go-module.sh; then
    if grep -q "https://github.com/SnapdragonPartners/maestro-demo.git" go.mod && grep -q "go 1.21" go.mod; then
        echo -e "${GREEN}✓ Test 4 passed: Module name and version corrected${NC}"
    else
        echo -e "${RED}✗ Test 4 failed: Module name/version not corrected${NC}"
        exit 1
    fi
else
    echo -e "${RED}✗ Test 4 failed: Script execution failed${NC}"
    exit 1
fi

# Cleanup
cd /
rm -rf "$TEST_DIR"

echo -e "\n${GREEN}All tests passed! Module initialization script is working correctly.${NC}"
