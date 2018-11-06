package nodes

import (
	"crypto/sha256"
	"encoding/hex"
	"math"
	"mime/multipart"

	"github.com/BjornGudmundsson/Peerster/data"
)

const chunkSize uint64 = 3

//HandleNewFile takes in a multipart.Fileheader and a multipart file and processes
//the file in such a way that is specified in the homework description. Adds all the
//chunks in a map that keeps track of which chunk corresponds to which part of text
//whithout respect for which file that bit of text belongs to.
func (g *Gossiper) HandleNewFile(fh *multipart.FileHeader, f multipart.File) {
	fSize := uint64(fh.Size)
	div := float64(fSize) / float64(chunkSize)
	metafile := make([]byte, 0)
	sizeInChunks := uint64(math.Ceil(div))
	for i := uint64(0); i < sizeInChunks; i++ {
		buf := make([]byte, chunkSize)
		n, _ := f.Read(buf)
		chunkString := string(buf)[0:n]
		hash := sha256.Sum256([]byte(chunkString))
		tempbs := make([]byte, 0)
		for _, b := range hash {
			tempbs = append(tempbs, b)
		}
		hxhash := hex.EncodeToString(tempbs)
		g.Chunks[hxhash] = chunkString
		metafile = append(metafile, tempbs...)
	}
	hashMF := sha256.Sum256(metafile)
	hashMFBs := make([]byte, 0)
	for _, b := range hashMF {
		hashMFBs = append(hashMFBs, b)
	}
	hexHash := hex.EncodeToString(hashMFBs)
	mf := &data.MetaData{
		FileName:       fh.Filename,
		FileSize:       fSize,
		MetaFile:       metafile,
		HashOfMetaFile: hexHash,
	}
	g.Files[mf.HashOfMetaFile] = *mf
}
