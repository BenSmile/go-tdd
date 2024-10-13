@PHONY: run
run:
	@go run ./cmd/api


@PHONY: test
test:
	@GOFLAG="-count=1" go test -v -cover -race ./...
