package transactions

import (
	"sync"
)

//TransactionBuffer is a buffer that
//allows for concurrent access to
//the transactions that this gossiper
//has received
type TransactionBuffer struct {
	mux    sync.Mutex
	Buffer []TxPublish
}

func NewBuffer() *TransactionBuffer {
	return &TransactionBuffer{
		Buffer: make([]TxPublish, 0),
	}
}

func (tb *TransactionBuffer) HasTx(tx TxPublish) bool {
	for _, t := range tb.Buffer {
		if tx.Compare(t) {
			return true
		}
	}
	return false
}

func (tb *TransactionBuffer) GetSize() int {
	return len(tb.Buffer)
}

func (tb *TransactionBuffer) AddTx(tx TxPublish) {
	tb.Buffer = append(tb.Buffer, tx)
}

//IsEmpty returns whether the buffer has any
//pending transactions or not.
func (tb *TransactionBuffer) IsEmpty() bool {
	return len(tb.Buffer) == 0
}

//Clear empties the transactionbuffer.
//This should only be done when the
//gossiper is done mining and will add the
//transactions to a block
func (tb *TransactionBuffer) Clear() {
	tb.mux.Lock()
	tb.Buffer = make([]TxPublish, 0)
	tb.mux.Unlock()
}

//GetTransactions returns the current list of transactions
//that this gossiper has received
func (tb *TransactionBuffer) GetTransactions() []TxPublish {
	return tb.Buffer
}

//SetNewTransactionBuffer sets a new buffer for the pending transactions.
func (tb *TransactionBuffer) SetNewTransactionBuffer(txs []TxPublish) {
	tb.Clear()
	for _, t := range txs {
		tb.AddTx(t)
	}
}

//FilterTx filters out the transactions that are present in a block
func (tb *TransactionBuffer) FilterTx(block Block) {
	temp := make([]TxPublish, 0)
	for _, tx := range tb.Buffer {
		if !block.HasTransaction(tx) {
			temp = append(temp, tx)
		}
	}
	tb.SetNewTransactionBuffer(temp)
}
