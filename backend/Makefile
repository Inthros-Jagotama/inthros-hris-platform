# HRIS Platform Makefile
# ============ Commands ============
# make build       - Build Go binary
# make run         - Run server locally
# make test        - Run all tests
# make race-test   - Run all tests with race detector in Docker (golang:1.24)
# make db-race-up  - Start MySQL + PostgreSQL containers for the migrator integration test
# make lint        - Run linter
# make check-openapi - Verify all registered endpoints are documented in openapi.json
# make docs        - Run check-openapi + regenerate docs/openapi-report.md
# make db-docs     - Regenerate docs/database-schema.md from SQL migrations
# make check-db-docs - Verify docs/database-schema.md is in sync with migrations (no overwrite)
# make docker      - Build Docker image
# make clean       - Clean build artifacts
# make migrate     - Run platform migrations
# make seed        - Run seeders
# make coverage    - Run tests with coverage

.PHONY: build build-installer run run-hot test test-verbose test-short test-pkg race-test db-race-up db-race-down test-employee test-organization test-cache test-cache-quick test-cache-integration test-cache-error bench-cache bench-cache-compare coverage cover-view cover-func lint vet check-openapi docs db-docs check-db-docs migrate seed seed-modules docker docker-compose-up docker-compose-down tidy clean help

BUILD_DIR := ./bin
BINARY := $(BUILD_DIR)/server

# Absolute path of the Makefile dir with forward slashes (docker -v source).
# Normalized so it works whether make is invoked from git-bash, cmd, or Linux.
BACKEND_DIR := $(subst \,/,$(CURDIR))

# Python interpreter for docs tooling (prefer python3, fallback to python)
PYTHON := $(shell command -v python3 2>/dev/null || command -v python 2>/dev/null)

# ============ Build ============

build:
	@echo "Building HRIS Platform..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 go build -ldflags="-s -w" -o $(BINARY) ./cmd/server
	@echo "Build complete: $(BINARY)"

build-installer:
	@echo "Building installer..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 go build -ldflags="-s -w" -o $(BUILD_DIR)/installer ./cmd/installer
	@echo "Build complete: $(BUILD_DIR)/installer"

# ============ Run ============

run:
	@echo "Starting server..."
	@go run ./cmd/server --config ./config/config.yaml

run-hot:
	@echo "Starting server with hot reload..."
	@which air || go install github.com/air-verse/air@latest
	@air

# ============ Test ============

# NOTE: `test` and `test-verbose` run -race locally, which needs cgo + a C
# compiler. On Windows without gcc these fail; use `make race-test` instead.
test:
	@echo "Running all tests..."
	@go test -v -race -count=1 ./...

test-verbose:
	@echo "Running all tests (verbose)..."
	@go test -v -race -count=1 ./... 2>&1 | tee test-output.log

test-short:
	@echo "Running short tests (non-integration)..."
	@go test -short -count=1 ./...

test-pkg:
	@echo "Running tests for package: $(pkg)"
	@go test -v -race -count=1 ./internal/$(pkg)/...

# Run the full test suite with the race detector inside a golang:1.24 container.
#
# Why Docker: on Windows the race detector requires cgo + a C compiler
# (CGO_ENABLED=1), which is not always installed locally. The official golang
# image ships with gcc, so -race works out of the box.
#
# Details:
#   - --network host  : integration tests (migrator) reach MySQL:3306 &
#                       PostgreSQL:5432 (started automatically via db-race-up).
#   - volume caches   : go-mod-cache & go-build-cache persist across runs so
#                       repeat invocations are fast.
#   - MSYS_NO_PATHCONV: prevents git-bash from mangling the -v mount path
#                       (Windows only; harmless no-op on Linux/macOS).
#   - PKG override    : e.g. make race-test PKG=./internal/modules/leave/...
# Default target set: the whole module tree.
GO_TEST_IMAGE ?= golang:1.24
RACE_PARALLEL ?= 4
PKG ?= ./...

# DB containers for the migrator integration test. The DSNs are hardcoded in
# internal/pkg/migrator/migrator_integration_test.go:
#   MySQL    : root:@tcp(localhost:3306)/hris_migrate_test
#   Postgres : postgres:password@localhost:5432/hris_migrate_test
RACE_MYSQL_CTR ?= hris-race-mysql
RACE_PG_CTR ?= hris-race-pg
RACE_MYSQL_IMAGE ?= mysql:8.0
RACE_PG_IMAGE ?= postgres:16-alpine

# Container creation commands (shared so db-race-up avoids duplication).
RACE_MYSQL_RUN := docker run -d --name $(RACE_MYSQL_CTR) --network host -e MYSQL_ALLOW_EMPTY_PASSWORD=yes -e MYSQL_DATABASE=hris_migrate_test $(RACE_MYSQL_IMAGE)
RACE_PG_RUN := docker run -d --name $(RACE_PG_CTR) --network host -e POSTGRES_PASSWORD=password -e POSTGRES_DB=hris_migrate_test $(RACE_PG_IMAGE)

