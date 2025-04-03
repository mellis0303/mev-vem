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

func EthToWei(eth int64) *big.Int {
	return new(big.Int).Mul(big.NewInt(eth), big.NewInt(1e18))
}

func main() {
	// Create a cancellable context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigChan
		log.Printf("Received signal %s - initiating shutdown...", sig)
		cancel()
	}()

	omega := mevomega.NewOmegaCore(5, EthToWei(1500), big.NewInt(1e12))

	// Main processing loop
	for {
		select {
		case <-ctx.Done():
			log.Println("Shutdown signal received, stopping MEV Omega.")
			return
		default:
			orderedTxs := omega.OptimizeTransactionOrdering()
			log.Printf("Optimized %d transactions...", len(orderedTxs))
			time.Sleep(1 * time.Second)
		}
	}
}