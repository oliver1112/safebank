package repository

import (
	"context"
	"safebank/internal/domain"
	"safebank/internal/repository/dao"
)

type AccRepository struct {
	dao *dao.AccountDAO
}

func (r *AccRepository) Create(ctx context.Context, u domain.User) error {
	return nil
}
