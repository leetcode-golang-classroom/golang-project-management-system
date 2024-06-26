.PHONY=build

build:
	@CGO_ENABLED=0 GOOS=linux go build -o bin/main cmd/main.go

run: build
	@./bin/main

test:
	@go test -v -cover ./test/...