package config

import "os"

type LoggerConfig struct {
	Level string
}

func NewLoggerConfig() LoggerConfig {
	l := new(LoggerConfig)
	l.Level = os.Getenv("LOG_LEVEL")
	return *l
}
