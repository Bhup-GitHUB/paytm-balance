package api

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"time"

	"paytm-balance/internal/store"
	"paytm-balance/internal/types"
	"paytm-balance/internal/util"

	"github.com/google/uuid"
)

func TransferHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("API Called: /transfer")
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req types.TransferRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		fmt.Println("API Crash Avoided: invalid json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.TransferResponse{Status: "error", Message: "invalid json"})
		return
	}

	if req.IdempotencyKey == "" {
		fmt.Println("API Error: Missing idempotency key")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.TransferResponse{Status: "error", Message: "idempotency_key required"})
		return
	}

	resp, exists := store.CheckIdempotency(req.IdempotencyKey)
	if exists {
		fmt.Printf("Idempotency hit for key: %s. Returning exact previous response.\n", req.IdempotencyKey)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		return
	}

	paise, err := util.ParseAmountToPaise(req.Amount)
	if err != nil {
		fmt.Println("API Error: invalid amount format")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.TransferResponse{Status: "error", Message: err.Error()})
		return
	}

	if paise <= 0 {
		fmt.Println("API Error: amount must be > 0")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.TransferResponse{Status: "error", Message: "amount must be > 0"})
		return
	}

	now := time.Now().UnixMilli()
	transactionID := uuid.New().String()

	debitEntry := types.Entry{
		ID:             uuid.New().String(),
		TransactionID:  transactionID,
		AccountID:      req.FromAccount,
		Amount:         paise,
		Type:           "debit",
		IdempotencyKey: req.IdempotencyKey,
		Timestamp:      now,
	}

	creditEntry := types.Entry{
		ID:             uuid.New().String(),
		TransactionID:  transactionID,
		AccountID:      req.ToAccount,
		Amount:         paise,
		Type:           "credit",
		IdempotencyKey: req.IdempotencyKey,
		Timestamp:      now,
	}

	finalResponse := types.TransferResponse{
		Status:  "success",
		Entries: []types.Entry{debitEntry, creditEntry},
		Message: "transfer complete",
	}

	store.ExecuteAtomicTransfer(req.IdempotencyKey, finalResponse, debitEntry, creditEntry)

	fmt.Printf("Ledger atomically appended 2 entries for amount %s (paise: %d)\n", req.Amount, paise)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(finalResponse)
}

func BalanceHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("API Called: /balance")
	accountID := r.URL.Query().Get("account_id")
	if accountID == "" {
		fmt.Println("API Error: account_id required")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sum := store.GetBalance(accountID)

	fmt.Printf("Dynamic balance calculated for %s: %d paise\n", accountID, sum)

	resp := types.BalanceResponse{
		AccountID: accountID,
		Balance:   util.FormatPaiseToAmount(sum),
		Paise:     sum,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func SimulateHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("API Called: /simulate")
	iterations := 1000000

	var floatSum float64 = 0
	for i := 0; i < iterations; i++ {
		floatSum += 0.01
	}

	var intSum int64 = 0
	for i := 0; i < iterations; i++ {
		intSum += 1
	}

	bigSum := new(big.Rat)
	addend := big.NewRat(1, 100)
	for i := 0; i < iterations; i++ {
		bigSum.Add(bigSum, addend)
	}

	expectedString := util.FormatPaiseToAmount(intSum)
	floatString := fmt.Sprintf("%.10f", floatSum)
	floatStringTrimmed := strings.TrimRight(floatString, "0")
	if strings.HasSuffix(floatStringTrimmed, ".") {
		floatStringTrimmed += "00"
	}
	floatString = floatStringTrimmed
	mathBigString := bigSum.FloatString(2)

	floatDrift, _ := new(big.Rat).SetString(floatString)
	expectedDrift, _ := new(big.Rat).SetString(expectedString)
	discrepancy := new(big.Rat).Sub(floatDrift, expectedDrift)

	fmt.Printf("Simulation Complete: Float sum %s, Int sum %s, Big sum %s\n", floatString, expectedString, mathBigString)

	resp := types.SimulationResponse{
		Iterations:            iterations,
		Float64Result:         floatString,
		Int64PaiseResult:      expectedString,
		MathBigDecimalResult:  mathBigString,
		FloatDriftDiscrepancy: strings.TrimRight(discrepancy.FloatString(10), "0"),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
