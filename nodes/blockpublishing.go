package nodes

import (
	"fmt"

	"github.com/BjornGudmundsson/Peerster/data/transactions"
)

//HandleBlockPublish handles blockpublish messages
func (g *Gossiper) HandleBlockPublish(msg GossipAddress) {
	blockPublish := msg.Msg.BlockPublish
	block := blockPublish.Block
	addr := msg.Addr
	haveSeenBlock := g.BlockChain.HasBlock(block)
	if haveSeenBlock {
		fmt.Println("Already had the block")
		//Dropping the block because I have seen it
		return
	}
	isValid := transactions.IsValidBlock(block)
	if !isValid {
		fmt.Println("Block was invalid so dropping it")
		return
	}
	isBlockValidOnFork := g.BlockChain.ValidateBlockFromFork(block)
	if !isBlockValidOnFork {
		fmt.Println("This block had duplicates")
		//Dropping fork since it has duplicate transactions
		return
	}
	prevHash := block.PrevHash
	hash := block.Hash()
	_, newChain, oldHead, commonBlock := g.BlockChain.AddBlock(hash, prevHash, block)
	if newChain {
		//We now have a new longest chain and must add the transactions back to the pool
		oldTransactions := g.BlockChain.GetTransactionsToCommonBlock(oldHead, commonBlock)
		//This method is really fucking inefficient but I don't really care
		for _, tx := range oldTransactions {
			isValid := g.BlockChain.HasTransaction(tx)
			if isValid {
				g.TransactionBuffer.AddTx(tx)
			}
		}
	}
	go g.BroadCastBlock(block, addr)
}
