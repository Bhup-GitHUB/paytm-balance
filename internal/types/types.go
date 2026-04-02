package types

type Entry struct {
	ID             string `json:"id"`
	TransactionID  string `json:"transaction_id"`
	AccountID      string `json:"account_id"`
	Amount         int64  `json:"amount"`
	Type           string `json:"type"`
	IdempotencyKey string `json:"idempotency_key"`
	Timestamp      int64  `json:"timestamp"`
}

type TransferRequest struct {
	Amount         string `json:"amount"`
	FromAccount    string `json:"from_account"`
	ToAccount      string `json:"to_account"`
	IdempotencyKey string `json:"idempotency_key"`
}

type TransferResponse struct {
	Status  string  `json:"status"`
	Entries []Entry `json:"entries,omitempty"`
	Message string  `json:"message,omitempty"`
}

type BalanceResponse struct {
	AccountID string `json:"account_id"`
	Balance   string `json:"balance"`
	Paise     int64  `json:"paise"`
}

type SimulationResponse struct {
	Iterations            int    `json:"iterations"`
	Float64Result         string `json:"float64_result"`
	Int64PaiseResult      string `json:"int64_paise_result"`
	MathBigDecimalResult  string `json:"mathbig_decimal_result"`
	FloatDriftDiscrepancy string `json:"float_drift_discrepancy"`
}
