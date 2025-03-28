package mevmax

import (
	"container/heap"
	"fmt"
	"sync"
)

// tasty profitable ETH transaction
type Transaction struct {
	Hash       string
	From       string
	To         string
	Value      uint64
	GasFee     uint64
	Profit     int64
	priority   int64
	index      int
}

// implements heap.Interface for sorting transactions by profitability.
type PriorityQueue []*Transaction

func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].priority > pq[j].priority
}
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}
func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	tx := x.(*Transaction)
	tx.index = n
	*pq = append(*pq, tx)
}
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	tx := old[n-1]
	old[n-1] = nil
	tx.index = -1
	*pq = old[0 : n-1]
	return tx
}

// manages and optimizes extraction via prioritized tx handling.
type MEVMempool struct {
	pq   PriorityQueue
	lock sync.Mutex
}

// initializes a new MEV maximizing transaction pool.
func NewMEVMempool() *MEVMempool {
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)
	return &MEVMempool{pq: pq}
}

// evaluates transaction profitability for MEV optimization.
func CalculatePriority(tx *Transaction) int64 {
	return int64(tx.GasFee + tx.Profit)
}

// adds a transaction to the mempool prioritized by profitability.
func (m *MEVMempool) AddTransaction(tx *Transaction) {
	m.lock.Lock()
	defer m.lock.Unlock()
	tx.priority = CalculatePriority(tx)
	heap.Push(&m.pq, tx)
}

// returns the most profitable bundle for block inclusion.
func (m *MEVMempool) GetOptimalBundle(maxTx int) []*Transaction {
	m.lock.Lock()
	defer m.lock.Unlock()

	bundleSize := maxTx
	if bundleSize > len(m.pq) {
		bundleSize = len(m.pq)
	}
	bundle := make([]*Transaction, bundleSize)
	for i := 0; i < bundleSize; i++ {
		bundle[i] = heap.Pop(&m.pq).(*Transaction)
	}
	return bundle
}

// Example demonstrates optimal extraction via transaction prioritization.
func Example() {
	mempool := NewMEVMempool()

	// Simulate adding profitable transactions
	mempool.AddTransaction(&Transaction{"0xTx1", "Alice", "DEX", 100, 50, 200, 0, 0})
	mempool.AddTransaction(&Transaction{"0xTx2", "Bob", "DEX", 200, 40, 250, 0, 0})
	mempool.AddTransaction(&Transaction{"0xTx3", "Carol", "DEX", 150, 60, 300, 0, 0})

	// Fetch optimal transaction bundle
	optimalBundle := mempool.GetOptimalBundle(2)
	for _, tx := range optimalBundle {
		fmt.Printf("TxHash: %s, From: %s, To: %s, Profit: %d, Priority: %d\n",
			tx.Hash, tx.From, tx.To, tx.Profit, tx.priority)
	}
}
