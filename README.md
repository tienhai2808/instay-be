# Instay Backend

Backend service for a multi-tenant hotel operations platform, built with Go.

## Table of Contents

- [Overview](#overview)
- [Tech Stack](#tech-stack)
- [Architecture and Directory Layout](#architecture-and-directory-layout)
- [System Requirements](#system-requirements)
- [Local Quick Start](#local-quick-start)
- [Migrations and Seeder](#migrations-and-seeder)
- [Run Individual Services](#run-individual-services)
- [Quality Checks](#quality-checks)
- [Docker](#docker)
- [CI](#ci)
- [Additional Documentation](#additional-documentation)
- [Troubleshooting](#troubleshooting)

## Overview

This repository builds multiple binaries:

- `server`: main HTTP API (Gin)
- `worker`: background worker
- `consumer`: message queue consumer (email, notifications, ...)
- `scheduler`: cron jobs
- `seeder`: tenant/admin/data seeding
- `healthcheck`: health binary (`/ping`)

The system uses a multi-tenant model based on `domain -> schema` resolution.
Most APIs require a valid tenant domain; otherwise they return `domain not allowed`.

## Tech Stack

- Go `1.25.12`
- HTTP: Gin
- ORM: GORM + PostgreSQL
- NoSQL: MongoDB
- Cache: Redis
- Message Queue: RabbitMQ
- AI integration: gRPC
- Storage: Google Cloud Storage
- Push notifications: Firebase Cloud Messaging
- Logging: Zap
- Migrations: Atlas

## Architecture and Directory Layout

```text
cmd/                              # entrypoints (server/worker/consumer/...)
internal/application              # use cases, DTOs, ports
internal/domain                   # models, repository contracts, exceptions
internal/infrastructure           # HTTP, persistence, providers, initialization
internal/container                # dependency wiring
migrations/                       # SQL migrations (Atlas)
proto/                            # protobuf contracts
scripts/                          # lint/security/container scan scripts
```

Organization follows a clean architecture style:

- `domain`: business models and contracts
- `application`: orchestration/use-case layer
- `infrastructure`: concrete adapters (DB, MQ, HTTP, external providers)

## System Requirements

Install these dependencies:

- Go `1.25.12`
- `protoc` (if you need to regenerate protobuf code)
- Atlas CLI (for migrations)
- PostgreSQL
- MongoDB
- Redis
- RabbitMQ
- AI gRPC service
- Google service account credentials (GCS/FCM)

Bash/WSL is recommended (Makefile commands are Unix shell-oriented).

## Local Quick Start

### 1) Clone and enter the repository

```bash
git clone <repo-url>
cd backend
```

### 2) Create runtime config file

The app reads `configs/config.yml` (not committed). Create it with:

```bash
cp configs/example.yml configs/config.yml
```

Then fill in real values in `configs/config.yml`.

### 3) Load environment variables

The repo ignores `.env*` files (except `.env.example`), so keep local env files private.

Recommended:

```bash
cp .env.local .env.local.override 2>/dev/null || true
set -a
source .env.local
set +a
```

Important environment variable groups (matching the current loader):

- Server: `SV_PORT`, `SV_WRITE_TIMEOUT`, `SV_READ_TIMEOUT`, `SV_IDLE_TIMEOUT`, `SV_ALLOW_METHODS`, `SV_ALLOW_HEADERS`, `SV_EXPOSE_HEADERS`, `SV_ALLOW_CREDENTIALS`, `SV_MAX_HEADER_BYTES`, `SV_MAX_AGE`
- SQL: `SQL_HOST`, `SQL_PORT`, `SQL_USER`, `SQL_PASSWORD`, `SQL_SSL_MODE`, `SQL_DB_NAME`
- MongoDB: `NOSQL_URI`, `NOSQL_DB_NAME`, `NOSQL_MAX_POOL_SIZE`, `NOSQL_MIN_POOL_SIZE`, `NOSQL_MAX_CONN_IDLE_TIME`
- Redis: `CA_HOST`, `CA_PORT`, `CA_PASSWORD`, `CA_USE_SSL`
- RabbitMQ: `MQ_HOST`, `MQ_PORT`, `MQ_USER`, `MQ_PASSWORD`, `MQ_VHOST`, `MQ_USE_SSL`
- Tokens: `TK_ACCESS_NAME`, `TK_REFRESH_NAME`, `TK_GUEST_NAME`, `TK_SECRET_KEY`, `TK_API_KEY_SECRET`, `TK_API_KEY_HEADER`, `TK_API_KEY_PREFIX`, `TK_ACCESS_EXPIRES_IN`, `TK_REFRESH_EXPIRES_IN`
- AI gRPC: `SVC_AI_HOST`, `SVC_AI_PORT`
- Mailer: `ML_HOST`, `ML_PORT`, `ML_USER`, `ML_PASSWORD`
- Storage/KYC: `ST_BUCKET`, `KYC_ACCESS_TOKEN`, `KYC_WORKFLOW_ID`, `KYC_WEBHOOK_SECRET`
- Logging: `LOG_LEVEL`, `LOG_ENCODING`
- Observability: HTTP metrics are exposed at `GET /metrics` for Prometheus; application logs are emitted as JSON to stdout for Loki/Grafana.

### 4) Google credentials

Set `GOOGLE_APPLICATION_CREDENTIALS` to your service account key JSON.

Current Makefile conventions:

- `configs/gcs.json` for `make run-sv`
- `configs/firebase.json` for `make run-csm`

You can use a single key file if that service account has both GCS and FCM permissions.

### 5) Verify dependencies are up

Server startup initializes and checks several dependencies (SQL/Mongo/Redis/RabbitMQ/GCS/FCM/AI gRPC).
Make sure they are running before starting the app.

## Migrations and Seeder

### Migrations (Atlas)

Build migration helper:

```bash
make migrate-build
```

Provide Atlas env values (`DATABASE_URL`, `DATABASE_DEV_URL`) via `.env.migrate`.

Common commands:

```bash
make migrate-up
make migrate-status
make migrate-down
make migrate-diff
```

### Seeder

`seeder` uses `SD_NAME` values:

- `tenant`: create a new tenant + schema + default permissions/sources/departments + pricing calendar (current year and next year) for that tenant only
- `admin`: create admin user for a schema
- `migrate`: auto-migrate all schemas + seed defaults

Run seeder:

```bash
set -a
source .env.local
set +a
SD_NAME=tenant make run-sd
SD_NAME=admin make run-sd
```

Additional variables per seed:

- `tenant`: `TA_DOMAIN`, `TA_SCHEMA`, `TA_NAME`, `TA_PHONE`, `TA_CITY` (optional `TA_LOGO`)
- `admin`: `AD_SCHEMA`, `AD_USERNAME`, `AD_EMAIL`, `AD_PASSWORD`

## Run Individual Services

### Server

```bash
set -a
source .env.local
set +a
make run-sv
```

Health check:

```bash
curl http://localhost:${SV_PORT}/ping
```

### Consumer

```bash
set -a
source .env.local
set +a
make run-csm
```

### Worker

```bash
set -a
source .env.local
set +a
make run-wk
```

### Scheduler

```bash
set -a
source .env.local
set +a
make run-sc
```

## Quality Checks

### Local scripts

```bash
make lint    # go fmt + golangci-lint
make vuln    # govulncheck
make sec     # gosec
make scan    # docker build + trivy image scan
```

### Tests

```bash
go test ./...
```

## Docker

Build image:

```bash
docker build -t instay-be:local .
```

This image does not define an `ENTRYPOINT`, so pass the target binary explicitly.
Example for API server:

```bash
docker run --rm \
  -p 8080:8080 \
  --env-file .env.local \
  -e GOOGLE_APPLICATION_CREDENTIALS=/app/configs/gcs.json \
  -v $(pwd)/configs:/app/configs \
  instay-be:local /app/server
```

Example for consumer:

```bash
docker run --rm \
  --env-file .env.local \
  -e GOOGLE_APPLICATION_CREDENTIALS=/app/configs/firebase.json \
  -v $(pwd)/configs:/app/configs \
  instay-be:local /app/consumer
```

## CI

Workflow: `.github/workflows/ci.yml`, using GitHub-hosted `ubuntu-24.04` runners.

Pushes to `main`, pull requests targeting `main` and manual runs execute the checks. Container scan and publish run only on a push to `main`.

1. `pre-checks`: gitleaks + golangci-lint.
2. `sast` and `sca` run in parallel: gosec + semgrep, and govulncheck.
3. `build-test`: `go test` with race detection and a `40%` coverage threshold, then binary build.
4. `container-scan`: Docker build → SPDX SBOM generation → Trivy SBOM scan.
5. `publish`: push `ghcr.io/<owner>/<repo>:<commit-sha>`, sign the image and attach the SBOM.

Binaries, image and SBOM move between jobs using workflow artifacts. Deployment selects an explicit commit SHA; the workflow does not publish a mutable branch tag. `make github CM="message"` requires the current branch to be `main` and pushes it via `git push origin main`.

## Additional Documentation

- [API notes](docs/api.md)
- [Router notes](docs/router.md)
- [Startup notes](docs/start.md)
- [Seed notes](docs/seed.md)
- [Dependency notes](docs/go_get.md)

## Troubleshooting

### `domain not allowed`

- Tenant resolution is based on request `Origin`/`Host`.
- Seed a tenant with the exact domain you are calling (for example `localhost`).

### `config file not found`

- Create `configs/config.yml` from `configs/example.yml`.

### AI gRPC connection errors

- Check `SVC_AI_HOST` and `SVC_AI_PORT`.
- Make sure the AI service is listening on that address.

### Google credential errors (GCS/FCM)

- Ensure the JSON key file exists.
- Ensure `GOOGLE_APPLICATION_CREDENTIALS` points to the correct path.

### Migration command failures

- Verify Atlas CLI is installed.
- Verify `.env.migrate` has `DATABASE_URL` (and `DATABASE_DEV_URL` for diff).
