package main

import (
	"log"
	"math/big"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mellis0303/mev-vem/pkg/mev-oraclex"
)

func main() {
	// Initialize OracleXEngine with minimum profit threshold
	oracleX := mevoraclex.NewOracleXEngine(1.5)

	// Example transactions
	txs := []*mevoraclex.Transaction{
		{
			Hash:      "0x123",
			Sender:    "0xAlice",
			Receiver:  "0xUniswap",
			GasPrice:  big.NewInt(100e9),
			Value:     big.NewInt(5e17),
			Timestamp: time.Now(),
		},
		{
			Hash:      "0x456",
			Sender:    "0xBob",
			Receiver:  "0xSushi",
			GasPrice:  big.NewInt(150e9),
			Value:     big.NewInt(7e17),
			Timestamp: time.Now(),
		},
		{
			Hash:      "0x789",
			Sender:    "0xEve",
			Receiver:  "0xBalancer",
			GasPrice:  big.NewInt(120e9),
			Value:     big.NewInt(9e17),
			Timestamp: time.Now(),
		},
	}

	// Add transactions to the engine
	for _, tx := range txs {
		oracleX.AddTransaction(tx)
	}

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Main event loop
	go func() {
		for {
			// Generate optimized Flashbots bundle
			bundle := oracleX.GenerateFlashbotsBundle(2)
			
			// Auction block space
			oracleX.AuctionBlockSpace(bundle)
			
			// Wait before next iteration
			time.Sleep(time.Second)
		}
	}()

	// Wait for shutdown signal
	<-sigChan
	log.Println("Shutting down MEV OracleX...")
}