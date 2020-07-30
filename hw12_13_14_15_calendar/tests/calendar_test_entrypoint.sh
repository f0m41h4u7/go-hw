#!/bin/sh
go get -u github.com/pressly/goose/cmd/goose
goose -dir ./migrations/ mysql "qwerty:pswd@(db:3306)/default?charset=utf8&parseTime=True&loc=Local" up
./main -config calendar_config.json
