package nodes

import (
	"crypto/rsa"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"

	"github.com/BjornGudmundsson/Peerster/data/peersterCrypto"

	"github.com/BjornGudmundsson/Peerster/data"
)

type pairBlockLen struct {
	Block *data.KeyBlock
	len   uint64
}

/**
Returns the PublicKey associate to a name in case that it exists in the blockChain
in case that the name has not a public key associated it returns nil.

*/
func (gossiper *Gossiper) GetPublicKey(name string) *peersterCrypto.PublicPair {
	var key *peersterCrypto.PublicPair

	// No blockChain initialised
	if gossiper.headBlock == nil {
		return key
	}

	found := false
	blockStruct := gossiper.headBlock
	hasNext := true

	for !found && hasNext {
		block := blockStruct.Block
		for _, transaction := range block.Transactions {
			if transaction.GetName() == name {
				found = true
				key = transaction.GetPublicKey()
			}
		}
		if !found {
			blockStruct, hasNext = gossiper.blocksMap[hex.EncodeToString(block.PrevHash[:])]
		}
	}

	return key
}

//HasTransaction check if a transaction is present in the longest
//chain of the gossiper
func (g *Gossiper) HasTransaction(tx *data.KeyTransaction) bool {
	g.blockChainMutex.Lock()
	defer g.blockChainMutex.Unlock()
	var has bool
	if g.headBlock == nil {
		return has
	}
	found := false
	blockStruct := g.headBlock
	hasNext := true
	for !found && hasNext {
		block := blockStruct.Block
		for _, transaction := range block.Transactions {
			pointerTx := &transaction
			if tx.Compare(pointerTx) {
				found = true
				has = true
			}
		}
		if !found {
			blockStruct, hasNext = g.blocksMap[hex.EncodeToString(block.PrevHash[:])]
		}
	}
	return has
}

/**
Checks whether a key associated to a name exists starting from a specific head of the chain
requires to have the mutex lock
*/
func (gossiper *Gossiper) existsTransactionFromBlock(tx *data.KeyTransaction, hash string) bool {

	// No blockChain initialised
	if gossiper.headBlock == nil {
		return false
	}

	found := false
	blockStruct, hasNext := gossiper.blocksMap[hash]

	for !found && hasNext {
		block := blockStruct.Block
		for _, transaction := range block.Transactions {
			pointerTx := &transaction
			if tx.Compare(pointerTx) {
				found = true
			}
		}
		if !found {
			blockStruct, hasNext = gossiper.blocksMap[hex.EncodeToString(block.PrevHash[:])]
		}
	}
	return found

}

func (gossiper *Gossiper) checkInsidePendingTransactions(tx *data.KeyTransaction) bool {
	found := false
	//fmt.Println("locking checkInsidePendingTransactions")
	gossiper.blockChainMutex.Lock()
	//fmt.Println("locked checkInsidePendingTransactions")
	for i := 0; i < len(gossiper.pendingTransactions) && !found; i++ {
		pendingTransaction := gossiper.pendingTransactions[i]
		found = tx.Compare(pendingTransaction)
		fmt.Println("Found: ", found)
	}
	//fmt.Println("unlocking checkInsidePendingTransactions")
	gossiper.blockChainMutex.Unlock()
	//fmt.Println("unlocked checkInsidePendingTransactions")

	return !found
}

func (gossiper *Gossiper) PublishPublicKey(name string, key rsa.PublicKey) bool {
	// check if the name already has an associated key
	// neither in the pending transactions
	newTransaction := data.NewEncryptionKeyTransaction(key, name)
	if gossiper.GetPublicKey(name) == nil && gossiper.checkInsidePendingTransactions(newTransaction) {
		// add the transaction to the pending transactions

		//fmt.Println("locking PublishPublicKey")
		gossiper.blockChainMutex.Lock()
		//fmt.Println("locked PublishPublicKey")
		gossiper.pendingTransactions = append(gossiper.pendingTransactions, newTransaction)
		//fmt.Println("unlocking PublishPublicKey")
		gossiper.blockChainMutex.Unlock()
		//fmt.Println("unlocked PublishPublicKey")

		// send the transaction to neighbors
		keyPublish := &data.KeyPublish{Transaction: newTransaction, HopLimit: hoplimit}
		packet := &data.GossipPacket{KeyPublish: keyPublish}
		gossiper.BroadCastPacket(packet)

		return true
	}
	return false

}

