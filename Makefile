SHELL       =   /bin/sh
PKG_PREFIX  :=  github.com/balerter/balerter
TAG         ?=  dev

.SUFFIXES:
.PHONY: help \
	build-dev push-dev \
	build push gobuild-balerter \
	build-tgtool push-tgtool \
	test-full

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build-dev: ## Build balerter/balerter and balerter/test docker images DEV version
	@echo Build Balerter dev
	docker build --build-arg version=dev -t ghcr.io/balerter/balerter:dev -t balerter/balerter:dev -f ./contrib/balerter.Dockerfile .
	docker build --build-arg version=dev -t ghcr.io/balerter/test:dev -t balerter/test:dev -f ./contrib/test.Dockerfile .

push-dev: ## Build balerter/balerter and balerter/test images to docker registry DEV version
	@echo Push Balerter dev
	docker push balerter/balerter:dev
	docker push ghcr.io/balerter/balerter:dev
	docker push balerter/test:dev
	docker push ghcr.io/balerter/test:dev

build: ## Build balerter/balerter and balerter/test docker images
	@echo Build Balerter $(TAG)
	docker build --build-arg version=$(TAG) -t ghcr.io/balerter/balerter:$(TAG) -t ghcr.io/balerter/balerter:latest -t balerter/balerter:$(TAG) -t balerter/balerter:latest -f ./contrib/balerter.Dockerfile .
	docker build --build-arg version=$(TAG) -t ghcr.io/balerter/test:$(TAG) -t ghcr.io/balerter/test:latest -t balerter/test:$(TAG) -t balerter/test:latest -f ./contrib/test.Dockerfile .

push: ## Build balerter/balerter and balerter/test images to docker registry
	@echo Push Balerter $(TAG)
	docker push balerter/balerter:$(TAG)
	docker push balerter/balerter:latest
	docker push ghcr.io/balerter/balerter:$(TAG)
	docker push ghcr.io/balerter/balerter:latest
	docker push balerter/test:$(TAG)
	docker push balerter/test:latest
	docker push ghcr.io/balerter/test:$(TAG)
	docker push ghcr.io/balerter/test:latest

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
	docker-compose up -d
	GO111MODULE=on go test -mod=vendor -coverprofile=coverage.txt -covermode=atomic ./internal/... ./cmd/...
	docker-compose down -v
