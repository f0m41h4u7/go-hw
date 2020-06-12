CMD := calendar

build:
	go build -o $(CMD) ./cmd/calendar/main.go

run:
	make build
	./$(CMD)

test:
	go test -v -count=1 -race -gcflags=-l -timeout=30s ./...

lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint
	golangci-lint run ./...

clean:
	rm $(CMD)

.PHONY: build run test lint clean
