ROOT := .
TMP_DIR := ./tmp
SV_BIN := $(TMP_DIR)/server
CSM_BIN := $(TMP_DIR)/consumer
SD_BIN := $(TMP_DIR)/seeder
SC_BIN := $(TMP_DIR)/scheduler
WK_BIN := $(TMP_DIR)/worker
MG_BIN := $(TMP_DIR)/migrate
SV_DIR := ./cmd/server
CSM_DIR := ./cmd/consumer
SD_DIR := ./cmd/seeder
SC_DIR := ./cmd/scheduler
WK_DIR := ./cmd/worker
MG_DIR := ./internal/infrastructure/tools/atlas.go
GCS_KEY_PATH := configs/gcs.json
FIREBASE_KEY_PATH := configs/firebase.json
PROTO_DIR := ./proto
GEN_DIR := ./proto/gen

.PHONY: sec scan vuln proto-gen build-mg migrate-build migrate-baseline migrate-up migrate-status migrate-down lint build-sv run-sv build-wk run-wk build-csm run-csm build-sd run-sd build-sc run-sc clean github

proto-gen:
	@protoc --go_out=$(GEN_DIR)/ai --go-grpc_out=$(GEN_DIR)/ai $(PROTO_DIR)/ai.proto

# Require GitBash, WSL, Ubuntu
lint:
	./scripts/lint.sh

vuln:
	./scripts/vuln.sh

sec:
	./scripts/sec.sh

scan:
	./scripts/container-scan.sh

migrate-build:
	go build -o $(MG_BIN) $(MG_DIR)

migrate-diff: migrate-build
	@set -a; . ./.env.migrate; set +a; \
	read -p "Migration name: " name; \
	name=$$(echo $$name | tr ' ' '_' | tr '[:upper:]' '[:lower:]'); \
	atlas migrate diff $$name --env local

migrate-baseline:
	@set -a; . ./.env.migrate; set +a; \
	read -p "Baseline: " name; \
	atlas migrate apply --env local --baseline $$name

migrate-up:
	@set -a; . ./.env.migrate; set +a; \
	atlas migrate apply --env local

migrate-status:
	@set -a; . ./.env.migrate; set +a; \
	atlas migrate status --env local

migrate-down:
	@set -a; . ./.env.migrate; set +a; \
	atlas migrate down --env local

build-sv:
	@echo "Building..."
	@mkdir -p $(TMP_DIR)
	go build -o $(SV_BIN) $(SV_DIR)

run-sv: build-sv
	@echo "Running..."
	@exec env GOOGLE_APPLICATION_CREDENTIALS=$(GCS_KEY_PATH) $(SV_BIN)

build-sc:
	@echo "Building..."
	@mkdir -p $(TMP_DIR)
	go build -o $(SC_BIN) $(SC_DIR)

run-sc: build-sc
	@echo "Running..."
	@$(SC_BIN)

build-csm:
	@echo "Building..."
	@mkdir -p $(TMP_DIR)
	go build -o $(CSM_BIN) $(CSM_DIR)

run-csm: build-csm
	@echo "Running..."
	@exec env GOOGLE_APPLICATION_CREDENTIALS=$(FIREBASE_KEY_PATH) $(CSM_BIN)

build-wk:
	@echo "Building..."
	@mkdir -p $(TMP_DIR)
	go build -o $(WK_BIN) $(WK_DIR)

run-wk: build-wk
	@echo "Running..."
	@$(WK_BIN)

build-sd:
	@echo "Building..."
	@mkdir -p $(TMP_DIR)
	go build -o $(SD_BIN) $(SD_DIR)

run-sd: build-sd
	@echo "Running..."
	@set -a && \
	. $(CURDIR)/.env.seeder && \
	set +a && \
	$(SD_BIN)

clean:
	@echo "Cleaning..."
	@rm -rf $(TMP_DIR)

github: lint vuln sec
	@if [ -z "$(CM)" ]; then \
		echo "Usage: make github CM=\"commit message\""; \
		exit 1; \
	fi
	@test "$$(git branch --show-current)" = main || { echo "Run make github from the main branch."; exit 1; }
	git add .
	git commit -m "$(CM)"
	git push origin main