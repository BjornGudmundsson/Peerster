package nodes

import (
	"encoding/hex"
	"fmt"
	"time"

	"github.com/BjornGudmundsson/Peerster/data"
	"github.com/BjornGudmundsson/Peerster/data/transactions"
)

//MiningThread keeps track of if it should be mining or not
func (g *Gossiper) MiningThread() {
	for {
		if !g.TransactionBuffer.IsEmpty() {
			g.StartMining()
		} else {
			//I sleep cause why not
			time.Sleep(time.Second)
		}
	}
}

//StartMining is a thread function that
//mines until they find a valid block
func (g *Gossiper) StartMining() {
	prev := g.BlockChain.LongestHash
	start := time.Now()
	for {
		head := g.BlockChain.LongestHash
		hxPrev := hex.EncodeToString(prev[:])
		hxHead := hex.EncodeToString(head[:])
		if hxPrev != hxHead {
			headBlock, _ := g.BlockChain.GetBlockByHash(head)
			g.TransactionBuffer.FilterTx(headBlock)
			if g.TransactionBuffer.IsEmpty() {
				break
			}
			prev = head
		}
		// fmt.Println("gett-start")
		newTransactions := g.TransactionBuffer.GetTransactions()
		// fmt.Println("gett-stop")

		block := transactions.MineABlock(newTransactions, prev)
		// fmt.Println("valid-start")
		if transactions.IsValidBlock(block) {
			hash := block.Hash()
			hx := hex.EncodeToString(hash[:])
			fmt.Println("FOUND BLOCK [", hx, "]")
			miningTime := time.Since(start)
			time.Sleep(miningTime)
			if g.BlockChain.HasBlock(block) {
				g.TransactionBuffer.Clear()
				//Because GetNonce is pseudo random can be a bitch
				break
			}
			if !g.BlockChain.ValidateBlockFromFork(block) {
				g.TransactionBuffer.Clear()
				//If this block has transactions that are already present on its fork
				break
			}
			//Broadcast to everyone. The empty string symbolizes that.
			g.BlockChain.AddBlock(hash, prev, block)
			g.TransactionBuffer.Clear()
			g.BroadCastBlock(block, "")
			break
		}
		// fmt.Println("valid-stop")
	}
}

//BroadCastBlock to every peer except for the peer that send you the block
func (g *Gossiper) BroadCastBlock(block transactions.Block, sender string) {
	peers := g.Neighbours.GetAllNeighboursWithException(sender)
	bp := &transactions.BlockPublish{
		Block:    block,
		HopLimit: hoplimit,
	}
	gp := &data.GossipPacket{
		BlockPublish: bp,
	}
	for _, peer := range peers {
		g.sendMessageToNeighbour(gp, peer)
	}
}
