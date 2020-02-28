TAG := 'latest' # $(shell git describe --tag)

.PHONY: build-balerter push-balerter gobuild-balerter

build-balerter:
	@echo Build Balerter $(TAG)
	docker build --build-arg revision=$(TAG) -t balerter/balerter:$(TAG) -f ./contrib/balerter.Dockerfile .
push-balerter:
	@echo Push Balerter $(TAG)
	docker push balerter/balerter:$(TAG)
gobuild-balerter:
	@echo Go Build Balerter
	go build -o ./.debug/balerter -ldflags "-X main.revision=${TAG} -s -w" ./cmd/balerter

build-tgtool:
	@echo Build tgtool
	docker build -t balerter/tgtool:$(TAG) -f ./contrib/tgtool.Dockerfile .
push-tgtool:
	@echo Push tgtool $(TAG)
	docker push balerter/tgtool:$(TAG)
