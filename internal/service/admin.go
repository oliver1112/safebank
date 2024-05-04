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
	InstituteDao   *dao.InstituteDAO
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
	instituteDao := dao.NewInstituteDao(db)
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
		InstituteDao:   instituteDao,
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

	if value, ok := updateData["apart"]; ok {
		accountData.Apart = cast.ToString(value)
	}

	a.AccountDao.Db.Where(&dao.Account{ID: accountData.ID}).Updates(&accountData)
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
		a.CheckingDao.Db.Where(&dao.Checking{AccountID: accountData.ID}).Updates(&checkingData)
	} else if accountData.AccountType == "S" {
		savingData, err := a.SavingDao.GetSaving(ctx, accountData.ID)
		if err != nil {
			return false, err
		}

		if value, ok := updateData["interest_rate"]; ok {
			savingData.InterestRate = cast.ToFloat64(value)
		}

		if value, ok := updateData["amount"]; ok {
			savingData.Amount = cast.ToInt64(value)
		}

		savingData.Account = accountData
		a.SavingDao.Db.Where(&dao.Saving{AccountID: accountData.ID}).Updates(&savingData)
	} else if accountData.AccountType == "L" {
		loanInfo, err := a.LoanDao.GetByAccountID(ctx, AccountID)
		if err != nil {
			return false, err
		}

		if value, ok := updateData["rate"]; ok {
			loanInfo.Rate = cast.ToFloat64(value)
		}

		if value, ok := updateData["amount"]; ok {
			loanInfo.Amount = cast.ToFloat64(value)
		}

		if value, ok := updateData["month"]; ok {
			loanInfo.Month = cast.ToInt(value)
		}

		if value, ok := updateData["payment"]; ok {
			loanInfo.Payment = cast.ToFloat64(value)
		}

		a.LoanDao.Db.Where(&dao.Loan{AccountID: accountData.ID}).Updates(&loanInfo)
		if loanInfo.Type == "H" {
			homeLoanInfo, err := a.HomeLoanDao.GetHomeLoan(ctx, AccountID)
			if err != nil {
				return false, err
			}

			if value, ok := updateData["build_year"]; ok {
				homeLoanInfo.BuildYear = cast.ToInt(value)
			}

			if value, ok := updateData["insur_acc_num"]; ok {
				homeLoanInfo.InsurAccNum = cast.ToInt(value)
			}

			if value, ok := updateData["insur_name"]; ok {
				homeLoanInfo.InsurName = cast.ToString(value)
			}

			if value, ok := updateData["insur_street"]; ok {
				homeLoanInfo.InsurStreet = cast.ToString(value)
			}

			if value, ok := updateData["insur_city"]; ok {
				homeLoanInfo.InsurCity = cast.ToString(value)
			}

			if value, ok := updateData["insur_state"]; ok {
				homeLoanInfo.InsurState = cast.ToString(value)
			}

			if value, ok := updateData["insur_zip"]; ok {
				homeLoanInfo.InsurZip = cast.ToInt(value)
			}

			if value, ok := updateData["year_insur_prm"]; ok {
				homeLoanInfo.YearInsurPrm = cast.ToFloat64(value)
			}

			a.HomeLoanDao.Db.Where(&dao.HomeLoan{AccountID: homeLoanInfo.AccountID}).Updates(&homeLoanInfo)
		} else if loanInfo.Type == "S" {
			studentLoanInfo, err := a.StuLoanDao.GetStuLoan(ctx, AccountID)
			if err != nil {
				return false, err
			}

			if value, ok := updateData["edu_institute"]; ok {
				instituteData := dao.Institute{
					InstituteName: cast.ToString(value),
				}
				institute, _ := a.InstituteDao.CreateOrUpdate(ctx, instituteData)
				studentLoanInfo.InstituteID = institute.InstituteID
			}

			if value, ok := updateData["student_id"]; ok {
				studentLoanInfo.StudentID = cast.ToInt(value)
			}

			if value, ok := updateData["grad_status"]; ok {
				studentLoanInfo.GradStatus = cast.ToString(value)
			}

			if value, ok := updateData["expect_grad_month"]; ok {
				studentLoanInfo.ExpectGradMonth = cast.ToInt(value)
			}

			if value, ok := updateData["expect_grad_year"]; ok {
				studentLoanInfo.ExpectGradYear = cast.ToInt(value)
			}

			a.StuLoanDao.Db.Where(&dao.StuLoan{AccountID: studentLoanInfo.AccountID}).Updates(&studentLoanInfo)
		} else if loanInfo.Type == "L" {
		}
	}

	return true, nil
}

func (a *AdminService) DeleteAccountByID(ctx *gin.Context, accountID int64) {
	// try to delete account
	a.HomeLoanDao.Db.Where(&dao.HomeLoan{AccountID: accountID}).Delete(&dao.HomeLoan{AccountID: accountID})
	a.StuLoanDao.Db.Where(&dao.StuLoan{AccountID: accountID}).Delete(&dao.StuLoan{AccountID: accountID})
	a.LoanDao.Db.Delete(&dao.Loan{AccountID: accountID})
	a.SavingDao.Db.Delete(&dao.Saving{AccountID: accountID})
	a.CheckingDao.Db.Delete(&dao.Checking{AccountID: accountID})
	a.AccountDao.Db.Delete(&dao.Account{ID: accountID})
}
