package dtos

import "github.com/google/uuid"


type BalanceResponseDto struct {
	WalletID uuid.UUID                 `json:"wallet_id"`
	Balance  int64                     `json:"balance"`
	User     SubUserBalanceResponseDto `json:"user"`
}

type SubUserBalanceResponseDto struct {
	UserID uuid.UUID `json:"user_id"`
	Name   string    `json:"name"`
}
