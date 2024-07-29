package pg

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"gophermart/internal/logger"
	"gophermart/internal/model"
	"time"
)

func (r PgRepo) AddAuthInfo(ctx context.Context, login, hashPass string) (int64, error) {
	ctx2, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var userID int64
	err := r.DB.QueryRow(ctx2, saveAuthInfoQuery, login, hashPass).Scan(&userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, model.ErrLoginAlreadyExist
		}
		return 0, fmt.Errorf("AddAuthInfo-Exec-err: %w", err)
	}

	return userID, nil
}

func (r PgRepo) AddOrder(ctx context.Context, orderID string, userID int64) error {
	ctx2, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	commandTag, err := r.DB.Exec(ctx2, addOrderQuery, orderID, userID)
	if err != nil {
		return fmt.Errorf("AddOrder-addOrderQuery-err: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		var userFromDB int64
		err := r.DB.QueryRow(ctx2, selectOrdersUserQuery, orderID).Scan(&userFromDB)
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

func (r PgRepo) Withdraw(ctx context.Context, withdraw model.Withdraw) error {
	ctx2, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	tx, err := r.DB.Begin(ctx2)
	if err != nil {
		tx.Rollback(ctx2)
		return fmt.Errorf("Withdraw-BeginTx-err: %w", err)
	}
	defer tx.Rollback(ctx2)

	var newBalance float64
	err = tx.QueryRow(ctx2, decreaseBalanceQuery, withdraw.UserID, withdraw.Sum).Scan(&newBalance)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.ErrNotEnoughMoney
		}
		return fmt.Errorf("Withdraw-decreaseBalanceQuery-err: %w", err)
	}

	if newBalance < 0 {
		return model.ErrNotEnoughMoney
	}

	commandTag, err := tx.Exec(ctx2, newWithdrawQuery, withdraw.UserID, withdraw.OrderID, withdraw.Sum)
	if err != nil {
		return fmt.Errorf("Withdraw-newWithdrawQuery-err: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return model.ErrOrderAlreadyUploaded
	}

	err = tx.Commit(ctx2)
	if err != nil {
		return fmt.Errorf("Withdraw-Commit-err: %w", err)
	}

	return nil
}

func (r PgRepo) SetAccrual(ctx context.Context, accrual model.Accrual) error {
	ctx2, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	tx, err := r.DB.Begin(ctx2)
	if err != nil {
		tx.Rollback(ctx2)
		return fmt.Errorf("SetAccrual-BeginTx-err: %w", err)
	}
	defer tx.Rollback(ctx2)

	var userID string
	err = tx.QueryRow(ctx2, setOrderStatusQuery, accrual.Order, accrual.Status, accrual.Accrual).Scan(&userID)
	if err != nil {
		return fmt.Errorf("SetAccrual-setOrderStatusQuery-err: %w", err)
	}

	accrualByte, err := json.Marshal(accrual)
	logger.Log.Info("SetAccrual info", zap.String("returned user", userID), zap.String("accrual", string(accrualByte)))

	if accrual.Accrual != 0 {
		_, err = tx.Exec(ctx2, increaseBalanceQuery, userID, accrual.Accrual)
		if err != nil {
			return fmt.Errorf("SetAccrual-increaseBalanceQuery-err: %w", err)
		}
	}

	err = tx.Commit(ctx2)
	if err != nil {
		return fmt.Errorf("SetAccrual-Commit-err: %w", err)
	}

	return nil
}
