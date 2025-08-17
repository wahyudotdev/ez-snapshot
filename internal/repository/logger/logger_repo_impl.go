package logger

import (
	"context"
	"ez-snapshot/internal/entity"
)

type impl struct {
}

func New() Repository {
	return &impl{}
}

func (r impl) SaveBackupLog(ctx context.Context, backup *entity.Backup) (*entity.Backup, error) {
	return nil, nil
}
