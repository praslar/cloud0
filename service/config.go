package service

import (
	"gitlab.com/goxp/cloud0/db"
)

type AppConfig struct {
	Port          int  `env:"PORT" envDefault:"8088"`
	DebugPort     int  `env:"DEBUG_PORT" envDefault:"7070"`
	ReadTimeout   int  `env:"READ_TIMEOUT" envDefault:"15"`
	EnableProfile bool `env:"ENABLE_PROFILE" envDefault:"true"` // enable profile listener

	db.Config
}
