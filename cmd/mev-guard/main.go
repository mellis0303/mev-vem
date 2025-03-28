package main

import (
	"crypto/rand"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mellis0303/mev-vem/pkg/mev-guard"
)

func main() {
	// Generate encryption key
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		log.Fatal("Failed to generate encryption key:", err)
	}

	// Initialize protected MEV mempool
	pool, err := mevguard.NewMEVMempool(key)
	if err != nil {
		log.Fatal("Failed to initialize MEV mempool:", err)
	}

	// Example transaction data
	txData := []byte(`{"to":"0xReceiver","value":"100ETH"}`)

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Main event loop
	go func() {
		nonce := uint64(1)
		for {
			// Add encrypted transaction
			err := pool.AddTransaction(txData, nonce, "0xSender")
			if err != nil {
				log.Printf("Error adding transaction: %v\n", err)
				continue
			}

			// Retrieve and decrypt transactions
			txs, err := pool.RetrieveTransactions()
			if err != nil {
				log.Printf("Error retrieving transactions: %v\n", err)
				continue
			}

			// Print decrypted transactions
			log.Println("Retrieved Transactions:")
			for _, tx := range txs {
				log.Printf("Decrypted Tx: %s, Nonce: %d, From: %s\n",
					tx.EncryptedData, tx.Nonce, tx.From)
			}

			nonce++
			time.Sleep(time.Second)
		}
	}()

	// Wait for shutdown signal
	<-sigChan
	log.Println("Shutting down MEV Guard...")
}