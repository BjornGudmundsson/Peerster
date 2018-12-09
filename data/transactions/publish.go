package transactions

import (
	"crypto/sha256"
	"encoding/binary"
)

//TxPublish is a struct
//representing the publishing
//of a File being indexed
type TxPublish struct {
	File     File
	HopLimit uint32
}

//Compare compares if two transactions are equal
func (t TxPublish) Compare(txp TxPublish) bool {
	return txp.File.Name == t.File.Name
}

//NewTransaction returns a new transaction of a given file.
func NewTransaction(f File, hopLimit uint32) TxPublish {
	return TxPublish{
		File:     f,
		HopLimit: hopLimit,
	}
}

//Hash returns the hash of a transaction
func (t *TxPublish) Hash() (out [32]byte) {
	h := sha256.New()
	binary.Write(h, binary.LittleEndian,
		uint32(len(t.File.Name)))
	h.Write([]byte(t.File.Name))
	h.Write(t.File.MetafileHash)
	copy(out[:], h.Sum(nil))
	return
}

//File is a struct representing
//a file in a transaction.
type File struct {
	Name         string
	Size         uint64
	MetafileHash []byte
}

//NewFile returns a new file with the function parameters as entries.
func NewFile(name string, size uint64, metafilehash []byte) File {
	return File{
		Name:         name,
		Size:         size,
		MetafileHash: metafilehash,
	}
}
