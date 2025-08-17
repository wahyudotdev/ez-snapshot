package config

import (
	"fmt"
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
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")      // optionally look for config in the working directory

	// read config
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

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