# Start MySQL + PostgreSQL as --network host containers matching the DSNs above.
# Idempotent: reuses an existing container; if it exists but fails to start, it is
# removed and recreated. NOTE: ports 3306/5432 must be free inside the Docker VM —
# stop `docker compose --profile mysql/postgres up` services first if they run.
db-race-up:
	@echo "Starting race-test DBs (MySQL + PostgreSQL, --network host)..."
	@if docker inspect $(RACE_MYSQL_CTR) >/dev/null 2>&1; then \
		docker start $(RACE_MYSQL_CTR) >/dev/null 2>&1 || { docker rm -f $(RACE_MYSQL_CTR) >/dev/null 2>&1; $(RACE_MYSQL_RUN); }; \
	else \
		$(RACE_MYSQL_RUN); \
	fi
	@if docker inspect $(RACE_PG_CTR) >/dev/null 2>&1; then \
		docker start $(RACE_PG_CTR) >/dev/null 2>&1 || { docker rm -f $(RACE_PG_CTR) >/dev/null 2>&1; $(RACE_PG_RUN); }; \
	else \
		$(RACE_PG_RUN); \
	fi
	@echo "Waiting for DBs to accept connections..."
	@for i in $$(seq 1 40); do \
		mysql_up=$$(docker exec $(RACE_MYSQL_CTR) mysqladmin ping -uroot --silent >/dev/null 2>&1 && echo up || echo down); \
		pg_up=$$(docker exec $(RACE_PG_CTR) pg_isready -U postgres -h localhost >/dev/null 2>&1 && echo up || echo down); \
		if [ "$$mysql_up" = up ] && [ "$$pg_up" = up ]; then echo "DBs ready (t=$$((i*3))s)"; exit 0; fi; \
		sleep 3; \
	done; \
	echo "ERROR: DBs not ready within 120s"; exit 1

# Stop and remove the race-test DB containers.
db-race-down:
	@echo "Stopping & removing race-test DB containers..."
	@docker rm -f $(RACE_MYSQL_CTR) 2>/dev/null || true
	@docker rm -f $(RACE_PG_CTR) 2>/dev/null || true
	@echo "DB containers removed."

race-test: db-race-up
	@echo "Running tests with race detector in Docker ($(GO_TEST_IMAGE))..."
	@echo "  NOTE: MySQL + PostgreSQL started automatically by db-race-up (migrator integration test)."
	@MSYS_NO_PATHCONV=1 docker run --rm \
		--network host \
		-v "$(BACKEND_DIR):/app" \
		-v go-mod-cache:/go/pkg/mod \
		-v go-build-cache:/root/.cache/go-build \
		-w /app \
		$(GO_TEST_IMAGE) go test -race -count=1 -p $(RACE_PARALLEL) -timeout 15m $(PKG)

test-employee:
	@echo "Running employee module tests..."
	@go test -v -count=1 ./internal/modules/employee/...

test-organization:
	@echo "Running organization module tests..."
	@go test -v -count=1 ./internal/modules/organization/...

test-cache:
	@echo "Running ALL cache package tests (unit + integration)..."
	@go test -v -count=1 -timeout 120s ./internal/pkg/cache/...

test-cache-quick:
	@echo "Running all cache tests except integration..."
	@go test -v -count=1 -timeout 60s -run '^(TestCache_|TestPubSub_|Benchmark)' ./internal/pkg/cache/...

test-cache-integration:
	@echo "Running cache integration tests only..."
	@go test -v -count=1 -timeout 120s -run TestCacheIntegration ./internal/pkg/cache/...

test-cache-error:
	@echo "Running cache error path tests..."
	@go test -v -count=1 -timeout 60s -run 'TestCache_New_Invalid|TestCache_New_Empty|TestCache_New_Closed|TestCache_New_Error|TestPubSub_NewPubSub_Invalid|TestPubSub_CallbackNotCalled' ./internal/pkg/cache/...

bench-cache:
	@echo "Running cache benchmarks..."
	@go test -bench=. -benchmem -count=1 -timeout 180s -run=^$ ./internal/pkg/cache/...

bench-cache-compare:
	@echo "Running cache benchmark comparison (2 iterations)..."
	@go test -bench=. -benchmem -count=2 -timeout 300s -run=^$ ./internal/pkg/cache/... 2>&1 | tee bench-cache-output.log
	@echo "Benchmark output saved to bench-cache-output.log"

coverage:
	@echo "Running tests with coverage..."
	@go test -v -count=1 -coverprofile=coverage.out -covermode=atomic ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

cover-view:
	@echo "Opening coverage report..."
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

cover-func:
	@echo "Coverage by function:"
	@go tool cover -func=coverage.out

# ============ Lint ============

lint:
	@echo "Running linter..."
	@which golangci-lint || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.59.0
	@golangci-lint run ./...

vet:
	@echo "Running go vet..."
	@go vet ./...

# ============ API Documentation ============

