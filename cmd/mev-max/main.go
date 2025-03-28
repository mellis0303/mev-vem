package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mellis0303/mev-vem/pkg/mev-max"
)

func main() {
	// Initialize MEV Mempool
	mempool := mevmax.NewMEVMempool()

	// Example transactions
	txs := []*mevmax.Transaction{
		{
			Hash:   "0xTx1",
			From:   "Alice",
			To:     "DEX",
			Value:  100,
			GasFee: 50,
			Profit: 200,
		},
		{
			Hash:   "0xTx2",
			From:   "Bob",
			To:     "DEX",
			Value:  200,
			GasFee: 40,
			Profit: 250,
		},
		{
			Hash:   "0xTx3",
			From:   "Carol",
			To:     "DEX",
			Value:  150,
			GasFee: 60,
			Profit: 300,
		},
	}

	// Add transactions to the mempool
	for _, tx := range txs {
		mempool.AddTransaction(tx)
	}

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Main event loop
	go func() {
		for {
			// Get optimal transaction bundle
			optimalBundle := mempool.GetOptimalBundle(2)
			
			// Print bundle details
			fmt.Println("Optimal Transaction Bundle:")
			for _, tx := range optimalBundle {
				fmt.Printf("TxHash: %s, From: %s, To: %s, Profit: %d, Priority: %d\n",
					tx.Hash, tx.From, tx.To, tx.Profit, tx.priority)
			}
			
			// Wait before next iteration
			time.Sleep(time.Second)
		}
	}()

	// Wait for shutdown signal
	<-sigChan
	log.Println("Shutting down MEV Max...")
}