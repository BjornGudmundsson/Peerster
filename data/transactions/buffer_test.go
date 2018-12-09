package transactions_test

import (
	"testing"

	"github.com/BjornGudmundsson/Peerster/data/transactions"
	"github.com/stretchr/testify/require"
)

func Test_buffer(t *testing.T) {
	nonce := transactions.GetNonce()
	buffer := transactions.NewBuffer()
	fileName1 := "Bjorn"
	fileName2 := "Bjorninn"

	file1 := transactions.NewFile(fileName1, 0, nonce[:])
	file2 := transactions.NewFile(fileName2, 0, nonce[:])

	tx1 := transactions.NewTransaction(file1, 10)
	tx2 := transactions.NewTransaction(file2, 10)
	buffer.AddTx(tx1)
	require.True(t, !buffer.IsEmpty(), "Added to the buffer")
	require.True(t, buffer.HasTx(tx1), "Making sure I can find tx1")
	buffer.AddTx(tx2)
	require.True(t, buffer.HasTx(tx2), "Making sure I can find tx2")
	require.Equal(t, 2, buffer.GetSize(), "Making sure that the transaction was added properly")
	//Clearing the buffer
	buffer.Clear()
	require.True(t, buffer.IsEmpty(), "Make sure that the buffer is empty")
	require.Equal(t, 0, buffer.GetSize(), "Make sure that the size of the buffer is 0")
}
