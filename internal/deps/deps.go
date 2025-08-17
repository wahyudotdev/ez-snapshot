package deps

import (
	"ez-snapshot/internal/config"
	"ez-snapshot/internal/repository/backup"
	"ez-snapshot/internal/repository/logger"
	"ez-snapshot/internal/repository/storage"
)

func NewBackupRepo() backup.Repository {
	cfg, err := config.LoadMySQLConfig()
	if err != nil {
		panic(err)
	}

	return backup.New(
		backup.WithDbType(backup.MYSQL),
		backup.WithDbHost(cfg.Host),
		backup.WithDbPort(cfg.Port),
		backup.WithDbUsername(cfg.Username),
		backup.WithDbPassword(cfg.Password),
		backup.WithDatabase(cfg.Database),
	)
}

func NewStorageRepo() storage.Repository {
	return storage.New()
}

func NewLoggerRepo() logger.Repository {
	return logger.New()
}
