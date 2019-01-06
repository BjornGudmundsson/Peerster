package nodes

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/BjornGudmundsson/Peerster/data/peersterCrypto"

	"github.com/BjornGudmundsson/Peerster/data"
)

const chunkSize uint64 = 8
const keySize int = 16

//HandleNewFile takes in a multipart.Fileheader and a multipart file and processes
//the file in such a way that is specified in the homework description. Adds all the
//chunks in a map that keeps track of which chunk corresponds to which part of text
//whithout respect for which file that bit of text belongs to.
func (g *Gossiper) HandleNewFile(fh *multipart.FileHeader, f multipart.File) {
	fSize := uint64(fh.Size)
	fileData := make([]byte, fSize)
	f.Read(fileData)
	mod := fSize % uint64(keySize)
	if mod != 0 {
		pad := keySize - int(mod)
		buf := make([]byte, pad)
		fileData = append(fileData, buf...)
	}
	IV := peersterCrypto.GetIV(keySize)
	key := peersterCrypto.GetKey(keySize)
	encryptedFileData, e := peersterCrypto.EncryptBlocks(fileData, key, IV)
	if e != nil {
		log.Fatal(e)
	}
	reader := bytes.NewReader(encryptedFileData)
	div := float64(fSize) / float64(chunkSize)
	metafile := make([]byte, 0)
	sizeInChunks := uint64(math.Ceil(div))
	for i := uint64(0); i < sizeInChunks; i++ {
		buf := make([]byte, chunkSize)
		n, _ := reader.Read(buf)
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
	//go g.SpreadMetaFile(metafile)
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

//HandleNewOSFile takes in a filename and gets
//the corresponding named file from the shared file of
//the gossiper and splits it into chunks and adds
//the chunks and hash to the memory of the Gossiper
func (g *Gossiper) HandleNewOSFile(fn string) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	temp := []string{dir, "_SharedFiles", fn}
	fp := strings.Join(temp, "/")
	f, e := os.Open(fp)
	if e != nil {
		log.Fatal(e)
	}
	fStat, e := f.Stat()
	if e != nil {
		log.Fatal()
	}
	fSize := uint64(fStat.Size())
	div := float64(fSize) / float64(chunkSize)
	metafile := make([]byte, 0)
	sizeInChunks := uint64(math.Ceil(div))
	fmt.Println("Size: ", fSize)
	fmt.Println("sizeInChunks", sizeInChunks)
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
	go g.SpreadMetaFile(metafile)
	hashMF := sha256.Sum256(metafile)
	hashMFBs := make([]byte, 0)
	for _, b := range hashMF {
		hashMFBs = append(hashMFBs, b)
	}
	hexHash := hex.EncodeToString(hashMFBs)
	mf := &data.MetaData{
		FileName:       fn,
		FileSize:       fSize,
		MetaFile:       metafile,
		HashOfMetaFile: hexHash,
	}
	g.Files[mf.HashOfMetaFile] = *mf
	/*
		file := transactions.NewFile(fn, fSize, hashMFBs)
		tx := transactions.NewTransaction(file, hoplimit)
		g.BroadCastTxPublish(tx, "")
		g.TransactionBuffer.AddTx(tx)
	*/
}

//CreateChunkmap takes in the metafile corresponding to a file
//and returns a chunkmap of all the chunks that this gossiper
//has for this file.
func (g *Gossiper) CreateChunkmap(metafile []byte) []uint64 {
	chunks := g.Chunks
	//I know the Metafile is a multiple of 32
	n := uint64(len(metafile)) / 32
	temp := make([]uint64, 0)
	for i := uint64(0); i < n; i++ {
		j := i + 1
		//Each chunks is exactly 32 bytes. That is why I have 32
		hx := hex.EncodeToString(metafile[i*32 : j*32])
		_, ok := chunks[hx]
		if ok {
			temp = append(temp, j)
		}
	}
	return temp
}
