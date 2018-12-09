package transactions_test

import (
	"fmt"
	"testing"

	"github.com/BjornGudmundsson/Peerster/data/transactions"
	"github.com/stretchr/testify/require"
)

func Test_MiningABlock(t *testing.T) {
	nonce := transactions.GetNonce()
	fileName1 := "Bjorn"
	fileName2 := "Bjorninn"
	fileName3 := "Bjorninn2"
	fileName4 := "Caehun"
	fileName5 := "Lefteris"
	file1 := transactions.NewFile(fileName1, 0, nonce[:])
	file2 := transactions.NewFile(fileName2, 0, nonce[:])
	file3 := transactions.NewFile(fileName3, 0, nonce[:])
	file4 := transactions.NewFile(fileName4, 0, nonce[:])
	file5 := transactions.NewFile(fileName5, 0, nonce[:])
	tx1 := transactions.NewTransaction(file1, 10)
	tx2 := transactions.NewTransaction(file2, 10)
	tx3 := transactions.NewTransaction(file3, 10)
	tx4 := transactions.NewTransaction(file4, 10)
	tx5 := transactions.NewTransaction(file5, 10)
	txs := []transactions.TxPublish{tx1, tx2, tx3, tx4, tx5}
	nonce2 := transactions.GetNonce()
	var k uint64
	for {
		k++
		newBlock := transactions.MineABlock(txs, nonce2)
		isValid := transactions.IsValidBlock(newBlock)
		if isValid {
			require.True(t, isValid)
			break
		}
	}
	fmt.Println("Took this many iterations: ", k)
	var i uint64
	for {
		nonce := transactions.GetNonce()
		newBlock := transactions.MineABlock(nil, nonce)
		i++
		if transactions.IsValidBlock(newBlock) {
			break
		}
	}
	fmt.Println("Took this many iterations: ", i)
	fmt.Println("Found a valid block")
}