check-openapi:
	@echo "Checking OpenAPI docs coverage against registered routes..."
	@test -n "$(PYTHON)" || { echo "ERROR: python3/python required to run this check"; exit 1; }
	@"$(PYTHON)" ../scripts/check_missing_openapi.py
	@echo "OpenAPI check passed: all registered endpoints are documented."

docs: check-openapi
	@echo "Regenerating OpenAPI report..."
	@cd .. && "$(PYTHON)" scripts/generate_openapi_report.py
	@echo "OpenAPI report regenerated: docs/openapi-report.md"

db-docs:
	@echo "Regenerating database schema documentation..."
	@test -n "$(PYTHON)" || { echo "ERROR: python3/python required to run this check"; exit 1; }
	@cd .. && "$(PYTHON)" scripts/generate_database_schema_doc.py
	@echo "Database schema doc regenerated: docs/database-schema.md"

check-db-docs:
	@echo "Checking database schema doc sync with SQL migrations..."
	@test -n "$(PYTHON)" || { echo "ERROR: python3/python required to run this check"; exit 1; }
	@cd .. && "$(PYTHON)" scripts/check_database_schema_doc.py
	@echo "Database schema doc check passed: docs/database-schema.md is up to date."

# ============ Database ============

migrate:
	@echo "Running platform migrations..."
	@go run ./cmd/server migrate --config ./config/config.yaml

seed:
	@echo "Running seeders..."
	@go run ./cmd/server seed --config ./config/config.yaml

seed-modules:
	@echo "Running module seeder..."
	@go run ./cmd/installer seed-modules --config ./config/config.yaml

# ============ Docker ============

docker:
	@echo "Building Docker image..."
	@docker build -t hris-platform:latest -f docker/Dockerfile .

docker-compose-up:
	@echo "Starting Docker Compose environment..."
	@docker compose -f ../docker/docker-compose.yml up -d

docker-compose-down:
	@echo "Stopping Docker Compose environment..."
	@docker compose -f ../docker/docker-compose.yml down

# ============ Utilities ============

clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html test-output.log
	@go clean

tidy:
	@echo "Tidying dependencies..."
	@go mod tidy
	@go mod verify

help:
	@echo "HRIS Platform Makefile"
	@echo "======================"
	@echo "=== Build ==="
	@echo "build          - Build Go binary"
	@echo "build-installer - Build CLI installer"
	@echo ""
	@echo "=== Run ==="
	@echo "run            - Run server locally"
	@echo "run-hot        - Run with hot reload (air)"
	@echo ""
	@echo "=== Test ==="
	@echo "test           - Run all tests (verbose + race + count=1)"
	@echo "test-verbose   - Run all tests, output to test-output.log"
	@echo "test-short     - Run short tests (non-integration)"
	@echo "test-pkg pkg=   - Run tests for a specific package (e.g. pkg=modules/employee)"
	@echo "race-test PKG=  - Run tests with race detector in Docker; starts DBs via db-race-up; PKG=./... default"
	@echo "db-race-up      - Start MySQL + PostgreSQL (--network host) for migrator integration test"
	@echo "db-race-down    - Stop & remove race-test DB containers"
	@echo "test-employee  - Run employee module tests"
	@echo "test-organization - Run organization module tests"
	@echo "test-cache     - Run ALL cache tests (unit + integration)"
	@echo "test-cache-quick - Run all cache tests except integration"
	@echo "test-cache-integration - Run cache integration tests only"
	@echo "test-cache-error - Run cache error path tests"
	@echo "bench-cache    - Run cache benchmarks"
	@echo "bench-cache-compare - Run 2 benchmark iterations + save to log"
	@echo "coverage       - Run tests with coverage report (HTML)"
	@echo "cover-view     - Open existing coverage report"
	@echo "cover-func     - Show coverage by function"
	@echo ""
	@echo "=== Lint ==="
	@echo "lint           - Run linter (golangci-lint)"
	@echo "vet            - Run go vet"
	@echo ""
	@echo "=== API Documentation ==="
	@echo "check-openapi  - Verify all registered endpoints are documented in openapi.json"
	@echo "docs           - Run check-openapi + regenerate docs/openapi-report.md"
	@echo "db-docs        - Regenerate docs/database-schema.md from SQL migrations"
	@echo "check-db-docs  - Verify docs/database-schema.md is in sync with migrations (no overwrite)"
	@echo ""
	@echo "=== Database ==="
	@echo "migrate        - Run platform migrations"
	@echo "seed           - Run seeders"
	@echo "seed-modules   - Seed/update modules table from module definitions"
	@echo ""
	@echo "=== Docker ==="
	@echo "docker         - Build Docker image"
	@echo "docker-compose-up   - Start Docker Compose environment"
	@echo "docker-compose-down - Stop Docker Compose environment"
	@echo ""
	@echo "=== Utilities ==="
	@echo "tidy           - Tidy and verify dependencies"
	@echo "clean          - Clean build artifacts"
	@echo "help           - Show this help"