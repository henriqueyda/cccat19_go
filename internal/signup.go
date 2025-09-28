package internal

import (
	"context"
	"fmt"
)

// ISP - Interface Segregation Principle - SignupData interface is declared next to the code that uses it and defines exactly what it needs to work.
// This turn this struct independent from external interfaces.
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

// Use case - Application Business Rules
type Signup struct {
	// DIP - Dependency Inversion Principle - Signup depends on abstractions, not on implementations. Using interfaces, we eliminate the signup dependencies from database.
	// The interface is more useful to the code that uses it, not the code that implements it.
	signupData    SignupData
	mailerGateway MailerGateway
}

func NewSignup(signupData SignupData, mailerGateway MailerGateway) *Signup {
	return &Signup{
		signupData:    signupData,
		mailerGateway: mailerGateway,
	}
}

// Renamed from Signup to Execute because the name of the struct is already Signup, so the method that performs the action should be named Execute.
func (s *Signup) Execute(ctx context.Context, input SignupInput) (string, error) {
	account, err := CreateAccount(input.Name, input.Email, input.CPF, input.CarPlate, input.Password, input.IsPassenger, input.IsDriver)
	if err != nil {
		return "", err
	}
	existingAccount, err := s.signupData.GetAccountByEmail(ctx, account.Email)
	// Rules that depend on external resources should be placed in use cases, not in entities, because the entities must be independent.
	if err != nil && err != ErrAccountNotFound {
		return "", fmt.Errorf("getting account by email: %w", err)
	}
	if existingAccount.AccountID != "" {
		return "", fmt.Errorf("%w: %s", ErrAccountAlreadyExists, existingAccount.AccountID)
	}

	err = s.signupData.SaveAccount(ctx, account)
	if err != nil {
		return "", fmt.Errorf("saving account: %w", err)
	}
	s.mailerGateway.send(account.Email, "Welcome", "...")
	return account.AccountID, nil

}
