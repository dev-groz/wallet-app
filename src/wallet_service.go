package main

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

func validateData(rBody RequestBody) error {
	if rBody.OperationType != DepositOperation && rBody.OperationType != WithdrawOperation {
		return errors.New("wrong operation type (valid: DEPOSIT, WITHDRAW)")
	}

	if rBody.ValletID == uuid.Nil {
		return errors.New("missing required field: valletId")
	}

	if rBody.Amount <= 0 {
		return errors.New("amount must be positive number")
	}

	return nil
}

func UpdateBalance(ctx context.Context, rBody RequestBody) error {
	err := validateData(rBody)
	if err != nil {
		return err
	}
	err = UpdateBalanceInDatabase(ctx, rBody)
	return err
}

func GetBalance(ctx context.Context, valletId uuid.UUID) (int, error) {
	balance, err := GetBalanceFromDatabase(ctx, valletId)
	return balance, err
}
