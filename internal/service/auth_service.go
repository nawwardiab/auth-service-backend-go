package service

import (
	"errors"
	"fmt"
	"server/internal/model"
	"server/internal/repo"

	"github.com/jackc/pgx"
	"golang.org/x/crypto/bcrypt"
)


var ErrInvalidCredentials = errors.New("service: invalid credentials")
var ErrUserExist = errors.New("service: can't register this user")


type AuthService struct {
	authRepo *repo.AuthRepo
}

func NewAuthService(authRepo *repo.AuthRepo) *AuthService {
	return &AuthService{authRepo: authRepo}
}

func (s *AuthService) Register(username, email, pwd string) (*model.User, error) {
	// Check if user exists
	_, exists  := s.authRepo.GetByEmail(email)
	if exists == nil {
		return nil, ErrUserExist
	}
	
	// hash password
	hashedPwd, hashErr := hashPassword(pwd)
	if hashErr != nil {
		return nil, fmt.Errorf("service: hash password: %w", hashErr)
	} else {
		usr := &model.User{
			Username: username,
			Email: email,
			PasswordHash: hashedPwd,
		}
		createUserErr := s.authRepo.CreateUser(usr)
		if createUserErr != nil{
			return nil, createUserErr
		} else {
			return usr, nil
		}
	}
}

// Login
func (s *AuthService) Login(email, password string) (*model.User, error){
	// Check if user exists
	usr, fetchingErr := s.authRepo.GetByEmail(email)
	if fetchingErr != nil {
		if errors.Is(fetchingErr, pgx.ErrNoRows) {
			return nil, ErrInvalidCredentials
		} else {
			return nil, fmt.Errorf("service: user lookup: %w", fetchingErr)
		}
	}

	// compare passwords
	pwdErr := checkPassword(usr.PasswordHash, password)
	if pwdErr != nil {
		return nil, ErrInvalidCredentials
	} else {
		return usr, nil
	}
}

// Helpers
// hashPassword
func hashPassword(password string) (string, error) {
	hashed, hashErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if hashErr != nil {
		return "", hashErr
	} else {
		return string(hashed), nil
	}
}

// checkPassword
func checkPassword(hashed, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
}