COVEROUT ?= coverage.html

.PHONY: test

test:
	go test -failfast -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o ${COVEROUT}

build:
	go build -o tripleworks main.go

run: build
	./tripleworks
