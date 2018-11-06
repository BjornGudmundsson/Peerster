package nodes

import (
	"encoding/hex"
	"fmt"

	"github.com/BjornGudmundsson/Peerster/data"
)

func (g *Gossiper) handleDataReplyMessage(msg GossipAddress) {
	fmt.Println("Got a data repy")
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
	if !g.dataReplyHandler.IsPending() {
		//I wasn't waiting for anything. Probably an old packet or something
		//some weird networking shit.
		return
	}
	nxtHop, ok := g.RoutingTable.Table[reply.Origin]
	if !ok {
		//I don't  know how this happened
		return
	}
	nxtChunk, isFinished, isMetaFile, wasValid := g.dataReplyHandler.Update(reply.HashValue, reply.Data)
	if isFinished {
		fmt.Println("This finished")
		hx := hex.EncodeToString(reply.HashValue)
		g.Chunks[hx] = string(reply.Data)
		return
	}
	fmt.Println("nxtChunk is", hex.EncodeToString(nxtChunk))
	dreq := &data.DataRequest{
		Origin:      g.Name,
		Destination: reply.Origin,
		HopLimit:    hoplimit,
		HashValue:   nxtChunk,
	}
	packet := &data.GossipPacket{
		DataRequest: dreq,
	}
	if isMetaFile {
		hx := hex.EncodeToString(reply.HashValue)
		md := data.MetaData{
			FileName:       g.dataReplyHandler.Name,
			FileSize:       7,
			HashOfMetaFile: hx,
			MetaFile:       reply.Data,
		}
		g.Files[hx] = md
		g.sendMessageToNeighbour(packet, nxtHop)
		return
	}
	if wasValid {
		hx := hex.EncodeToString(reply.HashValue)
		g.Chunks[hx] = string(reply.Data)
	}
	g.sendMessageToNeighbour(packet, nxtHop)
}
