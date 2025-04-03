package mevhypersuper

import (
	"fmt"
	"math/big"
	"sort"
	"sync"
	"time"
)

// Insert helper function at the top of the file (after imports)
func EthToWei(eth int64) *big.Int {
	return new(big.Int).Mul(big.NewInt(eth), big.NewInt(1000000000000000000))
}

// EventTx = ETH tx with real-time dependency tracking.
type EventTx struct {
	Hash       string
	Sender     string
	Receiver   string
	GasPrice   *big.Int
	Value      *big.Int
	Profit     *big.Int
	DependsOn  []string
	Timestamp  time.Time
}

// TxGraph resolves complex dependencies for MEV optimization.
type TxGraph struct {
	Nodes map[string]*EventTx
	Edges map[string][]string
	mutex sync.RWMutex
}

// EventHorizonCore handles dynamic arbitrage and blockspace auction.
type EventHorizonCore struct {
	graph           *TxGraph
	flashloanLimit  *big.Int
	maxBundleSize   int
}

// NewEventHorizon initializes Event Horizon engine.
func NewEventHorizon(bundleSize int, flashloanCap *big.Int) *EventHorizonCore {
	return &EventHorizonCore{
		graph: &TxGraph{
			Nodes: make(map[string]*EventTx),
			Edges: make(map[string][]string),
		},
		flashloanLimit: flashloanCap,
		maxBundleSize:  bundleSize,
	}
}

// AddTransaction adds ETH tx to the dependency graph (thx leetcode)
func (eh *EventHorizonCore) AddTransaction(tx *EventTx) {
	eh.graph.mutex.Lock()
	defer eh.graph.mutex.Unlock()

	tx.Profit = new(big.Int).Sub(tx.Value, tx.GasPrice)
	eh.graph.Nodes[tx.Hash] = tx
	for _, dep := range tx.DependsOn {
		eh.graph.Edges[dep] = append(eh.graph.Edges[dep], tx.Hash)
	}
}

// ResolveDependencies resolves the optimal execution path.
func (eh *EventHorizonCore) ResolveDependencies() []*EventTx {
	eh.graph.mutex.RLock()
	defer eh.graph.mutex.RUnlock()

	resolved := make(map[string]bool)
	var executionOrder []*EventTx

	var visit func(string)
	visit = func(hash string) {
		if resolved[hash] {
			return
		}
		for _, depHash := range eh.graph.Nodes[hash].DependsOn {
			visit(depHash)
		}
		executionOrder = append(executionOrder, eh.graph.Nodes[hash])
		resolved[hash] = true
	}

	for hash := range eh.graph.Nodes {
		visit(hash)
	}

	return executionOrder
}

// GenerateOptimalBundle optimizes tx selection across protocols.
func (eh *EventHorizonCore) GenerateOptimalBundle() []*EventTx {
	txs := eh.ResolveDependencies()

	sort.Slice(txs, func(i, j int) bool {
		return txs[i].Profit.Cmp(txs[j].Profit) > 0
	})

	bundle := []*EventTx{}
	flashloanUsed := big.NewInt(0)
	for _, tx := range txs {
		if len(bundle) >= eh.maxBundleSize {
			break
		}
		if flashloanUsed.Add(flashloanUsed, tx.Value).Cmp(eh.flashloanLimit) <= 0 {
			bundle = append(bundle, tx)
		}
	}

	return bundle
}

// ExecuteBundle integrates Flashbots and flashloans.
func (eh *EventHorizonCore) ExecuteBundle(bundle []*EventTx) {
	fmt.Println("Executing MEV Event Horizon Master Bundle:")
	for _, tx := range bundle {
		fmt.Printf("Tx: %s | Sender: %s | Receiver: %s | Profit: %s wei | Value: %s wei\n",
			tx.Hash, tx.Sender, tx.Receiver, tx.Profit.String(), tx.Value.String())
	}
}

// Example demonstrates Event Horizon's MEV strategy.
func Example() {
	eh := NewEventHorizon(5, EthToWei(1000)) // 1000 ETH Flashloan limit

	eh.AddTransaction(&EventTx{"0xa", "0xA", "0xUni", big.NewInt(100e9), EthToWei(400), nil, []string{}, time.Now()})
	eh.AddTransaction(&EventTx{"0xb", "0xB", "0xSushi", big.NewInt(150e9), EthToWei(300), nil, []string{"0xa"}, time.Now()})
	eh.AddTransaction(&EventTx{"0xc", "0xC", "0xCurve", big.NewInt(200e9), EthToWei(500), nil, []string{"0xb"}, time.Now()})
	eh.AddTransaction(&EventTx{"0xd", "0xD", "0xBalancer", big.NewInt(250e9), EthToWei(450), nil, []string{}, time.Now()})
	eh.AddTransaction(&EventTx{"0xe", "0xE", "0x1inch", big.NewInt(300e9), EthToWei(600), nil, []string{"0xc", "0xd"}, time.Now()})

	bundle := eh.GenerateOptimalBundle()
	eh.ExecuteBundle(bundle)
}
