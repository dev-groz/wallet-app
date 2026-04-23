package main

import (
	"testing"

	"github.com/google/uuid"
)

func TestValidateData(t *testing.T) {
	tests := []struct {
		name    string
		input   RequestBody
		wantErr bool
	}{
		{
			name: "valid deposit",
			input: RequestBody{
				ValletID:      uuid.New(),
				OperationType: DepositOperation,
				Amount:        100,
			},
			wantErr: false,
		},
		{
			name: "invalid operation",
			input: RequestBody{
				ValletID:      uuid.New(),
				OperationType: "INVALID",
				Amount:        100,
			},
			wantErr: true,
		},
		{
			name: "empty uuid",
			input: RequestBody{
				ValletID:      uuid.Nil,
				OperationType: DepositOperation,
				Amount:        100,
			},
			wantErr: true,
		},
		{
			name: "negative amount",
			input: RequestBody{
				ValletID:      uuid.New(),
				OperationType: DepositOperation,
				Amount:        -10,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateData(tt.input)

			if (err != nil) != tt.wantErr {
				t.Fatalf("expected error = %v, got %v", tt.wantErr, err)
			}
		})
	}
}
