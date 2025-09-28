package internal

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

var ErrAccountNotFound = fmt.Errorf("account not found")

// DAO - Data Access Object - direct mapping with the database, without any business rules.

// Repository - mediates the relationship between the domain and persistence layer.
// It applies all entities rules before inserting or updating data in the database.
type AccountRepository interface {
	SaveAccount(ctx context.Context, account Account) error
	GetAccountByEmail(ctx context.Context, email string) (Account, error)
	GetAccountByID(ctx context.Context, accountID string) (Account, error)
}

type AccountRepositoryDatabase struct {
}

func NewAccountRepositoryDatabase() *AccountRepositoryDatabase {
	return &AccountRepositoryDatabase{}
}

func (a *AccountRepositoryDatabase) GetAccountByEmail(ctx context.Context, email string) (Account, error) {
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
	var rAccountID, rName, rEmail, rCPF, rPassword, rCarPlate string
	var rIsPassenger, rIsDriver bool

	err = conn.QueryRow(ctx, "select account_id, name, email, cpf, password, car_plate, is_passenger, is_driver from ccca.account where email = $1", email).Scan(
		&rAccountID,
		&rName,
		&rEmail,
		&rCPF,
		&rPassword,
		&rCarPlate,
		&rIsPassenger,
		&rIsDriver,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return Account{}, ErrAccountNotFound
		}
		return Account{}, err
	}
	// apply business rules
	account, err := NewAccount(rAccountID, rName, rEmail, rCPF, rCarPlate, rPassword, rIsPassenger, rIsDriver)
	if err != nil {
		return Account{}, err
	}
	return account, nil
}

func (a *AccountRepositoryDatabase) GetAccountByID(ctx context.Context, accountID string) (Account, error) {
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
	var rAccountID, rName, rEmail, rCPF, rPassword, rCarPlate string
	var rIsPassenger, rIsDriver bool
	err = conn.QueryRow(ctx, "select account_id, name, email, cpf, password, car_plate, is_passenger, is_driver from ccca.account where account_id = $1", accountID).Scan(
		&rAccountID,
		&rName,
		&rEmail,
		&rCPF,
		&rPassword,
		&rCarPlate,
		&rIsPassenger,
		&rIsDriver,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return Account{}, ErrAccountNotFound
		}
		return Account{}, err
	}
	account, err := NewAccount(rAccountID, rName, rEmail, rCPF, rCarPlate, rPassword, rIsPassenger, rIsDriver)
	if err != nil {
		return Account{}, err
	}
	return account, nil
}

func (a *AccountRepositoryDatabase) SaveAccount(ctx context.Context, account Account) (err error) {
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

type AccountRepositoryMemory struct {
	accounts *[]Account
}

func NewAccountRepositoryMemory(accounts *[]Account) *AccountRepositoryMemory {
	return &AccountRepositoryMemory{accounts: accounts}
}

func (a *AccountRepositoryMemory) GetAccountByEmail(ctx context.Context, email string) (account Account, err error) {
	for _, account := range *a.accounts {
		if account.Email == email {
			return account, nil
		}
	}
	return Account{}, nil
}

func (a *AccountRepositoryMemory) GetAccountByID(ctx context.Context, accountID string) (account Account, err error) {
	for _, account := range *a.accounts {
		if account.AccountID == accountID {
			return account, nil
		}
	}
	return Account{}, nil
}

func (a *AccountRepositoryMemory) SaveAccount(ctx context.Context, account Account) (err error) {
	*a.accounts = append(*a.accounts, account)
	return nil
}
