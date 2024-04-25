package service

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"safebank/internal/repository"
	"safebank/internal/repository/dao"
)

type AccountService struct {
	userRepo    *repository.UserRepository
	AccountDao  *dao.AccountDAO
	CheckingDao *dao.CheckingDAO
	HomeLoanDao *dao.HomeLoanDAO
	LoanDao     *dao.LoanDAO
	SavingDao   *dao.SavingDAO
	StuLoanDao  *dao.StuLoanDAO
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

	return &AccountService{
		userRepo:    userRepo,
		AccountDao:  accountDao,
		CheckingDao: checkingDao,
		HomeLoanDao: homeLoanDao,
		LoanDao:     loanDao,
		SavingDao:   savingDao,
		StuLoanDao:  stuLoanDao,
	}
}

func (svc *AccountService) GetAccount(ctx *gin.Context) ([]dao.Account, error) {
	session := sessions.Default(ctx)
	id := session.Get("userId")
	accounts, err := svc.AccountDao.GetAccount(ctx, cast.ToInt64(id))
	return accounts, err
}

func (svc *AccountService) CreateOrUpdateAccount(ctx *gin.Context, data interface{}) ([]dao.Account, error) {
	session := sessions.Default(ctx)
	id := session.Get("userId")
	idInt64 := cast.ToInt64(id)
	accounts, err := svc.AccountDao.GetAccount(ctx, idInt64)
	return accounts, err
}
