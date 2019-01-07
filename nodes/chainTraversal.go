package nodes

import (
	"encoding/hex"
	"fmt"

	"github.com/BjornGudmundsson/Peerster/data/peersterCrypto"
)

//GetAllPublicKeyInLongestChain finds all public key transactions in the longest chain
func (g *Gossiper) GetAllPublicKeyInLongestChain() []peersterCrypto.PublicPair {
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

//GetSecretsForPeer gets the secrets that are meant for a peer.
func (g *Gossiper) GetSecretsForPeer(peer string) []peersterCrypto.EncryptedSecret {
	// No blockChain initialised
	list := make([]peersterCrypto.EncryptedSecret, 0)
	blockStruct := g.headBlock
	hasNext := true
	for hasNext {
		block := blockStruct.Block
		for _, transaction := range block.Transactions {
			tx := &transaction
			if !tx.IsKeyPublish() {
				fmt.Println("me?", tx.Secret.Destination, peer)
				secret := tx.Secret
				if tx.Secret.Destination == peer {
					list = append(list, *secret)
				}
			}
		}
		blockStruct, hasNext = g.blocksMap[hex.EncodeToString(block.PrevHash[:])]
	}
	return list
}
