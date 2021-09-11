BASEDIR = $(shell pwd)
APPNAME = $(shell basename $(BASEDIR))

PACKAGE = github.com/drvspw/$(APPNAME)

SOURCE_FILES?=./...
TEST_PATTERN?=.
TEST_OPTIONS?=-v

export GO111MODULE := on

LINTERS = fmtcheck lint tidy

.DEFAULT_GOAL := help

.PHONY: clean tools build build test fmt fmtcheck lint cover clean tidy

init: ## create module
	go mod init $(PACKAGE)

tools: ## setup
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.37.1

# Uncomment this rule for building executable
build: $(LINTERS) ## build
	GOOS=linux GOARCH=amd64 go build -v -o $(GOPATH)/bin/$(APPNAME)-linux-amd64 $(PACKAGE)

# If building a library, test should depend on build-lib
test: $(LINTERS) ## test
	go test $(TEST_OPTIONS) -failfast -race -coverpkg=./... -covermode=atomic -coverprofile=coverage.txt $(SOURCE_FILES) -run $(TEST_PATTERN) -timeout=2m

fmt: ## format
	@echo "==> Fixing source code with gofmt..."
	gofmt -s -w ./$(PKG_NAME)

fmtcheck: ## check format
	@sh -c "'$(CURDIR)/gofmtcheck.sh'"

lint: ## lint
	@echo "==> Checking source code against linters..."
	@GOGC=30 golangci-lint run ./...

tidy: ## tidy
	go mod tidy

cover: test ## generate coverage report
	go tool cover -html=coverage.txt

clean: ## clean
	rm $(GOPATH)/bin/$(APPNAME)*

run: build ## run
	$(GOPATH)/bin/$(APPNAME)-linux-amd64 $(filter-out $@,$(MAKECMDGOALS))
%:
	@:

docker-dev: ## run application in a dev container
	docker build -f Dockerfile -t $(APPNAME) ./
	docker run --rm -p 8090:8090 --volume $(PWD):/go/src/$(PACKAGE) --name $(APPNAME) $(APPNAME)

docker-prod: ## run application is a prod container
	docker build -f Dockerfile -t $(APPNAME) ./ --build-arg app_env=production
	docker run --rm -p 8090:8090 --name $(APPNAME) $(APPNAME)

help: ## help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
