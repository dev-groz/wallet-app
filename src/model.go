package main

import "github.com/google/uuid"

type Operation string

const (
	DepositOperation  Operation = "DEPOSIT"
	WithdrawOperation Operation = "WITHDRAW"
)

type RequestBody struct {
	ValletID      uuid.UUID `json:"valletId"`
	OperationType Operation `json:"operationType"`
	Amount        int64     `json:"amount"`
}
