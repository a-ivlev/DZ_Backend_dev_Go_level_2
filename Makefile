
lint:
	golangci-lint run -c golangci-lint.yml

test:
	go test ./...