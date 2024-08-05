package pg

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresRepository struct {
	DB        *pgxpool.Pool
	DBTimeout time.Duration
}

func NewConnect(ctx context.Context, dbDSN string, dbTimeout time.Duration) (PostgresRepository, error) {
	if dbDSN == "" {
		return PostgresRepository{}, nil
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	config, err := pgxpool.ParseConfig(dbDSN)
	if err != nil {
		return PostgresRepository{}, err
	}

	db, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return PostgresRepository{}, err
	}

	tx, err := db.Begin(ctx)
	if err != nil {
		tx.Rollback(ctx)
		return PostgresRepository{}, fmt.Errorf("NewConnect-BeginTx-err: %w", err)
	}
	defer tx.Rollback(ctx)

	var exists bool
	err = tx.QueryRow(ctx, existDBQuery).Scan(&exists)
	if err != nil {
		return PostgresRepository{}, err
	}

	if !exists {
		_, err = tx.Exec(ctx, createDBQuery)
		if err != nil {
			return PostgresRepository{}, err
		}
	}

	_, err = tx.Exec(ctx, createUserAuthTableQuery)
	if err != nil {
		return PostgresRepository{}, err
	}

	_, err = tx.Exec(ctx, createUserOrdersTableQuery)
	if err != nil {
		return PostgresRepository{}, err
	}

	_, err = tx.Exec(ctx, createUserOrdersUserIndexQuery)
	if err != nil {
		return PostgresRepository{}, err
	}

	_, err = tx.Exec(ctx, createUserBalanceTableQuery)
	if err != nil {
		return PostgresRepository{}, err
	}

	_, err = tx.Exec(ctx, createUserWithdrawalsTableQuery)
	if err != nil {
		return PostgresRepository{}, err
	}

	_, err = tx.Exec(ctx, createUserWithdrawalsUserIndexQuery)
	if err != nil {
		return PostgresRepository{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return PostgresRepository{}, fmt.Errorf("NewConnect-Commit-err: %w", err)
	}

	return PostgresRepository{db, dbTimeout}, nil
}
