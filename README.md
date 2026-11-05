# Instay Backend

Go backend for a multi-tenant hotel operations platform. Each tenant domain resolves to a database schema; most business endpoints require a registered tenant domain.

## Repository layout

```text
cmd/                     # server, consumer, scheduler, worker, seeder, healthcheck
internal/domain/         # models, repository contracts, exceptions
internal/application/    # use cases, DTOs, ports
internal/infrastructure/ # HTTP, persistence, providers and initialization
internal/container/      # dependency wiring
migrations/              # Atlas SQL migrations for shared tables
proto/                   # protobuf contracts and generated code
scripts/                 # lint, vulnerability and container scans
```

Stack: Go 1.26.10, Gin, GORM/PostgreSQL, MongoDB, Redis, RabbitMQ, gRPC, S3-compatible object storage, Firebase Cloud Messaging, Zap and Atlas.

## Requirements

- Go 1.26.10, GNU Make and Bash on Linux/WSL.
- PostgreSQL, MongoDB, Redis, RabbitMQ and the AI gRPC service, with reachable addresses and credentials.
- A Cloudflare R2 or AWS S3 bucket with a public read URL and credentials for storage operations.
- A Google development service account for Firebase Cloud Messaging.
- Atlas CLI for migrations; Docker if using Atlas's `docker://` development database or building images.
- `protoc`, `protoc-gen-go` and `protoc-gen-go-grpc` only when regenerating protobuf code.
- `golangci-lint`, `govulncheck`, `gosec` for their respective quality-check targets.

For a full local stack or VM deployment, use the separate **infrastructure repository and its README**. Its local Compose expects the backend, frontend, AI, CQA and landing source repositories in sibling folders named as documented there. Its production Compose pulls release images without requiring application source checkouts.

## Local quick start: run the backend on the host

Clone this repository into `backend`, then run all Make commands from that directory:

```bash
git clone <repo-url> backend
cd backend
go mod download
cp -n .env.example .env
chmod 600 .env
```

Edit `.env` before running services. Replace all `REPLACE_*` placeholders and configure the actual dependencies. The sample's `localhost` addresses assume those services are reachable on the host. When using the infrastructure Compose stack, use its own env templates and internal service hostnames instead.

The Go commands load `.env` from the working directory with godotenv. This works through Make, with `go run ./cmd/server`, or by running a built binary from the backend directory.

```bash
make help
make run-sv
```

This builds `tmp/server` and starts it. Apply migrations and create the local tenant/admin described below before using business endpoints.

In separate terminals, start the background processes you need:

```bash
make run-csm
make run-sc
make run-wk
```

Check health (8080 is the sample port; use your configured `SV_PORT`):

```bash
curl --fail http://localhost:8080/ping
```

`GET /metrics` exposes Prometheus metrics. Backend logs go to stdout; `LOG_ENCODING=json` is the recommended structured-log setting.

## Configuration and env precedence

`.env.example` lists the environment variables used by the backend, including `ST_*` and the AWS SDK's `AWS_*` variables for file storage.

- Process environment variables take precedence over dotenv values, including explicitly empty variables. Quote values with spaces or special characters. Use single quotes for literal values containing `$`.
- Runtime commands load `.env` from the working directory. Set `ENV_FILE` to select a different file. Direct runs with environment variables already injected, such as Docker Compose, can start without a dotenv file.
- For the seeder, `SEED_ENV_FILE` takes precedence over the runtime dotenv file. The process environment still takes priority, so `make run-sd SD_NAME=tenant` selects the tenant seed regardless of either file.
- Atlas Make targets source only `MIGRATE_ENV_FILE` (default `.env.migrate`) as trusted Bash assignments.
- Keep real env files and credential JSON outside Git.

Override env paths when necessary:

```bash
make run-sv ENV_FILE=/absolute/path/backend.env
make migrate-status MIGRATE_ENV_FILE=/absolute/path/migrate.env
make run-sd SD_NAME=tenant SEED_ENV_FILE=/absolute/path/tenant.env
```

An explicitly selected dotenv file must exist. Malformed dotenv syntax and invalid numbers, booleans or durations stop startup. The loader reports the affected file or variable without printing its value.

### File storage

Set the same bucket, API endpoint and region on the server and consumer. The server also needs the bucket's public read base URL for chat attachments.

