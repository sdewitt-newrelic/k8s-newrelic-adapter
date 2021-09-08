REGISTRY?=ghcr.io/kidk
IMAGE?=k8s-newrelic-adapter
OUT_DIR?=./_output
VENDOR_DOCKERIZED?=0

VERSION?:=latest
GOIMAGE=golang:1.13
GOFLAGS=-mod=vendor -tags=netgo

.PHONY: all docker-build push test build-local-image

all: test $(OUT_DIR)/adapter

src_deps=$(shell find pkg cmd -type f -name "*.go")
$(OUT_DIR)/adapter: $(src_deps)
	CGO_ENABLED=0 GOARCH=$* go build $(GOFLAGS) -o $(OUT_DIR)/$*/adapter cmd/adapter/adapter.go

build: test
	mkdir build/
	cp deploy/Dockerfile build/Dockerfile
	CGO_ENABLED=0 GO111MODULE=on go build -o build/adapter cmd/adapter/adapter.go

image:
	docker build -t $(REGISTRY)/$(IMAGE):$(VERSION) build/

push:
	docker push $(REGISTRY)/$(IMAGE):$(VERSION)

vendor: go.mod
ifeq ($(VENDOR_DOCKERIZED),1)
	docker run -it -v $(shell pwd):/src/k8s-newrelic-adapter -w /src/k8s-newrelic-adapter $(GOIMAGE) /bin/bash -c "\
		go mod vendor"
else
	go mod vendor
endif

test:
	CGO_ENABLED=0 GO111MODULE=on go test ./pkg/controller/...

clean:
	rm -rf ${OUT_DIR} vendor build

# Code gen helpers
gen-apis: vendor
	hack/update-codegen.sh

verify-apis: vendor
	hack/verify-codegen.sh
