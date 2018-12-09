package blockchain

import (
	"encoding/hex"
	"fmt"
	"sync"

	"github.com/BjornGudmundsson/Peerster/data/transactions"
)

//Link is a struct that holds
//all of the info about a link
//in a blockchain
type Link struct {
	PrevHash [32]byte
	Index    uint64
	Block    transactions.Block
}

//NewLink is an abstraction of creating a new link
//structure.
func NewLink(prevHash [32]byte, index uint64, block transactions.Block) Link {
	return Link{
		PrevHash: prevHash,
		Index:    index,
		Block:    block,
	}
}

var genesis [32]byte

//BlockChain contains all the info
//about the blockchain. LongestIndex
//is the length of the current longest chain,
//longest hash is the hash of the head of the
//longest chain.
type BlockChain struct {
	LongestIndex uint64
	LongestHash  [32]byte
	Chain        map[string]Link
	mux          sync.Mutex
}

//NewBlockChain returns a new empty BlockChain
func NewBlockChain() *BlockChain {
	chain := make(map[string]Link)
	return &BlockChain{
		LongestIndex: 0,
		Chain:        chain,
	}
}

//HasBlock check if a block is present in the blockchain
func (bc *BlockChain) HasBlock(block transactions.Block) bool {
	hash := block.Hash()
	hexHash := hex.EncodeToString(hash[:])
	_, ok := bc.Chain[hexHash]
	return ok
}

//GetBlockByHash returns the block that has this hash.
func (bc *BlockChain) GetBlockByHash(hash [32]byte) (transactions.Block, bool) {
	hx := hex.EncodeToString(hash[:])
	link, ok := bc.Chain[hx]
	return link.Block, ok
}

func (bc *BlockChain) FindCommonBlock(hash, prev [32]byte) (rewind uint64, common transactions.Block) {
	//var genesis [32]byte
	//genesisHex := hex.EncodeToString(genesis[:])
	prevHex := hex.EncodeToString(prev[:])
	hashHex := hex.EncodeToString(hash[:])
	_, hasPrev := bc.GetBlockByHash(prev)
	prevBlock, hasHash := bc.GetBlockByHash(hash)
	if !hasPrev || !hasHash {
		//Something went wrong
		fmt.Println("something went wrong")
		return 0, prevBlock
	}
	var i uint64
	if prevHex == hashHex {
		rewind = i
		common, _ = bc.GetBlockByHash(hash)
		return
	}
	link1, ok1 := bc.Chain[hashHex]
	link2, ok2 := bc.Chain[prevHex]
	for ok1 && ok2 {
		i++
		hash1 := link1.Block.Hash()
		hash2 := link2.Block.Hash()
		hx1 := hex.EncodeToString(hash1[:])
		hx2 := hex.EncodeToString(hash2[:])
		if hx1 == hx2 {
			rewind = i
			common = link1.Block
			return
		}
		nxtHash1 := link1.Block.PrevHash
		nxtHash2 := link2.Block.PrevHash
		hxHash1 := hex.EncodeToString(nxtHash1[:])
		hxHash2 := hex.EncodeToString(nxtHash2[:])
		link1, ok1 = bc.Chain[hxHash1]
		link2, ok2 = bc.Chain[hxHash2]
	}
	common = link1.Block
	rewind = i
	return
}

//AddBlock adds a blockhash to the Blockchain with the
//appropriate index.
func (bc *BlockChain) AddBlock(hash, prevHash [32]byte, block transactions.Block) (fork bool, newChain bool, oldHead [32]byte, commonBlock transactions.Block) {
	bc.mux.Lock()
	defer bc.mux.Unlock()
	prevHex := hex.EncodeToString(prevHash[:])
	hashHex := hex.EncodeToString(hash[:])
	link, ok := bc.Chain[prevHex]
	if !ok {
		l := NewLink(prevHash, 1, block)
		bc.Chain[hashHex] = l
		if bc.LongestIndex == 0 {
			nxt := bc.LongestIndex
			bc.LongestIndex = nxt + 1
			newChain = false
			oldHead = bc.LongestHash
			bc.LongestHash = hash
		} else {
			fmt.Println("FORK-SHORTER ", hashHex)
		}
		fork = true
		bc.PrintLongestChain()
		return
	}
	if prevHex != hex.EncodeToString(bc.LongestHash[:]) {
		fork = true
		if bc.HasBlockByHash(prevHash) {
			fmt.Println("FORK-SHORTER ", hashHex)
		}
	}
	nxtIndex := link.Index + 1
	if nxtIndex > bc.LongestIndex {
		if prevHex != hex.EncodeToString(bc.LongestHash[:]) {
			newChain = true
			oldHead = bc.LongestHash
			var rewind uint64
			rewind, commonBlock = bc.FindCommonBlock(oldHead, prevHash)
			fmt.Printf("\nFORK-LONGER rewind %v blocks\n", rewind)
		}
		bc.LongestIndex = nxtIndex
		bc.LongestHash = hash
	}
	l := NewLink(prevHash, nxtIndex, block)
	bc.Chain[hashHex] = l
	bc.PrintLongestChain()
	return
}

