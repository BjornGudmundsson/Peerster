package nodes

import (
	"encoding/hex"

	"github.com/BjornGudmundsson/Peerster/data/peersterCrypto"
)

//GetAllPublicKeyInLongestChain finds all public key transactions in the longest chain
func (g *Gossiper) GetAllPublicKeyInLongestChain() []peersterCrypto.PublicPair {
	g.blockChainMutex.Lock()
	defer g.blockChainMutex.Unlock()

	// No blockChain initialised
	if g.headBlock == nil {
		return nil
	}
	list := make([]peersterCrypto.PublicPair, 0)
	blockStruct := g.headBlock
	hasNext := true

	for hasNext {
		block := blockStruct.Block
		for _, transaction := range block.Transactions {
			tx := &transaction
			if tx.IsKeyPublish() {
				pair := transaction.KeyPublish
				list = append(list, *pair)
			}
		}
		blockStruct, hasNext = g.blocksMap[hex.EncodeToString(block.PrevHash[:])]
	}
	return list
}
