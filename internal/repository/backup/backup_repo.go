package backup

import (
	"context"
	"io"
)

type Repository interface {
	Dump(ctx context.Context) (string, error)
	Restore(ctx context.Context, reader io.ReadCloser) error
	DropAllTables(ctx context.Context) error
}
