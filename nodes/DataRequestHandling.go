package nodes

import (
	"encoding/hex"
	"fmt"
	"time"

	"github.com/BjornGudmundsson/Peerster/data"
)

//I'll wait for this amout of seconds for the
// next chunk. If this many seconds pass I'll stop the download
const waitForChunk int = 5

func (g *Gossiper) DownLoadAFile(fn string, mf []byte, dst string) {
	//This is the current next chunk
	nxtChunk := mf
	//count how many seconds have passed
	counter := 0
	g.dataReplyHandler.Start(fn, hex.EncodeToString(mf))
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
		fmt.Println("This destination did not exist")
		return
	}
	g.sendMessageToNeighbour(gp, nxtHop)
	fmt.Println("Sending datarequest")
	for {
		finished := g.dataReplyHandler.IsFinished()
		if finished {
			fmt.Println("Finished downloading the file")
			fmt.Println("I got this many chunks: ", len(g.dataReplyHandler.GetCurrentChunks()))
			g.dataReplyHandler.Clear()
			return
		}
		if counter > waitForChunk {
			fmt.Println("Timed out waitinf for file")
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
	fmt.Println("Got a data request")
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
	fmt.Println("sending this text")
	if ok {
		dr := data.DataReply{}
		hashBytes, e := hex.DecodeString(hexHash)
		if e != nil {
			fmt.Println("Could not convrt hex string to byte slice")
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
