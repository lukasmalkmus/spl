# TOOLCHAIN
GO				?= CGO_ENABLED=0 GOFLAGS=-mod=vendor GOBIN=$(CURDIR)/bin go
GO_BIN_IN_PATH	?= CGO_ENABLED=0 GOFLAGS=-mod=vendor go
GO_NO_VENDOR	?= CGO_ENABLED=0 GOBIN=$(CURDIR)/bin go
GOFMT			?= $(GO)fmt

# ENVIRONMENT
VERBOSE 	=
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
BENCHSTAT				:= bin/benchstat
GOLANGCI_LINT   		:= bin/golangci-lint
GOTESTSUM   		    := bin/gotestsum

# FLAGS
GOFLAGS			?= -buildmode=exe -tags=netgo -installsuffix=cgo -trimpath \
					-ldflags='-s -w -extldflags "-static" \
					-X $(MOD_NAME)/pkg/version.buildTime=$(BUILD_TIME) \
					-X $(MOD_NAME)/pkg/version.commit=$(COMMIT) \
					-X $(MOD_NAME)/pkg/version.release=$(RELEASE) \
					-X $(MOD_NAME)/pkg/version.user=$(USER)'
GOTESTSUM_FLAGS	?= --jsonfile tests.json --junitfile junit.xml
GO_TEST_FLAGS 	?= -race -coverprofile=$(COVERPROFILE)

# Enable verbose test output if explicitly set.
ifdef VERBOSE
	GOTESTSUM_FLAGS	+= --format=standard-verbose
endif

# MISC
COVERPROFILE	:= coverage.out
DIRTY			:= $(shell git diff-index --quiet HEAD || echo "untracked")

# FUNCS
# func go-list-pkg-sources(package)
go-list-pkg-sources = $(GO) list $(GOFLAGS) -f '{{ range $$index, $$filename := .GoFiles }}{{ $$.Dir }}/{{ $$filename }}{{end}}' $(1)
# func go-pkg-sourcefiles(package)
go-pkg-sourcefiles = $(shell $(call go-list-pkg-sources,$(strip $1)))

.PHONY: all
all: dep fmt lint test build config ## Runs dep, fmt, lint, test, build and config.

.PHONY: bench
bench: $(BENCHSTAT) ## Runs all benchmarks and compares the benchmark results of a dirty workspace to the ones of a clean workspace if present.
	@echo ">> running benchmarks"
	@mkdir -p benchmarks
ifneq ($(DIRTY),untracked)
	@$(GO) test -run='$^' -bench=. -benchmem -count=5 ./... > benchmarks/$(COMMIT).txt
	@$(BENCHSTAT) benchmarks/$(COMMIT).txt
else
	@$(GO) test -run='$^' -bench=. -benchmem -count=5 ./... > benchmarks/$(COMMIT)-dirty.txt
	@$(BENCHSTAT) benchmarks/$(COMMIT).txt benchmarks/$(COMMIT)-dirty.txt
endif

.PHONY: build
build: .build/spl-$(GOOS)-$(GOARCH) ## Builds all binaries.

.PHONY: clean
clean: ## Removes build and test artifacts.
	@echo ">> cleaning up artifacts"
	@rm -rf .build $(COVERPROFILE) tests.json junit.xml

.PHONY: config
config: build ## Generates default configuration files.
	@echo ">> generating default configuration file"
	@rm -rf configs/spl.toml
	@./.build/spl-$(GOOS)-$(GOARCH) config > configs/spl.toml

.PHONY: cover
cover: | $(COVERPROFILE) ## Calculates the code coverage score.
	@echo ">> calculating code coverage"
	@$(GO) tool cover -func=$(COVERPROFILE)

.PHONY: dep-clean
dep-clean: ## Removes obsolete dependencies.
	@echo ">> cleaning dependencies"
	@$(GO_NO_VENDOR) mod tidy

.PHONY: dep-upgrade
dep-upgrade: ## Upgrades all direct dependencies to their latest version.
	@echo ">> upgrading dependencies"
	@$(GO_NO_VENDOR) get $(shell $(GO) list -f '{{if not (or .Main .Indirect)}}{{.Path}}{{end}}' -m all)
	@make dep

.PHONY: dep
dep: dep-clean dep.stamp ## Installs and verifies dependencies and removes obsolete ones.

dep.stamp: $(GOMODDEPS)
	@echo ">> installing dependencies"
	@$(GO) mod download
	@$(GO) mod verify
	@$(GO) mod vendor
	@touch $@

.PHONY: fmt
fmt: ## Formats and simplifies the source code using `gofmt`.
	@echo ">> formatting code"
	@! $(GOFMT) -s -w $(shell find . -path ./vendor -prune -o -name '*.go' -print) | grep '^'

.PHONY: install
install: $(GOPATH)/bin/spl ## Install all binaries into the $GOPATH/bin directory.

.PHONY: lint
lint: | $(GOLANGCI_LINT) ## Lints the source code.
	@echo ">> linting code"
	@GO111MODULE=on $(GOLANGCI_LINT) run

.PHONY: test
test: $(GOTESTSUM) ## Runs all tests. Run with VERBOSE=1 to get verbose test output ('-v' flag).
	@echo ">> running tests"
	@$(GOTESTSUM) $(GOTESTSUM_FLAGS) -- $(GO_TEST_FLAGS) ./...

.PHONY: tools
tools: | $(BENCHSTAT) $(GOLANGCI_LINT) $(GOTESTSUM) ## Installs all tools into the projects local $GOBIN directory.

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

# BUILD TARGETS

.build/spl-$(GOOS)-$(GOARCH): $(GOMODDEPS) $(call go-pkg-sourcefiles, $(MOD_NAME)/cmd/spl)
	@echo ">> building spl production binary for $(GOOS)/$(GOARCH)"
	@GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO) build $(GOFLAGS) -o .build/spl-$(GOOS)-$(GOARCH) ./cmd/spl

# INSTALL TARGETS

$(GOPATH)/bin/spl: $(GOMODDEPS) $(call go-pkg-sourcefiles, $(MOD_NAME)/cmd/spl)
	@echo ">> installing spl binary"
	@$(GO_BIN_IN_PATH) install $(GOFLAGS) ./cmd/spl

# TEST TARGETS

$(COVERPROFILE):
	@make test

# TOOLS

$(BENCHSTAT): $(GOMODDEPS)
	@echo ">> installing benchstat"
	@$(GO) install golang.org/x/perf/cmd/benchstat

$(GOLANGCI_LINT): $(GOMODDEPS)
	@echo ">> installing golangci-lint"
	@$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint

$(GOTESTSUM): $(GOMODDEPS)
	@echo ">> installing gotestsum"
	@$(GO) install gotest.tools/gotestsum
