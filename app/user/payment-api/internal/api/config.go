package api

import (
	"fmt"
	"os"
)

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	Params   string
}

func GetDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "3306"),
		User:     getEnv("DB_USER", "root"),
		Password: getEnv("DB_PASSWORD", ""),
		DBName:   getEnv("DB_NAME", "payments"),
		Params:   getEnv("DB_PARAMS", "parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci"),
	}
}

func (c DatabaseConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.DBName,
		c.Params,
	)
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}