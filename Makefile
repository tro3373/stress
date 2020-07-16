NAME := go-cobra-example
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
	@go build

.PHONY: help
help:
	@go run ./main.go --help

.PHONY: run
run:
	@go run ./main.go hello

# [Goでhttpリクエストを送信する方法](https://qiita.com/taizo/items/c397dbfed7215969b0a5)
# https://qiita.com/so-heee/items/c739687ed2a609196a33
###############################################################################
# https://qiita.com/tkit/items/3cdeafcde2bd98612428
# git@github.com:tkit/go-cmd-example.git
###############################################################################
# NAME            := cmd-test
# VERSION         := v0.0.1
# REVISION        := $(shell git rev-parse --short HEAD)
# OSARCH          := "darwin/amd64 linux/amd64"
# PROJECTROOT			:= "./"
#
# ifndef GOBIN
# GOBIN := $(shell echo "$${GOPATH%%:*}/bin")
# endif
#
# LINT := $(GOBIN)/golint
# GOX := $(GOBIN)/gox
# ARCHIVER := $(GOBIN)/archiver
# DEP := $(GOBIN)/dep
# JUNITREPORT := $(JUNITREPORT)/dep
#
# $(LINT): ; @go get github.com/golang/lint/golint
# $(GOX): ; @go get github.com/mitchellh/gox
# $(ARCHIVER): ; @go get github.com/mholt/archiver/cmd/archiver
# $(DEP): ; @go get github.com/golang/dep/cmd/dep
# $(JUNITREPORT): ; @go get github.com/jstemmer/go-junit-report
#
# .DEFAULT_GOAL := build
#
# .PHONY: deps
# deps: $(DEP)
# 		dep ensure
#
# .PHONY: build
# build: deps
# 		go build -o bin/$(NAME) ${PROJECTROOT}
#
# .PHONY: install
# install: deps
# 		go install ${PROJECTROOT}
#
# .PHONY: cross-build
# cross-build: deps $(GOX)
# 		rm -rf ./out && \
#  		gox -osarch $(OSARCH) -output "./out/${NAME}_${VERSION}_{{.OS}}_{{.Arch}}/{{.Dir}}" ${PROJECTROOT}
#
# .PHONY: package
# package: cross-build $(ARCHIVER)
# 		rm -rf ./pkg && mkdir ./pkg && \
# 		cd out && \
# 		find * -type d -exec archiver make ../pkg/{}.tar.gz {}/$(NAME) \; && \
# 		cd ../
#
# .PHONY: lint
# lint: $(LINT)
# 		@golint $(PROJECTROOT)/cmd/...
#
# .PHONY: vet
# vet:
# 		@go vet $(PROJECTROOT)/cmd/...
#
# .PHONY: test
# test: deps
# 		@go test -v $(PROJECTROOT)/cmd/...
#
# .PHONY: test-junit
# test-junit: deps $(JUNITREPORT)
# 		@go test -v $(PROJECTROOT)/cmd/... 2>&1 | go-junit-report > report.xml
#
# .PHONY: check
# check: lint vet test build
#
# .PHONY: check-job
# check-job: lint vet test-junit build
###############################################################################
#
#
#
###############################################################################
# https://qiita.com/minamijoyo/items/cfd22e9e6d3581c5d81f
# git@github.com:minamijoyo/api-cli-go-example.git
###############################################################################
# NAME				:= hoge
# VERSION			:= v0.0.1
# REVISION		:= $(shell git rev-parse --short HEAD)
# LDFLAGS			:= "-X github.com/minamijoyo/api-cli-go-example/cmd.Version=${VERSION} -X github.com/minamijoyo/api-cli-go-example/cmd.Revision=${REVISION}"
# OSARCH			:= "darwin/amd64 linux/amd64"
# GITHUB_USER	:= minamijoyo
#
# ifndef GOBIN
# GOBIN := $(shell echo "$${GOPATH%%:*}/bin")
# endif
#
# LINT := $(GOBIN)/golint
# GOX := $(GOBIN)/gox
# ARCHIVER := $(GOBIN)/archiver
# GHR := $(GOBIN)/ghr
#
# $(LINT): ; @go get github.com/golang/lint/golint
# $(GOX): ; @go get github.com/mitchellh/gox
# $(ARCHIVER): ; @go get github.com/mholt/archiver/cmd/archiver
# $(GHR): ; @go get github.com/tcnksm/ghr
#
# .DEFAULT_GOAL := build
#
# .PHONY: deps
# deps:
# 	go get -d -v .
#
# .PHONY: build
# build: deps
# 	go build -ldflags $(LDFLAGS) -o bin/$(NAME)
#
# .PHONY: install
# install: deps
# 	go install -ldflags $(LDFLAGS)
#
# .PHONY: cross-build
# cross-build: deps $(GOX)
# 	rm -rf ./out && \
# 	gox -ldflags $(LDFLAGS) -osarch $(OSARCH) -output "./out/${NAME}_${VERSION}_{{.OS}}_{{.Arch}}/{{.Dir}}"
#
# .PHONY: package
# package: cross-build $(ARCHIVER)
# 	rm -rf ./pkg && mkdir ./pkg && \
# 	pushd out && \
# 	find * -type d -exec archiver make ../pkg/{}.tar.gz {}/$(NAME) \; && \
# 	popd
#
# .PHONY: release
# release: $(GHR)
# 	ghr -u $(GITHUB_USER) $(VERSION) pkg/
#
# .PHONY: lint
# lint: $(LINT)
# 	@golint ./...
#
# .PHONY: vet
# vet:
# 	@go vet ./...
#
# .PHONY: test
# test:
# 	@go test ./...
#
# .PHONY: check
# check: lint vet test build
