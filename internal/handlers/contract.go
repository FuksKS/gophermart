package handlers

import (
	"context"
	"gophermart/internal/model"
)

type gmService interface {
	AddAuthInfo(ctx context.Context, login, pass string, passKey []byte) (int64, error)
	GetAuthInfo(ctx context.Context, login, pass string, passKey []byte) (int64, error)
	AddOrder(ctx context.Context, orderID string, userID int64) error
	GetOrders(ctx context.Context, userID int64) ([]model.Order, error)
	GetBalance(ctx context.Context, userID int64) (model.Balance, error)
	Withdraw(ctx context.Context, withdraw model.Withdraw) error
	GetWithdrawals(ctx context.Context, userID int64) ([]model.Withdraw, error)
}
