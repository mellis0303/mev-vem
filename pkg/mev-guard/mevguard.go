package mevguard

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"sync"
)

// ETH transaction
type Transaction struct {
	EncryptedData string
	Nonce         uint64
	From          string
}

// protected + encrypted transactional pool
type MEVMempool struct {
	transactions []*Transaction
	mutex        sync.RWMutex
	key          []byte
}

// NewMEVMempool initializes a new MEV-protected mempool
func NewMEVMempool(key []byte) (*MEVMempool, error) {
	if len(key) != 32 {
		return nil, errors.New("encryption key must be 32 bytes long")
	}
	return &MEVMempool{
		transactions: []*Transaction{},
		key:          key,
	}, nil
}

// Encrypt encrypts transaction data using AES-GCM 
func (mp *MEVMempool) Encrypt(txData []byte) (string, error) {
	block, err := aes.NewCipher(mp.key)
	if err != nil {
		return "", err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := aesGCM.Seal(nonce, nonce, txData, nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}


func (mp *MEVMempool) Decrypt(encodedCipher string) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encodedCipher)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(mp.key)
	if err != nil {
		return nil, err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return aesGCM.Open(nil, nonce, ciphertext, nil)
}

// encrypt + add transaction to protected pool
func (mp *MEVMempool) AddTransaction(txData []byte, nonce uint64, from string) error {
	mp.mutex.Lock()
	defer mp.mutex.Unlock()

	encTx, err := mp.Encrypt(txData)
	if err != nil {
		return err
	}

	tx := &Transaction{
		EncryptedData: encTx,
		Nonce:         nonce,
		From:          from,
	}

	mp.transactions = append(mp.transactions, tx)
	return nil
}

// decrypt + retrieves transaction 
func (mp *MEVMempool) RetrieveTransactions() ([]*Transaction, error) {
	mp.mutex.RLock()
	defer mp.mutex.RUnlock()

	decryptedTxs := make([]*Transaction, len(mp.transactions))
	for i, tx := range mp.transactions {
		decData, err := mp.Decrypt(tx.EncryptedData)
		if err != nil {
			return nil, err
		}
		decryptedTxs[i] = &Transaction{
			EncryptedData: string(decData),
			Nonce:         tx.Nonce,
			From:          tx.From,
		}
	}

	return decryptedTxs, nil
}

// Example usage
func Example() {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		panic(err)
	}
	pool, err := NewMEVMempool(key)
	if err != nil {
		panic(err)
	}

	txData := []byte(`{"to":"0xReceiver","value":"100ETH"}`)
	err = pool.AddTransaction(txData, 1, "0xSender")
	if err != nil {
		fmt.Println("Error adding transaction:", err)
		return
	}

	txs, err := pool.RetrieveTransactions()
	if err != nil {
		fmt.Println("Error retrieving transactions:", err)
		return
	}

	for _, tx := range txs {
		fmt.Printf("Decrypted Tx: %s, Nonce: %d, From: %s\n", tx.EncryptedData, tx.Nonce, tx.From)
	}
}
