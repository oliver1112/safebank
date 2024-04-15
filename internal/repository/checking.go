package repository

import (
	"context"
	"safebank/internal/domain"
	"safebank/internal/repository/dao"
)

type CheckingRepository struct {
	dao *dao.CheckingDAO
}

func (r *CheckingRepository) Create(ctx context.Context, u domain.Checking) error {
	return r.dao.Insert(ctx, dao.Checking{})
}
