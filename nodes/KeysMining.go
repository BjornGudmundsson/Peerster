package nodes

import (
	"crypto/rsa"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/BjornGudmundsson/Peerster/data"
	"math/rand"
	"time"
)

type BlockRequest struct {
	Origin string
	Destination string
	HopLimit uint32
	HashValue [32]byte
}

type BlockReply struct {
	Origin string
	Destination string
	HopLimit uint32
	Block KeyBlock
}

type KeyTransaction struct {
	Name string
	Key rsa.PublicKey
}

type KeyBlock struct {
	PrevHash [32]byte
	Nonce [32]byte
	Transactions []KeyTransaction
}

type KeyPublish struct {
	Transaction *KeyTransaction
	HopLimit uint32
}

type KeyBlockPublish struct {
	Origin string
	Block *KeyBlock
	HopLimit uint32
}

type pairBlockLen struct {
	Block *KeyBlock
	len uint64
}


/**
	Returns the PublicKey associate to a name in case that it exists in the blockChain
	in case that the name has not a public key associated it returns nil.

 */
func (gossiper *Gossiper) GetPublicKey(name string) *rsa.PublicKey  {
	var key *rsa.PublicKey
	gossiper.blockChainMutex.Lock()

	// No blockChain initialised
	if gossiper.headBlock == nil {
		return key
	}

	found := false
	blockStruct := gossiper.headBlock
	hasNext := true

	for !found && hasNext{
		block := blockStruct.Block
		for _, transaction := range block.Transactions{
			if transaction.Name == name{
				found = true
				key = &transaction.Key
			}
		}
		if !found{
			blockStruct, hasNext = gossiper.blocksMap[hex.EncodeToString(block.PrevHash[:])]
		}
	}
	gossiper.blockChainMutex.Unlock()

	return key
}

/**
	Checks whether if a key associated to a name exists starting from a specific head of the chain
	requires to have the mutex lock
 */
func (gossiper *Gossiper) existsPublicKeyFromBlock(name string, hash string) bool  {

	// No blockChain initialised
	if gossiper.headBlock == nil {
		return false
	}

	found := false
	blockStruct,_ := gossiper.blocksMap[hash]
	hasNext := true

	for !found && hasNext{
		block := blockStruct.Block
		for _, transaction := range block.Transactions{
			if transaction.Name == name{
				found = true
			}
		}
		if !found{
			blockStruct, hasNext = gossiper.blocksMap[hex.EncodeToString(block.PrevHash[:])]
		}
	}

	return found

}

func (gossiper *Gossiper) checkInsidePendingTransactions(name string) bool  {
	found := false
	gossiper.blockChainMutex.Lock()
	for i := 0; i < len(gossiper.pendingTransactions) && !found; i++ {
		pendingTransaction := gossiper.pendingTransactions[i]
		found = pendingTransaction.Name == name
	}
	gossiper.blockChainMutex.Unlock()

	return !found
}

func (gossiper *Gossiper) PublishPublicKey(name string, key rsa.PublicKey) bool{
	// check if the name already has an associated key
	// neither in the pending transactions
	if gossiper.GetPublicKey(name) == nil && gossiper.checkInsidePendingTransactions(name	){
		// add the transaction to the pending transactions
		gossiper.blockChainMutex.Lock()
		newTransaction := &KeyTransaction{Name: name, Key: key}
		gossiper.pendingTransactions = append(gossiper.pendingTransactions, newTransaction)
		gossiper.blockChainMutex.Unlock()

		// send the transaction to neighbors
		keyPublish := &KeyPublish{Transaction: newTransaction, HopLimit: hoplimit}
		packet := &data.GossipPacket{KeyPublish: keyPublish}
		gossiper.BroadCastPacket(packet)

		return true
	}
	return false

}

