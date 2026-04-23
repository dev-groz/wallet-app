package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleWallet_BadJSON(t *testing.T) {
	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/wallet",
		bytes.NewBuffer([]byte(`invalid json`)),
	)

	w := httptest.NewRecorder()

	HandleWallet(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", w.Code)
	}
}

func TestHandleBalance_InvalidUUID(t *testing.T) {
	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/wallets/invalid",
		nil,
	)

	req.SetPathValue("WALLET_UUID", "invalid")

	w := httptest.NewRecorder()

	HandleBalance(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", w.Code)
	}
}
