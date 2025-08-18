package usecase

import (
	"context"
	"ez-snapshot/internal/repository/backup"
	"ez-snapshot/internal/repository/storage"
	"os"
)

type BackupDatabaseUseCase struct {
	backup  backup.Repository
	storage storage.Repository
}

func NewBackupDatabaseUseCase(
	backup backup.Repository,
	storage storage.Repository,
) *BackupDatabaseUseCase {
	return &BackupDatabaseUseCase{
		backup:  backup,
		storage: storage,
	}
}

func (uc *BackupDatabaseUseCase) Execute(ctx context.Context) error {
	dumpPath, err := uc.backup.Dump(ctx)
	if err != nil {
		return err
	}
	f, err := os.Open(dumpPath)
	if err != nil {
		return err
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return err
	}

	_, err = uc.storage.Upload(ctx, stat.Name(), f)
	if err != nil {
		return err
	}

	return nil
}
