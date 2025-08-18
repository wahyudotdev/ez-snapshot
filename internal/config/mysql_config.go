package config

import (
	"github.com/spf13/viper"
)

type MySQLConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

func LoadMySQLConfig() (*MySQLConfig, error) {
	// set defaults (in case values are missing)
	viper.SetDefault("mysql.port", "3306")

	cfg := &MySQLConfig{
		Host:     viper.GetString("mysql.host"),
		Port:     viper.GetString("mysql.port"),
		Username: viper.GetString("mysql.username"),
		Password: viper.GetString("mysql.password"),
		Database: viper.GetString("mysql.database"),
	}

	return cfg, nil
}
