package config

import (
	"github.com/spf13/viper"
)

type RCloneConfig struct {
	Host   string // rclone host, eg : http://localhost:5572
	Fs     string // file system, eg : s3:my-aws-bucket
	Remote string // remote path
}

func LoadRCloneConfig() (*RCloneConfig, error) {
	viper.SetDefault("rclone.host", "http://localhost:5572")

	cfg := &RCloneConfig{
		Host:   viper.GetString("rclone.host"),
		Fs:     viper.GetString("rclone.fs"),
		Remote: viper.GetString("rclone.remote"),
	}

	return cfg, nil
}
