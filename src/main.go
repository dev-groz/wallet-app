package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	http.HandleFunc("POST /api/v1/wallet", HandleWallet)
	http.HandleFunc("GET /api/v1/wallets/{WALLET_UUID}", HandleBalance)

	err := godotenv.Load("config.env")
	if err != nil {
		slog.Error("error loading enviroment vars", "error", err)
		return
	}

	dbEnv := PostgresEnv{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		UserName: os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DbName:   os.Getenv("POSTGRES_DB"),
	}

	err = ConnectToDatabase(dbEnv)
	if err != nil {
		slog.Error("error connecting to db", "error", err)
		return
	}
	defer CloseDbConnection()

	srv := &http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: nil,

		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,

		ReadHeaderTimeout: 2 * time.Second,
	}

	slog.Info("Server listening on", "PORT", srv.Addr)
	err = srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		slog.Error("server failed", "error", err)
	}
}
