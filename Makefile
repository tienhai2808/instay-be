SHELL := /bin/bash
.SHELLFLAGS := -eu -o pipefail -c
.DEFAULT_GOAL := help

TMP_DIR := ./tmp
SV_BIN := $(TMP_DIR)/server
CSM_BIN := $(TMP_DIR)/consumer
SD_BIN := $(TMP_DIR)/seeder
SC_BIN := $(TMP_DIR)/scheduler
WK_BIN := $(TMP_DIR)/worker
MG_BIN := $(TMP_DIR)/migrate
MG_SRC := ./internal/infrastructure/tools/atlas.go
PROTO_DIR := ./proto
GEN_DIR := ./proto/gen

ENV_FILE ?= .env.local
MIGRATE_ENV_FILE ?= .env.migrate
SEED_ENV_FILE ?= .env.seeder
GCS_KEY_PATH ?= configs/gcs.json
NAME ?=
BASELINE ?=
CM ?=
export ENV_FILE MIGRATE_ENV_FILE SEED_ENV_FILE GCS_KEY_PATH NAME BASELINE CM
# Export command-line/environment seed selection without embedding it in shell code.
export SD_NAME

# Env files are trusted Bash assignments, not parsed as Make syntax.
# File values override inherited environment values unless explicitly noted below.
define LOAD_ENV
load_env() { \
  local file="$$1"; \
  if [[ ! -f "$$file" ]]; then \
    printf 'Missing env file: %s\nSee README.md for setup.\n' "$$file" >&2; \
    return 1; \
  fi; \
  set -a; source "$$file"; set +a; \
}
endef

.PHONY: help sec scan vuln lint test proto-gen migrate-build migrate-diff \
 migrate-baseline migrate-up migrate-status migrate-down build-sv run-sv \
 build-wk run-wk build-csm run-csm build-sd run-sd build-sc run-sc clean github

help:
	@printf '%s\n' \
	  'Setup: copy .env.example to .env.local and fill credentials (see README.md).' \
	  'Run:   make run-sv | run-csm | run-wk | run-sc' \
	  'Seed:  make run-sd SD_NAME=tenant  (loads runtime env + .env.seeder)' \
	  'Atlas: make migrate-up | migrate-status | migrate-down' \
	  '       make migrate-diff NAME=add_room | migrate-baseline BASELINE=<version>' \
	  'Check: make test | lint | vuln | sec | scan' \
	  'Paths: ENV_FILE=.env.local MIGRATE_ENV_FILE=.env.migrate SEED_ENV_FILE=.env.seeder'

proto-gen:
	@mkdir -p "$(GEN_DIR)/ai"
	protoc --go_out="$(GEN_DIR)/ai" --go-grpc_out="$(GEN_DIR)/ai" "$(PROTO_DIR)/ai.proto"

lint:
	bash -eu -o pipefail ./scripts/lint.sh

vuln:
	bash -eu -o pipefail ./scripts/vuln.sh

sec:
	bash -eu -o pipefail ./scripts/sec.sh

scan:
	bash -eu -o pipefail ./scripts/container-scan.sh

test:
	go test ./...

migrate-build:
	@mkdir -p "$(TMP_DIR)"
	go build -o "$(MG_BIN)" "$(MG_SRC)"

# atlas.hcl's local environment refers to ./tmp/migrate.
migrate-up migrate-status migrate-down: migrate-build
	@$(LOAD_ENV); \
	load_env "$$MIGRATE_ENV_FILE"; \
	: "$${DATABASE_URL:?Set DATABASE_URL in the migration env file}"; \
	atlas migrate $(if $(filter migrate-up,$@),apply,$(patsubst migrate-%,%,$@)) --env local

migrate-diff: migrate-build
	@$(LOAD_ENV); \
	load_env "$$MIGRATE_ENV_FILE"; \
	: "$${DATABASE_URL:?Set DATABASE_URL in the migration env file}"; \
	: "$${DATABASE_DEV_URL:?Set DATABASE_DEV_URL to a disposable development database}"; \
	name="$$NAME"; \
	if [[ -z "$$name" ]]; then read -r -p 'Migration name: ' name; fi; \
	name=$$(printf '%s' "$$name" | tr ' ' '_' | tr '[:upper:]' '[:lower:]'); \
	[[ -n "$$name" ]] || { echo 'Migration name is required.' >&2; exit 1; }; \
	atlas migrate diff "$$name" --env local

migrate-baseline: migrate-build
	@$(LOAD_ENV); \
	load_env "$$MIGRATE_ENV_FILE"; \
	: "$${DATABASE_URL:?Set DATABASE_URL in the migration env file}"; \
	baseline="$$BASELINE"; \
	if [[ -z "$$baseline" ]]; then read -r -p 'Baseline version: ' baseline; fi; \
	[[ -n "$$baseline" ]] || { echo 'Baseline version is required.' >&2; exit 1; }; \
	atlas migrate apply --env local --baseline "$$baseline"

build-sv: BINARY = server
build-csm: BINARY = consumer
build-sc: BINARY = scheduler
build-wk: BINARY = worker
build-sd: BINARY = seeder

build-sv build-csm build-sc build-wk build-sd:
	@mkdir -p "$(TMP_DIR)"
	go build -o "$(TMP_DIR)/$(BINARY)" "./cmd/$(BINARY)"

run-sv: build-sv
	@$(LOAD_ENV); \
	load_env "$$ENV_FILE"; \
	export GOOGLE_APPLICATION_CREDENTIALS="$${GOOGLE_APPLICATION_CREDENTIALS:-$$GCS_KEY_PATH}"; \
	exec "$(SV_BIN)"

run-csm: build-csm
	@$(LOAD_ENV); \
	load_env "$$ENV_FILE"; \
	export GOOGLE_APPLICATION_CREDENTIALS="$${GOOGLE_APPLICATION_CREDENTIALS:-$$GCS_KEY_PATH}"; \
	exec "$(CSM_BIN)"

run-sc: build-sc
	@$(LOAD_ENV); \
	load_env "$$ENV_FILE"; \
	exec "$(SC_BIN)"

run-wk: build-wk
	@$(LOAD_ENV); \
	load_env "$$ENV_FILE"; \
	exec "$(WK_BIN)"

run-sd: build-sd
	@$(LOAD_ENV); \
	seed_override="$${SD_NAME-}"; \
	load_env "$$ENV_FILE"; \
	load_env "$$SEED_ENV_FILE"; \
	if [[ -n "$$seed_override" ]]; then export SD_NAME="$$seed_override"; fi; \
	case "$${SD_NAME-}" in tenant|admin|migrate) ;; \
	  *) echo 'SD_NAME must be tenant, admin or migrate.' >&2; exit 1 ;; esac; \
	exec "$(SD_BIN)"

clean:
	rm -rf -- "$(TMP_DIR)"

github: lint vuln sec
	@[[ -n "$$CM" ]] || { echo 'Usage: make github CM="commit message"' >&2; exit 1; }
	@test "$$(git branch --show-current)" = main || { echo 'Run make github from main.' >&2; exit 1; }
	git add .
	git commit -m "$$CM"
	git push origin main
