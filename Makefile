###
# make dev       # untuk menjalankan
# make run       # untuk generate documentation & menjalankan
# make docs      # untuk generate OpenAPI docs
# make build     # untuk build binary ke bin/mcp-server
# make tidy      # untuk membersihkan dan sinkronkan dependensi

APP_NAME = api-server
OUTPUT_DIR = bin
PACKAGE := apiserver
# Docker image to generate OAS3 specs
OAS3_GENERATOR_DOCKER_IMAGE = openapitools/openapi-generator-cli

VERSION := $(shell git describe --tags --always 2>/dev/null || echo "dev")
GIT_COMMIT := $(shell git rev-parse --short HEAD)
BUILD_DATE := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -X $(PACKAGE)/configs.Version=$(VERSION) -X $(PACKAGE)/configs.GitCommit=$(GIT_COMMIT) -X $(PACKAGE)/configs.BuildDate=$(BUILD_DATE)
.PHONY: docs dev

dev:
	go run cmd/api/main.go

seed:
	go run cmd/api/main.go --seed

run:
	make docs
	go run cmd/api/main.go

docs:
	rm -rf docs/swagger.json
	rm -rf docs/swagger.yaml
	swag init -g cmd/api/main.go -o docs

docs-noexample:
	# With build tags
	swag init -g cmd/api/main.go -o docs --tags "!Example,!Group,!Permission"

docs-v2:
	rm -rf docs/swagger.json
	rm -rf docs/swagger.yaml
	swag init --dir . \
	  --generalInfo ./cmd/api/main.go \
	  --output docs \
	  --outputTypes json,yaml \
	  --parseDependency false \
	  --parseInternal false \
	  --parseVendor false \
	  --exclude ./data,./bin,sync/atomic

docs-openapi-scalar:
	# First generate v2 docs
	make docs
	# Then convert to v3 using openapi-generator in Docker
	docker run -v ${PWD}:/local $(OAS3_GENERATOR_DOCKER_IMAGE) \
		generate \
		-i /local/docs/swagger.json \
		-g openapi \
		-o /local/docs \
		--skip-validate-spec \
		--additional-properties=outputFileName=openapi.json

build:
	go build -o $(OUTPUT_DIR)/$(APP_NAME) -ldflags "$(LDFLAGS)" ./cmd/api

buildx:
	# make docs
	go mod tidy
	go build -o $(OUTPUT_DIR)/$(APP_NAME) -ldflags "$(LDFLAGS)" ./cmd/api

tidy:
	go mod tidy

