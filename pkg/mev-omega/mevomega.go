package mevomega

import (
	"fmt"
	"math/big"
	"sort"
	"sync"
	"time"
)

// OmegaTx = advanced ETH transactions.
type OmegaTx struct {
	Hash         string
	Sender       string
	Receiver     string
	GasPrice     *big.Int
	Value        *big.Int
	Profit       *big.Int
	Dependencies []string
	Timestamp    time.Time
}

// OmegaGraph resolves dynamic transaction dependencies.
type OmegaGraph struct {
	Nodes map[string]*OmegaTx
	Edges map[string][]string
	mutex sync.RWMutex
}

// OmegaCore executes real-time MEV strategies.
type OmegaCore struct {
	graph          *OmegaGraph
	maxFlashloan   *big.Int
	maxBundleTxs   int
	gasCap         *big.Int
}

// Initialize OmegaCore.
func NewOmegaCore(bundleTxLimit int, flashloanCap, gasCap *big.Int) *OmegaCore {
	return &OmegaCore{
		graph: &OmegaGraph{
			Nodes: make(map[string]*OmegaTx),
			Edges: make(map[string][]string),
		},
		maxFlashloan: flashloanCap,
		maxBundleTxs: bundleTxLimit,
		gasCap:       gasCap,
	}
}

// AddTx dynamically integrates ETH transactions.
func (oc *OmegaCore) AddTx(tx *OmegaTx) {
	oc.graph.mutex.Lock()
	oc.graph.Nodes[tx.Hash] = tx
	for _, dep := range tx.Dependencies {
		oc.graph.Edges[dep] = append(oc.graph.Edges[dep], tx.Hash)
	}
	oc.graph.mutex.Unlock()
}

// OptimizeTransactionOrdering solves transaction graphs dynamically.
func (oc *OmegaCore) OptimizeTransactionOrdering() []*OmegaTx {
	oc.graph.mutex.RLock()
	defer oc.graph.mutex.RUnlock()

	resolved := map[string]bool{}
	var orderedTxs []*OmegaTx

	var dfsResolve func(string)
	dfsResolve = func(txHash string) {
		if resolved[txHash] {
			return
		}
		for _, dep := range oc.graph.Nodes[txHash].Dependencies {
			dfsResolve(dep)
		}
		orderedTxs = append(orderedTxs, oc.graph.Nodes[txHash])
		resolved[txHash] = true
	}

	for txHash := range oc.graph.Nodes {
		dfsResolve(txHash)
	}
	return orderedTxs
}

// SelectOptimalBundle picks the absolute highest-profit transactions.
func (oc *OmegaCore) SelectOptimalBundle(txs []*OmegaTx) []*OmegaTx {
	sort.Slice(txs, func(i, j int) bool {
		return txs[i].Profit.Cmp(txs[j].Profit) > 0
	})

	bundle := []*OmegaTx{}
	usedFlashloan := big.NewInt(0)
	usedGas := big.NewInt(0)

	for _, tx := range txs {
		if len(bundle) >= oc.maxBundleTxs || usedGas.Cmp(oc.gasCap) >= 0 {
			break
		}
		if usedFlashloan.Add(usedFlashloan, tx.Value).Cmp(oc.maxFlashloan) <= 0 {
			bundle = append(bundle, tx)
			usedGas.Add(usedGas, tx.GasPrice)
		}
	}

	return bundle
}

// ExecuteStrategicBundle manages Flashbots-style auctions dynamically.
func (oc *OmegaCore) ExecuteStrategicBundle(bundle []*OmegaTx) {
	fmt.Println("Executing MEV Omega Strategic Bundle:")
	for _, tx := range bundle {
		fmt.Printf("Tx: %s | From: %s | To: %s | Profit: %s wei | Value: %s wei\n",
			tx.Hash, tx.Sender, tx.Receiver, tx.Profit.String(), tx.Value.String())
	}
}

// Example demonstrates MEV Omega's stuff.
func Example() {
	omega := NewOmegaCore(5, big.NewInt(1500e18), big.NewInt(1e12)) // 1500 ETH Flashloan limit, gas cap

	omega.AddTx(&OmegaTx{"0x1", "0xA", "0xUniswap", big.NewInt(200e9), big.NewInt(500e18), big.NewInt(300e18), []string{}, time.Now()})
	omega.AddTx(&OmegaTx{"0x2", "0xB", "0xCurve", big.NewInt(250e9), big.NewInt(400e18), big.NewInt(200e18), []string{"0x1"}, time.Now()})
	omega.AddTx(&OmegaTx{"0x3", "0xC", "0xSushi", big.NewInt(300e9), big.NewInt(600e18), big.NewInt(400e18), []string{"0x1"}, time.Now()})
	omega.AddTx(&OmegaTx{"0x4", "0xD", "0xBalancer", big.NewInt(350e9), big.NewInt(700e18), big.NewInt(500e18), []string{"0x2", "0x3"}, time.Now()})
	omega.AddTx(&OmegaTx{"0x5", "0xE", "0x1inch", big.NewInt(400e9), big.NewInt(800e18), big.NewInt(600e18), []string{"0x4"}, time.Now()})

	optimalOrder := omega.OptimizeTransactionOrdering()
	bundle := omega.SelectOptimalBundle(optimalOrder)
	omega.ExecuteStrategicBundle(bundle)
}