func (gossiper *Gossiper) HandleBlockReply(reply *data.BlockReply) {
	if reply.Destination == gossiper.Name {
		//fmt.Println("i got a Block reply")
		// the reply is for me
		gossiper.HandleNewBlock(reply.KeyBlockPublish)

	} else {
		// the reply is for another peer
		reply.KeyBlockPublish.HopLimit -= 1
		if reply.KeyBlockPublish.HopLimit > 0 {
			// create gossip packet with request and send it to Destination
			packet := &data.GossipPacket{BlockReply: reply}
			gossiper.SendPacketViaRoutingTable(packet, reply.Destination)

		}
	}
}

func (gossiper *Gossiper) HandleBlockRequest(request *data.BlockRequest) {
	if request.Destination == gossiper.Name {
		//fmt.Println("I got a block request")
		// the request is for me
		blockHashBytes := request.HashValue
		blockHashString := hex.EncodeToString(blockHashBytes[:])

		// block to be send
		//fmt.Println("locking HandleBlockRequest")
		gossiper.blockChainMutex.Lock()
		//fmt.Println("locked HandleBlockRequest")
		block, found := gossiper.blocksMap[blockHashString]
		//fmt.Println("unlocking HandleBlockRequest")
		gossiper.blockChainMutex.Unlock()
		//fmt.Println("unlocked HandleBlockRequest")

		if found {
			publish := &data.KeyBlockPublish{Origin: gossiper.Name, HopLimit: hoplimit, Block: block.Block}
			reply := &data.BlockReply{Destination: request.Origin, KeyBlockPublish: publish}
			// create gossip packet with reply and send it
			packet := &data.GossipPacket{BlockReply: reply}

			gossiper.SendPacketViaRoutingTable(packet, request.Origin)

			//hashBlcok := block.Block.Hash()
			//fmt.Println("i sended the block with hash", hex.EncodeToString(hashBlcok[:]))

			if len(block.Block.Transactions) > 0 {
				//fmt.Println(*block.Block.Transactions[0].KeyPublish)
				//fmt.Println(*block.Block.Transactions[0].Secret)
			}

		}

	} else {
		// the request is for another peer
		request.HopLimit -= 1
		if request.HopLimit > 0 {
			// create gossip packet with request and send it to Destination
			packet := &data.GossipPacket{BlockRequest: request}
			gossiper.SendPacketViaRoutingTable(packet, request.Destination)

		}
	}
}

func (gossiper *Gossiper) HandleKeyTransaction(publish *data.KeyPublish) {
	fmt.Println("New transaction")
	transaction := publish.Transaction
	// Check it is not published in the main chain
	valid := gossiper.GetPublicKey(transaction.GetName()) == nil

	//fmt.Println("locking HandleKeyTransaction")
	gossiper.blockChainMutex.Lock()
	//fmt.Println("locked HandleKeyTransaction")
	// Check it is not in the pending transactions
	for _, pending := range gossiper.pendingTransactions {
		valid = valid && transaction.GetName() != pending.GetName()
	}

	if valid {
		// Add it
		gossiper.pendingTransactions = append(gossiper.pendingTransactions, transaction)
	}
	//fmt.Println("unlocking HandleKeyTransaction")
	gossiper.blockChainMutex.Unlock()
	//fmt.Println("unlocked HandleKeyTransaction")

}

