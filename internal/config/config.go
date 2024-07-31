package config

import (
	"fmt"
	"os"
)

const (
	// names of required envs
	DB_HOST     = "DB_HOST"
	DB_PORT     = "DB_PORT"
	DB_USER     = "DB_USER"
	DB_PASSWORD = "DB_PASSWORD"
	DB_NAME     = "DB_NAME"
	REDIS_HOST  = "REDIS_HOST"
	REDIS_PORT  = "REDIS_PORT"
)

type DbConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type CacheConfig struct {
	Host string
	Port string
}

type Config struct {
	DB    DbConfig
	Cache CacheConfig
}

func LoadConfig() (*Config, error) {
	err := checkAllEnvVariables()
	if err != nil {
		return nil, err
	}

	return &Config{
		DB: DbConfig{
			Host:     os.Getenv(DB_HOST),
			Port:     os.Getenv(DB_PORT),
			User:     os.Getenv(DB_USER),
			Password: os.Getenv(DB_PASSWORD),
			Name:     os.Getenv(DB_NAME),
		},
		Cache: CacheConfig{
			Host: os.Getenv(REDIS_HOST),
			Port: os.Getenv(REDIS_PORT),
		},
	}, nil
}

// if at least one env params missing - error
func checkAllEnvVariables() error {
	required := []string{DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME, REDIS_HOST, REDIS_PORT}
	for _, str := range required {
		if _, ok := os.LookupEnv(str); !ok {
			return fmt.Errorf("%s is missing", str)
		}
	}
	return nil
}
