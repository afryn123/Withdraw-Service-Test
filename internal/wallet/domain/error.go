package domain

import "errors"

var (
	ErrInsufficientBalance   = errors.New("insufficient balance")
	ErrInvalidWithdrawAmount = errors.New("withdraw amount must be positive")
)
