package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresEnv struct {
	Host     string
	Port     string
	UserName string
	Password string
	DbName   string
}

var ErrWalletNotFound = errors.New("wallet not found")
var ErrInsufficientFunds = errors.New("insufficient funds")
var ErrUnknownOperation = errors.New("unknown operation")

var db *pgxpool.Pool

func BuildConnString(dbEnv PostgresEnv) string {
	return fmt.Sprintf("postgres://%s:%s@%s:5432/%s",
		dbEnv.UserName,
		dbEnv.Password,
		dbEnv.Host,
		dbEnv.DbName,
	)
}

func ConnectToDatabase(dbEnv PostgresEnv) error {
	ctx := context.Background()
	connString := BuildConnString(dbEnv)

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return err
	}

	config.MaxConns = 50
	config.MinConns = 10
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 10 * time.Minute
	config.HealthCheckPeriod = 30 * time.Second

	config.ConnConfig.ConnectTimeout = 5 * time.Second

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return err
	}

	db = pool

	createTableQuery := `
        CREATE TABLE IF NOT EXISTS balance (
            valletId UUID PRIMARY KEY,
            amount BIGINT NOT NULL
        );
    `

	_, err = db.Exec(ctx, createTableQuery)
	return err
}

func GetBalanceFromDatabase(ctx context.Context, valletId uuid.UUID) (int, error) {
	var balance int
	selectQuery := `SELECT amount FROM balance WHERE valletId = $1;`
	err := db.QueryRow(ctx, selectQuery, valletId).Scan(&balance)

	if errors.Is(err, pgx.ErrNoRows) {
		return 0, ErrWalletNotFound
	}

	return balance, err
}

func UpdateBalanceInDatabase(ctx context.Context, rBody RequestBody) error {
	switch rBody.OperationType {
	case DepositOperation:
		query := `INSERT INTO balance(valletId, amount) 
					VALUES ($1, $2) 
					ON CONFLICT (valletId) DO UPDATE
					SET amount = excluded.amount + balance.amount;`

		_, err := db.Exec(ctx, query, rBody.ValletID, rBody.Amount)
		return err
	case WithdrawOperation:
		query := `UPDATE balance
					SET amount = amount - $2
					WHERE valletId = $1 
					AND amount >= $2
					RETURNING amount;`
		var updatedAmount int
		err := db.QueryRow(ctx, query, rBody.ValletID, rBody.Amount).Scan(&updatedAmount)

		if err == nil {
			return nil
		}

		if errors.Is(err, pgx.ErrNoRows) {
			var exists bool

			checkQuery := `SELECT EXISTS(
						SELECT 1 FROM balance WHERE valletId = $1
						);`

			err = db.QueryRow(ctx, checkQuery, rBody.ValletID).Scan(&exists)
			if err != nil {
				return err
			}

			if !exists {
				return ErrWalletNotFound
			}

			return ErrInsufficientFunds

		}
		return err
	}
	return ErrUnknownOperation
}

func CloseDbConnection() {
	db.Close()
}
