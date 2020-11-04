#@IgnoreInspection BashAddShebang

export APP=openflag

export ROOT=$(realpath $(dir $(lastword $(MAKEFILE_LIST))))

export LDFLAGS="-w -s"

all: format lint build

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