func (gossiper *Gossiper) HandleBlockRequest(request *BlockRequest)  {
	if request.Destination == gossiper.Name{
		// the request is for me
		blockHashBytes := request.HashValue
		blockHashString := hex.EncodeToString(blockHashBytes[:])

		// block to be send
		gossiper.blockChainMutex.Lock()
		block, found := gossiper.blocksMap[blockHashString]
		gossiper.blockChainMutex.Unlock()

		if found{
			reply := &BlockReply{Destination: request.Origin, Origin: gossiper.Name, HopLimit: hoplimit, Block: *block.Block}
			// create gossip packet with reply and send it
			packet := &data.GossipPacket{BlockReply: reply}
			gossiper.SendPacketViaRoutingTable(packet, request.Origin)

		}

	} else{
		// the request is for another peer
		request.HopLimit -= 1
		if request.HopLimit > 0 {
			// create gossip packet with request and send it to Destination
			packet := &data.GossipPacket{BlockRequest: request}
			gossiper.SendPacketViaRoutingTable(packet, request.Destination)

		}
	}
}

func (gossiper *Gossiper)HandleNewBlock(blockPublish *KeyBlockPublish)  {
	newBlock := blockPublish.Block
	// check the prove of work
	newBlockHash := newBlock.Hash()
	if newBlockHash[0] == 0 && newBlockHash[1] == 0 {

		// if the prev block is the origin block, we must accept the new block.
		allZeros := true
		for _, b := range newBlock.PrevHash{
			allZeros = allZeros && b == 0
		}

		if allZeros {
			valid := true
			// check for duplications in the block
			for i, transaction := range newBlock.Transactions{
				for j := i + 1; j < len(newBlock.Transactions); j++ {
					valid = valid && transaction.Name != newBlock.Transactions[j].Name
				}
			}

			if valid {
				gossiper.blockChainMutex.Lock()
				newBlockStruct := &pairBlockLen{Block: newBlock, len: 1}
				// add to the chain
				gossiper.blocksMap[hex.EncodeToString(newBlockHash[:])] = newBlockStruct
				gossiper.blockChainMutex.Unlock()

				// broadcast the block
				blockPublish.HopLimit -= 1
				if blockPublish.HopLimit > 0 {
					packet := &data.GossipPacket{KeyBlockPublish: blockPublish}
					gossiper.BroadCastPacket(packet)
				}

			}

		}


		gossiper.blockChainMutex.Lock()

		prevBlockStruct, haveIt := gossiper.blocksMap[hex.EncodeToString(newBlock.PrevHash[:])]
		// if i don't have the prev block, ask for it.
		if !haveIt{
			gossiper.RequestBlock(newBlock.PrevHash, blockPublish.Origin)
		}

		askedTime := time.Now()
		// wait until i have the prev block (max 5 seg)
		for !haveIt && time.Now().Sub(askedTime) < 5 * time.Second{
			time.Sleep(10 * time.Millisecond)
			prevBlockStruct, haveIt = gossiper.blocksMap[hex.EncodeToString(newBlock.PrevHash[:])]

		}

		// if i haven't received it, do not accept the new block
		if haveIt{
			prevHashString := hex.EncodeToString(prevBlockStruct.Block.PrevHash[:])
			valid := true
			// check validity of all the transactions of the block (in the chain and no repeated)

			// check for duplications in the chain
			for _, transaction := range newBlock.Transactions{
				valid = valid && !gossiper.existsPublicKeyFromBlock(transaction.Name, prevHashString)
			}

			// check for duplications in the block
			for i, transaction := range newBlock.Transactions{
				for j := i + 1; j < len(newBlock.Transactions); j++ {
					valid = valid && transaction.Name != newBlock.Transactions[j].Name
				}
			}

			if valid{
				// count new length (prev + 1)
				newBlockStruct := &pairBlockLen{Block: newBlock, len: prevBlockStruct.len + 1}

				// add to the chain
				gossiper.blocksMap[hex.EncodeToString(newBlockHash[:])] = newBlockStruct

				if gossiper.headBlock.len < newBlockStruct.len{
					// change head
					gossiper.headBlock = newBlockStruct

					// as head has been changed, we may need to delete some transactions
					newTransactions := make([]*KeyTransaction, 0)
					for _, transaction := range gossiper.pendingTransactions{
						if gossiper.GetPublicKey(transaction.Name) == nil{
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
		}
		gossiper.blockChainMutex.Unlock()
	}
}


func (gossiper *Gossiper) RequestBlock(hash [32] byte, dest string)  {
	request := &BlockRequest{Origin: gossiper.Name, Destination: dest, HopLimit: hoplimit, HashValue: hash}
	packet := &data.GossipPacket{BlockRequest: request}
	gossiper.SendPacketViaRoutingTable(packet, dest)
}



func (gossiper *Gossiper)SendPacketViaRoutingTable(packet *data.GossipPacket, name string){
	nxtHop, ok := gossiper.RoutingTable.Table[name]
	//If I don't know the next hop, discard the message
	if ok {
		gossiper.sendMessageToNeighbour(packet, nxtHop)
	}
}

func (gossiper *Gossiper)BroadCastPacket(packet *data.GossipPacket)  {
	peers := gossiper.Neighbours.GetAllNeighboursWithException("")
	for _, peer := range peers {
		gossiper.sendMessageToNeighbour(packet, peer)
	}
}

func (gossiper *Gossiper)KeyMiningThread() {

	for true {

		gossiper.blockChainMutex.Lock()
		headHash := [32]byte{}

		thereWasChain := gossiper.headBlock != nil
		if thereWasChain {
			// A block has already been added
			headHash = gossiper.headBlock.Block.Hash()
		}

		// Save in listToPublish all pending transactions to create a new Block
		listToPublish := make([]KeyTransaction, len(gossiper.pendingTransactions))
		for i, pointer := range gossiper.pendingTransactions {
			listToPublish[i] = *pointer
		}

		// Generating random nonce
		var nonce [32]byte
		rand.Read(nonce[:])

		// Creating new Block
		newBlock := &KeyBlock{PrevHash: headHash, Nonce: nonce, Transactions: listToPublish}
		newHash := newBlock.Hash()

		// Check validity, Prove of work is 16 bits (first 2 bytes) equals to 0
		valid := newHash[0] == 0 && newHash[1] == 0

		if valid {

			// new valid block found!
			fmt.Println("FOUND-KEY-BLOCK", hex.EncodeToString(newHash[:]))

			// count length (prev + 1)
			newBlockStruct := &pairBlockLen{Block: newBlock, len: 1}
			prevStruct, prevExists := gossiper.blocksMap[hex.EncodeToString(newBlock.PrevHash[:])]
			if prevExists {
				newBlockStruct.len += prevStruct.len
			}

			// add to the chain
			gossiper.blocksMap[hex.EncodeToString(newHash[:])] = newBlockStruct


			// change head
			gossiper.headBlock = newBlockStruct
			s := "KEY-CHAIN" + hex.EncodeToString(newHash[:])
			fmt.Println(s)

			// all pending transactions have been added, removing them
			gossiper.pendingTransactions = make([]*KeyTransaction, 0)
			gossiper.blockChainMutex.Unlock()


			// publish
			fmt.Println("publishing new key block")
			blockPublish := &KeyBlockPublish{Block: newBlock, HopLimit: hoplimit, Origin: gossiper.Name}
			packet := &data.GossipPacket{KeyBlockPublish: blockPublish}
			gossiper.BroadCastPacket(packet)

		}
	}
}


func (b *KeyBlock) Hash() (out [32]byte) {
	h := sha256.New()
	h.Write(b.PrevHash[:])
	h.Write(b.Nonce[:])
	binary.Write(h,binary.LittleEndian, uint32(len(b.Transactions)))
	for _, t := range b.Transactions {
		th := t.Hash()
		h.Write(th[:])
	}
	copy(out[:], h.Sum(nil))
	return
}

func (t *KeyTransaction) Hash() (out [32]byte) {
	h := sha256.New()
	h.Write([]byte(t.Name))
	binary.Write(h,binary.LittleEndian, uint32(t.Key.E))
	h.Write(t.Key.N.Bytes())
	copy(out[:], h.Sum(nil))
	return
}
