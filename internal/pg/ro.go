package pg

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"gophermart/internal/model"
	"time"
)

func (r PgRepo) GetAuthInfo(ctx context.Context, login string) (int64, string, error) {
	ctx2, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var userID int64
	var pass string
	err := r.DB.QueryRow(ctx2, getAuthInfoQuery, login).Scan(&userID, &pass)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, "", model.ErrWrongLogin
		}
		return 0, "", fmt.Errorf("GetAuthInfo-Query-err: %w", err)
	}

	return userID, pass, nil
}

func (r PgRepo) GetOrders(ctx context.Context, userID int64) ([]model.Order, error) {
	ctx2, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	rows, err := r.DB.Query(ctx2, getUserOrdersQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("GetOrders-getUserOrdersQuery-err: %w", err)
	}
	defer rows.Close()

	orders, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[model.Order])
	if err != nil {
		return nil, fmt.Errorf("PgRepo-GetOrders-CollectRows-err: %w", err)
	}

	return orders, nil
}

func (r PgRepo) GetBalance(ctx context.Context, userID int64) (model.Balance, error) {
	ctx2, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	rows, err := r.DB.Query(ctx2, getUserBalanceQuery, userID)
	if err != nil {
		return model.Balance{}, fmt.Errorf("GetBalance-getUserBalanceQuery-err: %w", err)
	}
	defer rows.Close()

	balance, err := pgx.CollectOneRow(rows, pgx.RowToStructByPos[model.Balance])
	if err != nil && !errors.Is(err, pgx.ErrNoRows) { // если запись по юзеру не найдена, то вместо ошибки вернем нулевой баланс
		return model.Balance{}, fmt.Errorf("PgRepo-GetBalance-CollectRows-err: %w", err)
	}

	return balance, nil
}

func (r PgRepo) GetWithdrawals(ctx context.Context, userID int64) ([]model.Withdraw, error) {
	ctx2, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	rows, err := r.DB.Query(ctx2, getUserWithdrawalsQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("GetWithdrawals-getUserWithdrawalsQuery-err: %w", err)
	}
	defer rows.Close()

	withdrawals, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[model.Withdraw])
	if err != nil {
		return nil, fmt.Errorf("PgRepo-GetWithdrawals-CollectRows-err: %w", err)
	}

	return withdrawals, nil
}

func (r PgRepo) GetOrderForAccrual(ctx context.Context) (string, error) {
	ctx2, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var orderID string
	err := r.DB.QueryRow(ctx2, getOrderForAccrualQuery).Scan(&orderID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", err
		}
		return "", fmt.Errorf("GetOrderForAccrual-Query-err: %w", err)
	}

	return orderID, nil
}
