package logger

import (
	"context"
	"ez-snapshot/internal/entity"
)

type Repository interface {
	SaveBackupLog(ctx context.Context, backup *entity.Backup) (*entity.Backup, error)
}
