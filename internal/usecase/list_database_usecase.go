package usecase

import (
	"context"
	"ez-snapshot/internal/entity"
	"ez-snapshot/internal/repository/storage"
)

type ListDatabaseUseCase struct {
	storage storage.Repository
}

func NewListDatabaseUseCase(storage storage.Repository) ListDatabaseUseCase {
	return ListDatabaseUseCase{
		storage: storage,
	}
}

func (r ListDatabaseUseCase) Execute(ctx context.Context) ([]*entity.Backup, error) {
	list, err := r.storage.List(ctx)
	if err != nil {
		return nil, err
	}
	return list, nil
}