func (gossiper *Gossiper) HandleNewBlock(blockPublish *data.KeyBlockPublish) {
	//fmt.Println("New Block")
	newBlock := blockPublish.Block
	// check the prove of work
	newBlockHash := newBlock.Hash()
	//fmt.Println("new block hash", hex.EncodeToString(newBlockHash[:]))
	if newBlockHash[0] == 0 && newBlockHash[1] == 0 {

		//fmt.Println("proof of work correct")
		// if the prev block is the origin block, we must accept the new block.
		allZeros := true
		for _, b := range newBlock.PrevHash {
			allZeros = allZeros && b == 0
		}

		if allZeros {
			//fmt.Println("first block")
			valid := true
			// check for duplications in the block
			for i, transaction := range newBlock.Transactions {
				for j := i + 1; j < len(newBlock.Transactions); j++ {
					pointerTx := &transaction
					otherTransaction := &newBlock.Transactions[j]
					valid = valid && !pointerTx.Compare(otherTransaction)
				}
			}

			if valid {
				//fmt.Println("locking HandleNewBlock")
				gossiper.blockChainMutex.Lock()
				//fmt.Println("locked HandleNewBlock")
				newBlockStruct := &pairBlockLen{Block: newBlock, len: 1}
				// add to the chain
				//fmt.Println("ADDING BLOCK", hex.EncodeToString(newBlockHash[:]))
				gossiper.blocksMap[hex.EncodeToString(newBlockHash[:])] = newBlockStruct
				//fmt.Println("unlocking HandleNewBlock")
				gossiper.blockChainMutex.Unlock()
				//fmt.Println("unlocked HandleNewBlock")

				// broadcast the block
				blockPublish.HopLimit -= 1
				if blockPublish.HopLimit > 0 {
					packet := &data.GossipPacket{KeyBlockPublish: blockPublish}
					gossiper.BroadCastPacket(packet)
				}

			}

		} else {
			//fmt.Println("not first block")
			//fmt.Println("locking HandleNewBlock")
			gossiper.blockChainMutex.Lock()
			//fmt.Println("locked HandleNewBlock")

			prevBlockStruct, haveIt := gossiper.blocksMap[hex.EncodeToString(newBlock.PrevHash[:])]
			gossiper.blockChainMutex.Unlock()

			// if i don't have the prev block, ask for it.
			if !haveIt {
				//fmt.Println("i don't have the prev", hex.EncodeToString(newBlock.PrevHash[:]))
				if len(blockPublish.Block.Transactions) > 0 {
					//fmt.Println(*blockPublish.Block.Transactions[0].KeyPublish)
					//fmt.Println(*blockPublish.Block.Transactions[0].Secret)
				}
				gossiper.RequestBlock(newBlock.PrevHash, blockPublish.Origin)
			}

			askedTime := time.Now()
			// wait until i have the prev block (max 5 seg)
			for !haveIt && time.Now().Sub(askedTime) < 5*time.Second {
				time.Sleep(10 * time.Millisecond)
				gossiper.blockChainMutex.Lock()
				prevBlockStruct, haveIt = gossiper.blocksMap[hex.EncodeToString(newBlock.PrevHash[:])]
				gossiper.blockChainMutex.Unlock()

			}

			// if i haven't received it, do not accept the new block
			if haveIt {
				gossiper.blockChainMutex.Lock()

				//fmt.Println("I have the prev")
				prevHashString := hex.EncodeToString(prevBlockStruct.Block.PrevHash[:])
				valid := true
				// check validity of all the transactions of the block (in the chain and no repeated)

				// check for duplications in the chain
				for _, transaction := range newBlock.Transactions {
					valid = valid && !gossiper.existsTransactionFromBlock(&transaction, prevHashString)
				}

				// check for duplications in the block
				for i, transaction := range newBlock.Transactions {
					for j := i + 1; j < len(newBlock.Transactions); j++ {
						pointerTx := &transaction
						otherTransaction := &newBlock.Transactions[j]
						valid = valid && !pointerTx.Compare(otherTransaction)
					}
				}

				//fmt.Println("valid", valid)

				if valid {
					// count new length (prev + 1)
					newBlockStruct := &pairBlockLen{Block: newBlock, len: prevBlockStruct.len + 1}

					// add to the chain
					//fmt.Println("ADDING BLOCK", hex.EncodeToString(newBlockHash[:]))
					gossiper.blocksMap[hex.EncodeToString(newBlockHash[:])] = newBlockStruct

					if gossiper.headBlock == nil ||
						gossiper.headBlock.len < newBlockStruct.len {
						// change head
						gossiper.headBlock = newBlockStruct

						// as head has been changed, we may need to delete some transactions
						newTransactions := make([]*data.KeyTransaction, 0)
						for _, transaction := range gossiper.pendingTransactions {
							if gossiper.GetPublicKey(transaction.GetName()) == nil {
								newTransactions = append(newTransactions, transaction)
							}
						}
						gossiper.pendingTransactions = newTransactions
					}

					// broadcast the block
					blockPublish.HopLimit -= 1
					if blockPublish.HopLimit > 0 {
						packet := &data.GossipPacket{KeyBlockPublish: blockPublish}
						gossiper.BroadCastPacket(packet)
					}
				}
				//fmt.Println("unlocking HandleNewBlock")
				gossiper.blockChainMutex.Unlock()
				//fmt.Println("unlocked HandleNewBlock")
			}
		}

	} else {
		//fmt.Println("wrong block?")
		if len(blockPublish.Block.Transactions) > 0 {
			//fmt.Println(*blockPublish.Block.Transactions[0].KeyPublish)
			//fmt.Println(*block.Block.Transactions[0].Secret)
		}

	}
}

