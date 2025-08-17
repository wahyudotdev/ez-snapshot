package backup

import "context"

type Repository interface {
	Dump(ctx context.Context, opts ...DumpDbOpts) (string, error)
}
