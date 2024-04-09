package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"safebank/internal/domain"
	"safebank/internal/repository"
)

var ErrUserDuplicateEmail = repository.ErrUserDuplicateEmail
var ErrInvalidUserOrPassword = errors.New("email or password is wrong")

type UserService struct {
	userRepo    *repository.UserRepository
	accountRepo *repository.AccRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		userRepo: repo,
	}
}

func NewAccService(repo *repository.UserRepository) *UserService {
	return &UserService{
		userRepo: repo,
	}
}

func (svc *UserService) SignUp(ctx context.Context, u domain.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)

	return svc.userRepo.Create(ctx, u)
}

func (svc *UserService) Login(ctx context.Context, email, password string) (domain.User, error) {
	// find the user
	u, err := svc.userRepo.FindByEmail(ctx, email)
	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrInvalidUserOrPassword
	}

	if err != nil {
		return domain.User{}, err
	}

	//compare the password
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		// DEBUG log
		return domain.User{}, ErrInvalidUserOrPassword
	}

	return u, nil
}

func (svc *UserService) GetProfile(userID int64) (gin.H, error) {

	//savingAcc := FindSavingByUserID(userID)
	//checkingAcc := FindCheckingByUserID(userID)
	//loanAcc := FindLoanAccByUserID(userID)
	//studentLoanAcc := FindStudentLoanByUserID(userID)
	//homeLoanAcc := FindHomeLoanByUserID(userID)

	//return gin.H{
	//
	//
	//	"Account" : gin.H{
	//		"savingAcc":      savingAcc,
	//		"checkingAcc":    checkingAcc,
	//		"loanAcc":        loanAcc,
	//		"studentLoanAcc": studentLoanAcc,
	//		"homeLoanAcc":    homeLoanAcc,
	//	},
	//}, nil

	user := svc.userRepo.FindUserByUserID(userID)
	jsonData, _ := json.Marshal(user)

	//for _, account := range user.Account {
	//	if account.AccountType == "S": {
	//		saving := svc.accountRepo.FindSavingByAccID(account.ID)
	//	} else if account.AccountType == "C": {
	//		checking := svc.accountRepo.FindCheckByAccID(account.ID)
	//	} else if
	//}

	return gin.H{"data": string(jsonData)}, nil

}
