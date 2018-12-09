package nodes

import (
	"github.com/BjornGudmundsson/Peerster/data"
	"github.com/BjornGudmundsson/Peerster/data/transactions"
)

//HandleTxPublish handles an incoming transaction publish
func (g *Gossiper) HandleTxPublish(msg GossipAddress) {
	tx := *msg.Msg.TxPublish
	addr := msg.Addr
	hasTransaction := g.BlockChain.HasTransaction(tx) || g.TransactionBuffer.HasTx(tx)
	if hasTransaction {
		//Droppping the transaction
		return
	}
	g.TransactionBuffer.AddTx(tx)
	//Broadcast to everyone except for the person the gossiper received the transaction from
	g.BroadCastTxPublish(tx, addr)
}

//BroadCastTxPublish broadcasts a Transactions  to every peer except for the sender
func (g *Gossiper) BroadCastTxPublish(tx transactions.TxPublish, sender string) {
	peers := g.Neighbours.GetAllNeighboursWithException(sender)
	gp := &data.GossipPacket{
		TxPublish: &tx,
	}
	for _, peer := range peers {
		g.sendMessageToNeighbour(gp, peer)
	}
}
