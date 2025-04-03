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

func EthToWei(eth int64) *big.Int {
	return new(big.Int).Mul(big.NewInt(eth), big.NewInt(1000000000000000000))
}

func main() {
	// Initialize Event Horizon Core
	eh := mevhypersuper.NewEventHorizon(5, EthToWei(1000)) // 1000 ETH Flashloan limit

	// Example transactions with dependencies
	txs := []*mevhypersuper.EventTx{
		{
			Hash:      "0xa",
			Sender:    "0xA",
			Receiver:  "0xUni",
			GasPrice:  big.NewInt(100e9),
			Value:     EthToWei(400),
			DependsOn: []string{},
			Timestamp: time.Now(),
		},
		{
			Hash:      "0xb",
			Sender:    "0xB",
			Receiver:  "0xSushi",
			GasPrice:  big.NewInt(150e9),
			Value:     EthToWei(300),
			DependsOn: []string{"0xa"},
			Timestamp: time.Now(),
		},
		{
			Hash:      "0xc",
			Sender:    "0xC",
			Receiver:  "0xCurve",
			GasPrice:  big.NewInt(200e9),
			Value:     EthToWei(500),
			DependsOn: []string{"0xb"},
			Timestamp: time.Now(),
		},
		{
			Hash:      "0xd",
			Sender:    "0xD",
			Receiver:  "0xBalancer",
			GasPrice:  big.NewInt(250e9),
			Value:     EthToWei(450),
			DependsOn: []string{},
			Timestamp: time.Now(),
		},
		{
			Hash:      "0xe",
			Sender:    "0xE",
			Receiver:  "0x1inch",
			GasPrice:  big.NewInt(300e9),
			Value:     EthToWei(600),
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