# Containerized Development Environment

This project includes a comprehensive containerized development environment that supports multiple programming languages and development workflows.

## Features

- **Base OS**: Ubuntu 22.04 LTS
- **Languages**: Node.js (LTS), Python 3, Go 1.21
- **Development Tools**: Git, Vim, Nano, curl, wget, jq, tree, htop
- **Build Tools**: GCC, G++, Make, CMake, pkg-config
- **Network Tools**: ping, telnet, netcat
- **User**: Non-root developer user with sudo access

## Quick Start

### Building the Container

```bash
# Build the container image
docker build -t maestro-demo-dev .

# Or use docker-compose
docker-compose build
```

### Running the Container

```bash
# Run with docker-compose (recommended)
docker-compose up -d dev
docker-compose exec dev bash

# Or run directly with Docker
docker run -it --rm -v $(pwd):/workspace -p 3000:3000 maestro-demo-dev
```

### Development Workflow

1. Start the development environment
2. Your code is mounted at `/workspace`
3. Available tools: Node.js, Python 3, Go
4. Exposed ports: 3000, 8000, 8080, 9000

## Container Management

The container includes all necessary development tools pre-installed and configured for immediate use.

### Environment Variables

- `NODE_ENV=development`
- `PYTHONPATH=/workspace`
- `GOPATH=/home/developer/go`

### Security

- Non-root user (developer) with sudo access
- Secure defaults for development workflow

For more information, see the Dockerfile and docker-compose.yml configuration files.
