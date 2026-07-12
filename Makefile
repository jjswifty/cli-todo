.PHONY: run fmt lint lint-fix test cover cover-html

run:
	go run .

fmt:
	golangci-lint fmt

lint:
	golangci-lint run

lint-fix:
	golangci-lint run --fix

test:
	go test ./...

cover:
	go test -coverprofile=cover.out ./...
	go tool cover -func=cover.out

cover-html: cover
	go tool cover -html=cover.out
