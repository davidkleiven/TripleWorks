COVEROUT ?= coverage.html

.PHONY: build-templates e2e-cleanup release

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
	go-semantic-release --provider github --token $$(gh auth token) --provider-opt "slug=davidkleiven/Tripleworks"

e2e: build
	rm -f tripleworks-e2e.db
	TRIPLE_WORKS_CONFIG="e2e_sqlite" TRIPLEWORKS_E2E="1" ./tripleworks

e2e-cleanup:
	rm tripleworks-e2e.db