func (gossiper *Gossiper) RequestBlock(hash [32]byte, dest string) {
	request := &data.BlockRequest{Origin: gossiper.Name, Destination: dest, HopLimit: hoplimit, HashValue: hash}
	packet := &data.GossipPacket{BlockRequest: request}
	gossiper.SendPacketViaRoutingTable(packet, dest)
}

func (gossiper *Gossiper) SendPacketViaRoutingTable(packet *data.GossipPacket, name string) {
	nxtHop, ok := gossiper.RoutingTable.Table[name]
	//If I don't know the next hop, discard the message
	if ok {
		gossiper.sendMessageToNeighbour(packet, nxtHop)
	}
}

func (gossiper *Gossiper) BroadCastPacket(packet *data.GossipPacket) {
	peers := gossiper.Neighbours.GetAllNeighboursWithException("")
	for _, peer := range peers {
		gossiper.sendMessageToNeighbour(packet, peer)
	}
}

func (gossiper *Gossiper) KeyMiningThread() {
	for true {
		//fmt.Println("locking KeyMiningThread")
		gossiper.blockChainMutex.Lock()
		//fmt.Println("locked KeyMiningThread")
		headHash := [32]byte{}
		thereWasChain := gossiper.headBlock != nil
		if thereWasChain {
			// A block has already been added
			headHash = gossiper.headBlock.Block.Hash()
		}

		// Save in listToPublish all pending transactions to create a new Block
		listToPublish := make([]data.KeyTransaction, len(gossiper.pendingTransactions))
		/*
			if len(listToPublish) == 0 {
				gossiper.blockChainMutex.Unlock()
				continue
			}*/
		for i, pointer := range gossiper.pendingTransactions {
			listToPublish[i] = *pointer
		}

		// Generating random nonce
		var nonce [32]byte
		rand.Read(nonce[:])

		// Creating new Block
		newBlock := &data.KeyBlock{PrevHash: headHash, Nonce: nonce, Transactions: listToPublish}
		newHash := newBlock.Hash()

		// Check validity, Prove of work is 16 bits (first 2 bytes) equals to 0
		valid := newHash[0] == 0 && newHash[1] == 0

		if valid {

			// new valid block found!
			//fmt.Println("FOUND-KEY-BLOCK", hex.EncodeToString(newHash[:]))

			if len(newBlock.Transactions) > 0 {
				//fmt.Println(*newBlock.Transactions[0].KeyPublish)
				//fmt.Println(*block.Block.Transactions[0].Secret)
			}

			// count length (prev + 1)
			newBlockStruct := &pairBlockLen{Block: newBlock, len: 1}
			prevStruct, prevExists := gossiper.blocksMap[hex.EncodeToString(newBlock.PrevHash[:])]
			if prevExists {
				newBlockStruct.len += prevStruct.len
			}

			// add to the chain
			//fmt.Println("ADDING BLOCK", hex.EncodeToString(newHash[:]))
			gossiper.blocksMap[hex.EncodeToString(newHash[:])] = newBlockStruct

			// change head
			gossiper.headBlock = newBlockStruct
			//s := "KEY-CHAIN " + hex.EncodeToString(newHash[:])
			//fmt.Println(s)

			// all pending transactions have been added, removing them
			gossiper.pendingTransactions = make([]*data.KeyTransaction, 0)

			// publish
			//fmt.Println("publishing new key block")
			blockPublish := &data.KeyBlockPublish{Block: newBlock, HopLimit: hoplimit, Origin: gossiper.Name}
			packet := &data.GossipPacket{KeyBlockPublish: blockPublish}
			gossiper.BroadCastPacket(packet)
			//fmt.Println("published")

			gossiper.printBlockChain()

			time.Sleep(100 * time.Millisecond)
		}
		//fmt.Println("unlocking KeyMiningThread")
		gossiper.blockChainMutex.Unlock()
		//fmt.Println("unlocked KeyMiningThread")

	}
}

func (gossiper *Gossiper) printBlockChain() {
	found := false
	blockStruct := gossiper.headBlock
	hasNext := blockStruct != nil
	s := "BLOCKCHAIN: "

	for !found && hasNext {
		block := blockStruct.Block
		blockHash := block.Hash()
		s += hex.EncodeToString(blockHash[:])
		for _, transaction := range block.Transactions {
			if transaction.IsKeyPublish() {
				s += ":" + transaction.GetName()
			} else {
				s += "|" + transaction.Secret.Origin
			}
		}
		s += " "
		blockStruct, hasNext = gossiper.blocksMap[hex.EncodeToString(block.PrevHash[:])]

	}
	if false {
		fmt.Println(s)
	}
}
