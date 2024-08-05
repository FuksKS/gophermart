package pg

import (
	"context"
	"errors"
	"fmt"
	"gophermart/internal/model"

	"github.com/jackc/pgx/v5"
)

func (r PostgresRepository) AddAuthInfo(ctx context.Context, login, hashPass string) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.DBTimeout)
	defer cancel()

	var userID int64
	err := r.DB.QueryRow(ctx, saveAuthInfoQuery, login, hashPass).Scan(&userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, model.ErrLoginAlreadyExist
		}
		return 0, fmt.Errorf("AddAuthInfo-Exec-err: %w", err)
	}

	return userID, nil
}

func (r PostgresRepository) AddOrder(ctx context.Context, orderID string, userID int64) error {
	ctx, cancel := context.WithTimeout(ctx, r.DBTimeout)
	defer cancel()

	commandTag, err := r.DB.Exec(ctx, addOrderQuery, orderID, userID)
	if err != nil {
		return fmt.Errorf("AddOrder-addOrderQuery-err: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		var userFromDB int64
		err := r.DB.QueryRow(ctx, selectOrdersUserQuery, orderID).Scan(&userFromDB)
		if err != nil {
			return fmt.Errorf("AddOrder-selectOrdersUserQuery-err: %w", err)
		}

		if userFromDB == userID {
			return model.ErrAlreadyUploadedByThisUser
		} else {
			return model.ErrAlreadyUploadedByAnotherUser
		}
	}

	return nil
}

func (r PostgresRepository) Withdraw(ctx context.Context, withdraw model.Withdraw) error {
	ctx, cancel := context.WithTimeout(ctx, r.DBTimeout)
	defer cancel()

	tx, err := r.DB.Begin(ctx)
	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("Withdraw-BeginTx-err: %w", err)
	}
	defer tx.Rollback(ctx)

	var newBalance float64
	err = tx.QueryRow(ctx, decreaseBalanceQuery, withdraw.UserID, withdraw.Sum).Scan(&newBalance)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.ErrNotEnoughMoney
		}
		return fmt.Errorf("Withdraw-decreaseBalanceQuery-err: %w", err)
	}

	if newBalance < 0 {
		return model.ErrNotEnoughMoney
	}

	commandTag, err := tx.Exec(ctx, newWithdrawQuery, withdraw.UserID, withdraw.OrderID, withdraw.Sum)
	if err != nil {
		return fmt.Errorf("Withdraw-newWithdrawQuery-err: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return model.ErrOrderAlreadyUploaded
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("Withdraw-Commit-err: %w", err)
	}

	return nil
}

func (r PostgresRepository) SetAccrual(ctx context.Context, accrual model.Accrual) error {
	ctx, cancel := context.WithTimeout(ctx, r.DBTimeout)
	defer cancel()

	tx, err := r.DB.Begin(ctx)
	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("SetAccrual-BeginTx-err: %w", err)
	}
	defer tx.Rollback(ctx)

	var userID string
	err = tx.QueryRow(ctx, setOrderStatusQuery, accrual.Order, accrual.Status, accrual.Accrual).Scan(&userID)
	if err != nil {
		return fmt.Errorf("SetAccrual-setOrderStatusQuery-err: %w", err)
	}

	if accrual.Accrual != 0 {
		_, err = tx.Exec(ctx, increaseBalanceQuery, userID, accrual.Accrual)
		if err != nil {
			return fmt.Errorf("SetAccrual-increaseBalanceQuery-err: %w", err)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("SetAccrual-Commit-err: %w", err)
	}

	return nil
}
