package repository

import (
	"context"
	"safebank/internal/domain"
	"safebank/internal/repository/dao"
)

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
	ErrUserNotFound       = dao.ErrUserNotFound
)

type UserRepository struct {
	dao *dao.UserDAO
}

func NewUserRepository(dao *dao.UserDAO) *UserRepository {
	return &UserRepository{
		dao: dao,
	}
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		ID:       u.ID,
		Email:    u.Email,
		Password: u.Password,
	}, nil
}

func (r *UserRepository) FindUserByUserID(ctx context.Context, userID int64) (domain.User, error) {
	u, err := r.dao.FindById(ctx, userID)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		ID:      u.ID,
		Email:   u.Email,
		FName:   u.FName,
		LName:   u.LName,
		Country: u.Country,

		State:  u.State,
		City:   u.City,
		Street: u.Street,
		Apart:  u.Apart,
		Zip:    u.Zip,
	}, nil
}

func (r *UserRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
		FName:    u.FName,
		LName:    u.LName,
		Country:  u.Country,

		State:  u.State,
		City:   u.City,
		Street: u.Street,
		Apart:  u.Apart,
		Zip:    u.Zip,
	})
}
