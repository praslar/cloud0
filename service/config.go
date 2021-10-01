package service

import (
	"gitlab.com/goxp/cloud0/db"
)

// AppConfig presents some basic app configuration
type AppConfig struct {
	Port          int      `env:"PORT" envDefault:"8088"`
	Env           string   `env:"ENV" envDefault:"stg"`
	DebugPort     int      `env:"DEBUG_PORT" envDefault:"7070"`
	ReadTimeout   int      `env:"READ_TIMEOUT" envDefault:"15"`
	EnableProfile bool     `env:"ENABLE_PROFILE" envDefault:"true"` // enable profile listener
	EnableDB      bool     `env:"ENABLE_DB" envDefault:"false"`
	TrustedProxy  []string `env:"TRUSTED_PROXY" envSeparator:"," envDefault:"127.0.0.1,10.0.0.0/8,192.168.0.0/16"`
	Debug         bool     `env:"DEBUG" envDefault:"false"`
	DB            *db.Config
}

func NewAppConfig() *AppConfig {
	return &AppConfig{
		DB: &db.Config{},
	}
}
