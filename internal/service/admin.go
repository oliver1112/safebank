package service

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"safebank/internal/repository"
	"safebank/internal/repository/dao"
)

type AdminService struct {
	userRepo    *repository.UserRepository
	AccountDao  *dao.AccountDAO
	CheckingDao *dao.CheckingDAO
	HomeLoanDao *dao.HomeLoanDAO
	LoanDao     *dao.LoanDAO
	SavingDao   *dao.SavingDAO
	StuLoanDao  *dao.StuLoanDAO
	EmployeeDao *dao.EmployeeDAO
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

	return &AdminService{
		userRepo:    userRepo,
		AccountDao:  accountDao,
		CheckingDao: checkingDao,
		HomeLoanDao: homeLoanDao,
		LoanDao:     loanDao,
		SavingDao:   savingDao,
		StuLoanDao:  stuLoanDao,
		EmployeeDao: employeeDao,
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

func (a *AdminService) GetAccountInfo(ctx *gin.Context, userID int64) (interface{}, error) {
	// get user info
	// TODO
	accountInfo := make(map[string]interface{}, 0)

	return accountInfo, nil
}

func (a *AdminService) UpdateAccountInfo(ctx *gin.Context, AccountID int64, updateData map[string]interface{}) (interface{}, error) {

	accountInfo, err := a.AccountDao.GetAccountByID(ctx, AccountID)
	if err != nil {
		return nil, err
	}

	// update account info
	if value, ok := updateData["name"]; ok {
		accountInfo.Name = cast.ToString(value)
	}

	if value, ok := updateData["street"]; ok {
		accountInfo.Street = cast.ToString(value)
	}

	if value, ok := updateData["city"]; ok {
		accountInfo.City = cast.ToString(value)
	}

	if value, ok := updateData["state"]; ok {
		accountInfo.State = cast.ToString(value)
	}

	if value, ok := updateData["zip"]; ok {
		accountInfo.Zip = cast.ToString(value)
	}

	//a.AccountDao.db.Save(&accountInfo)

	if accountInfo.AccountType == "C" {
		//checkingInfo, _ := a.CheckingDao.GetChecking(ctx, accountInfo.ID)

		if value, ok := updateData["zip"]; ok {
			accountInfo.Zip = cast.ToString(value)
		}

	} else if accountInfo.AccountType == "S" {
		// TODO update saving
	} else if accountInfo.AccountType == "L" {
		loanInfo, err := a.LoanDao.GetByAccountID(ctx, AccountID)
		if err != nil {
			return nil, err
		}
		if loanInfo.Type == "H" {
			// TODO update homeloan
		} else if loanInfo.Type == "S" {
			// TODO update student loan
		} else if loanInfo.Type == "L" {
			// TODO update general loan
		}
	}

	return accountInfo, nil
}
