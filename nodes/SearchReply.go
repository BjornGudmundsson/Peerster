package nodes

import (
	"encoding/hex"
	"time"

	"github.com/BjornGudmundsson/Peerster/data"
)

//HandleSearchReply handles all incoming SearchReply messages
//coming to the Gossiper
func (g *Gossiper) HandleSearchReply(msg GossipAddress) {
	sr := msg.Msg.SearchReply
	dst := sr.Destination
	g.RoutingTable.UpdateRoutingTable(sr.Origin, msg.Addr)
	if g.Name == dst {
		//Handle the datareply
		g.ProcessDataReply(*sr)
	} else {
		hoplimit := sr.HopLimit
		if hoplimit == 1 {
			//The hoplimit has been exceeded
			return
		}
		nxtHop, ok := g.RoutingTable.Table[dst]
		if !ok {
			//This is a weird situation but whatever. I won't handle it.
			return
		}
		sr.HopLimit = sr.HopLimit - 1
		gp := &data.GossipPacket{
			SearchReply: sr,
		}
		g.sendMessageToNeighbour(gp, nxtHop)
	}
}

//ProcessDataReply handles the local processing of data
//if the data is not supposed to be forwarded any further
func (g *Gossiper) ProcessDataReply(msg data.SearchReply) {
	src := msg.Origin
	results := msg.Results
	for _, result := range results {
		metafilehash := result.MetafileHash
		hexHash := hex.EncodeToString(metafilehash)
		hasMetafile := g.StateFileFinder.HasMetaFile(hexHash)
		if !hasMetafile {
			go g.RequestMetaFile(src, *result)
		}
	}
}

//RequestMetaFile is a method bound to a Gossiper that takes in an origin
//node name and a SearchResult. It will request the metafile of the file
//from the source node until it has the metafile. Can in theory have come
//from a different searchresult but it will continue until it gets it
func (g *Gossiper) RequestMetaFile(src string, result data.SearchResult) {
	hxHash := hex.EncodeToString(result.MetafileHash)
	for {
		datarequest := &data.DataRequest{
			Origin:      g.Name,
			Destination: src,
			HopLimit:    hoplimit,
			HashValue:   result.MetafileHash,
		}
		gp := &data.GossipPacket{
			DataRequest: datarequest,
		}
		g.sendMessageToNeighbour(gp, src)
		//Making the goroutine sleep for half a second before sending again
		//This number was kind of chosen arbitrarily.
		time.Sleep(500 * time.Millisecond)
		foundMetafile := g.StateFileFinder.HasMetaFile(hxHash)
		if foundMetafile {
			break
		}
	}
	g.StateFileFinder.AddOrigin(result.FileName, src, result.ChunkMap)
}
