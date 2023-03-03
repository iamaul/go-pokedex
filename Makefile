.PHONY: docker local test

local:
	echo "Starting local environment"
	docker-compose -f docker-compose.yml up --build


run-linter:
	echo "Starting linters"
	golangci-lint run ./...

run:
	go run ./cmd/app/main.go

build:
	go build ./cmd/app/main.go

test:
	go test -cover ./...


deps-reset:
	git checkout -- go.mod
	go mod tidy
	go mod vendor

tidy:
	go mod tidy
	go mod vendor

deps-upgrade:
	# go get $(go list -f '{{if not (or .Main .Indirect)}}{{.Path}}{{end}}' -m all)
	go get -u -t -d -v ./...
	go mod tidy
	go mod vendor

deps-cleancache:
	go clean -modcache

FILES := $(shell docker ps -aq)

down-local:
	docker stop $(FILES)
	docker rm $(FILES)

clean:
	docker system prune -f

logs-local:
	docker logs -f $(FILES)