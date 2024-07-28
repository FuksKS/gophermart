package pg

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"time"
)

type PgRepo struct {
	DB *pgxpool.Pool
}

func NewConnect(ctx context.Context, dbDSN string) (PgRepo, error) {
	if dbDSN == "" {
		return PgRepo{}, nil
	}

	ctx2, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	config, err := pgxpool.ParseConfig(dbDSN)
	if err != nil {
		return PgRepo{}, err
	}

	db, err := pgxpool.NewWithConfig(ctx2, config)
	if err != nil {
		return PgRepo{}, err
	}

	var exists bool
	err = db.QueryRow(ctx2, existDBQuery).Scan(&exists)
	if err != nil {
		return PgRepo{}, err
	}

	if !exists {
		_, err = db.Exec(ctx2, createDBQuery)
		if err != nil {
			return PgRepo{}, err
		}
	}

	_, err = db.Exec(ctx2, createUserAuthTableQuery)
	if err != nil {
		return PgRepo{}, err
	}

	_, err = db.Exec(ctx2, createUserOrdersTableQuery)
	if err != nil {
		return PgRepo{}, err
	}

	_, err = db.Exec(ctx2, createUserOrdersUserIndexQuery)
	if err != nil {
		return PgRepo{}, err
	}

	_, err = db.Exec(ctx2, createUserBalanceTableQuery)
	if err != nil {
		return PgRepo{}, err
	}

	_, err = db.Exec(ctx2, createUserWithdrawalsTableQuery)
	if err != nil {
		return PgRepo{}, err
	}

	_, err = db.Exec(ctx2, createUserWithdrawalsUserIndexQuery)
	if err != nil {
		return PgRepo{}, err
	}

	return PgRepo{db}, nil
}
