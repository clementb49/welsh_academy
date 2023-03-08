// Package config contains the configuration settings for the application.
// It imports packages for the Gin web framework, the Viper configuration library, and the Zap logging library.
package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// DbConfig is the struct that store welsh academy database configuration.
type DbConfig struct {
	Host     string
	User     string
	Password string
	DbName   string
	Port     int
	SslMode  string
	TimeZone string
}

// DbConfigFromViper constructs a new DbConfig instance from viper configuration settings
func DbConfigFromViper() *DbConfig {
	// Set default values for optional settings
	viper.SetDefault("dbPort", 5432)
	viper.SetDefault("dbSslMode", "disable")
	viper.SetDefault("timezone", "UTC")
	// Build and return the DbConfig instance
	return &DbConfig{
		Host:     viper.GetString("dbHost"),
		User:     viper.GetString("dbUser"),
		Password: viper.GetString("dbPassword"),
		DbName:   viper.GetString("dbName"),
		Port:     viper.GetInt("dbPort"),
		TimeZone: viper.GetString("timeZone"),
		SslMode:  viper.GetString("dbSslMode"),
	}
}

// BuildDsn builds a PostgreSQL DSN string from the DbConfig instance
func (cfg *DbConfig) BuildDsn() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s timezone=%s", cfg.Host, cfg.User, cfg.Password, cfg.DbName, cfg.Port, cfg.SslMode, cfg.TimeZone)
}
