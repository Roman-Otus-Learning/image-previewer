BIN="./bin/image-previewer"
DOCKER_IMG="image-previewer:develop"

build:
	go build -v -o $(BIN) ./cmd

run: build
	$(BIN) --config ./configs/config.yaml image-previewer


install-linter:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.50.1

lint: install-linter
	golangci-lint run ./...

lint-fix:
	golangci-lint run ./... --fix

build-img:
	docker build \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

test:
	go test -race ./internal/...

.PHONY: build run install-linter lint lint-fix build-img