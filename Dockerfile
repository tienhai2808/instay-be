ARG GO_VERSION=1.25.12
ARG ALPINE_VERSION=3.23
ARG DEBIAN_VERSION=13

FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS builder

ARG VERSION=v2.0.0
ARG COMMIT_SHA=abc123xyz789
ARG BUILD_DATE=2026-01-06

RUN apk add --no-cache tzdata

WORKDIR /app

COPY go.mod go.sum ./

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download && go mod verify

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    set -e; \
    for cmd in healthcheck consumer seeder scheduler worker; do \
        CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
            -ldflags="-w -s" -trimpath \
            -o $cmd ./cmd/$cmd; \
    done && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
        -ldflags="-w -s -X main.Version=${VERSION} -X main.CommitSHA=${COMMIT_SHA} -X main.BuildDate=${BUILD_DATE}" \
        -trimpath -o server ./cmd/server

FROM alpine:${ALPINE_VERSION} AS tools

RUN wget -q https://release.ariga.io/atlas/atlas-linux-amd64-latest \
    -O /atlas && chmod +x /atlas

FROM gcr.io/distroless/static-debian${DEBIAN_VERSION}:nonroot AS production

WORKDIR /app 

COPY --from=tools /atlas ./atlas
COPY --from=builder --chown=nonroot:nonroot /app/healthcheck .
COPY --from=builder --chown=nonroot:nonroot /app/consumer .
COPY --from=builder --chown=nonroot:nonroot /app/seeder .
COPY --from=builder --chown=nonroot:nonroot /app/scheduler .
COPY --from=builder --chown=nonroot:nonroot /app/worker .
COPY --from=builder --chown=nonroot:nonroot /app/server .
COPY --chown=nonroot:nonroot migrations/ ./migrations/
COPY --chown=nonroot:nonroot atlas.hcl .

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["/app/healthcheck"]

ENV GIN_MODE=release

EXPOSE 8080
