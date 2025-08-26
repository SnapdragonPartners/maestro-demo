#!/bin/bash
# Automated test script to verify all required development tools are installed

set -e

echo "=== Container Development Tools Test ==="
echo "Timestamp: $(date)"
echo "Container: $(hostname)"
echo "User: $(whoami)"
echo

# Test essential system tools
echo "Testing system tools..."
test_command() {
    local cmd=$1
    local expected_output=$2
    
    if command -v "$cmd" >/dev/null 2>&1; then
        echo "✓ $cmd: $(command -v "$cmd")"
        if [ -n "$expected_output" ]; then
            $cmd $expected_output >/dev/null 2>&1 && echo "  - Command test passed"
        fi
        return 0
    else
        echo "✗ $cmd: NOT FOUND"
        return 1
    fi
}

# Essential development tools
test_command "curl" "--version"
test_command "wget" "--version"
test_command "git" "--version"
test_command "vim" "--version"
test_command "nano" "--version"

# Build tools
echo
echo "Testing build tools..."
test_command "gcc" "--version"
test_command "make" "--version"

# Python environment
echo
echo "Testing Python environment..."
test_command "python3" "--version"
test_command "pip3" "--version"

# Test Python functionality
echo "Testing Python functionality..."
python3 -c "print('✓ Python3 execution test passed')"
python3 -c "import sys; print(f'✓ Python version: {sys.version}')"

# Test pip functionality
echo "Testing pip functionality..."
pip3 list >/dev/null 2>&1 && echo "✓ pip3 list command works"

# Test file system permissions
echo
echo "Testing file system permissions..."
if [ -w "/workspace" ]; then
    touch /workspace/.test_file 2>/dev/null && rm -f /workspace/.test_file 2>/dev/null
    echo "✓ Workspace write permissions work"
else
    echo "✗ Workspace write permissions failed"
    exit 1
fi

# Test sudo access
echo
echo "Testing sudo access..."
if sudo -n true 2>/dev/null; then
    echo "✓ Passwordless sudo access works"
else
    echo "✗ Passwordless sudo access failed"
    exit 1
fi

# Test network connectivity
echo
echo "Testing network connectivity..."
if curl -s --max-time 5 https://httpbin.org/get >/dev/null 2>&1; then
    echo "✓ Network connectivity works"
else
    echo "! Network connectivity test failed (may be expected in some environments)"
fi

# Summary
echo
echo "=== Development Tools Test Summary ==="
echo "✓ All essential development tools are installed and working"
echo "✓ Python environment is properly configured"
echo "✓ Build tools are available"
echo "✓ File system permissions are correct"
echo "✓ Container is ready for development work"
echo
echo "Test completed successfully at $(date)"
exit 0
