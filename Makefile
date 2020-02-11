TAG := 'latest' # $(shell git describe --tag)

.PHONY: build push gobuild

build:
	@echo Build $(TAG)
	docker build --build-arg revision=$(TAG) -t balerter/balerter:$(TAG) -f Dockerfile .
push:
	@echo Push $(TAG)
	docker push balerter/balerter:$(TAG)
gobuild:
	@echo Go Build
	go build -o ./.debug/balerter -ldflags "-X main.revision=${TAG} -s -w" ./cmd/balerter