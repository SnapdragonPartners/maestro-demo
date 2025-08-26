#!/bin/bash
# Health check script for container validation

set -e

echo "=== Container Health Check ==="
echo "Timestamp: $(date)"

# Check if essential commands are available
echo "Checking essential tools..."
commands=("curl" "wget" "git" "vim" "python3" "pip3" "gcc" "make")
for cmd in "${commands[@]}"; do
    if command -v "$cmd" >/dev/null 2>&1; then
        echo "✓ $cmd is available"
    else
        echo "✗ $cmd is missing"
        exit 1
    fi
done

# Check Python installation
echo "Checking Python environment..."
python3 --version
pip3 --version

# Check if workspace is accessible
echo "Checking workspace accessibility..."
if [ -d "/workspace" ] && [ -w "/workspace" ]; then
    echo "✓ Workspace is accessible and writable"
else
    echo "✗ Workspace is not accessible or not writable"
    exit 1
fi

# Check if user is correct
echo "Checking user configuration..."
if [ "$(whoami)" = "developer" ]; then
    echo "✓ Running as developer user"
else
    echo "✗ Not running as developer user"
    exit 1
fi

# Check if ports are available for binding
echo "Checking port availability..."
ports=(3000 8000 8080 9000)
for port in "${ports[@]}"; do
    if ! netstat -ln 2>/dev/null | grep ":$port " >/dev/null; then
        echo "✓ Port $port is available"
    else
        echo "! Port $port is in use (this may be expected)"
    fi
done

echo "=== Health Check Complete ==="
echo "All checks passed successfully!"
exit 0
