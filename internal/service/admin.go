package service

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"safebank/internal/domain"
	"safebank/internal/repository"
	"safebank/internal/repository/dao"
)

type AdminService struct {
	userRepo       *repository.UserRepository
	AccountDao     *dao.AccountDAO
	CheckingDao    *dao.CheckingDAO
	HomeLoanDao    *dao.HomeLoanDAO
	LoanDao        *dao.LoanDAO
	SavingDao      *dao.SavingDAO
	StuLoanDao     *dao.StuLoanDAO
	EmployeeDao    *dao.EmployeeDAO
	AccountService *AccountService
}

func NewAdminService(db *gorm.DB) *AdminService {
	// init repo
	userDao := dao.NewUserDao(db)
	userRepo := repository.NewUserRepository(userDao)

	accountDao := dao.NewAccountDao(db)
	checkingDao := dao.NewCheckingDao(db)
	homeLoanDao := dao.NewHomeLoanDao(db)
	loanDao := dao.NewLoanDao(db)
	savingDao := dao.NewSavingDao(db)
	stuLoanDao := dao.NewStuLoanDao(db)
	employeeDao := dao.NewEmployeeDao(db)
	accountService := NewAccountService(db)

	return &AdminService{
		userRepo:       userRepo,
		AccountDao:     accountDao,
		CheckingDao:    checkingDao,
		HomeLoanDao:    homeLoanDao,
		LoanDao:        loanDao,
		SavingDao:      savingDao,
		StuLoanDao:     stuLoanDao,
		EmployeeDao:    employeeDao,
		AccountService: accountService,
	}
}

func (a *AdminService) Login(ctx *gin.Context, email, password string) (dao.Employee, error) {
	// find the user
	employee, err := a.EmployeeDao.FindByEmail(ctx, email)
	if err == repository.ErrUserNotFound {
		return dao.Employee{}, ErrInvalidUserOrPassword
	}

	if err != nil {
		return dao.Employee{}, err
	}

	//compare the password
	err = bcrypt.CompareHashAndPassword([]byte(employee.Password), []byte(password))
	if err != nil {
		// DEBUG log
		return dao.Employee{}, ErrInvalidUserOrPassword
	}

	return employee, nil
}

func (a *AdminService) GetAccountInfo(ctx *gin.Context, userID int64) (domain.UserCenter, error) {
	// get user info
	return a.AccountService.GetAccountsByUserID(ctx, userID)
}

func (a *AdminService) UpdateAccountInfo(ctx *gin.Context, AccountID int64, updateData map[string]interface{}) (bool, error) {

	accountData, err := a.AccountDao.GetAccountByID(ctx, AccountID)
	if err != nil {
		return false, err
	}

	if accountData.ID <= 0 {
		return false, nil
	}

	// update account info
	if value, ok := updateData["name"]; ok {
		accountData.Name = cast.ToString(value)
	}

	if value, ok := updateData["street"]; ok {
		accountData.Street = cast.ToString(value)
	}

	if value, ok := updateData["city"]; ok {
		accountData.City = cast.ToString(value)
	}

	if value, ok := updateData["state"]; ok {
		accountData.State = cast.ToString(value)
	}

	if value, ok := updateData["zip"]; ok {
		accountData.Zip = cast.ToString(value)
	}

	if accountData.AccountType == "C" {
		checkingData, err := a.CheckingDao.GetChecking(ctx, accountData.ID)
		if err != nil {
			return false, err
		}

		if value, ok := updateData["service_charge"]; ok {
			checkingData.ServiceCharge = cast.ToFloat64(value)
		}

		if value, ok := updateData["amount"]; ok {
			checkingData.Amount = cast.ToInt64(value)
		}

		checkingData.Account = accountData
		a.CheckingDao.Db.Save(&checkingData)
	} else if accountData.AccountType == "S" {
		savingData, err := a.SavingDao.GetSaving(ctx, accountData.ID)
		if err != nil {
			return false, err
		}

		if value, ok := updateData["interest_rate"]; ok {
			savingData.InterestRate = cast.ToFloat64(value)
		}

		if value, ok := updateData["amount"]; ok {
			savingData.Amount = cast.ToFloat64(value)
		}

		savingData.Account = accountData
		a.SavingDao.Db.Save(&savingData)
	} else if accountData.AccountType == "L" {
		loanInfo, err := a.LoanDao.GetByAccountID(ctx, AccountID)
		if err != nil {
			return false, err
		}
		loanInfo.Account = accountData

		if loanInfo.Type == "H" {
			homeLoanInfo, err := a.HomeLoanDao.GetHomeLoan(ctx, AccountID)
			if err != nil {
				return false, err
			}
			homeLoanInfo.Loan = loanInfo
			a.HomeLoanDao.Db.Save(&homeLoanInfo)
		} else if loanInfo.Type == "S" {
			studentLoanInfo, err := a.StuLoanDao.GetStuLoan(ctx, AccountID)
			if err != nil {
				return false, err
			}
			studentLoanInfo.Loan = loanInfo
			a.StuLoanDao.Db.Save(&studentLoanInfo)
		} else if loanInfo.Type == "L" {
			a.LoanDao.Db.Save(&loanInfo)
		}
	}

	return true, nil
}
