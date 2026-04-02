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

func ExecuteAtomicTransfer(key string, response types.TransferResponse, entries ...types.Entry) {
	ledgerMutex.Lock()
	idempotencyLock.Lock()
	defer ledgerMutex.Unlock()
	defer idempotencyLock.Unlock()

	ledger = append(ledger, entries...)
	idempotencyMap[key] = response
}

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

func CheckIdempotency(key string) (types.TransferResponse, bool) {
	idempotencyLock.RLock()
	defer idempotencyLock.RUnlock()
	resp, exists := idempotencyMap[key]
	return resp, exists
}
