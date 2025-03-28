package main

import (
	"log"
	"math/big"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mellis0303/mev-vem/pkg/mev-hypersuper"
)

func main() {
	// Initialize Event Horizon Core
	eh := mevhypersuper.NewEventHorizon(5, big.NewInt(1000e18)) // 1000 ETH Flashloan limit

	// Example transactions with dependencies
	txs := []*mevhypersuper.EventTx{
		{
			Hash:      "0xa",
			Sender:    "0xA",
			Receiver:  "0xUni",
			GasPrice:  big.NewInt(100e9),
			Value:     big.NewInt(400e18),
			DependsOn: []string{},
			Timestamp: time.Now(),
		},
		{
			Hash:      "0xb",
			Sender:    "0xB",
			Receiver:  "0xSushi",
			GasPrice:  big.NewInt(150e9),
			Value:     big.NewInt(300e18),
			DependsOn: []string{"0xa"},
			Timestamp: time.Now(),
		},
		{
			Hash:      "0xc",
			Sender:    "0xC",
			Receiver:  "0xCurve",
			GasPrice:  big.NewInt(200e9),
			Value:     big.NewInt(500e18),
			DependsOn: []string{"0xb"},
			Timestamp: time.Now(),
		},
		{
			Hash:      "0xd",
			Sender:    "0xD",
			Receiver:  "0xBalancer",
			GasPrice:  big.NewInt(250e9),
			Value:     big.NewInt(450e18),
			DependsOn: []string{},
			Timestamp: time.Now(),
		},
		{
			Hash:      "0xe",
			Sender:    "0xE",
			Receiver:  "0x1inch",
			GasPrice:  big.NewInt(300e9),
			Value:     big.NewInt(600e18),
			DependsOn: []string{"0xc", "0xd"},
			Timestamp: time.Now(),
		},
	}

	// Add transactions to the core
	for _, tx := range txs {
		eh.AddTransaction(tx)
	}

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Main event loop
	go func() {
		for {
			// Generate optimal bundle
			bundle := eh.GenerateOptimalBundle()
			
			// Execute bundle
			eh.ExecuteBundle(bundle)
			
			// Wait before next iteration
			time.Sleep(time.Second)
		}
	}()

	// Wait for shutdown signal
	<-sigChan
	log.Println("Shutting down MEV HyperSuper...")
}