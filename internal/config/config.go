package config

import (
	"fmt"
	"os"
)

type Postgres struct {
	Host     string
	User     string
	Pass     string
	Port     string
	Database string
}

var PostgresConfig Postgres

func LoadPostgresConfig() {
	PostgresConfig.Host = os.Getenv("DB_HOST")
	PostgresConfig.User = os.Getenv("DB_USER")
	PostgresConfig.Pass = os.Getenv("DB_PASSWORD")
	PostgresConfig.Port = os.Getenv("DB_PORT")
	PostgresConfig.Database = os.Getenv("DB_DATABASE")
}

func GetDBUrl() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		PostgresConfig.User,
		PostgresConfig.Pass,
		PostgresConfig.Host,
		PostgresConfig.Port,
		PostgresConfig.Database)
}
