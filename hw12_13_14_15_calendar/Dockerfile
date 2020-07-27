FROM    golang:1.14.6-alpine3.12 AS build
#ARG     APP
WORKDIR /Calendar
COPY    . .
RUN     go build ./cmd/calendar/main.go && pwd && ls
