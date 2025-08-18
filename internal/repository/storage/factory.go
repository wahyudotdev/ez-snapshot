package storage

import (
	"context"
	"ez-snapshot/internal/config"
)

func New(_ context.Context, cfg *config.RCloneConfig) Repository {
	return newRCloneImpl(cfg.Host, cfg.Fs, cfg.Remote)
}
