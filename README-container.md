# Go Development Container

This container provides a secure, Go-optimized development environment with strict security constraints.

## Container Features

### Base Image
- **golang:1.21-alpine3.18** - Official Go image optimized for development
- Alpine Linux for minimal footprint and security
- Go 1.21 with full toolchain

### Security Features
- **Rootless execution**: Runs as `nobody` user (UID: 65534)
- **Read-only filesystem**: Entire filesystem is read-only except /tmp
- **Network isolation**: No network access (`--network=none`)
- **Minimal attack surface**: Only essential packages installed

### Development Capabilities
- Full Go toolchain (go build, go run, go test, etc.)
- Git for version control
- CA certificates for HTTPS operations
- Writable /tmp directory for build artifacts

## Usage

### Build Container
```bash
# Using container tools (recommended)
container_build maestro-demo-dev

# Using docker-compose
docker-compose build
```

### Test Container
```bash
# Boot test
container_test maestro-demo-dev

# Validation script
container_test maestro-demo-dev "sh scripts/validate-go-container.sh"

# Go compilation test
container_test maestro-demo-dev "cd /tmp && echo 'package main; import \"fmt\"; func main() { fmt.Println(\"Hello Go!\") }' > hello.go && go run hello.go"
```

## Security Constraints

- **User**: nobody:nobody (65534:65534)
- **Filesystem**: Read-only except /tmp
- **Network**: Completely disabled
- **Privileges**: No root access

## DevOps Story

This container implementation satisfies the acceptance criteria:
✓ Go-optimized base image (golang:1.21-alpine3.18)
✓ Rootless execution (--user=nobody)
✓ Read-only filesystem except /tmp
✓ Network disabled (--network=none)
✓ Successfully builds and runs Go applications
