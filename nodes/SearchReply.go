package nodes

import (
	"fmt"

	"github.com/BjornGudmundsson/Peerster/data"
)

//HandleSearchReply handles all incoming SearchReply messages
//coming to the Gossiper
func (g *Gossiper) HandleSearchReply(msg GossipAddress) {
	fmt.Println("Got a search reply")
	sr := msg.Msg.SearchReply
	dst := sr.Destination
	if g.Name == dst {
		//Handle the datareply
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
