.PHONY: run build tidy

## run: start the API server locally
run:
	go run ./cmd/api

## build: compile to bin/stories-go
build:
	go build -o bin/stories-go ./cmd/api

## tidy: sync go.mod / go.sum
tidy:
	go mod tidy

## dev: tidy + run
dev: tidy run
