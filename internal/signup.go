package internal

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/google/uuid"
)

var validNameRegex = regexp.MustCompile("[a-zA-Z] [a-zA-Z]+")
var validEmailRegex = regexp.MustCompile("^(.+)@(.+)$")
var validCarPlateRegex = regexp.MustCompile("[A-Z]{3}[0-9]{4}")

var (
	ErrAccountAlreadyExists = errors.New("account already exists")
	ErrInvalidName          = errors.New("invalid name")
	ErrInvalidEmail         = errors.New("invalid email")
	ErrInvalidCarPlate      = errors.New("invalid car plate")
	ErrInvalidCpf           = errors.New("invalid cpf")
)

// ISP - Interface Segregation Principle
type SignupData interface {
	SaveAccount(ctx context.Context, account Account) error
	GetAccountByEmail(ctx context.Context, email string) (Account, error)
}

type SignupInput struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	CPF         string `json:"cpf"`
	Password    string `json:"password"`
	CarPlate    string `json:"car_plate"`
	IsPassenger bool   `json:"is_passenger"`
	IsDriver    bool   `json:"is_driver"`
}

type Signup struct {
	signupData    SignupData
	mailerGateway MailerGateway
}

func NewSignup(signupData SignupData, mailerGateway MailerGateway) *Signup {
	return &Signup{
		signupData:    signupData,
		mailerGateway: mailerGateway,
	}
}

func (s *Signup) isValidName(name string) bool {
	return validNameRegex.MatchString(name)
}

func (s *Signup) isValidEmail(email string) bool {
	return validEmailRegex.MatchString(email)
}

func (s *Signup) isValidCarPlate(carPlate string) bool {
	return validCarPlateRegex.MatchString(carPlate)
}

func (s *Signup) Signup(ctx context.Context, input SignupInput) (string, error) {
	account := Account{
		AccountID:   uuid.NewString(),
		Name:        input.Name,
		Email:       input.Email,
		CPF:         input.CPF,
		Password:    input.Password,
		CarPlate:    input.CarPlate,
		IsPassenger: input.IsPassenger,
		IsDriver:    input.IsDriver,
	}

	existingAccount, err := s.signupData.GetAccountByEmail(ctx, account.Email)
	if err != nil {
		return "", fmt.Errorf("getting account by email: %w", err)
	}
	if existingAccount.AccountID != "" {
		return "", fmt.Errorf("%w: %s", ErrAccountAlreadyExists, existingAccount.AccountID)
	}
	if !s.isValidName(account.Name) {
		return "", ErrInvalidName
	}
	if !s.isValidEmail(account.Email) {
		return "", ErrInvalidEmail
	}
	if !validateCpf(account.CPF) {
		return "", ErrInvalidCpf
	}
	if account.IsDriver && !s.isValidCarPlate(account.CarPlate) {
		return "", ErrInvalidCarPlate
	}
	err = s.signupData.SaveAccount(ctx, account)
	if err != nil {
		return "", fmt.Errorf("saving account: %w", err)
	}
	s.mailerGateway.send(account.Email, "Welcome", "...")
	return account.AccountID, nil

}
