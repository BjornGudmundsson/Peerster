package nodes

import (
	"encoding/hex"

	"github.com/BjornGudmundsson/Peerster/data"
)

func (g *Gossiper) handleDataReplyMessage(msg GossipAddress) {
	reply := msg.Msg.DataReply
	//This if statement handles if I am supposed to forward the reply
	//further. Check if it is for me, if not, continue forwarding
	//according the the stated criteria for forwarding.
	if reply.Destination != g.Name {
		reply.HopLimit = reply.HopLimit - 1
		if reply.HopLimit == 0 {
			//Dropping this packet because it has exceeded the hoplimit
			return
		}
		nxtHop, ok := g.RoutingTable.Table[reply.Destination]
		if !ok {
			//Dropping this packet since it has no forwarding destination
			return
		}
		gp := &data.GossipPacket{
			DataReply: reply,
		}
		g.sendMessageToNeighbour(gp, nxtHop)
		return
	}
	g.HandlerDataReplies.PassReplyToChannel(*reply)
}

//HasAllChunksOfFile checks if all of the chunks in a metafile
//are present in the chunk repository of the gossiper. It also
//returns the index of the next chunk it needs in order.
func (g *Gossiper) HasAllChunksOfFile(metafile []byte) (bool, uint64) {
	//The length of the metafile is always a multiple of 32
	n := len(metafile) / 32
	for i := 0; i < n; i++ {
		j := i + 1
		hx := hex.EncodeToString(metafile[i*32 : j*32])
		_, ok := g.Chunks[hx]
		if !ok {
			return false, uint64(i * 32)
		}
	}
	return true, uint64(n * 32)
}

//SendDataRequest is just an abstraction of sending a datarequest to a destination.
func (g *Gossiper) SendDataRequest(dst string, chunk []byte) {
	datarequest := &data.DataRequest{
		Destination: dst,
		HopLimit:    hoplimit,
		Origin:      g.Name,
		HashValue:   chunk,
	}
	gp := &data.GossipPacket{
		DataRequest: datarequest,
	}
	nxtHop, ok := g.RoutingTable.Table[dst]
	if !ok {
		return
	}
	g.sendMessageToNeighbour(gp, nxtHop)
}
