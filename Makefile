.DEFAULT_GOAL := all
.PHONY := clean all fmt coverage build docker

go_files := $(shell find . -path ./vendor -prune -o -path '*/testdata' -prune -o -path './tools' -prune -o -type f -name '*.go' -print)
commands := $(notdir $(shell find cmd/* -maxdepth 0 -type d))
local_bins := $(addprefix bin/,$(commands))

clean:
	# Removing all generated files...
	@rm -rf bin/ || true

bin/goimports: $(shell find ./tools/goimports -type f)
	# Building goimports
	@cd ./tools/goimports && go generate ./tools.go

bin/ko: $(shell find ./tools/ko -type f)
	# Building ko
	@cd ./tools/ko && go generate ./tools.go

bin/golangci-lint: $(shell find ./tools/golangci-lint -type f)
	# Building golangci-lint
	@cd ./tools/golangci-lint && go generate ./tools.go

bin/.vendor: go.mod go.sum
	# Downloading modules...
	@go mod download
	@mkdir -p bin/
	@touch bin/.vendor

bin/.generate: $(go_files) bin/.vendor
	@go generate ./...
	@touch bin/.generate

fmt: bin/.generate $(go_files) bin/goimports
	# Formatting files...
	@bin/goimports -w $(go_files)

bin/.vet: bin/golangci-lint bin/.generate $(go_files)
	bin/golangci-lint run --enable bodyclose --enable goimports
	@touch bin/.vet

bin/.fmtcheck: bin/goimports bin/.generate $(go_files)
	# Checking format of Go files...
	@GOIMPORTS=$$(bin/goimports -l $(go_files)) && \
	if [ "$$GOIMPORTS" != "" ]; then \
		bin/goimports -d $(go_files); \
		exit 1; \
	fi
	@touch bin/.fmtcheck

bin/.coverage.out: bin/.generate $(go_files)
	@go test -cover -v -count=1 ./... -coverprofile bin/.coverage.tmp
	@mv bin/.coverage.tmp bin/.coverage.out

coverage: bin/.coverage.out
	@go tool cover -html=bin/.coverage.out

$(local_bins): bin/.fmtcheck bin/.vet bin/.coverage.out $(go_files)
	CGO_ENABLED=0 go build -trimpath -o $@ ./cmd/$(basename $(@F))

docker: bin/ko bin/.fmtcheck bin/.vet bin/.coverage.out $(go_files)
	KO_DOCKER_REPO=github.com/wjam bin/ko publish --base-import-paths --push=false --platform linux/amd64,linux/arm64 ./cmd/policy

build: $(local_bins)

all: build
