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
)

type DbConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type Config struct {
	DB DbConfig
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
	}, nil
}

// if at least one env params missing - error
func checkAllEnvVariables() error {
	if _, ok := os.LookupEnv(DB_HOST); !ok {
		return fmt.Errorf("%s is missing", DB_HOST)
	}

	if _, ok := os.LookupEnv(DB_PORT); !ok {
		return fmt.Errorf("%s is missing", DB_PORT)
	}

	if _, ok := os.LookupEnv(DB_USER); !ok {
		return fmt.Errorf("%s is missing", DB_USER)
	}

	if _, ok := os.LookupEnv(DB_PASSWORD); !ok {
		return fmt.Errorf("%s is missing", DB_PASSWORD)
	}

	if _, ok := os.LookupEnv(DB_NAME); !ok {
		return fmt.Errorf("%s is missing", DB_NAME)
	}

	return nil
}
