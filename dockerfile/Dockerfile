FROM alpine:latest

# Install required dependencies
RUN apk add --no-cache \
    bash \
    gcc \
    g++ \
    musl-dev \
    libc6-compat \
    wget \
    git

# Specify the Go version
ENV GO_VERSION=1.21.1

# Install Go
RUN wget https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz && \
    rm go${GO_VERSION}.linux-amd64.tar.gz

# Add Go to PATH
ENV PATH="/usr/local/go/bin:${PATH}"

# Verify Go installation
RUN go version
