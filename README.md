# maestro-demo
Testing Maestro development

## Development Environment

This project includes a containerized development environment with the following features:

### Container Specifications
- **Base Image**: Ubuntu 22.04 LTS
- **Languages**: Python 3 with pip
- **Development Tools**: Git, Vim, Nano, curl, wget
- **Build Tools**: GCC, G++, build-essential
- **User**: Non-root developer user with sudo access

### Quick Start

Build and run the development container:

```bash
# Build the container
docker build -t maestro-demo-dev .

# Or use docker-compose
docker-compose up -d dev
docker-compose exec dev bash

# Direct run
docker run -it --rm -v $(pwd):/workspace maestro-demo-dev
```

### Available Ports
- 3000: Web development
- 8000: Python development server  
- 8080: Alternative web server
- 9000: Additional services

### Files Created
- `Dockerfile`: Container definition
- `docker-compose.yml`: Container orchestration
- `.dockerignore`: Build optimization
- `README-container.md`: Detailed container documentation

The development environment is ready for application code deployment and testing.
