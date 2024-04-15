package repository

import (
	"context"
	"safebank/internal/domain"
	"safebank/internal/repository/dao"
)

type SavingRepository struct {
	dao *dao.SavingDAO
}

func (r *SavingRepository) Create(ctx context.Context, u domain.Saving) error {
	return r.dao.Insert(ctx, dao.Saving{})
}
