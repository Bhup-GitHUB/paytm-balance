package main

import (
	"fmt"
	"log"
	"net/http"

	"paytm-balance/internal/api"
)

func main() {
	http.HandleFunc("/transfer", api.TransferHandler)
	http.HandleFunc("/balance", api.BalanceHandler)
	http.HandleFunc("/simulate", api.SimulateHandler)

	fmt.Println("Server starting on :8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Server crashed: %v", err)
	}
}
