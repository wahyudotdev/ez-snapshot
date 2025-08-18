package storage

import (
	"context"
	"ez-snapshot/internal/entity"
	"io"
)

type Repository interface {
	Upload(ctx context.Context, key string, r io.Reader) (string, error)
	Download(ctx context.Context, key string) (io.ReadCloser, error)
	Delete(ctx context.Context, key string) error
	List(ctx context.Context) ([]*entity.Backup, error)
}
