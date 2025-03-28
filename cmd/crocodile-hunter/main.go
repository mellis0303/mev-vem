package main

import (
	"log"
	"math/big"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mellis0303/mev-vem/pkg/crocodile-hunter"
)

func main() {
	// Initialize FlashHunter with gas and profit thresholds
	maxGas := big.NewInt(50000000000)    // 50 Gwei
	minProf := big.NewInt(50000000000000000) // 0.05 ETH
	hunter := crocodilehunter.NewFlashHunter(maxGas, minProf)

	// Example transactions
	txs := []*crocodilehunter.Tx{
		{
			Hash:      "0xabc",
			From:      "ArbBot",
			To:        "Uniswap",
			GasPrice:  big.NewInt(40000000000),
			Profit:    big.NewInt(60000000000000000),
			Timestamp: time.Now(),
		},
		{
			Hash:      "0xdef",
			From:      "FrontRunner",
			To:        "SushiSwap",
			GasPrice:  big.NewInt(45000000000),
			Profit:    big.NewInt(70000000000000000),
			Timestamp: time.Now(),
		},
		{
			Hash:      "0xghi",
			From:      "BackRunner",
			To:        "Curve",
			GasPrice:  big.NewInt(30000000000),
			Profit:    big.NewInt(80000000000000000),
			Timestamp: time.Now(),
		},
	}

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Main event loop
	go func() {
		for {
			// Add transactions to mempool
			for _, tx := range txs {
				hunter.AddTx(tx)
			}

			// Analyze and create bundles
			hunter.AnalyzeAndBundle()

			// Submit bundles for MEV extraction
			hunter.SubmitBundles()

			// Wait before next iteration
			time.Sleep(time.Second)
		}
	}()

	// Wait for shutdown signal
	<-sigChan
	log.Println("Shutting down Crocodile Hunter...")
}