package getaccrual

import (
	"context"
	"gophermart/internal/model"
)

type accrualService interface {
	GetAccrual(orderID string) (int, model.Accrual, error)
}

type storager interface {
	GetOrderForAccrual(ctx context.Context) (string, error)
	SetAccrual(ctx context.Context, accrual model.Accrual) error
}
