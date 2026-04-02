package store

import (
	"paytm-balance/internal/types"
	"sync"
)

var (
	ledger          []types.Entry
	ledgerMutex     sync.RWMutex
	idempotencyMap  = make(map[string]types.TransferResponse)
	idempotencyLock sync.RWMutex
)

// AddEntries appends new entries to the core ledger.
func AddEntries(entries ...types.Entry) {
	ledgerMutex.Lock()
	defer ledgerMutex.Unlock()
	ledger = append(ledger, entries...)
}

// GetBalance calculates the balance of an account from the ledger directly.
func GetBalance(accountID string) int64 {
	var sum int64 = 0
	ledgerMutex.RLock()
	defer ledgerMutex.RUnlock()

	for _, entry := range ledger {
		if entry.AccountID == accountID {
			if entry.Type == "credit" {
				sum += entry.Amount
			} else if entry.Type == "debit" {
				sum -= entry.Amount
			}
		}
	}
	return sum
}

// CheckIdempotency checks if a transfer key already exists.
func CheckIdempotency(key string) (types.TransferResponse, bool) {
	idempotencyLock.RLock()
	defer idempotencyLock.RUnlock()
	resp, exists := idempotencyMap[key]
	return resp, exists
}

// SaveIdempotency saves the transfer response for a key.
func SaveIdempotency(key string, response types.TransferResponse) {
	idempotencyLock.Lock()
	defer idempotencyLock.Unlock()
	idempotencyMap[key] = response
}
