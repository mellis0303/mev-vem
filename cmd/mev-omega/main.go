package main

import (
	"context"
	"log"
	"math/big"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mellis0303/mev-vem/pkg/mev-omega"
)

// EthToWei converts an integer ETH value to Wei.
func EthToWei(eth int64) *big.Int {
	wei := new(big.Int).Mul(big.NewInt(eth), big.NewInt(1e18))
	return wei
}

func main() {
	// Create a cancellable context so we can gracefully shut down the goroutines.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("Shutting down MEV Omega...")
		cancel()
	}()

	// Initialize OmegaCore with corrected ETH-to-Wei conversions.
	omega := mevomega.NewOmegaCore(5, EthToWei(1500), big.NewInt(1e12)) // gas cap assumed to be within range

	// Example transaction with proper conversions.
	tx := &mevomega.OmegaTx{
		Hash:         "0x123",
		Sender:       "0xAlice",
		Receiver:     "0xUniswap",
		GasPrice:     big.NewInt(200e9),
		Value:        EthToWei(500),
		Profit:       EthToWei(300),
		Dependencies: []string{},
		Timestamp:    time.Now(),
	}
	omega.AddTx(tx)

	// Main event loop using context cancellation
	for {
		select {
		case <-ctx.Done():
			return
		default:
			// Optimize transaction ordering
			orderedTxs := omega.OptimizeTransactionOrdering()
			// Select optimal bundle
			bundle := omega.SelectOptimalBundle(orderedTxs)
			// Execute strategic bundle
			omega.ExecuteStrategicBundle(bundle)
			time.Sleep(time.Second)
		}
	}
}