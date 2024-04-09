package repository

import "safebank/internal/repository/dao"

type AccRepository struct {
	dao *dao.AccountDAO
}
