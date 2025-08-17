package storage

import (
	"context"
	"io"
)

type Repository interface {
	Upload(ctx context.Context, key string, r io.Reader) (string, error)
	Download(ctx context.Context, key string) (io.ReadCloser, error)
	Delete(ctx context.Context, key string) error
}
