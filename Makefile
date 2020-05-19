SHELL       =   /bin/sh
PKG_PREFIX  :=  github.com/balerter/balerter
TAG         ?=  latest

.SUFFIXES:
.PHONY: help \
	build-balerter push-balerter gobuild-balerter \
	build-tgtool push-tgtool \
	test-full test-integration

build-balerter: ## Build balerter docker image
	@echo Build Balerter $(TAG)
	docker build --build-arg version=$(TAG) -t balerter/balerter:$(TAG) -f ./contrib/balerter.Dockerfile .

push-balerter: ## Build balerter image to docker registry
	@echo Push Balerter $(TAG)
	docker push balerter/balerter:$(TAG)

gogenerate: ## Call go generate
	@echo Calling go generate
	go generate ./...

gobuild-balerter: ## Build balerter binary file
	@echo Go Build Balerter
	go build -o ./.debug/balerter -ldflags "-X main.revision=${TAG} -s -w" ./cmd/balerter

build-tgtool: ## Build tgtool docker image
	@echo Build tgtool
	docker build -t balerter/tgtool:$(TAG) -f ./contrib/tgtool.Dockerfile .

push-tgtool: ## Build tgtool image to docker registry
	@echo Push tgtool $(TAG)
	docker push balerter/tgtool:$(TAG)

test-full: ## Run full tests
	GO111MODULE=on go test -mod=vendor -coverprofile=coverage.txt -covermode=atomic ./internal/... ./cmd/...

test-integration: ## Run integration tests
	go build -race -o ./integration/balerter ./cmd/balerter
	go test ./integration

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
