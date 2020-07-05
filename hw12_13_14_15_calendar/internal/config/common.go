package config

import "errors"

const (
	Debug LogLevel = "debug"
	Error LogLevel = "error"
	Info  LogLevel = "info"
	Warn  LogLevel = "warn"
)

var (
	ErrWrongDBHost       = errors.New("wrong database host")
	ErrWrongRMQHost      = errors.New("wrong rabbitmq host")
	ErrCannotReadConfig  = errors.New("cannot read config file")
	ErrCannotParseConfig = errors.New("cannot parse config file")
	ErrWrongServerHost   = errors.New("wrong server host")
	ErrWrongFile         = errors.New("wrong log file")
)

type DBConfiguration struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type ServerConfiguration struct {
	Host string
	Port string
}

type LogConfiguration struct {
	File  string
	Level LogLevel
}
