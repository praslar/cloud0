package db

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	// dbDefault presents current active database,
	// should be initialized on starting app by calling MustOpenDefault or OpenDefault
	dbDefault *gorm.DB
)

// Config presents configuration that's necessary to work with database
type Config struct {
	Driver          string `env:"DB_DRIVER" envDefault:"postgres"`
	DSN             string `env:"DB_DSN"`
	MaxOpenConns    int    `env:"DB_MAX_OPEN_CONNS" envDefault:"25"`
	MaxIdleConns    int    `env:"DB_MAX_IDLE_CONNS" envDefault:"25"`
	ConnMaxLifetime int    `env:"DB_CONN_MAX_LIFETIME" envDefault:"600"`

	Host   string `env:"DB_HOST"`
	Port   string `env:"DB_PORT" envDefault:"5432"`
	User   string `env:"DB_USER"`
	Pass   string `env:"DB_PASS"`
	Name   string `env:"DB_NAME"`
	Schema string `env:"DB_SCHEMA" envDefault:"public"`
	Tz     string `env:"DB_TZ" envDefault:"UTC"`
}

// GetDSN returns a dsn that is read from ENV or built from separated env DB_*
func (c Config) GetDSN() string {
	if c.DSN != "" {
		return c.DSN
	}

	c.DSN = fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable connect_timeout=5 TimeZone=%s",
		c.Host,
		c.Port,
		c.User,
		c.Name,
		c.Pass,
		c.Tz,
	)

	return c.DSN
}

// Open open a DB connection
//  dbDefault, err := Open(config)
func Open(config *Config) (*gorm.DB, error) {

	var dialector gorm.Dialector
	switch config.Driver {
	case "sqlite", "sqlite3":
		dialector = sqlite.Open(config.GetDSN())
	case "postgres":
		dialector = postgres.Open(config.GetDSN())
	default:
		return nil, fmt.Errorf("unsupported driver %s", config.Driver)
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
			TablePrefix:   config.Schema + ".",
		},
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	theDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if config.MaxIdleConns > 0 {
		theDB.SetMaxIdleConns(config.MaxIdleConns)
	}
	if config.MaxOpenConns > 0 {
		theDB.SetMaxOpenConns(config.MaxOpenConns)
	}
	if config.ConnMaxLifetime > 0 {
		theDB.SetConnMaxLifetime(time.Duration(config.ConnMaxLifetime) * time.Second)
	}

	return db, nil
}

// Close release a DB instance
func Close(db *gorm.DB) {
	if db != nil {
		if dbInstance, err := db.DB(); err == nil {
			_ = dbInstance.Close()
		}
	}
}

// inMemorySqliteCfg presents configuration for quick testing
// this is lightweight database, should consider to user a real DB
// in more advanced testing like concurrency writing
var inMemorySqliteCfg = &Config{
	Driver:          "sqlite3",
	DSN:             ":memory:",
	MaxOpenConns:    1, // should be 1 cuz sqlite doesn't support concurrency writing operation.
	MaxIdleConns:    1,
	ConnMaxLifetime: 600,
}

// MustSetupTest setups an in-memory DB for testing and set to default
// it'll panic if errors occur
func MustSetupTest() {
	cfg := new(Config)
	*cfg = *inMemorySqliteCfg

	db, err := Open(cfg)
	if err != nil {
		panic(err)
	}

	dbDefault = db
}

// GetDB gets default database connection
func GetDB() *gorm.DB {
	if dbDefault == nil {
		panic(errors.New("uninitialized database. Please connect first"))
	}
	return dbDefault
}

// OpenDefault opens default database connection and assign to default
func OpenDefault(config *Config) error {
	db, err := Open(config)
	if err != nil {
		return err
	}
	dbDefault = db

	return nil
}

// MustOpenDefault open connection & assign to dbDefault, this will panic application if failed
func MustOpenDefault(config *Config) {
	if err := OpenDefault(config); err != nil {
		panic(err)
	}
}

// CloseDB closes default database
func CloseDB() {
	Close(dbDefault)
	dbDefault = nil
}
