SERVICE_NAME ?= go-admin
REPO_FOLDER ?= $(shell pwd)
BUILD_DIR = $(REPO_FOLDER)/build
GIT_HASH = $(shell git rev-parse --short HEAD)
GIT_TAG = $(shell git describe --tags --exact-match HEAD 2>/dev/null)
CONTAINER_IMAGE_TAG ?= $(if $(GIT_TAG),$(GIT_TAG),$(GIT_HASH))
CONTAINER_IMAGE = $(SERVICE_NAME):$(CONTAINER_IMAGE_TAG)
CONTAINER_IMAGE_LATEST = $(SERVICE_NAME):latest

# Entry point defaults based on service
ENTRYPOINT = cmd/$(SERVICE_NAME)/main.go
DEBUG_ENTRYPOINT = cmd/debug/debug.go

# Goose / DB
GOOSE = go run github.com/pressly/goose/v3/cmd/goose
DB_DRIVER = postgres
DB_NAME = go-admin
DB_DSN = postgres://$(DB_NAME):$(DB_NAME)@localhost:5432/$(DB_NAME)?sslmode=disable
MIGRATIONS_DIR = ./internal/database/sqlc/schema

ENABLE_PPROF = true
PROFILE_OUTPUT = top

######################
# Go Clean & Setup
######################
clean:
	go clean
# go clean -modcache

mod-update:
	go mod tidy
# go mod vendor

######################
# Run Modes
######################
run-debug:
	ENABLE_PPROF=$(ENABLE_PPROF) ENVIRONMENT=debug go run $(ENTRYPOINT)

run-dev:
	ENABLE_PPROF=$(ENABLE_PPROF) ENVIRONMENT=dev go run $(ENTRYPOINT)

run-prod:
	ENABLE_PPROF=$(ENABLE_PPROF) ENVIRONMENT=prod go run $(ENTRYPOINT)

debug:
	go run $(DEBUG_ENTRYPOINT)

######################
# Service-Specific Run Shortcuts
######################
run-go-admin:
	$(MAKE) run-debug SERVICE_NAME=go-admin

run-gatekeeper:
	$(MAKE) run-debug SERVICE_NAME=gate-keeper

######################
# Testing
######################
test:
	go test -v ./...

test-summary:
	go test -v -json ./... | tparse

######################
# Profiling
######################
mem-profile:
	echo "$(PROFILE_OUTPUT)" | go tool pprof mem.pprof

cpu-profile:
	echo "$(PROFILE_OUTPUT)" | go tool pprof cpu.pprof

######################
# DB Migration With Goose
######################
migrate-up:
	go get github.com/pressly/goose/v3/cmd/goose
	$(GOOSE) -dir $(MIGRATIONS_DIR) $(DB_DRIVER) "$(DB_DSN)" up

migrate-down:
	$(GOOSE) -dir $(MIGRATIONS_DIR) $(DB_DRIVER) "$(DB_DSN)" down

migrate-status:
	$(GOOSE) -dir $(MIGRATIONS_DIR) $(DB_DRIVER) "$(DB_DSN)" status

sqlc-generate:
	sqlc generate

######################
# Build
######################
compile-artifacts: pre_compile compile

pre_compile: clean mod-update

compile:
	mkdir -p $(BUILD_DIR)
	cd $(REPO_FOLDER)/cmd/$(SERVICE_NAME) && GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(SERVICE_NAME) main.go

######################
# Container
######################
build-container: compile-artifacts
	docker build --build-arg BINARY_NAME=$(SERVICE_NAME) -t $(CONTAINER_IMAGE) . && \
	docker tag $(CONTAINER_IMAGE) $(CONTAINER_IMAGE_LATEST)

remove-container:
	docker images --format "{{.Repository}}:{{.Tag}}" | grep '^$(SERVICE_NAME)' | xargs -r docker rmi

start-container:
	docker compose -f docker-compose.yaml up -d

stop-container:
	docker compose -f docker-compose.yaml down


######################
# Proto buf
######################

generate-gatekeeper-proto:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	export PATH="$PATH:$(go env GOPATH)/bin"
	protoc \
	--go_out=internal/gatekeeper \
	--go-grpc_out=internal/gatekeeper \
	--go_opt=paths=source_relative \
	--go-grpc_opt=paths=source_relative \
	proto/gatekeeper.proto


######################
# Clean build and Start
######################

clean-build-run:
	$(MAKE) stop-container
	$(MAKE) remove-container SERVICE_NAME=go-admin
	$(MAKE) remove-container SERVICE_NAME=gate-keeper
	$(MAKE) remove-container SERVICE_NAME=gatari-db-init
	$(MAKE) build-container SERVICE_NAME=go-admin
	$(MAKE) build-container SERVICE_NAME=gate-keeper
	$(MAKE) build-container SERVICE_NAME=gatari-db-init
	$(MAKE) start-container
