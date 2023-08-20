.PHONY: lint test drone

lint:
	test -z $$(gofmt -s -l .)
	go vet ./...

test:
	go test -v ./...

drone:
	drone starlark --format
	drone lint .drone.yml --trusted
	drone --server https://drone.grafana.net sign --save grafana/kindsys
