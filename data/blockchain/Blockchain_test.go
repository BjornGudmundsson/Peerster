package blockchain_test

import (
	"encoding/hex"
	"testing"

	"github.com/BjornGudmundsson/Peerster/data/blockchain"
	"github.com/BjornGudmundsson/Peerster/data/transactions"
	"github.com/stretchr/testify/require"
)

//This test tests the workflow when adding a block to the blockchain.
//It makes sure that the blockchain is in its intened state
func Test_AddToEmptyChain(t *testing.T) {
	blockchain := blockchain.NewBlockChain()
	nonce := transactions.GetNonce()
	prev := transactions.GetNonce()
	fileName1 := "Bjorn"
	fileName2 := "Bjorninn"
	fileName3 := "Bjorninn2"
	fileName4 := "Caehun"
	fileName5 := "Lefteris"
	file1 := transactions.NewFile(fileName1, 0, nil)
	file2 := transactions.NewFile(fileName2, 0, nil)
	file3 := transactions.NewFile(fileName3, 0, nil)
	file4 := transactions.NewFile(fileName4, 0, nil)
	file5 := transactions.NewFile(fileName5, 0, nil)
	tx1 := transactions.NewTransaction(file1, 10)
	tx2 := transactions.NewTransaction(file2, 10)
	tx3 := transactions.NewTransaction(file3, 10)
	tx4 := transactions.NewTransaction(file4, 10)
	tx5 := transactions.NewTransaction(file5, 10)
	//This is a completely random block
	block1 := transactions.NewBlock(prev, nonce, []transactions.TxPublish{tx1})
	hashBlock1 := block1.Hash()
	block2 := transactions.NewBlock(hashBlock1, nonce, []transactions.TxPublish{tx2})
	hashBlock2 := block2.Hash()
	hxBlock2 := hex.EncodeToString(hashBlock2[:])
	block3 := transactions.NewBlock(hashBlock2, nonce, []transactions.TxPublish{tx3})
	hashBlock3 := block3.Hash()

	//I'll use this to get a fork and create a new longest chain
	block4 := transactions.NewBlock(hashBlock1, nonce, []transactions.TxPublish{tx4})
	hashBlock4 := block4.Hash()
	//hxBlock4 := hex.EncodeToString(hashBlock4[:])
	block5 := transactions.NewBlock(hashBlock4, nonce, []transactions.TxPublish{tx5})
	hashBlock5 := block5.Hash()
	hxBlock5 := hex.EncodeToString(hashBlock5[:])

	//hxBlock3 := hex.EncodeToString(hashBlock3[:])
	hxBlock1 := hex.EncodeToString(hashBlock1[:])
	//Adding the first block
	blockchain.AddBlock(hashBlock1, prev, block1)
	hxHead := hex.EncodeToString(blockchain.LongestHash[:])
	require.Equal(t, blockchain.LongestIndex, uint64(1))
	require.Equal(t, hxHead, hxBlock1)
	require.True(t, blockchain.HasBlock(block1), "Making sure that block 1 was added")
	//Adding a block not in the blockchain and is not the first block
	blockchain.AddBlock(hashBlock3, prev, block3)
	hasBlock := blockchain.HasBlock(block3)
	require.True(t, hasBlock, "Block 3 should not have been added")
	//Extending the blockchain
	blockchain.AddBlock(hashBlock2, hashBlock1, block2)
	require.True(t, blockchain.HasBlock(block2), "Making sure block 2 was added")
	hxHead = hex.EncodeToString(blockchain.LongestHash[:])
	require.Equal(t, hxHead, hxBlock2, "making sure that the head was updated to block 2")
	//Forking the blockchain
	blockchain.AddBlock(hashBlock4, hashBlock1, block4)
	require.True(t, blockchain.HasBlock(block4), "Making sure that block 4 was added but to a fork")
	require.Equal(t, hxHead, hxBlock2, "Require that the head has not changed from block 2")
	//Creating a new longest chain
	_, _, oldHead, commonBlock := blockchain.AddBlock(hashBlock5, hashBlock4, block5)
	hxHead = hex.EncodeToString(blockchain.LongestHash[:])
	require.Equal(t, blockchain.LongestIndex, uint64(3), "Making sure the the length of the longest chain has been updated")
	require.Equal(t, hxHead, hxBlock5, "Making sure that the longest chain has been updated")
	require.NotEqual(t, hxHead, hxBlock2)
	oldTransActions := blockchain.GetTransactionsToCommonBlock(oldHead, commonBlock)
	require.Equal(t, len(oldTransActions), 1, "Testing the length of the oldTransactions")

}

func Test_FindTransAction(t *testing.T) {
	blockchain := blockchain.NewBlockChain()
	nonce := transactions.GetNonce()
	prev := transactions.GetNonce()
	fileName1 := "Bjorn"
	fileName2 := "Bjorninn"
	fileName3 := "Bjorninn2"
	file1 := transactions.NewFile(fileName1, 0, nil)
	file2 := transactions.NewFile(fileName2, 0, nil)
	file3 := transactions.NewFile(fileName3, 0, nil)
	tx1 := transactions.NewTransaction(file1, 10)
	tx2 := transactions.NewTransaction(file2, 10)
	tx3 := transactions.NewTransaction(file3, 10)
	block1 := transactions.NewBlock(prev, nonce, []transactions.TxPublish{tx1})
	hashBlock1 := block1.Hash()
	block2 := transactions.NewBlock(hashBlock1, nonce, []transactions.TxPublish{tx2})
	hashBlock2 := block2.Hash()
	var empty [32]byte
	blockchain.AddBlock(hashBlock1, empty, block1)
	//Making sure that I can find the transaction
	require.True(t, blockchain.HasTransaction(tx1), "Check that I can find a transaction in a newly made blockchain")
	blockchain.AddBlock(hashBlock2, hashBlock1, block2)
	require.True(t, blockchain.HasTransaction(tx1), "Check that I can find a transaction after a new head has been appended")
	require.True(t, blockchain.HasTransaction(tx2), "Check that I can find a transaction that has been appended to the longest chain")
	require.True(t, !blockchain.HasTransaction(tx3), "Checking that if a transaction is not there I can process it")
}
