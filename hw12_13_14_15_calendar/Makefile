CALENDAR_CMD := calendar
SCHEDULER_CMD := scheduler
SENDER_CMD := sender

build:
	go build -o $(CALENDAR_CMD) ./cmd/calendar/main.go
	go build -o $(SCHEDULER_CMD) ./cmd/scheduler/main.go
	go build -o $(SENDER_CMD) ./cmd/sender/main.go

run:
	make build
	./$(CALENDAR_CMD)

up:
	docker-compose -f ./deployments/docker-compose.yml up #--build -d

down:
	docker-compose -f ./deployments/docker-compose.yml down

restart:
	make down
	make up

test:
	go test -v -count=1 -race -gcflags=-l -timeout=30s ./...

lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint
	golangci-lint run ./...

clean:
	rm $(CALENDAR_CMD) $(SCHEDULER_CMD) $(SENDER_CMD) *.out

migrate:
	goose -dir migrations/ mysql "qwerty:pswd@(127.0.0.1:1338)/default?charset=utf8&parseTime=True&loc=Local" up

gen:
	go generate ./...

bdd:
	set -e ;\
	docker-compose -f ./tests/docker-compose-test.yml up --build -d ;\
	test_status_code=0 ;\
	docker-compose -f ./tests/docker-compose-test.yml run integration_tests go test || test_status_code=$$? ;\
	docker-compose -f ./tests/docker-compose-test.yml down ;\
	exit $$test_status_code ;\

test-cleanup:
	docker-compose down \
		--rmi local \
		--volumes \
		--remove-orphans \
		--timeout 60; \
  	docker-compose rm -f

.PHONY: build run test lint clean migrate gen up down restart bdd test-cleanup
