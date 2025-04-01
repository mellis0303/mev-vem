// This file contains unit tests for the OmegaCore implementation including
// transaction addition, DFS dependency resolution, and cycle detection.

package mevomega

import (
	"math/big"
	"testing"
	"time"
)

// EthToWei is a helper for converting an integer ETH value to Wei.
func EthToWei(eth int64) *big.Int {
	wei := new(big.Int)
	wei.Mul(big.NewInt(eth), big.NewInt(1e18))
	return wei
}

func TestAddTxAndOrder(t *testing.T) {
	// Use helper to avoid float literals (which can be imprecise for big.Int)
	omega := NewOmegaCore(5, EthToWei(1500), big.NewInt(1e12))

	tx1 := &OmegaTx{
		Hash:         "tx1",
		Sender:       "A",
		Receiver:     "B",
		GasPrice:     big.NewInt(200_000_000_000), // 200 Gwei
		Value:        EthToWei(1),
		Profit:       EthToWei(1), // simplified
		Dependencies: []string{},
		Timestamp:    time.Now(),
	}
	tx2 := &OmegaTx{
		Hash:         "tx2",
		Sender:       "B",
		Receiver:     "C",
		GasPrice:     big.NewInt(200_000_000_000),
		Value:        EthToWei(1),
		Profit:       EthToWei(1),
		Dependencies: []string{"tx1"},
		Timestamp:    time.Now(),
	}

	// Add tx2 first then tx1 (the ordering logic should resolve the dependency order)
	omega.AddTx(tx2)
	omega.AddTx(tx1)

	ordered := omega.OptimizeTransactionOrdering()
	if len(ordered) != 2 {
		t.Fatalf("Expected 2 transactions, got %d", len(ordered))
	}
	if ordered[0].Hash != "tx1" {
		t.Errorf("Expected first transaction to be tx1, got %s", ordered[0].Hash)
	}
	if ordered[1].Hash != "tx2" {
		t.Errorf("Expected second transaction to be tx2, got %s", ordered[1].Hash)
	}
}

func TestCycleDetection(t *testing.T) {
	// Create a cycle in the dependency graph: tx1 depends on tx2 and vice versa.
	omega := NewOmegaCore(5, EthToWei(1500), big.NewInt(1e12))

	tx1 := &OmegaTx{
		Hash:         "tx1",
		Sender:       "A",
		Receiver:     "B",
		GasPrice:     big.NewInt(200_000_000_000),
		Value:        EthToWei(1),
		Profit:       EthToWei(1),
		Dependencies: []string{"tx2"},
		Timestamp:    time.Now(),
	}
	tx2 := &OmegaTx{
		Hash:         "tx2",
		Sender:       "B",
		Receiver:     "C",
		GasPrice:     big.NewInt(200_000_000_000),
		Value:        EthToWei(1),
		Profit:       EthToWei(1),
		Dependencies: []string{"tx1"},
		Timestamp:    time.Now(),
	}

	omega.AddTx(tx1)
	omega.AddTx(tx2)

	ordered := omega.OptimizeTransactionOrdering()
	if len(ordered) != 2 {
		t.Fatalf("Expected 2 transactions, got %d", len(ordered))
	}
	// Even if a cycle is detected, both transactions should appear.
	foundTx1, foundTx2 := false, false
	for _, tx := range ordered {
		if tx.Hash == "tx1" {
			foundTx1 = true
		}
		if tx.Hash == "tx2" {
			foundTx2 = true
		}
	}
	if !foundTx1 || !foundTx2 {
		t.Errorf("Cycle detection test failed: missing transactions in ordering")
	}
}
