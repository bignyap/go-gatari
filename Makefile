SERVICE_NAME = go-admin
GIT_HASH = $(shell git rev-parse --short HEAD)
GIT_TAG = $(shell git describe --tags --exact-match HEAD 2>/dev/null)
CONTAINER_IMAGE_TAG ?= $(if $(GIT_TAG),$(GIT_TAG),$(GIT_HASH))
CONTAINER_IMAGE = $(SERVICE_NAME):$(CONTAINER_IMAGE_TAG)
CONTAINER_IMAGE_LATEST = $(SERVICE_NAME):latest
REPO_FOLDER ?= $(shell pwd)
BUILD_DIR = $(REPO_FOLDER)/build
ENTRYPOINT = cmd/go-admin/main.go
DEBUG_ENTRYPOINT = cmd/debug/debug.go
PROFILE_OUTPUT = top

######################
# Go Clean & Setup
######################
clean:
	go clean
	go clean -modcache

mod-update:
	go mod tidy
	go mod vendor

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
# Build
######################

compile-artifacts: pre_compile compile

pre_compile: clean mod-setup mod-update

compile:
	mkdir -p $(BUILD_DIR)
	cd $(REPO_FOLDER)/cmd/$(SERVICE_NAME) && GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(SERVICE_NAME) main.go

######################
# Container
######################

build-container:
	cd $(REPO_FOLDER) && \
	docker build -t $(CONTAINER_IMAGE) . && \
	docker tag $(CONTAINER_IMAGE) $(CONTAINER_IMAGE_LATEST)

start-container:
	docker compose -f docker-compose.yaml up -d

stop-container:
	docker compose -f docker-compose.yaml down