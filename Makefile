lint:
	golangci-lint run ./... --timeout=120s

test:
	go test ./...