BIN := "./bin/"

build:
	go build -v -o $(BIN) ./cmd/image_previewer

run: build
	$(BIN)/image_previewer --config="./configs/config.toml"

up:
	docker-compose -f ./deployments/docker-compose.yaml up -d

down:
	docker-compose -f ./deployments/docker-compose.yaml down

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v2.3.0

lint: install-lint-deps
	golangci-lint run \
		--new-from-rev=origin/main \
        --config=.golangci.yml \
        --sort-results \
        --max-issues-per-linter=1000 \
        --max-same-issues=1000 \
        ./...


test:
	go test -race -v ./internal/...

integration-tests:
	@docker-compose -f ./test/e2e/docker-compose.test.yaml up -d
	@go test ./test/e2e/ -v; TEST_EXIT=$$?; \
	docker-compose -f ./test/e2e/docker-compose.test.yaml down -v; \
	if [ $$TEST_EXIT -ne 0 ]; then exit 1; else exit 0; fi

.PHONY: build run lint run up down test integration-tests
