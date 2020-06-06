CMD := calendar

build:
	go build -o $(CMD) .

run:
	make build
	./$(CMD)

test:
	go test -v -count=1 -race -gcflags=-l -timeout=30s ./tests

lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint
	golangci-lint run ./...

clean:
	rm $(CMD)

.PHONY: build run test lint clean
