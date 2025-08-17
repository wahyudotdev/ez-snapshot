package usecase

import (
	"context"
	"ez-snapshot/internal/entity"
	"ez-snapshot/internal/repository/backup"
	"ez-snapshot/internal/repository/logger"
	"ez-snapshot/internal/repository/storage"
	"os"
	"time"
)

type BackupDatabaseUseCase struct {
	backup  backup.Repository
	storage storage.Repository
	logger  logger.Repository
}

func NewBackupDatabaseUseCase(
	backup backup.Repository,
	storage storage.Repository,
	logger logger.Repository,
) *BackupDatabaseUseCase {
	return &BackupDatabaseUseCase{
		backup:  backup,
		storage: storage,
		logger:  logger,
	}
}

func (uc *BackupDatabaseUseCase) Execute(ctx context.Context) (*entity.Backup, error) {
	dumpPath, err := uc.backup.Dump(ctx)
	if err != nil {
		return nil, err
	}
	f, err := os.Open(dumpPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}

	uploadPath, err := uc.storage.Upload(ctx, stat.Name(), f)
	if err != nil {
		return nil, err
	}

	result, err := uc.logger.SaveBackupLog(ctx, &entity.Backup{
		CreatedAt: time.Now(),
		Name:      stat.Name(),
		Path:      uploadPath,
		Size:      stat.Size(),
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}
