package internal

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

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

type AccountDAO interface {
	SaveAccount(ctx context.Context, account Account) error
	GetAccountByEmail(ctx context.Context, email string) (Account, error)
	GetAccountByID(ctx context.Context, accountID string) (Account, error)
}

type AccountDAODatabase struct {
}

func NewAccountDAODatabase() *AccountDAODatabase {
	return &AccountDAODatabase{}
}

func (a *AccountDAODatabase) GetAccountByEmail(ctx context.Context, email string) (account Account, err error) {
	conn, err := pgx.Connect(ctx, "postgres://postgres:123456@localhost:5432/app")
	defer func() {
		closeErr := conn.Close(ctx)
		if err != nil {
			if closeErr != nil {
				fmt.Printf("failed to close connection: %s", closeErr)
			}
			return
		}
		err = closeErr
	}()
	if err != nil {
		return Account{}, err
	}
	conn.QueryRow(ctx, "select account_id, name, email, cpf, password, car_plate, is_passenger, is_driver from ccca.account where email = $1", email).Scan(
		&account.AccountID,
		&account.Name,
		&account.Email,
		&account.CPF,
		&account.Password,
		&account.CarPlate,
		&account.IsPassenger,
		&account.IsDriver,
	)
	return account, nil
}

func (a *AccountDAODatabase) GetAccountByID(ctx context.Context, accountID string) (account Account, err error) {
	conn, err := pgx.Connect(ctx, "postgres://postgres:123456@localhost:5432/app")
	defer func() {
		closeErr := conn.Close(ctx)
		if err != nil {
			if closeErr != nil {
				fmt.Printf("failed to close connection: %s", closeErr)
			}
			return
		}
		err = closeErr
	}()
	if err != nil {
		return Account{}, err
	}
	err = conn.QueryRow(ctx, "select account_id, name, email, cpf, password, car_plate, is_passenger, is_driver from ccca.account where account_id = $1", accountID).Scan(
		&account.AccountID,
		&account.Name,
		&account.Email,
		&account.CPF,
		&account.Password,
		&account.CarPlate,
		&account.IsPassenger,
		&account.IsDriver,
	)
	if err != nil {
		return Account{}, err
	}
	return account, nil
}

func (a *AccountDAODatabase) SaveAccount(ctx context.Context, account Account) (err error) {
	conn, err := pgx.Connect(ctx, "postgres://postgres:123456@localhost:5432/app")
	defer func() {
		closeErr := conn.Close(ctx)
		if err != nil {
			if closeErr != nil {
				fmt.Printf("failed to close connection: %s", closeErr)
			}
			return
		}
		err = closeErr
	}()
	if err != nil {
		return err
	}
	_, err = conn.Exec(ctx, "insert into ccca.account (account_id, name, email, cpf, car_plate, is_passenger, is_driver, password) values ($1, $2, $3, $4, $5, $6, $7, $8)",
		&account.AccountID,
		&account.Name,
		&account.Email,
		&account.CPF,
		&account.CarPlate,
		&account.IsPassenger,
		&account.IsDriver,
		&account.Password,
	)
	if err != nil {
		return err
	}
	return nil
}

type AccountDAOMemory struct {
	accounts *[]Account
}

func NewAccountDAOMemory(accounts *[]Account) *AccountDAOMemory {
	return &AccountDAOMemory{accounts: accounts}
}

func (a *AccountDAOMemory) GetAccountByEmail(ctx context.Context, email string) (account Account, err error) {
	for _, account := range *a.accounts {
		if account.Email == email {
			return account, nil
		}
	}
	return Account{}, nil
}

func (a *AccountDAOMemory) GetAccountByID(ctx context.Context, accountID string) (account Account, err error) {
	for _, account := range *a.accounts {
		if account.AccountID == accountID {
			return account, nil
		}
	}
	return Account{}, nil
}

func (a *AccountDAOMemory) SaveAccount(ctx context.Context, account Account) (err error) {
	*a.accounts = append(*a.accounts, account)
	return nil
}
