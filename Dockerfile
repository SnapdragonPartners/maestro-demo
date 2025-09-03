# Go-optimized development environment with security constraints
FROM golang:1.21-alpine3.18

# Set environment variables for Go
ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:/usr/local/go/bin:$PATH

# Install minimal required packages for Go development including make, gh, and golangci-lint
RUN apk add --no-cache \
    git \
    ca-certificates \
    make \
    curl \
    wget \
    gnupg \
    && rm -rf /var/cache/apk/* \
    && wget -q -O /tmp/gh.tar.gz https://github.com/cli/cli/releases/download/v2.40.1/gh_2.40.1_linux_amd64.tar.gz \
    && cd /tmp \
    && tar -xzf gh.tar.gz \
    && cp gh_2.40.1_linux_amd64/bin/gh /usr/local/bin/ \
    && chmod +x /usr/local/bin/gh \
    && rm -rf /tmp/gh* \
    && wget -q -O /tmp/golangci-lint.tar.gz https://github.com/golangci/golangci-lint/releases/download/v1.55.2/golangci-lint-1.55.2-linux-amd64.tar.gz \
    && cd /tmp \
    && tar -xzf golangci-lint.tar.gz \
    && cp golangci-lint-1.55.2-linux-amd64/golangci-lint /usr/local/bin/ \
    && chmod +x /usr/local/bin/golangci-lint \
    && rm -rf /tmp/golangci-lint*

# Create necessary directories with proper permissions
RUN mkdir -p /tmp && chmod 1777 /tmp && \
    mkdir -p /go/src /go/bin /go/pkg && \
    mkdir -p /workspace && \
    chown -R nobody:nobody /go /workspace /tmp

# Set working directory
WORKDIR /workspace

# Switch to nobody user for security (rootless execution)
USER nobody

# No network ports exposed for security
# No EXPOSE directives - networking will be disabled

# Health check that works with nobody user and no network
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD go version || exit 1

# Keep container running for development work
CMD ["tail", "-f", "/dev/null"]
