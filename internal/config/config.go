package config

import (
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// Config ...
type Config struct {
	DBConfig *DBConfig
}

// DBConfig ...
type DBConfig struct {
	Dialect  string
	Name     string
	Username string
	Password string
	Host     string
	Port     string
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); len(value) != 0 {
		return value
	}

	return fallback
}

// GetConf ...
func GetConf() *Config {
	return &Config{
		DBConfig: &DBConfig{
			Dialect:  getenv("TODO_DB_DIALECT", "mysql"),
			Name:     getenv("TODO_DB_NAME", "1pJU3DlSp7"),
			Username: getenv("TODO_DB_USERNAME", "1pJU3DlSp7"),
			Password: getenv("TODO_DB_PASSWORD", "L276GsXLFa"),
			Host:     getenv("TODO_DB_HOST", "remotemysql.com"),
			Port:     getenv("TODO_DB_PORT", "3306"),
		},
	}
}
