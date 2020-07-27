#!/bin/sh
go get -u github.com/pressly/goose/cmd/goose
goose -dir ./migrations/ mysql "$MYSQL_USER:$MYSQL_PASSWORD@($MYSQL_HOST:$MYSQL_PORT)/$MYSQL_DATABASE?charset=utf8&parseTime=True&loc=Local" up
./main -config calendar_config.json
