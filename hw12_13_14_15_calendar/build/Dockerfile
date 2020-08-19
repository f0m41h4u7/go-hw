FROM    golang:1.14.6 AS build
ARG     APP
ENV     GO111MODULE=on
WORKDIR /Calendar
COPY    go.mod .
COPY    go.sum .
RUN     go mod download
COPY    . .
RUN     go build ./cmd/$APP/main.go