| Setting | Cloudflare R2 | AWS S3 |
| --- | --- | --- |
| `ST_BUCKET` | R2 bucket name | S3 bucket name |
| `ST_ENDPOINT` | `https://<account-id>.r2.cloudflarestorage.com` (use the jurisdiction-specific API endpoint if applicable) | Leave empty to use the SDK's regional endpoint |
| `ST_PUBLIC_BASE_URL` | The public bucket domain, such as `https://files.example.com` or the bucket's `r2.dev` development URL | The public bucket or CDN base URL |
| `AWS_REGION` | `auto` | Bucket region, for example `ap-southeast-1` |
| Credentials | R2 S3 Access Key ID and Secret Access Key via `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` | An IAM role on AWS, or an SDK credential source such as `AWS_PROFILE` locally |

Chat attachments are read directly at `ST_PUBLIC_BASE_URL/<object-key>`, including the attachment URLs returned to CQA. Use the R2 S3 API endpoint for the existing presigned upload and view endpoints. Configure bucket CORS for the app origins, `PUT` and `GET`, and the `Content-Type` request header. Send the file body directly with `PUT` using the same `Content-Type` supplied to the presign API.

Before switching an existing environment, stop issuing upload URLs, allow outstanding uploads to finish, drain the delete queue and sync objects to the destination with their exact keys and metadata. Deploy the server and consumer against the destination together. Check a browser upload, a download of an existing file, a CQA attachment download and a queued delete before resuming normal traffic. Keep the source for rollback and reconcile writes and deletes if rolling back after traffic has resumed.

### Firebase credentials

Set `GOOGLE_APPLICATION_CREDENTIALS` in `.env` to the Firebase service account credential file. Server and consumer both honor this value.

`PORT` is read by the standalone healthcheck binary; keep it equal to `SV_PORT` if you change the default 8080.

## Migrations: Atlas

Create a private `.env.migrate` file with Bash-compatible assignments, for example:

```dotenv
DATABASE_URL='postgres://instay:REPLACE_WITH_PASSWORD@localhost:5432/instay?sslmode=disable'
DATABASE_DEV_URL='docker://postgres/17/dev?search_path=public'
```

Percent-encode special characters in URL credentials. `DATABASE_URL` is the target database. `DATABASE_DEV_URL` is required for diffing and must refer to a disposable development database, not production.

```bash
chmod 600 .env.migrate
make migrate-status
make migrate-up
make migrate-diff NAME=add_room
```

- `migrate-up` runs `atlas migrate apply --env local`.
- `migrate-status` runs `atlas migrate status --env local`.
- `migrate-diff` generates a migration; `NAME` is lowercased and spaces become underscores. Omit it for an interactive prompt.
- `migrate-down` runs the Atlas down command against the selected database; review the intended rollback first.
- `migrate-baseline BASELINE=<existing-migration-version>` marks the baseline for a database already matching that version. Omit `BASELINE` for a prompt; this is not part of every fresh setup.

Atlas targets automatically build the GORM schema helper at `tmp/migrate`, referenced by `atlas.hcl`'s `local` environment. Keep that path aligned if you customize `TMP_DIR`. Atlas manages shared tables; tenant-schema creation/migration is performed by the seeders below. The Docker image's Atlas `prod` environment selects `/app/migrations`; it does not choose a database automatically—the target still comes from `DATABASE_URL`.

## Seeder

Create a private `.env.seeder` containing only seed-specific values; keep shared SQL/log settings in `.env`. Make selects `.env.seeder` by default. For a direct run, set `SEED_ENV_FILE=.env.seeder` to load it. Example for a new local tenant:

```dotenv
TA_DOMAIN=localhost
TA_SCHEMA=local_hotel
TA_NAME='Local Hotel'
TA_PHONE=0900000000
TA_CITY='Ho Chi Minh City'
TA_LOGO=
```

```bash
chmod 600 .env.seeder
make run-sd SD_NAME=tenant
```

Valid seed names:

| `SD_NAME` | Purpose | Required seed variables |
| --- | --- | --- |
| `tenant` | Create tenant, schema and defaults including pricing calendars | `TA_DOMAIN`, `TA_SCHEMA`, `TA_NAME`, `TA_PHONE`, `TA_CITY`; optional `TA_LOGO` |
| `admin` | Create an admin in an existing schema | `AD_SCHEMA`, `AD_USERNAME`, `AD_EMAIL`, `AD_PASSWORD` |
| `migrate` | Auto-migrate tenant schemas and seed defaults across schemas | Existing tenant database; review scope before running |

For an admin, create another private env file, for example `.env.admin`, containing:

```dotenv
AD_SCHEMA=local_hotel
AD_USERNAME=admin
AD_EMAIL=admin@example.com
AD_PASSWORD='REPLACE_WITH_A_STRONG_PASSWORD'
```

