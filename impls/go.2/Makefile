.PHONY: lint test

lint:
	@golangci-lint run ./src/...\
	  --max-issues-per-linter 0\
	  --max-same-issues 0\
	  --allow-parallel-runners

test:
	@go test ./src/...

fmt:
	@go fmt ./...
