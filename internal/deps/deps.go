package deps

import (
	"context"
	"ez-snapshot/internal/config"
	"ez-snapshot/internal/repository/backup"
	"ez-snapshot/internal/repository/storage"
)

func NewBackupRepo(_ context.Context) backup.Repository {
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

func NewStorageRepo(ctx context.Context) storage.Repository {
	cfg, err := config.LoadRCloneConfig()
	if err != nil {
		panic(err)
	}
	return storage.New(
		ctx,
		cfg,
	)
}
