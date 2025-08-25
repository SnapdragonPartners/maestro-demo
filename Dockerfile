# Minimal development environment
FROM ubuntu:22.04

# Set environment variables
ENV DEBIAN_FRONTEND=noninteractive
ENV TZ=UTC
ENV LANG=C.UTF-8
ENV LC_ALL=C.UTF-8

# Create non-root user
RUN useradd -m -s /bin/bash -u 1000 developer && \
    usermod -aG sudo developer && \
    echo 'developer ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers

# Install essential tools only
RUN apt-get update && apt-get install -y \
    curl wget git vim nano \
    ca-certificates build-essential \
    python3 python3-pip \
    && rm -rf /var/lib/apt/lists/*

# Set working directory and ownership
WORKDIR /workspace
RUN chown -R developer:developer /workspace

# Switch to developer user
USER developer

# Expose common development ports
EXPOSE 3000 8000 8080 9000

# Keep container running - use tail to keep alive for boot test
CMD ["tail", "-f", "/dev/null"]
