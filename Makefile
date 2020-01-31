WITH_COVERAGE?=false

ifeq ($(WITH_COVERAGE),true)
GOCMD_TEST?=go test -mod=vendor -coverpkg=./... -coverprofile=coverage.txt -covermode=atomic ./...
else
GOCMD_TEST?=go test -mod=vendor
endif

.PHONY: lint
lint:
	docker run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:v1.23.1 golangci-lint run -c .golangci.yaml

.PHONY: tests
tests:
	$(GOCMD_TEST) ./...
