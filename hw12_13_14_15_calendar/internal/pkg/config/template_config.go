package config

import (
	"errors"
	"net"
)

type LogLevel string

var (
	ErrWrongServerHost = errors.New("wrong server host")
	ErrWrongDBHost     = errors.New("wrong database host")
	ErrWrongFile       = errors.New("wrong log file")
)

const (
	Debug LogLevel = "debug"
	Error LogLevel = "error"
	Info  LogLevel = "info"
	Warn  LogLevel = "warn"
)

type Config struct {
	Server struct {
		Host string
		Port string
	}
	Logger struct {
		File  string
		Level LogLevel
	}
	SQL      bool
	Database struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
	}
}

func (c Config) validate() error {
	if _, err := net.LookupHost(c.Server.Host); err != nil {
		return ErrWrongServerHost
	}

	if c.SQL {
		if _, err := net.LookupHost(c.Database.Host); err != nil {
			return ErrWrongDBHost
		}
	}

	return nil
}
