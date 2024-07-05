.PHONY: default proto hackcenter all clean docker

GOBIN = $(shell pwd)/build/bin
TAG ?= latest
GOFILES_NOVENDOR := $(shell go list -f "{{.Dir}}" ./...)

default: hackcenter

all: proto hackcenter

hackcenter:
	go build $(BUILD_FLAGS) -o=${GOBIN}/$@ -gcflags "all=-N -l" .
	@echo "Done building."

proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./hackcenter/center.proto


clean:
	rm -fr build/*

docker:
	docker build -t hackcenter:${TAG} .
