package dtos

type WithdrawRequestDTO struct {
	Amount int64  `json:"amount" validate:"required,gt=0"`
	Remark *string `json:"remark"`
}

type WithdrawResponseDTO struct {
	TransactionID   string                            `json:"transaction_id"`
	ReferenceNumber string                            `json:"reference_number"`
	Transaction     SubTransactionWithdrawResponseDto `json:"transaction"`
	User            SubUserResponseDto                `json:"user"`
}

type SubUserResponseDto struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
}

type SubTransactionWithdrawResponseDto struct {
	Amount     int64  `json:"amount"`
	Type       string `json:"type"`
	BalanceNow int64  `json:"balance_now"`
	Remark     *string `json:"remark,omitempty"`
}
