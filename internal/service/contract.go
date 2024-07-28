package service

import (
	"context"
	"gophermart/internal/model"
)

type gophermartRepo interface {
	AddAuthInfo(ctx context.Context, login, hashPass string) (int64, error)
	GetAuthInfo(ctx context.Context, login string) (int64, string, error)
	AddOrder(ctx context.Context, orderID string, userID int64) error
	GetOrders(ctx context.Context, userID int64) ([]model.Order, error)
	GetBalance(ctx context.Context, userID int64) (model.Balance, error)
	Withdraw(ctx context.Context, withdraw model.Withdraw) error
	GetWithdrawals(ctx context.Context, userID int64) ([]model.Withdraw, error)
	GetOrderForAccrual(ctx context.Context) (string, error)
	SetAccrual(ctx context.Context, accrual model.Accrual) error
}
