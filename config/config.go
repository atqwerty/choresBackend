package config

import "os"

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
			Name:     getenv("TODO_DB_NAME", "choresdb"),
			Username: getenv("TODO_DB_USERNAME", "root"),
			Password: getenv("TODO_DB_PASSWORD", "Savanna1"),
			Host:     getenv("TODO_DB_HOST", "localhost"),
			Port:     getenv("TODO_DB_PORT", "3306"),
		},
	}
}
