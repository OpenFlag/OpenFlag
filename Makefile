#@IgnoreInspection BashAddShebang

export ROOT=$(realpath $(dir $(lastword $(MAKEFILE_LIST))))

export APP=openflag

export APP_VERSION=v0.1.1

export BUILD_INFO_PKG="github.com/OpenFlag/OpenFlag/pkg/version"

export LDFLAGS="-w -s -X '$(BUILD_INFO_PKG).AppVersion=$(APP_VERSION)' -X '$(BUILD_INFO_PKG).Date=$$(date)' -X '$(BUILD_INFO_PKG).BuildVersion=$$(git rev-parse HEAD | cut -c 1-8)' -X '$(BUILD_INFO_PKG).VCSRef=$$(git rev-parse --abbrev-ref HEAD)'"

export PROTO_DIR = pkg/proto

all: format lint build

run-version:
	go run -ldflags $(LDFLAGS) ./cmd/openflag version

run-migrate:
	go run -ldflags $(LDFLAGS) ./cmd/openflag migrate

run-server:
	go run -ldflags $(LDFLAGS) ./cmd/openflag server

build:
	go build -ldflags $(LDFLAGS)  ./cmd/openflag

install:
	go install -ldflags $(LDFLAGS) ./cmd/openflag

release:
	./scripts/release.sh $(APP) $(APP_VERSION) $(LDFLAGS)

check-formatter:
	which goimports || GO111MODULE=off go get -u golang.org/x/tools/cmd/goimports

format: check-formatter
	find $(ROOT) -type f -name "*.go" -not -path "$(ROOT)/vendor/*" | xargs -n 1 -I R goimports -w R
	find $(ROOT) -type f -name "*.go" -not -path "$(ROOT)/vendor/*" | xargs -n 1 -I R gofmt -s -w R

check-linter:
	which golangci-lint || GO111MODULE=off go get -u github.com/golangci/golangci-lint/cmd/golangci-lint@v1.23.8

lint: check-linter
	golangci-lint run $(ROOT)/...

check-go-bindata:
	which go-bindata || GO111MODULE=off go get -u github.com/jteeuwen/go-bindata/...

bindata: check-go-bindata
	cd internal/app/openflag/migrations/postgres && go-bindata -pkg postgres -o ./../bindata/postgres/bindata.go .

install-protoc-gen-go:
	which protoc-gen-go || GO111MODULE=off go get -u github.com/golang/protobuf/protoc-gen-go

code-gen:
	protoc --go_out=plugins=grpc:pkg $(PROTO_DIR)/*.proto

test:
	go test -ldflags $(LDFLAGS) -v -race -p 1 `go list ./... | grep -v integration`

ci-test:
	go test -ldflags $(LDFLAGS) -v -race -p 1 -coverprofile=coverage.txt -covermode=atomic ./...
	go tool cover -func coverage.txt

integration-tests:
	go test -ldflags $(LDFLAGS) -v -race -p 1 `go list ./... | grep integration`

up:
	docker-compose up -d redis postgres

down:
	docker-compose down
