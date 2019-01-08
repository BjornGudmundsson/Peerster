package nodes

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/BjornGudmundsson/Peerster/data/peersterCrypto"

	"github.com/BjornGudmundsson/Peerster/data"
)

//DownloadingFile is a process that handles downloading a file.
//It will not stop sending requests till the file has
//completed its download. When starting a new
//DownloadingFile thread it will assume that this is a file
//you do not have.
func (g *Gossiper) DownloadingFile(filename string) {
	var lastKnownDestination string
	metadata := g.Files[filename]
	metafile := metadata.MetaFile
	metafileHash := metadata.HashOfMetaFile
	//All chunks will go through this channel
	chunkChannel := make(chan data.DataReply)
	pmf, has := g.Chunks[metafileHash]
	if has {
		v := []byte(pmf)
		metadata.MetaFile = v
		g.Files[filename] = metadata
		g.PopulateFromMetafile(metadata.MetaFile, filename, metafileHash)
	}
	var nxtChunk []byte
	if metadata.MetaFile == nil {
		lastKnownDestination = g.ChunkToPeer.GetRandomOwnerOfMetafile(metafileHash)
		fmt.Println("Last known destination: ", lastKnownDestination)
		nxtChunk, _ = hex.DecodeString(metafileHash)
		mfh := metafileHash
		g.HandlerDataReplies.AddChunk(mfh, chunkChannel)
		g.SendDataRequest(lastKnownDestination, nxtChunk)
	} else {
		fmt.Println("Had the metafile 2: ", hex.EncodeToString(metadata.MetaFile[0:32]))
		isFullyDownloaded, nxtIndex := g.HasAllChunksOfFile(metadata.MetaFile)
		if isFullyDownloaded {
			g.ReconstructFile(filename, metadata.MetaFile)
			return
		}
		nxtChunk = metadata.MetaFile[nxtIndex : nxtIndex+32]
		lastKnownDestination = g.ChunkToPeer.GetRandomOwnerOfChunk(metafileHash, (nxtIndex/32)+1)
		g.HandlerDataReplies.AddMetafile(metafile, chunkChannel)
		g.SendDataRequest(lastKnownDestination, nxtChunk)
		g.Chunks[metadata.HashOfMetaFile] = string(metadata.MetaFile)
		fmt.Println("Populating")
		g.PopulateFromMetafile(metadata.MetaFile, filename, metafileHash)
	}
	for {
		ticker := time.NewTicker(2 * time.Second)
		select {
		case datareply := <-chunkChannel:
			hash := datareply.HashValue
			hashHex := hex.EncodeToString(hash)
			if hashHex == metafileHash {
				fmt.Println("Got metafile 2222")
				metadata.MetaFile = datareply.Data
				fmt.Println("Metafilehash: ", metafileHash)
				fmt.Println("reply hash: ", hex.EncodeToString(datareply.HashValue))
				fmt.Println("MEtafile: ", string(datareply.Data[0:32]))
				g.Chunks[metafileHash] = string(datareply.Data)
				g.HandlerDataReplies.AddMetafile(metadata.MetaFile, chunkChannel)
				g.Files[filename] = metadata
				g.PopulateFromMetafile(metadata.MetaFile, filename, metafileHash)
				_, i := g.HasAllChunksOfFile(metadata.MetaFile)
				for {
					lastKnownDestination = g.ChunkToPeer.GetRandomOwnerOfChunk(metafileHash, (i/32)+1)
					if lastKnownDestination == g.Name {
						i = i + 32
						continue
					}
					break
				}
				fmt.Println(lastKnownDestination)
				nxtChunk = metadata.MetaFile[i : i+32]
				g.SendDataRequest(lastKnownDestination, nxtChunk)
			} else {
				fmt.Println("Getting a chunk")
				md := g.Files[filename]
				g.Chunks[hashHex] = string(datareply.Data)
				done, index := g.HasAllChunksOfFile(md.MetaFile)
				if done {
					g.ReconstructFile(filename, md.MetaFile)
					g.HandlerDataReplies.DeleteMetafile(md.MetaFile)
					return
				}
				for {
					lastKnownDestination = g.ChunkToPeer.GetRandomOwnerOfChunk(metafileHash, (index/32)+1)
					if lastKnownDestination == g.Name {
						index = index + 32
						continue
					}
					break
				}
				nxtChunk = md.MetaFile[index : index+32]
				g.SendDataRequest(lastKnownDestination, nxtChunk)
			}
		case <-ticker.C:
			//I'll keep persisting until I get a datarequest
			//g.SendDataRequest(lastKnownDestination, nxtChunk)
			//lastKnownDestination = g.ChunkToPeer.GetRandomOwnerOfChunk()
			fmt.Println("Ticker expired")
			fmt.Println(lastKnownDestination)
			g.SendDataRequest(lastKnownDestination, nxtChunk)
		}
	}
}

//ReconstructFile takes in the filename and metafile and reconstructs the
//file from that and adds it to the downloaded function.
func (g *Gossiper) ReconstructFile(filename string, metafile []byte) {
	fmt.Println("Reconstructing file")
	//This reconstructs the file. Don't feel like writing it right now.
	n := len(metafile) / 32
	fmt.Println("Length of the metafile: ", n)
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	temp := []string{dir, "/_Downloads", filename}
	fp := strings.Join(temp, "/")
	f, e := os.Create(fp)
	if e != nil {
		fmt.Println("os create")
		log.Fatal(e)
	}
	md := g.Files[filename]
	fmt.Println("metadata reconstruction: ")
	IV := md.IV
	Key := md.Key
	buffer := make([]byte, 0)
	for i := 0; i < n; i++ {
		j := i + 1
		chunk := hex.EncodeToString(metafile[i*32 : j*32])
		buffer = append(buffer, []byte(g.Chunks[chunk])...)
	}
	fmt.Println("created the buffer")
	fmt.Println("multiple of blocksize", len(buffer)%keySize == 0, len(buffer))
	decryptedFile, e := peersterCrypto.DecryptCiphertext(buffer, Key, IV)
	if e != nil {
		fmt.Println("Got an error")
		log.Fatal(e)
	}
	fmt.Fprintf(f, string(decryptedFile))
	fmt.Printf("Reconstructed file %v", filename)
}
