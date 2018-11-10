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

//I'll wait for this amout of seconds for the
// next chunk. If this many seconds pass I'll stop the download
const waitForChunk int = 8 * 1024

func (g *Gossiper) DownLoadAFile(fn string, mf []byte, dst string) {
	//This is the current next chunk
	nxtChunk := mf
	//count how many seconds have passed
	counter := 0
	fs, ok := g.DownloadState[fn]
	if ok {
		l := len(fs.CurrentChunks) * 32
		g.dataReplyHandler.Start(fn, hex.EncodeToString(mf), fs.MetaFile, l)
	} else {
		g.dataReplyHandler.Start(fn, hex.EncodeToString(mf), nil, 0)
	}
	dr := &data.DataRequest{
		Origin:      g.Name,
		Destination: dst,
		HopLimit:    hoplimit,
		HashValue:   mf,
	}
	gp := &data.GossipPacket{
		DataRequest: dr,
	}
	nxtHop, ok := g.RoutingTable.Table[dst]
	if !ok {
		//This destination is not registered
		return
	}
	g.sendMessageToNeighbour(gp, nxtHop)
	for {
		finished := g.dataReplyHandler.IsFinished()
		if finished {
			fmt.Println("")
			fmt.Printf("Reconstructed file %v", g.dataReplyHandler.Name)
			chunks := g.DownloadState[fn].CurrentChunks
			dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
			if err != nil {
				log.Fatal(err)
			}
			temp := []string{dir, "/_Downloads", g.dataReplyHandler.Name}
			fp := strings.Join(temp, "/")
			f, e := os.Create(fp)
			if e != nil {
				log.Fatal(e)
			}
			for _, chunk := range chunks {
				fmt.Fprintf(f, g.Chunks[chunk])
			}
			g.dataReplyHandler.Clear()
			return
		}
		if counter > waitForChunk {
			g.dataReplyHandler.Clear()
			return
		}
		latestChunk := g.dataReplyHandler.NextChunk()
		d, _ := hex.DecodeString(latestChunk)
		if !data.Compare(nxtChunk, d) {
			nxtChunk = d
			counter = 0
			continue
		}
		counter++
		time.Sleep(1 * time.Second)
	}
}

func (g *Gossiper) handleDataRequestMessage(msg GossipAddress) {
	//refactor this function at a good opportunity. This is too long
	gp := data.GossipPacket{}
	req := msg.Msg.DataRequest
	g.RoutingTable.UpdateRoutingTable(req.Origin, msg.Addr)
	hash := req.HashValue
	hexHash := hex.EncodeToString(hash)
	metafile, ok := g.Files[hexHash]
	if ok {
		//Send the metafile
		mf := metafile.HashOfMetaFile
		mfbs, e := hex.DecodeString(mf)
		if e != nil {
			//something went wrong
			return
		}
		dr := data.DataReply{}
		dr.HashValue = mfbs
		dr.HopLimit = hoplimit
		dr.Destination = req.Origin
		dr.Origin = req.Destination
		dr.Data = metafile.MetaFile
		gp.DataReply = &dr
		nxtHop, ok := g.RoutingTable.Table[req.Origin]
		//If I don't know the next hop, discard the message
		if !ok {
			//this really should not happen
			return
		}
		g.sendMessageToNeighbour(&gp, nxtHop)
		return
	}
	txt, ok := g.Chunks[hexHash]
	if ok {
		dr := data.DataReply{}
		hashBytes, e := hex.DecodeString(hexHash)
		if e != nil {
			//This really should not happen and I have no idea
			//what to do if it happens so I am doing nothing
			return
		}
		dr.HashValue = hashBytes
		dr.HopLimit = hoplimit
		dr.Origin = req.Destination
		dr.Destination = req.Origin
		dr.Data = []byte(txt)
		gp.DataReply = &dr
		nxtHop, ok := g.RoutingTable.Table[req.Origin]
		//If I don't know the next hop, discard the message
		if !ok {
			//this really should not happen
			return
		}
		g.sendMessageToNeighbour(&gp, nxtHop)
		return
	}

	nxtHop, ok := g.RoutingTable.Table[req.Destination]
	//If I don't know the next hop, discard the message
	if !ok {
		//I do not have this in my routing table
		return
	}
	req.HopLimit = req.HopLimit - 1
	if req.HopLimit == 0 {
		//Dropping packet because it has exceeded the hoplimit
		return
	}
	gp.DataRequest = req
	g.sendMessageToNeighbour(&gp, nxtHop)
}
