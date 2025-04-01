// This file contains tests for the MEVGuardianEngine functionality,
// including transaction submission and decryption.

package mevgrandmotherguardia

import (
	"math/big"
	"testing"
	"time"
)

func TestSubmitAndDecryptTransactions(t *testing.T) {
	engine := NewMEVGuardianEngine()

	tx1 := &Tx{
		Hash:      "tx1",
		Sender:    "0xAlice",
		Receiver:  "0xDEX",
		GasPrice:  big.NewInt(100_000_000_000), // 100 Gwei
		Value:     big.NewInt(500_000_000_000_000_000),
		Encrypted: true,
		Timestamp: time.Now(),
	}

	tx2 := &Tx{
		Hash:      "tx2",
		Sender:    "0xBob",
		Receiver:  "0xDEX",
		GasPrice:  big.NewInt(120_000_000_000),
		Value:     big.NewInt(300_000_000_000_000_000),
		Encrypted: false,
		Timestamp: time.Now(),
	}

	engine.SubmitTransaction(tx1)
	engine.SubmitTransaction(tx2)

	decrypted := engine.DecryptTransactions()
	for _, tx := range decrypted {
		if tx.Encrypted {
			t.Errorf("Expected transaction %s to be decrypted", tx.Hash)
		}
	}
}
