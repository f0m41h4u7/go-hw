package config

type LogLevel string

const (
	Debug      LogLevel = "debug"
	Production          = "production"
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
	Sql      bool
	Database struct {
		Host     string
		Port     string
		User     string
		Password string
	}
}
