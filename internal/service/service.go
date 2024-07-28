package service

import (
	"context"
	"fmt"
	"gophermart/internal/crypto"
	"gophermart/internal/model"
)

type service struct {
	gmRepo gophermartRepo
}

func New(gmRepo gophermartRepo) *service {
	return &service{gmRepo: gmRepo}
}

func (s service) AddAuthInfo(ctx context.Context, login, pass string, passKey []byte) (int64, error) {
	encodedPass, err := crypto.PassEncrypt(passKey, pass)
	if err != nil {
		return 0, fmt.Errorf("AddAuthInfo-PassEncrypt-err: %w", err)
	}

	return s.gmRepo.AddAuthInfo(ctx, login, encodedPass)
}

func (s service) GetAuthInfo(ctx context.Context, login, pass string, passKey []byte) (int64, error) {
	userID, encryptPassFromDB, err := s.gmRepo.GetAuthInfo(ctx, login)
	if err != nil {
		return 0, fmt.Errorf("GetAuthInfo-GetAuthInfo-err: %w", err)
	}

	passFromDB, err := crypto.PassDecrypt(passKey, encryptPassFromDB)
	if err != nil {
		return 0, fmt.Errorf("GetAuthInfo-PassDecrypt-err: %w", err)
	}

	if passFromDB != pass {
		return 0, model.ErrWrongPas
	}

	return userID, nil
}

func (s service) AddOrder(ctx context.Context, orderID string, userID int64) error {
	return s.gmRepo.AddOrder(ctx, orderID, userID)
}

func (s service) GetOrders(ctx context.Context, userID int64) ([]model.Order, error) {
	return s.gmRepo.GetOrders(ctx, userID)
}

func (s service) GetBalance(ctx context.Context, userID int64) (model.Balance, error) {
	return s.gmRepo.GetBalance(ctx, userID)
}

func (s service) Withdraw(ctx context.Context, withdraw model.Withdraw) error {
	return s.gmRepo.Withdraw(ctx, withdraw)
}

func (s service) GetWithdrawals(ctx context.Context, userID int64) ([]model.Withdraw, error) {
	return s.gmRepo.GetWithdrawals(ctx, userID)
}

func (s service) GetOrderForAccrual(ctx context.Context) (string, error) {
	return s.gmRepo.GetOrderForAccrual(ctx)
}

func (s service) SetAccrual(ctx context.Context, accrual model.Accrual) error {
	return s.gmRepo.SetAccrual(ctx, accrual)
}
