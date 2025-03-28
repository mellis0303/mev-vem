package main

import (
	"log"
	"math/big"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mellis0303/mev-vem/pkg/mev-nexus"
)

func main() {
	// Initialize MEV Simulation engine
	nexus := mevnexus.NewMEVSimulation()

	// Example transactions
	txs := []*mevnexus.Transaction{
		{
			Hash:          "0xtx1",
			Sender:        "0xA",
			Receiver:      "0xB",
			GasPrice:      big.NewInt(50e9),
			Value:         big.NewInt(3e17),
			BlockIncluded: 0,
		},
		{
			Hash:          "0xtx2",
			Sender:        "0xB",
			Receiver:      "0xC",
			GasPrice:      big.NewInt(60e9),
			Value:         big.NewInt(2e17),
			BlockIncluded: 0,
		},
		{
			Hash:          "0xtx3",
			Sender:        "0xC",
			Receiver:      "0xA",
			GasPrice:      big.NewInt(70e9),
			Value:         big.NewInt(1e17),
			BlockIncluded: 0,
		},
	}

	// Add transactions to the simulation
	for _, tx := range txs {
		nexus.AddTransaction(tx)
	}

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Main event loop
	go func() {
		for {
			// Run simulations for next 5 blocks
			currentBlock := uint64(19000000)
			nexus.RunSimulations(currentBlock, currentBlock+5)

			// Find optimal block for extraction
			optimalBlock := nexus.OptimizeExtraction()
			log.Printf("Optimal block identified: %d\n", optimalBlock)

			// Execute optimized bundle
			nexus.ExecuteOptimizedBundle(optimalBlock)

			// Wait before next iteration
			time.Sleep(time.Second)
		}
	}()

	// Wait for shutdown signal
	<-sigChan
	log.Println("Shutting down MEV Nexus...")
}