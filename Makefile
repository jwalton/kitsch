# Adapted from https://github.com/vincentbernat/hellogopher

MODULE   = $(shell env GO111MODULE=on $(GO) list -m)
DATE    ?= $(shell date +%FT%T%z)
VERSION ?= $(shell git describe --tags --always --dirty --match=v* | cut -c2- 2> /dev/null || \
			cat $(CURDIR)/.version 2> /dev/null || echo v0)
COMMIT  ?= $(shell git rev-parse HEAD)
PKGS     = $(or $(PKG),$(shell env GO111MODULE=on $(GO) list ./...))
TESTPKGS = $(shell env GO111MODULE=on $(GO) list -f \
			'{{ if or .TestGoFiles .XTestGoFiles }}{{ .ImportPath }}{{ end }}' \
			$(PKGS))
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/)
RELEASE  = $(CURDIR)/dist

GO           = go
GORELEASER   = goreleaser
GOLINT       = golint
GOLANGCILINT = golangci-lint
GOCOV        = gocov
GOCOVXML     = gocov-xml
GO2XUNIT     = go2xunit

TIMEOUT = 15
V = 0
Q = $(if $(filter 1,$V),,@)
# Prompt shown before each item
M = $(shell printf "\033[34;1m▶\033[0m")

export GO111MODULE=on

.PHONY: all
all: generate fmt lint test build ## Format, test, and build

.PHONY: generate
generate: ; $(info $(M) generating source…) @ ## Generate source code
	$Q go generate ./...

.PHONY: build
build: ; $(info $(M) building executable…) @ ## Build program binary
	$Q $(GO) build \
		-tags release \
		-ldflags '-X $(MODULE)/cmd.version=$(VERSION) -X $(MODULE)/cmd.commit=$(COMMIT)'

.PHONY: goreleaser
goreleaser: ; $(info $(M) running…) @ ## Dry run of goreleaser
	$Q goreleaser build --skip-validate --rm-dist

.PHONY: install
install: ; $(info $(M) installing executable…) @ ## Install program binary
	$Q $(GO) install \
		-tags release \
		-ldflags '-X $(MODULE)/cmd.version=$(VERSION) -X $(MODULE)/cmd.commit=$(COMMIT)'

.PHONY: snapshot
snapshot: ; $(info $(M) creating snapshot…) @ ## Run goreleaser snapshot
	$Q goreleaser release --snapshot --rm-dist

.PHONY: docs
docs: build; $(info $(M) generating documentation…) @ ## Generate documentation
	$Q ./install.sh --dir .
	$Q cd docs && npm install && KITSCH=$(CURDIR)/kitsch npm run build
	$Q cp $(CURDIR)/install.sh $(CURDIR)/docs/build/install.sh
	$Q $(CURDIR)/kitsch version > $(CURDIR)/docs/build/latest-version.txt

# Tests

TEST_TARGETS := test-default test-bench test-short test-verbose test-race
.PHONY: $(TEST_TARGETS) test-xml check test tests
test-bench:   ARGS=-run=__absolutelynothing__ -bench=. ## Run benchmarks
test-short:   ARGS=-short        ## Run only short tests
test-verbose: ARGS=-v            ## Run tests in verbose mode with coverage reporting
test-race:    ARGS=-race         ## Run tests with race detector
$(TEST_TARGETS): NAME=$(MAKECMDGOALS:test-%=%)
$(TEST_TARGETS): test
test: lint ; $(info $(M) running $(NAME:%=% )tests…) @ ## Run tests
	$Q $(GO) test -timeout $(TIMEOUT)s $(ARGS) $(TESTPKGS)

test-ci: golint ; $(info $(M) running $(NAME:%=% )tests…) @ ## Run tests
	$Q $(GO) test -timeout $(TIMEOUT)s $(ARGS) $(TESTPKGS)

test-xml: lint | $(GO2XUNIT) ; $(info $(M) running xUnit tests…) @ ## Run tests with xUnit output
	$Q mkdir -p test
	$Q 2>&1 $(GO) test -timeout $(TIMEOUT)s -v $(TESTPKGS) | tee test/tests.output
	$(GO2XUNIT) -fail -input test/tests.output -output test/tests.xml

COVERAGE_MODE    = atomic
COVERAGE_PROFILE = $(COVERAGE_DIR)/profile.out
COVERAGE_XML     = $(COVERAGE_DIR)/coverage.xml
COVERAGE_HTML    = $(COVERAGE_DIR)/index.html
.PHONY: test-coverage test-coverage-tools
test-coverage-tools: | $(GOCOV) $(GOCOVXML)
test-coverage: COVERAGE_DIR := $(CURDIR)/test/coverage.$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
test-coverage: lint test-coverage-tools ; $(info $(M) running coverage tests…) @ ## Run coverage tests
	$Q mkdir -p $(COVERAGE_DIR)
	$Q $(GO) test \
		-coverpkg=$$($(GO) list -f '{{ join .Deps "\n" }}' $(TESTPKGS) | \
					grep '^$(MODULE)/' | \
					tr '\n' ',' | sed 's/,$$//') \
		-covermode=$(COVERAGE_MODE) \
		-coverprofile="$(COVERAGE_PROFILE)" $(TESTPKGS)
	$Q $(GO) tool cover -html=$(COVERAGE_PROFILE) -o $(COVERAGE_HTML)
	$Q $(GOCOV) convert $(COVERAGE_PROFILE) | $(GOCOVXML) > $(COVERAGE_XML)

.PHONY: lint
lint: golint golangci-lint ## Run linters

golint: | ; $(info $(M) running golint…) @ ## Run golint
	$Q $(GOLINT) -set_exit_status ./...

.PHONY: golangci-lint
golangci-lint: | ; $(info $(M) running golangci-lint…) @ ## Run golangci-lint
	$Q $(GOLANGCILINT) run

.PHONY: fmt
fmt: ; $(info $(M) running gofmt…) @ ## Run gofmt on all source files
	$Q $(GO) fmt $(PKGS)

# Misc

.PHONY: clean
clean: ; $(info $(M) cleaning…)	@ ## Cleanup everything
	@rm -rf dist $(basename $(MODULE))
	@rm -rf test/tests.* test/coverage.*

.PHONY: help
help:
	@grep -hE '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-17s\033[0m %s\n", $$1, $$2}'

.PHONY: version
version:
	@echo $(VERSION)
