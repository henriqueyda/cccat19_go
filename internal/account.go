package internal

import (
	"errors"
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

// Entity - Independent Business Rules. We can reuse these rules in other parts of the application or even in other applications.
type Account struct {
	AccountID   string `json:"account_id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	CPF         string `json:"cpf"`
	Password    string `json:"password"`
	CarPlate    string `json:"car_plate"`
	IsPassenger bool   `json:"is_passenger"`
	IsDriver    bool   `json:"is_driver"`
}

func NewAccount(accountID, name, email, cpf, carPlate, password string, isPassenger, isDriver bool) (Account, error) {
	account := Account{
		AccountID:   accountID,
		Name:        name,
		Email:       email,
		CPF:         cpf,
		CarPlate:    carPlate,
		Password:    password,
		IsPassenger: isPassenger,
		IsDriver:    isDriver,
	}
	if !account.isValidName(name) {
		return Account{}, ErrInvalidName
	}
	if !account.isValidEmail(email) {
		return Account{}, ErrInvalidEmail
	}
	if !validateCpf(cpf) {
		return Account{}, ErrInvalidCpf
	}
	if isDriver && !account.isValidCarPlate(carPlate) {
		return Account{}, ErrInvalidCarPlate
	}
	return account, nil
}

// Factory function to create a new account with a generated UUID.
func CreateAccount(name, email, cpf, carPlate, password string, isPassenger, isDriver bool) (Account, error) {
	return NewAccount(uuid.NewString(), name, email, cpf, carPlate, password, isPassenger, isDriver)
}

func (a *Account) isValidName(name string) bool {
	return validNameRegex.MatchString(name)
}

func (a *Account) isValidEmail(email string) bool {
	return validEmailRegex.MatchString(email)
}

func (a *Account) isValidCarPlate(carPlate string) bool {
	return validCarPlateRegex.MatchString(carPlate)
}
