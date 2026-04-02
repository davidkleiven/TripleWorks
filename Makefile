COVEROUT ?= coverage.html

.PHONY: build-templates

build-templates:
	templ generate

test: build-templates
	go test -failfast -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o ${COVEROUT}

build: build-templates
	go build -o tripleworks main.go

run: build
	./tripleworks

release:
	go-semantic-release --provider github --allow-initial-releases