//PrintLongestChain prints the longest chain
func (bc *BlockChain) PrintLongestChain() {
	temp := "CHAIN"
	head := bc.LongestHash
	hxHead := hex.EncodeToString(head[:])
	link, ok := bc.Chain[hxHead]
	for ok {
		block := link.Block
		hyphen := " - "
		printBlock := block.PrintBlock()
		temp = temp + hyphen + printBlock
		nxtBlock := hex.EncodeToString(link.PrevHash[:])
		link, ok = bc.Chain[nxtBlock]
	}
	fmt.Println(temp)
}

//HasTransaction check if a transaction is present on the blockchain.
func (bc *BlockChain) HasTransaction(t transactions.TxPublish) bool {
	if bc.LongestIndex == 0 {
		return false
	}
	currentHead := bc.LongestHash
	hexHead := hex.EncodeToString(currentHead[:])
	link, ok := bc.Chain[hexHead]
	for ok {
		block := link.Block
		hasTransaction := block.HasTransaction(t)
		if hasTransaction {
			return true
		}
		nxtHash := hex.EncodeToString(link.PrevHash[:])
		link, ok = bc.Chain[nxtHash]
	}
	return false
}

//HasBlockByHash allows for finding a block on the blockhain by its hash
func (bc *BlockChain) HasBlockByHash(hash [32]byte) bool {
	if bc.LongestIndex == 0 {
		return false
	}
	hexHash := hex.EncodeToString(hash[:])
	currentHead := bc.LongestHash
	hexHead := hex.EncodeToString(currentHead[:])
	link, ok := bc.Chain[hexHead]
	for ok {
		block := link.Block
		nHash := block.Hash()
		hexnHash := hex.EncodeToString(nHash[:])
		if hexHash == hexnHash {
			return true
		}
		nxtHash := hex.EncodeToString(link.PrevHash[:])
		link, ok = bc.Chain[nxtHash]
	}
	return false
}

//GetTransactionsFromFork allows for finding all the transaction on a given chain or fork
func (bc *BlockChain) GetTransactionsFromFork(hash [32]byte) []transactions.TxPublish {
	temp := make([]transactions.TxPublish, 0)
	if bc.LongestIndex == 0 {
		return nil
	}
	hexHash := hex.EncodeToString(hash[:])
	link, ok := bc.Chain[hexHash]
	for ok {
		block := link.Block
		temp = append(temp, block.Transactions...)
		nxtHash := hex.EncodeToString(link.PrevHash[:])
		link, ok = bc.Chain[nxtHash]
	}
	return temp
}

//ValidateBlockFromFork validates if a block is valid on the chain that it extends
func (bc *BlockChain) ValidateBlockFromFork(block transactions.Block) bool {
	if bc.LongestIndex == 0 {
		return true
	}
	prevHex := hex.EncodeToString(block.PrevHash[:])
	link, ok := bc.Chain[prevHex]
	for ok {
		nxtBlock := link.Block
		if block.CompareBlock(nxtBlock) {
			return false
		}
	}
	return true
}

//GetTransactionsToCommonBlock gets all the lost transactions from a fork up to the common head
func (bc *BlockChain) GetTransactionsToCommonBlock(oldhead [32]byte, common transactions.Block) []transactions.TxPublish {
	commonHash := common.Hash()
	commonHex := hex.EncodeToString(commonHash[:])
	temp := make([]transactions.TxPublish, 0)
	hx := hex.EncodeToString(oldhead[:])
	link, ok := bc.Chain[hx]
	for ok {
		block := link.Block
		hash := block.Hash()
		hashHex := hex.EncodeToString(hash[:])
		if hashHex == commonHex {
			return temp
		}
		temp = append(temp, block.Transactions...)
		prevHex := hex.EncodeToString(block.PrevHash[:])
		link, ok = bc.Chain[prevHex]
	}
	return temp
}
