ARG GO_VERSION
FROM mcr.microsoft.com/devcontainers/go:1-${GO_VERSION}-bookworm
ARG GOLANGCI_VERSION=2.1.5


RUN go install github.com/mitranim/gow@latest \
    && go install github.com/swaggo/swag/cmd/swag@latest \
    && go install github.com/nats-io/nats-top@latest \
    && go install github.com/nats-io/natscli/nats@latest \
    && go install mvdan.cc/sh/v3/cmd/shfmt@latest \
    && go install github.com/jstemmer/go-junit-report@latest \
    && curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(go env GOPATH)/bin v${GOLANGCI_VERSION}
