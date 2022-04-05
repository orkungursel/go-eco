lint:
	golangci-lint run

coverage:
	go test -cover ./...