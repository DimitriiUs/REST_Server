package config

import (
	"fmt"
	"os"
	"sync"
)

type Postgres struct {
	Host     string
	User     string
	Pass     string
	Port     string
	Database string
	DBURL    string
}

type Config struct {
	Postgres Postgres
	LogLevel string
}

func GetConfig() Config {
var (
	cfg  Config
	once sync.Once
)
	once.Do(func() {
		cfg.Postgres.Host = os.Getenv("DB_HOST")
		cfg.Postgres.User = os.Getenv("DB_USER")
		cfg.Postgres.Pass = os.Getenv("DB_PASSWORD")
		cfg.Postgres.Port = os.Getenv("DB_PORT")
		cfg.Postgres.Database = os.Getenv("DB_DATABASE")
		cfg.Postgres.DBURL = fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s",
			cfg.Postgres.User,
			cfg.Postgres.Pass,
			cfg.Postgres.Host,
			cfg.Postgres.Port,
			cfg.Postgres.Database)
		cfg.LogLevel = os.Getenv("LOG_LEVEL")
	})

	return cfg
}
