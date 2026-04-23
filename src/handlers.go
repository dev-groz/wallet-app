package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

func HandleWallet(w http.ResponseWriter, r *http.Request) {
	var rBody RequestBody
	err := json.NewDecoder(r.Body).Decode(&rBody)

	if err != nil {
		slog.Error("failed to unmarshal json", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	err = UpdateBalance(ctx, rBody)

	switch {
	case errors.Is(err, ErrWalletNotFound):
		slog.Error("wallet not found", "error", err, "valletId", rBody.ValletID)
		http.Error(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, ErrInsufficientFunds):
		slog.Error("insufficient funds", "error", err, "valletId", rBody.ValletID, "amount", rBody.Amount)
		http.Error(w, err.Error(), http.StatusConflict)
	case err != nil:
		slog.Error("internal server error", "error", err, "method", r.Method, "path", r.URL.Path)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	default:
		fmt.Fprintf(w, `{"status":"success"}`)
		slog.Info("POST wallet balance succesfully updated", "valletId", rBody.ValletID, "operationType", rBody.OperationType, "amount", rBody.Amount)
	}
}

func HandleBalance(w http.ResponseWriter, r *http.Request) {
	valletId := r.PathValue("WALLET_UUID")

	valId, err := uuid.Parse(valletId)

	if err != nil {
		slog.Error("valletId is not UUID", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	balance, err := GetBalance(ctx, valId)
	if err != nil {
		slog.Error("no wallet with this valletId", "error", err)
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	fmt.Fprintf(w, `{"balance": %v}`, balance)
	slog.Info("GET wallet balance", "valletId", valId, "amount", balance)

}
