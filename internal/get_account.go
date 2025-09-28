package internal

import "context"

// ISP - Interface Segregation Principle
type GetAccountData interface {
	GetAccountByID(ctx context.Context, accountID string) (Account, error)
}

type GetAccount struct {
	getAccountData GetAccountData
}

func NewGetAccount(getAccountData GetAccountData) *GetAccount {
	return &GetAccount{
		getAccountData: getAccountData,
	}
}

func (g *GetAccount) Execute(ctx context.Context, accountID string) (Account, error) {
	return g.getAccountData.GetAccountByID(ctx, accountID)
}
