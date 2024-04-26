package service

import (
	"github.com/gin-gonic/gin"
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
