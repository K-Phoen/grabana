WITH_COVERAGE?=false

ifeq ($(WITH_COVERAGE),true)
GOCMD_TEST?=go test -coverpkg=./... -coverprofile=coverage.txt -covermode=atomic ./...
else
GOCMD_TEST?=go test
endif

.PHONY: lint
lint:
	docker run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:v1.43.0 golangci-lint run -c .golangci.yaml

.PHONY: tests
tests:
	$(GOCMD_TEST) ./...

.PHONY: up
up:
	docker run -d \
        --name=grabana_prometheus \
        -p 9090:9090 \
        -v $(shell pwd)/testdata/prometheus.yml:/etc/prometheus/prometheus.yml \
        prom/prometheus
	docker run -d \
      -p 3000:3000 \
      --name=grabana_grafana \
      -e "GF_SECURITY_ADMIN_PASSWORD=secret" \
      grafana/grafana:8.3.2

.PHONY: up-all
up-all: up
	docker run -d\
		 --name grabana_graphite\
		 --restart=always\
		 -p 8081:80\
		 -p 2003-2004:2003-2004\
		 -p 2023-2024:2023-2024\
		 -p 8125:8125/udp\
		 -p 8126:8126\
		 graphiteapp/graphite-statsd
	docker run -d \
		-p 8086:8086 \
		--name=grabana_influxdb \
		-e DOCKER_INFLUXDB_INIT_MODE=setup \
		-e DOCKER_INFLUXDB_INIT_USERNAME=my-user \
		-e DOCKER_INFLUXDB_INIT_PASSWORD=my-password \
		-e DOCKER_INFLUXDB_INIT_ORG=my-org \
		-e DOCKER_INFLUXDB_INIT_BUCKET=my-bucket \
		-e DOCKER_INFLUXDB_INIT_RETENTION=1w \
		-e DOCKER_INFLUXDB_INIT_ADMIN_TOKEN=my-super-secret-auth-token \
        influxdb:2.0

.PHONY: down
down:
	docker rm -f grabana_grafana
	docker rm -f grabana_prometheus

.PHONY: down-all
down-all: down
	docker rm -f grabana_graphite
	docker rm -f grabana_influxdb

install_goreleaser:
	go install github.com/goreleaser/goreleaser@latest
