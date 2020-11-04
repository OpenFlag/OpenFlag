#@IgnoreInspection BashAddShebang

export APP=openflag

export ROOT=$(realpath $(dir $(lastword $(MAKEFILE_LIST))))

export BUILD_INFO_PKG="github.com/OpenFlag/OpenFlag/pkg/version"

export LDFLAGS="-w -s -X $(BUILD_INFO_PKG).Date=$$(TZ=Asia/Tehran date '+%FT%T') -X $(BUILD_INFO_PKG).Version=$$(git rev-parse HEAD | cut -c 1-8) -X $(BUILD_INFO_PKG).VCSRef=$$(git rev-parse --abbrev-ref HEAD)"

all: format lint build

run-version:
	go run -ldflags $(LDFLAGS) ./cmd/openflag version

run-migrate:
	go run -ldflags $(LDFLAGS) ./cmd/openflag migrate --path migrations

run-server:
	go run -ldflags $(LDFLAGS) ./cmd/openflag server

build:
	CGO_ENABLED=1 go build -ldflags $(LDFLAGS)  ./cmd/openflag

install:
	CGO_ENABLED=1 go install -ldflags $(LDFLAGS) ./cmd/openflag

check-formatter:
	which goimports || GO111MODULE=off go get -u golang.org/x/tools/cmd/goimports

format: check-formatter
	find $(ROOT) -type f -name "*.go" -not -path "$(ROOT)/vendor/*" | xargs -n 1 -I R goimports -w R
	find $(ROOT) -type f -name "*.go" -not -path "$(ROOT)/vendor/*" | xargs -n 1 -I R gofmt -s -w R

check-linter:
	which golangci-lint || GO111MODULE=off go get -u github.com/golangci/golangci-lint/cmd/golangci-lint

lint: check-linter
	golangci-lint run $(ROOT)/...

test:
	go test -v -race -p 1 ./...

ci-test:
	go test -v -race -p 1 -coverprofile=coverage.txt -covermode=atomic ./...
	go tool cover -func coverage.txt

up:
	docker-compose up -d

down:
	docker-compose down
	rm -rf data
