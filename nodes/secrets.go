package nodes

import (
	"encoding/hex"
	"errors"
	"fmt"
	"log"

	"github.com/BjornGudmundsson/Peerster/data"
	"github.com/BjornGudmundsson/Peerster/data/peersterCrypto"
)

//This file is for handlign publishing and creating secrets

//ShareSecret handles sharing secrets on the blockchain
func (g *Gossiper) ShareSecret(fn, peer string) {
	publicKey := g.GetPublicKey(peer)
	if publicKey == nil {
		log.Fatal(errors.New("This person was not on the longest chain. Not logged"))
	}
	metadata, ok := g.Files[fn]
	if !ok {
		fmt.Println("This file has not been indexed yet")
		return
	}
	IV := metadata.IV
	Key := metadata.Key
	hashofmetafile, e := hex.DecodeString(metadata.HashOfMetaFile)
	if e != nil {
		log.Fatal(e)
	}
	pub := g.PublicKey
	secret := peersterCrypto.NewSecret(fn, g.Name, peer, pub, hashofmetafile, IV, Key)
	priv := g.PrivateKey
	fmt.Println("Secret destination: ", secret.Destination)
	encryptedSecret, e := priv.EncryptSecret(secret, publicKey.PublicKey)
	if e != nil {
		log.Fatal(e)
	}
	fmt.Println("Encrypted secret destination: ", encryptedSecret.Destination)
	g.PublishSecret(encryptedSecret)
}

//PublishSecret publishes a secret to the blockchain
func (g *Gossiper) PublishSecret(es *peersterCrypto.EncryptedSecret) {
	newTransaction := data.NewSecretKeyTransaction(es)
	fmt.Println("Has: ", g.HasTransaction(newTransaction))
	fmt.Println("Pending: ", g.checkInsidePendingTransactions(newTransaction))
	if !g.HasTransaction(newTransaction) && g.checkInsidePendingTransactions(newTransaction) {
		fmt.Println("passed condition")
		// add the transaction to the pending transactions
		g.blockChainMutex.Lock()
		g.pendingTransactions = append(g.pendingTransactions, newTransaction)
		g.blockChainMutex.Unlock()

		// send the transaction to neighbors
		keyPublish := &data.KeyPublish{Transaction: newTransaction, HopLimit: hoplimit}
		packet := &data.GossipPacket{KeyPublish: keyPublish}
		g.BroadCastPacket(packet)
	}
}
