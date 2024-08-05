package pg

import (
	"context"
	"errors"
	"fmt"
	"gophermart/internal/model"

	"github.com/jackc/pgx/v5"
)

func (r PostgresRepository) GetAuthInfo(ctx context.Context, login string) (int64, string, error) {
	ctx, cancel := context.WithTimeout(ctx, r.DBTimeout)
	defer cancel()

	var userID int64
	var pass string
	err := r.DB.QueryRow(ctx, getAuthInfoQuery, login).Scan(&userID, &pass)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, "", model.ErrWrongLogin
		}
		return 0, "", fmt.Errorf("GetAuthInfo-Query-err: %w", err)
	}

	return userID, pass, nil
}

func (r PostgresRepository) GetOrders(ctx context.Context, userID int64) ([]model.Order, error) {
	ctx, cancel := context.WithTimeout(ctx, r.DBTimeout)
	defer cancel()

	rows, err := r.DB.Query(ctx, getUserOrdersQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("GetOrders-getUserOrdersQuery-err: %w", err)
	}
	defer rows.Close()

	orders, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[model.Order])
	if err != nil {
		return nil, fmt.Errorf("GetOrders-CollectRows-err: %w", err)
	}

	return orders, nil
}

func (r PostgresRepository) GetBalance(ctx context.Context, userID int64) (model.Balance, error) {
	ctx, cancel := context.WithTimeout(ctx, r.DBTimeout)
	defer cancel()

	rows, err := r.DB.Query(ctx, getUserBalanceQuery, userID)
	if err != nil {
		return model.Balance{}, fmt.Errorf("GetBalance-getUserBalanceQuery-err: %w", err)
	}
	defer rows.Close()

	balance, err := pgx.CollectOneRow(rows, pgx.RowToStructByPos[model.Balance])
	if err != nil && !errors.Is(err, pgx.ErrNoRows) { // если запись по юзеру не найдена, то вместо ошибки вернем нулевой баланс
		return model.Balance{}, fmt.Errorf("GetBalance-CollectRows-err: %w", err)
	}

	return balance, nil
}

func (r PostgresRepository) GetWithdrawals(ctx context.Context, userID int64) ([]model.Withdraw, error) {
	ctx, cancel := context.WithTimeout(ctx, r.DBTimeout)
	defer cancel()

	rows, err := r.DB.Query(ctx, getUserWithdrawalsQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("GetWithdrawals-getUserWithdrawalsQuery-err: %w", err)
	}
	defer rows.Close()

	withdrawals, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[model.Withdraw])
	if err != nil {
		return nil, fmt.Errorf("GetWithdrawals-CollectRows-err: %w", err)
	}

	return withdrawals, nil
}

func (r PostgresRepository) GetOrderForAccrual(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, r.DBTimeout)
	defer cancel()

	var orderID string
	err := r.DB.QueryRow(ctx, getOrderForAccrualQuery).Scan(&orderID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", err
		}
		return "", fmt.Errorf("GetOrderForAccrual-Query-err: %w", err)
	}

	return orderID, nil
}
