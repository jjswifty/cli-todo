.PHONY: run fmt lint lint-fix

run:
	go run .

fmt:
	golangci-lint fmt

lint:
	golangci-lint run

lint-fix:
	golangci-lint run --fix

