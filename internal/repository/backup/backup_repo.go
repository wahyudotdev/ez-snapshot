package backup

import "context"

type Repository interface {
	Dump(ctx context.Context) (string, error)
}