```bash
chmod 600 .env.admin
make run-sd SD_NAME=admin SEED_ENV_FILE=.env.admin
```

Alternatively, set `SD_NAME` in the seed env file and run `make run-sd`. A nonempty command-line selection wins over that default. Make validates the seed name before starting the binary. The current Go seeder logs some execution failures and returns without a nonzero exit status; inspect its result/logs rather than treating exit code alone as proof of successful seeding.

Tenant resolution uses request `Origin`/`Host`. Register the domain used by your local frontend or API client; a healthy `/ping` does not mean a tenant exists.

## Quality checks and builds

```bash
make test       # go test ./...
make lint       # go fmt (modifies formatting), then golangci-lint
make vuln       # govulncheck
make sec        # gosec
make scan       # Docker build, then Trivy image scan
make build-sv   # tmp/server
make build-csm  # tmp/consumer
make build-sc   # tmp/scheduler
make build-wk   # tmp/worker
make build-sd   # tmp/seeder
make proto-gen # regenerate Go protobuf code; requires protoc and Go plugins
make clean      # remove tmp/ build output
```

The script targets run Bash with failure checking enabled. They do not load application env files. `make github CM="message"` runs lint/vulnerability/security checks, then requires `main`, stages changes, commits and pushes `main`; this is an explicit publishing command, not a routine build check.

## Docker and deployment

```bash
docker build -t instay-be:local .
```

The image targets Linux amd64 and runs as a non-root user. It has no default application command: select `/app/server`, `/app/consumer`, `/app/worker`, `/app/scheduler` or `/app/seeder` explicitly. Worker/consumer/scheduler/seeder must not inherit the image's HTTP server healthcheck; disable it for those containers, as the infrastructure Compose does.

For a complete stack, use the infrastructure repository. Do not reuse the host `.env` unchanged inside a container: `localhost` refers to that container. Supply a container env file with service hostnames on a shared Docker network, or the correct external service addresses.

An API-container example, after preparing `.env.container`, an existing network `instay-local_backend`, and credentials readable by container UID 65532:

```bash
docker run --rm --init \
  --network instay-local_backend \
  -p 127.0.0.1:8080:8080 \
  --env-file .env.container \
  -e GOOGLE_APPLICATION_CREDENTIALS=/app/configs/firebase.json \
  --mount "type=bind,source=$(pwd)/configs/firebase.json,target=/app/configs/firebase.json,readonly" \
  instay-be:local /app/server
```

This example does not start dependencies. Ensure port 8080 is free if the Compose server is already running. For VM production, use the infrastructure README's root-owned configuration under `/etc/instay/`, release images and restricted deploy command; the host Make targets here are for development.

## CI

Workflow: `.github/workflows/ci.yml`, on GitHub-hosted `ubuntu-24.04` runners.

Pushes to `main`, PRs targeting `main` and manual runs execute checks. Only a push to `main` runs container scan/publish:

1. `pre-checks`: Gitleaks + golangci-lint.
2. `sast` and `sca` in parallel: Gosec/Semgrep and govulncheck.
3. `build-test`: race tests, the existing 40% coverage gate, then binary build.
4. `container-scan`: Docker build → SPDX SBOM → Trivy SBOM scan.
5. `publish`: push `tienhai2808/instay-be:<commit-sha>` to Docker Hub, sign the image and attach the SBOM.

Set the repository Actions secret `DOCKERHUB_TOKEN` to a Docker Hub access token with write access for `tienhai2808`.

Workflow artifacts transfer binaries, image and SBOM between jobs. Select an explicitly published image SHA when deploying. CI does not automatically deploy a VM.

## Additional documentation

- [API notes](docs/api.md)
- [Router notes](docs/router.md)
- [Startup notes](docs/start.md)
- [Seed notes](docs/seed.md)
- [Dependency notes](docs/go_get.md)

## Troubleshooting

- **Missing env file:** create it as described above or pass the correct `ENV_FILE`, `MIGRATE_ENV_FILE` or `SEED_ENV_FILE` to Make.
- **Invalid configuration:** compare env names and types with `.env.example`. Use the environment templates for your run mode.
- **Domain not allowed:** create/register the tenant matching the request domain.
- **Google credential error:** check the path, file permissions, project and service-account permissions; container paths differ from host paths.
- **AI gRPC errors:** check `SVC_AI_HOST`, `SVC_AI_PORT` and actual reachability. Constructing the gRPC client at startup does not prove a successful RPC.
- **Migration failure:** check Atlas installation, target DB credentials and the dev database used for diffing. Baseline and rollback commands require review of the current schema/history.
