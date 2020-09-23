NAME := stress
AWS_PROFILE=exp
VERSION := v0.0.1
# REVISION := $(shell git rev-parse --short HEAD)
OSARCH := "darwin/amd64 linux/amd64"
PACKAGE := github.com/tro3373/$(NAME)

ifndef GOBIN
GOBIN := $(shell echo "$${GOPATH%%:*}/bin")
endif

COBRA := $(GOBIN)/cobra

$(COBRA): ; @go get -v -u github.com/spf13/cobra/cobra

.DEFAULT_GOAL := run

.PHONY: init-gen
init-gen: $(COBRA)
	@go mod init $(PACKAGE) \
	&& $(COBRA) init --pkg-name $(PACKAGE)

.PHONY: add-hello
add-hello: $(COBRA)
	@$(COBRA) add hello

.PHONY: deps
deps:
	@go list -m all

.PHONY: tidy
tidy:
	@go mod tidy

.PHONY: build
build:
	@env GOOS=linux go build -ldflags="-s -w"

.PHONY: build-lambda
build-lambda:
	@env GOOS=linux go build -ldflags="-s -w" -o bin/front ./lambda/front.go

.PHONY: clean
clean:
	rm -rf ./bin logs stress

.PHONY: _do_sls
_do_sls:
	@if [[ -e package.json && ! -e node_modules ]]; then \
		npm install; \
	fi; \
	export AWS_PROFILE=${PROFILE}; \
	export AWS_DEFAULT_PROFILE=${PROFILE}; \
	export SLS_DEBUG=*; \
	sls ${SUB_CMD} -v --stage ${STAGE}; \
	echo done!;

.PHONY: deploy
deploy: clean build-lambda
	$(MAKE) PROFILE=${AWS_PROFILE} STAGE=dev SUB_CMD=deploy _do_sls
.PHONY: remove
remove:
	$(MAKE) PROFILE=${AWS_PROFILE} STAGE=dev SUB_CMD=remove _do_sls

.PHONY: help
help:
	@go run ./main.go --help

.PHONY: front
front:
	@go run ./main.go front

.PHONY: back
back:
	@go run ./main.go back

.PHONY: reguser
reguser:
	@go run ./main.go reguser

.PHONY: run
run: reguser

