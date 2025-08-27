#!/bin/bash

# Go Container Validation Script
# This script validates the Go development container meets security requirements

echo "=== Go Container Security Validation ==="

# Test 1: Verify Go is installed and working
echo "Test 1: Checking Go installation..."
if go version; then
    echo "✓ Go is properly installed"
else
    echo "✗ Go installation failed"
    exit 1
fi

# Test 2: Verify running as nobody user
echo -e "\nTest 2: Checking user permissions..."
current_user=$(whoami)
current_uid=$(id -u)
if [ "$current_user" = "nobody" ] && [ "$current_uid" = "65534" ]; then
    echo "✓ Running as nobody user (UID: $current_uid)"
else
    echo "✗ Not running as nobody user (current: $current_user, UID: $current_uid)"
    exit 1
fi

# Test 3: Verify filesystem restrictions
echo -e "\nTest 3: Checking filesystem permissions..."
if ! touch /root/test 2>/dev/null; then
    echo "✓ Root filesystem is properly restricted"
else
    echo "✗ Root filesystem is writable (security risk)"
    rm -f /root/test 2>/dev/null
    exit 1
fi

# Test 4: Verify /tmp is writable
echo -e "\nTest 4: Checking /tmp writability..."
if touch /tmp/test && rm /tmp/test; then
    echo "✓ /tmp is writable as expected"
else
    echo "✗ /tmp is not writable"
    exit 1
fi

# Test 5: Verify workspace is accessible
echo -e "\nTest 5: Checking workspace access..."
if [ -d "/workspace" ] && [ -r "/workspace" ]; then
    echo "✓ Workspace is accessible"
else
    echo "✗ Workspace is not accessible"
    exit 1
fi

# Test 6: Test Go compilation capability
echo -e "\nTest 6: Testing Go compilation..."
cat > /tmp/hello.go << 'GOEOF'
package main

import "fmt"

func main() {
    fmt.Println("Hello, secure Go container!")
}
GOEOF

if cd /tmp && go build -o hello hello.go && ./hello; then
    echo "✓ Go compilation and execution successful"
    rm -f /tmp/hello /tmp/hello.go
else
    echo "✗ Go compilation failed"
    exit 1
fi

# Test 7: Verify network isolation (should fail)
echo -e "\nTest 7: Checking network isolation..."
if ! ping -c 1 8.8.8.8 2>/dev/null; then
    echo "✓ Network access is properly disabled"
else
    echo "✗ Network access is available (security risk)"
    exit 1
fi

echo -e "\n=== All security validations passed! ==="
echo "Container is ready for secure Go development."
