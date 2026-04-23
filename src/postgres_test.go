package main

import "testing"

func TestBuildConnString(t *testing.T) {
	env := PostgresEnv{
		Host:     "localhost",
		UserName: "user",
		Password: "pass",
		DbName:   "wallet",
	}

	conn := BuildConnString(env)

	expected := "postgres://user:pass@localhost:5432/wallet"

	if conn != expected {
		t.Fatalf("expected %s, got %s", expected, conn)
	}
}
