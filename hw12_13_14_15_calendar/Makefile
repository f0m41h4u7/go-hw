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

migrate:
	goose -dir migrations/ mysql "qwerty:pswd@(127.0.0.1:1338)/default?charset=utf8&parseTime=True&loc=Local" up

gen:
	go generate ./...

.PHONY: build run test lint clean migrate generate
