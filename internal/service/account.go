package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"safebank/internal/domain"
	"safebank/internal/lib"
	"safebank/internal/repository"
	"safebank/internal/repository/dao"
)

type AccountService struct {
	UserRepo       *repository.UserRepository
	AccountDao     *dao.AccountDAO
	CheckingDao    *dao.CheckingDAO
	HomeLoanDao    *dao.HomeLoanDAO
	LoanDao        *dao.LoanDAO
	SavingDao      *dao.SavingDAO
	StuLoanDao     *dao.StuLoanDAO
	InstituteDao   *dao.InstituteDAO
	TransactionDao *dao.TransactionDAO
}

func NewAccountService(db *gorm.DB) *AccountService {
	// init repo
	userDao := dao.NewUserDao(db)
	userRepo := repository.NewUserRepository(userDao)

	accountDao := dao.NewAccountDao(db)
	checkingDao := dao.NewCheckingDao(db)
	homeLoanDao := dao.NewHomeLoanDao(db)
	loanDao := dao.NewLoanDao(db)
	savingDao := dao.NewSavingDao(db)
	stuLoanDao := dao.NewStuLoanDao(db)
	instituteDao := dao.NewInstituteDao(db)
	transactionDao := dao.NewTransactionDao(db)

	return &AccountService{
		UserRepo:       userRepo,
		AccountDao:     accountDao,
		CheckingDao:    checkingDao,
		HomeLoanDao:    homeLoanDao,
		LoanDao:        loanDao,
		SavingDao:      savingDao,
		StuLoanDao:     stuLoanDao,
		InstituteDao:   instituteDao,
		TransactionDao: transactionDao,
	}
}

func (a *AccountService) GetAccount(ctx *gin.Context) ([]dao.Account, error) {
	userID, _ := ctx.Get("userID")
	accounts, err := a.AccountDao.GetAccountList(ctx, cast.ToInt64(userID))
	return accounts, err
}

func (a *AccountService) CreateOrUpdateAccount(ctx *gin.Context, data interface{}) ([]dao.Account, error) {
	id, _ := ctx.Get("userID")
	accounts, err := a.AccountDao.GetAccountList(ctx, cast.ToInt64(id))
	return accounts, err
}

func (a *AccountService) GetAccountsByUserID(ctx *gin.Context, userID int64) (domain.UserCenter, error) {
	var info domain.UserCenter

	userInfo, err := a.UserRepo.FindUserByUserID(ctx, userID)
	if err != nil {
		return info, fmt.Errorf("system error 1")
	}

	accountList, err := a.AccountDao.GetAccountList(ctx, userID)
	if err != nil {
		return info, fmt.Errorf("system error 2")
	}

	var data domain.AccountData

	for _, account := range accountList {
		if account.ID <= 0 {
			continue
		}
		accountDetail := make(map[string]interface{})

		if account.AccountType == "S" {
			savingData, err := a.SavingDao.GetSaving(ctx, account.ID)
			if err != nil {
				return info, fmt.Errorf("system error 3")
			}
			savingData.Account = account
			err = lib.StructToMapSingleD2(savingData, "json", &accountDetail)
			if err != nil {
				return info, fmt.Errorf("system error 4")
			}
			data.SavingAccount = accountDetail

		} else if account.AccountType == "C" {
			checkingData, err := a.CheckingDao.GetChecking(ctx, account.ID)
			if err != nil {
				return info, fmt.Errorf("system error 5")
			}
			checkingData.Account = account
			err = lib.StructToMapSingleD2(checkingData, "json", &accountDetail)
			if err != nil {
				return info, fmt.Errorf("system error 6")
			}
			data.CheckingAccount = accountDetail
		} else if account.AccountType == "L" {
			loanData, err := a.LoanDao.GetLoan(ctx, account.ID)
			if err != nil {
				return info, fmt.Errorf("system error 7")
			}
			loanData.Account = account

			if loanData.Type == "L" {
				err = lib.StructToMapSingleD2(loanData, "json", &accountDetail)
				if err != nil {
					return info, fmt.Errorf("system error 13")
				}
				data.PersonalLoanAccount = accountDetail
			} else if loanData.Type == "H" {
				homeLoanData, err := a.HomeLoanDao.GetHomeLoan(ctx, account.ID)
				if err != nil {
					return info, fmt.Errorf("system error 8")
				}
				homeLoanData.Loan = loanData
				err = lib.StructToMapSingleD2(homeLoanData, "json", &accountDetail)
				if err != nil {
					return info, fmt.Errorf("system error 9")
				}
				data.HomeLoanAccount = accountDetail

			} else if loanData.Type == "S" {
				studentLoanData, err := a.StuLoanDao.GetStuLoan(ctx, account.ID)
				if err != nil {
					return info, fmt.Errorf("system error 10")
				}
				studentLoanData.Loan = loanData
				err = lib.StructToMapSingleD2(studentLoanData, "json", &accountDetail)
				if err != nil {
					return info, fmt.Errorf("system error 11")
				}

				instituteData, err := a.InstituteDao.GetByID(ctx, studentLoanData.InstituteID)
				err = lib.StructToMapSingleD2(instituteData, "json", &accountDetail)
				if err != nil {
					return info, fmt.Errorf("system error 12")
				}
				data.StudentLoanAccount = accountDetail
			}
		}
	}

	info.UserInfo = userInfo
	info.AccountInfo = data

	return info, err
}

func (a *AccountService) GetAccountsByEmail(ctx *gin.Context, email string) (domain.UserCenter, error) {
	var info domain.UserCenter

	userInfo, err := a.UserRepo.FindByEmail(ctx, email)
	if err != nil {
		return info, fmt.Errorf("system error 30")
	}

	return a.GetAccountsByUserID(ctx, userInfo.ID)
}

func (a *AccountService) GetAccountsByAccountID(ctx *gin.Context, accountID int64) (domain.UserCenter, error) {
	var info domain.UserCenter

	accountData, err := a.AccountDao.GetAccountByID(ctx, accountID)
	if err != nil {
		return info, fmt.Errorf("system error 30")
	}

	if accountData.ID <= 0 {
		return info, nil
	}

	allAccountData, err := a.GetAccountsByUserID(ctx, accountData.UserID)
	if err != nil {
		return info, nil
	}
	// only reserve account
	if allAccountData.AccountInfo.SavingAccount["account_id"] != accountID {
		allAccountData.AccountInfo.SavingAccount = nil
	}

	if allAccountData.AccountInfo.CheckingAccount["account_id"] != accountID {
		allAccountData.AccountInfo.CheckingAccount = nil
	}

	if allAccountData.AccountInfo.StudentLoanAccount["account_id"] != accountID {
		allAccountData.AccountInfo.StudentLoanAccount = nil
	}

	if allAccountData.AccountInfo.HomeLoanAccount["account_id"] != accountID {
		allAccountData.AccountInfo.HomeLoanAccount = nil
	}

	if allAccountData.AccountInfo.PersonalLoanAccount["account_id"] != accountID {
		allAccountData.AccountInfo.PersonalLoanAccount = nil
	}

	return allAccountData, nil
}
