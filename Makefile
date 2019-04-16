# TOOLCHAIN
GO			?= CGO_ENABLED=0 GO111MODULE=on go
CGO			?= CGO_ENABLED=1 GO111MODULE=on go
GOFMT		?= $(GO)fmt

# ENVIRONMENT
FAST		=
GOFILES		:= $(shell find . -name '*.go')
GOPATH		:= $(GOPATH)
GOOS		:= $(shell echo $(shell uname -s) | tr A-Z a-z)
GOARCH		:= amd64
MOD_NAME	:= github.com/lukasmalkmus/spl

# APPLICATION INFORMATION
BUILD_TIME	?= $(shell date -u '+%Y-%m-%d_%H:%M:%S')
COMMIT		?= $(shell git rev-parse --short HEAD)
RELEASE		?= $(shell cat VERSION)
USER		?= $(shell whoami)

# TOOLS
GOLANGCI_LINT	:= $(GOPATH)/bin/golangci-lint

# TEST MISC
COVERPROFILE 	:= coverage.out

# BUILD INFORMATION
GOFLAGS		?= -mod=vendor -buildmode=exe -tags=netgo -installsuffix=cgo \
				-gcflags=-trimpath=$(GOPATH)/src \
				-asmflags=-trimpath=$(GOPATH)/src \
				-ldflags='-s -w -extldflags "-static" \
				-X $(MOD_NAME)/pkg/version.BuildTime=$(BUILD_TIME) \
				-X $(MOD_NAME)/pkg/version.Commit=$(COMMIT) \
				-X $(MOD_NAME)/pkg/version.Release=$(RELEASE) \
				-X $(MOD_NAME)/pkg/version.User=$(USER)'

# Disable caching if FAST is not explicitly set.
ifndef FAST
	GOFLAGS	+= -a
endif

# Add msan target if system is linux/amd64 or linux/arm64.
ifeq ($(GOOS),linux)
	ifeq ($(shell uname -m),x86_64)
		all	= lint test race bench msan build config
	endif
	ifeq ($(shell uname -m),arm64)
		all	= lint test race bench msan build config
	endif
endif

.PHONY: all bench build config cover clean dep dep-update fmt install lint msan race test tools

all: fmt lint test race bench build config clean

bench: $(GOFILES)
	@echo ">> running benchmarks"
	@$(GO) test -mod=vendor -run=NONE -bench=. -benchmem ./...

build: .build/spl-$(GOOS)-$(GOARCH)

config: build
	@echo ">> generating default configuration file"
	@rm -rf configs/spl.toml
	@./.build/spl-$(GOOS)-$(GOARCH) config > configs/spl.toml

cover: | $(COVERPROFILE)
	@echo ">> calculating code coverage"
	@$(GO) tool cover -func $(COVERPROFILE)

clean:
	@echo ">> cleaning up artifacts"
	@rm -rf .build
	@rm -rf $(COVERPROFILE)

dep:
	@echo ">> installing dependencies"
	@$(GO) mod tidy
	@$(GO) mod download
	@$(GO) mod verify
	@$(GO) mod vendor

fmt: $(GOFILES)
	@echo ">> formatting code"
	@! $(GOFMT) -s -w $(shell find . -path ./vendor -prune -o -name '*.go' -print) | grep '^'

install: $(GOPATH)/bin/spl

lint: $(GOFILES) | $(GOLANGCI_LINT)
	@echo ">> linting code"
	@GO111MODULE=on $(GOLANGCI_LINT) run

msan: $(GOFILES)
	@echo ">> running memory sanitizer test"
	@$(CGO) test -mod=vendor -short -msan ./...

race: $(GOFILES)
	@echo ">> running race detector test"
	@$(CGO) test -mod=vendor -short -race ./...

test: $(GOFILES)
	@echo ">> running tests"
	@$(GO) test -mod=vendor -short -covermode=count -coverprofile=$(COVERPROFILE) ./...

tools: | $(GOLANGCI_LINT)

# BUILD TARGETS

.build/spl-$(GOOS)-$(GOARCH): $(GOFILES)
	@echo ">> building spl production binary for $(GOOS)/$(GOARCH)"
	@GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO) build $(GOFLAGS) -o .build/spl-$(GOOS)-$(GOARCH) ./cmd/spl

# INSTALL TARGETS

$(GOPATH)/bin/spl: $(GOFILES)
	@echo ">> installing spl binary"
	@$(GO) install -mod=vendor $(GOFLAGS) ./cmd/spl

# TEST TARGETS

$(COVERPROFILE):
	@make test

# TOOLS

$(GOLANGCI_LINT):
	@echo ">> installing golangci-lint"
	@$(GO) install -mod=vendor github.com/golangci/golangci-lint/cmd/golangci-lint
