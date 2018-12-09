package nodes

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/BjornGudmundsson/Peerster/data"
)

//DownloadingFile is a process that handles downloading a file.
//It will not stop sending requests till the file has
//completed its download. When starting a new
//DownloadingFile thread it will assume that this is a file
//you do not have.
func (g *Gossiper) DownloadingFile(filename string) {
	var meta []byte
	var lastKnownDestination string
	metadata := g.Files[filename]
	metafile := metadata.MetaFile
	metafileHash := metadata.HashOfMetaFile
	//All chunks will go through this channel
	chunkChannel := make(chan data.DataReply)
	var gotMetaFile bool
	var nxtChunk []byte
	if metadata.MetaFile == nil {
		lastKnownDestination = g.ChunkToPeer.GetRandomOwnerOfMetafile(metafileHash)
		meta = make([]byte, len(metadata.MetaFile))
		copy(meta, metadata.MetaFile)
		nxtChunk, _ = hex.DecodeString(metafileHash)
		mfh := metafileHash
		g.HandlerDataReplies.AddChunk(mfh, chunkChannel)
		g.SendDataRequest(lastKnownDestination, nxtChunk)
	} else {
		isFullyDownloaded, nxtIndex := g.HasAllChunksOfFile(metadata.MetaFile)
		if isFullyDownloaded {
			g.ReconstructFile(filename, metadata.MetaFile)
			return
		}
		nxtChunk = metafile[nxtIndex : nxtIndex+32]
		lastKnownDestination = g.ChunkToPeer.GetRandomOwnerOfChunk(metafileHash, (nxtIndex/32)+1)
		g.HandlerDataReplies.AddMetafile(metafile, chunkChannel)
		g.SendDataRequest(lastKnownDestination, nxtChunk)
	}
	for {
		ticker := time.NewTicker(5 * time.Second)
		select {
		case datareply := <-chunkChannel:
			hash := datareply.HashValue
			hashHex := hex.EncodeToString(hash)
			if hashHex == metafileHash {
				fmt.Println("DOWNLOADING metafile")
				mf := datareply.Data
				meta = make([]byte, len(mf))
				copy(meta, mf)
				_, i := g.HasAllChunksOfFile(meta)
				lastKnownDestination = g.ChunkToPeer.GetRandomOwnerOfChunk(metafileHash, (i/32)+1)
				if g.Files[filename].MetaFile != nil {
					nxtChunk = meta[i : i+32]
					g.SendDataRequest(lastKnownDestination, nxtChunk)
					continue
				}
				if !gotMetaFile {
					metadata.MetaFile = mf
					g.Files[filename] = metadata
					g.HandlerDataReplies.AddMetafile(mf, chunkChannel)
					nxtChunk = mf[i : i+32]
					g.SendDataRequest(lastKnownDestination, nxtChunk)
					gotMetaFile = !gotMetaFile
				}
			} else {
				g.Chunks[hashHex] = string(datareply.Data)
				done, index := g.HasAllChunksOfFile(meta)
				if done {
					g.ReconstructFile(filename, meta)
					g.HandlerDataReplies.DeleteMetafile(meta)
					return
				}
				lastKnownDestination = g.ChunkToPeer.GetRandomOwnerOfChunk(metafileHash, (index/32)+1)
				fmt.Println("DOWNLOADING chunk ", index/32)
				nxtChunk = meta[index : index+32]
				g.SendDataRequest(lastKnownDestination, nxtChunk)
			}
		case <-ticker.C:
			//I'll keep persisting until I get a datarequest
			//g.SendDataRequest(lastKnownDestination, nxtChunk)
			g.SendDataRequest(lastKnownDestination, nxtChunk)
		}
	}
}

//ReconstructFile takes in the filename and metafile and reconstructs the
//file from that and adds it to the downloaded function.
func (g *Gossiper) ReconstructFile(filename string, metafile []byte) {
	//This reconstructs the file. Don't feel like writing it right now.
	n := len(metafile) / 32
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	temp := []string{dir, "/_Downloads", filename}
	fp := strings.Join(temp, "/")
	f, e := os.Create(fp)
	if e != nil {
		log.Fatal(e)
	}
	for i := 0; i < n; i++ {
		j := i + 1
		chunk := hex.EncodeToString(metafile[i*32 : j*32])
		fmt.Fprintf(f, g.Chunks[chunk])
	}
	fmt.Printf("Reconstructed file %v", filename)
}
