package mevoraclex

import (
	"fmt"
	"math/big"
	"sync"
	"time"
)

// Transaction = an ETH transaction ripe for MEV
type Transaction struct {
	Hash        string
	Sender      string
	Receiver    string
	GasPrice    *big.Int
	Value       *big.Int
	ProfitScore float64
	Timestamp   time.Time
}

// MEVMempool = an optimized mempool for MEV
type MEVMempool struct {
	transactions []*Transaction
	mutex        sync.RWMutex
}

// OracleXEngine integrates real-time prediction and bundle creation
type OracleXEngine struct {
	mempool          *MEVMempool
	minProfitScore   float64
	frontRunDetector map[string]bool
	mutex            sync.Mutex
}

// NewOracleXEngine initializes OracleX with MEV mempool
func NewOracleXEngine(minProfit float64) *OracleXEngine {
	return &OracleXEngine{
		mempool: &MEVMempool{
			transactions: []*Transaction{},
		},
		minProfitScore:   minProfit,
		frontRunDetector: make(map[string]bool),
	}
}

// adds transactions with advanced MEV analytics
func (ox *OracleXEngine) AddTransaction(tx *Transaction) {
	ox.mutex.Lock()
	defer ox.mutex.Unlock()

	tx.ProfitScore = ox.predictProfitScore(tx)
	if tx.ProfitScore >= ox.minProfitScore && !ox.isFrontRun(tx) {
		ox.mempool.mutex.Lock()
		defer ox.mempool.mutex.Unlock()
		ox.mempool.transactions = append(ox.mempool.transactions, tx)
		ox.frontRunDetector[tx.Hash] = true
	}
}

// predictProfitScore predicts transaction profitability intelligently (duh)
func (ox *OracleXEngine) predictProfitScore(tx *Transaction) float64 {
	baseScore := float64(tx.Value.Int64()) / 1e18
	gasFactor := float64(tx.GasPrice.Int64()) / 1e9
	return baseScore*0.6 + gasFactor*0.4
}

// isFrontRun detects potential front-running attempts
func (ox *OracleXEngine) isFrontRun(tx *Transaction) bool {
	_, exists := ox.frontRunDetector[tx.Hash]
	return exists
}

// GenerateFlashbotsBundle intelligently generates transaction bundles
func (ox *OracleXEngine) GenerateFlashbotsBundle(maxTx int) []*Transaction {
	ox.mempool.mutex.RLock()
	defer ox.mempool.mutex.RUnlock()

	var bundle []*Transaction
	var highestProfit float64

	for _, tx := range ox.mempool.transactions {
		if len(bundle) < maxTx && tx.ProfitScore > highestProfit {
			bundle = append(bundle, tx)
			highestProfit += tx.ProfitScore
		}
	}

	return bundle
}

// AuctionBlockSpace simulates optimal block-space auctioning for MEV
func (ox *OracleXEngine) AuctionBlockSpace(bundle []*Transaction) {
	fmt.Println("Auctioning Optimized MEV Bundle to Flashbots:")
	for _, tx := range bundle {
		fmt.Printf("Tx: %s Sender: %s Receiver: %s ProfitScore: %.2f\n",
			tx.Hash, tx.Sender, tx.Receiver, tx.ProfitScore)
	}
}

// Example demonstrates full OracleX functionality
func Example() {
	oracleX := NewOracleXEngine(1.5)

	// Add realistic Ethereum transactions
	oracleX.AddTransaction(&Transaction{"0x123", "0xAlice", "0xUniswap", big.NewInt(100e9), big.NewInt(5e17), 0, time.Now()})
	oracleX.AddTransaction(&Transaction{"0x456", "0xBob", "0xSushi", big.NewInt(150e9), big.NewInt(7e17), 0, time.Now()})
	oracleX.AddTransaction(&Transaction{"0x789", "0xEve", "0xBalancer", big.NewInt(120e9), big.NewInt(9e17), 0, time.Now()})

	// Generate and auction optimized MEV bundle
	bundle := oracleX.GenerateFlashbotsBundle(2)
	oracleX.AuctionBlockSpace(bundle)
}
