VERSION := "v0.1.0"
REVISION := $(shell git rev-parse --short HEAD)
LDFLAGS := -X 'main.version=$(VERSION)' \
        -X 'main.revision=$(REVISION)'

.PHONY: docker-build
docker-build:
	docker build \
        --build-arg VERSION=$(VERSION) \
        -t ohsawa0515/gpu-cloudwatch-reporting:$(VERSION) ./

.PHONY: docker-push
docker-push:
	docker push ohsawa0515/gpu-cloudwatch-reporting:$(VERSION)
