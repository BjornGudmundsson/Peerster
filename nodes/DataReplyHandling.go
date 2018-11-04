package nodes

import (
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
	//Here I start processing datareplies
	//Write shit here in the comments since I don't have a clear idea yet
	/*
		Check if I am waiting for a hashvalue
		Am I waiting for a hashvalue
		Yes?
			Check if this is the hashvalue I am waiting for
			Is this the hashvalue I am waiting for?
			Yes?
				Update the next hashvalue I am waiting for
				Check if there are any more hashvalues to receive.reply
				Start reconstructing the file
				Possibly by keeping all of the current hashvalues in a
				"sorted" list. Basically a list that has the hashvalues in order.reply
				reconstruct file from that.
			No?
				Discard this packet
				Send another request with the hashvalue I am waiting for.reply
				back to the "source" of the file.
		No?
			Some networking delay nonesense, discard the packet
	*/
}
