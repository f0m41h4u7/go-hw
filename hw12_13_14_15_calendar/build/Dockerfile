FROM    golang:1.14.6 AS build
ARG     APP
ENV     GO111MODULE=on
WORKDIR /Calendar
COPY    . .
#RUN     apk update && apk upgrade && apk add gcc
RUN     go build ./cmd/$APP/main.go
