package transactions

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"strings"
)

//BlockPublish is a struct
//representing publishing a
//block to the blockchain
type BlockPublish struct {
	Block    Block
	HopLimit uint32
}

//Block is a struct
//representing a block
//on the blockchain
type Block struct {
	PrevHash     [32]byte
	Nonce        [32]byte
	Transactions []TxPublish
}

//NewBlock returns a new block with the given parameters
func NewBlock(prev, nonce [32]byte, tx []TxPublish) Block {
	return Block{
		PrevHash:     prev,
		Nonce:        nonce,
		Transactions: tx,
	}
}

//PrintBlock returns the print string for a block
func (b *Block) PrintBlock() string {
	temp := "["
	h := b.Hash()
	hash := hex.EncodeToString(h[:])
	prevHash := hex.EncodeToString(b.PrevHash[:])
	temp = temp + hash + ":" + prevHash + ":"
	slice := make([]string, 0)
	for _, tx := range b.Transactions {
		slice = append(slice, tx.File.Name)
	}
	files := strings.Join(slice, ",")
	temp = temp + files + "]"
	return temp
}

//CompareBlock compares if two block share any
//transactions.
func (b *Block) CompareBlock(block Block) bool {
	for _, tx := range b.Transactions {
		if block.HasTransaction(tx) {
			return true
		}
	}
	return false
}

//HasTransaction checks if a block has this
//particular transaction logged.
func (b *Block) HasTransaction(t TxPublish) bool {
	for _, tx := range b.Transactions {
		equal := t.Compare(tx)
		if equal {
			return true
		}
	}
	return false
}

//Hash hashes a block
func (b *Block) Hash() (out [32]byte) {
	h := sha256.New()
	h.Write(b.PrevHash[:])
	h.Write(b.Nonce[:])
	binary.Write(h, binary.LittleEndian, uint32(len(b.Transactions)))
	for _, t := range b.Transactions {
		th := t.Hash()
		h.Write(th[:])
	}
	copy(out[:], h.Sum(nil))
	return
}

//IsValidBlock checks if a block is valid
func IsValidBlock(block Block) bool {
	// fmt.Println("Started hashing")
	hash := block.Hash()
	// fmt.Println("stopped hashing")
	if hash[0] == byte(0) && hash[1] == byte(0) {
		return true
	}
	return false
}

//MineABlock takes in the paremeters of a block and generates a new
//block with a new nonce
func MineABlock(transactions []TxPublish, prevHash [32]byte) Block {
	nonce := GetNonce()
	minedBlock := NewBlock(prevHash, nonce, transactions)
	return minedBlock
}
