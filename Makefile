.PHONY: lint test

lint:
	@which golangci-lint > /dev/null || (echo "golangci-lint is not installed. Please install it first:" && echo " visit: https://golangci-lint.run/welcome/install/" && exit 1)
	golangci-lint run -v ./...

test:
	go test -v ./...
