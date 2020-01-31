.PHONY: lint
lint:
	docker run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:v1.23.1 golangci-lint run -c .golangci.yaml

.PHONY: tests
tests:
	go test -mod=vendor ./...
