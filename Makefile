PKG_PREFIX := github.com/balerter/balerter
TAG := latest # $(shell git describe --tag)

.PHONY: build-balerter
build-balerter: ## Build balerter docker image
	@echo Build Balerter $(TAG)
	docker build --build-arg version=$(TAG) -t balerter/balerter:$(TAG) -f ./contrib/balerter.Dockerfile .

.PHONY: push-balerter
push-balerter: ## Build balerter image to docker registry
	@echo Push Balerter $(TAG)
	docker push balerter/balerter:$(TAG)

.PHONY: gobuild-balerter
gobuild-balerter: ## Build balerter binary file
	@echo Go Build Balerter
	go build -o ./.debug/balerter -ldflags "-X main.revision=${TAG} -s -w" ./cmd/balerter

.PHONY: build-tgtool
build-tgtool: ## Build tgtool docker image
	@echo Build tgtool
	docker build -t balerter/tgtool:$(TAG) -f ./contrib/tgtool.Dockerfile .

.PHONY: push-tgtool
push-tgtool: ## Build tgtool image to docker registry
	@echo Push tgtool $(TAG)
	docker push balerter/tgtool:$(TAG)

.PHONY: test-full
test-full: ## Run full tests
	GO111MODULE=on go test -mod=vendor -coverprofile=coverage.txt -covermode=atomic ./internal/... ./cmd/...

.PHONY: test-integration
test-integration: ## Run integration tests
	go build -race -o ./balerter ./cmd/balerter
	docker-compose -f ./test/docker-compose.yml up -d
	sleep 5
	./balerter -config ./test/config.yml -once > out.txt
	diff out.txt ./test/out.etalon.txt
	docker-compose -f ./test/docker-compose.yml down -v


# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := build-balerter
