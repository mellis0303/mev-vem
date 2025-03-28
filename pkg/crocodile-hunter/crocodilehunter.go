package crocodilehunter

import (
	"fmt"
	"math/big"
	"sync"
	"time"
)

// Tx = a profitable ETH transaction.
type Tx struct {
	Hash       string
	From       string
	To         string
	GasPrice   *big.Int
	Profit     *big.Int
	Timestamp  time.Time
}

// Bundle = a group of transactions to submit for MEV extraction.
type Bundle struct {
	Transactions []*Tx
	TotalProfit  *big.Int
}

// dynamically identifies and bundles transactions for MEV optimization.
type FlashHunter struct {
	mempool      []*Tx
	bundles      []*Bundle
	mutex        sync.RWMutex
	maxGasPrice  *big.Int
	minProfit    *big.Int
}

// initializes a FlashHunter MEV engine instance.
func NewFlashHunter(maxGasPrice, minProfit *big.Int) *FlashHunter {
	return &FlashHunter{
		mempool:     []*Tx{},
		bundles:     []*Bundle{},
		maxGasPrice: maxGasPrice,
		minProfit:   minProfit,
	}
}

// adds transactions to the mempool with profitability evaluation.
func (fh *FlashHunter) AddTx(tx *Tx) {
	fh.mutex.Lock()
	defer fh.mutex.Unlock()

	if tx.GasPrice.Cmp(fh.maxGasPrice) <= 0 && tx.Profit.Cmp(fh.minProfit) >= 0 {
		fh.mempool = append(fh.mempool, tx)
	}
}

// identifies profitable MEV opportunities dynamically.
func (fh *FlashHunter) AnalyzeAndBundle() {
	fh.mutex.Lock()
	defer fh.mutex.Unlock()

	profitMap := map[string]*Bundle{}

	for _, tx := range fh.mempool {
		key := fmt.Sprintf("%s->%s", tx.From, tx.To)
		if _, exists := profitMap[key]; !exists {
			profitMap[key] = &Bundle{
				Transactions: []*Tx{},
				TotalProfit:  big.NewInt(0),
			}
		}

		bundle := profitMap[key]
		bundle.Transactions = append(bundle.Transactions, tx)
		bundle.TotalProfit.Add(bundle.TotalProfit, tx.Profit)
	}

	// Only keep bundles surpassing minimum threshold
	for _, bundle := range profitMap {
		if bundle.TotalProfit.Cmp(fh.minProfit) >= 0 {
			fh.bundles = append(fh.bundles, bundle)
		}
	}

	// Clear mempool after analysis
	fh.mempool = []*Tx{}
}

// simulates bundle submission for MEV extraction.
func (fh *FlashHunter) SubmitBundles() {
	fh.mutex.RLock()
	defer fh.mutex.RUnlock()

	for i, bundle := range fh.bundles {
		fmt.Printf("Submitting Bundle %d: Profit=%s\n", i+1, bundle.TotalProfit.String())
		for _, tx := range bundle.Transactions {
			fmt.Printf("\tTx: %s From: %s To: %s Gas: %s Profit: %s\n", tx.Hash, tx.From, tx.To, tx.GasPrice.String(), tx.Profit.String())
		}
	}

	// Clear bundles after submission
	fh.bundles = []*Bundle{}
}

// Example demonstrates FlashHunter engine functionality.
func Example() {
	maxGas := big.NewInt(50000000000) // 50 Gwei
	minProf := big.NewInt(50000000000000000) // 0.05 ETH

	fh := NewFlashHunter(maxGas, minProf)

	fh.AddTx(&Tx{"0xabc", "ArbBot", "Uniswap", big.NewInt(40000000000), big.NewInt(60000000000000000), time.Now()})
	fh.AddTx(&Tx{"0xdef", "FrontRunner", "SushiSwap", big.NewInt(45000000000), big.NewInt(70000000000000000), time.Now()})
	fh.AddTx(&Tx{"0xghi", "BackRunner", "Curve", big.NewInt(30000000000), big.NewInt(80000000000000000), time.Now()})

	// Dynamic analysis and bundling
	fh.AnalyzeAndBundle()

	// Simulate submission to Flashbots for MEV profit extraction
	fh.SubmitBundles()
}
