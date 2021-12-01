export GOBIN ?= $(shell pwd)/bin

GOLINT = $(GOBIN)/golint
STATICCHECK = $(GOBIN)/staticcheck

GO_FILES := $(shell \
	find . '(' -path '*/.*' -o -path './vendor' ')' -prune \
	-o -name '*.go' -print | cut -b3-)

.PHONY: all
all: lint test

$(GOLINT):
	go install golang.org/x/lint/golint

$(STATICCHECK):
	go install honnef.co/go/tools/cmd/staticcheck

.PHONY: dependencies
dependencies:
	go mod download

.PHONY: lint
lint: $(GOLINT) $(STATICCHECK)
	@rm -rf lint.log
	@echo "Checking gofmt"
	@gofmt -e -s -l $(GO_FILES) 2>&1 | tee -a lint.log
	@echo "Checking govet"
	@go vet ./... 2>&1 | grep -v '^#' |  tee -a lint.log
	@echo "Checking golint"
	@$(GOLINT) ./... 2>&1 | tee -a lint.log
	@echo "Checking staticcheck"
	@$(STATICCHECK) ./... 2>&1 |  tee -a lint.log
	@echo "Checking for license headers"
	@scripts/check_license.sh | tee -a lint.log
	@[ ! -s lint.log ]

.PHONY: test
test: verifyversion
	go test -race ./...

.PHONY: cover
cover:
	go test -coverprofile=cover.out -coverpkg=./... -v ./...
	go tool cover -html=cover.out -o cover.html


.PHONY: verifyversion
verifyversion:
	$(eval CHANGELOG_VERSION := $(shell perl -ne '/^## (\S+)/ && print "$$1\n"' CHANGELOG.md | head -n1))
	$(eval INTHECODE_VERSION := $(shell perl -ne '/^const Version.*"([^"]+)".*$$/ && print "v$$1\n"' version.go))
	@if [ "$(INTHECODE_VERSION)" = "$(CHANGELOG_VERSION)" ]; then \
		echo "net/metrics: $(CHANGELOG_VERSION)"; \
	elif [ "$(CHANGELOG_VERSION)" = "vUnreleased" ]; then \
		echo "net/metrics (development): $(INTHECODE_VERSION)"; \
	else \
		echo "Version number in version.go does not match CHANGELOG.md"; \
		echo "version.go: $(INTHECODE_VERSION)"; \
		echo "CHANGELOG : $(CHANGELOG_VERSION)"; \
		exit 1; \
	fi

