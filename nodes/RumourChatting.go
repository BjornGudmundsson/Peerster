package nodes

import (
	"time"

	"github.com/BjornGudmundsson/Peerster/data"
)

//RumourChatting is a function that sends rumour
//chat messages of sorts to help update the routing
//table of every node
func (g *Gossiper) RumourChatting(rtTimer int) {
	if rtTimer == 0 {
		return
	}
	peers := g.Neighbours.Neighbours
	startingRumour := &data.RumourMessage{
		Origin: g.Name,
		ID:     1,
		Text:   "",
	}
	startingGossipPacket := &data.GossipPacket{
		Rumour: startingRumour,
	}
	for peer := range peers {
		g.sendMessageToNeighbour(startingGossipPacket, peer)
	}
	for {
		//sendToEveryone
		counter := g.Counter.ReturnCounter()
		rm := &data.RumourMessage{
			Origin: g.Name,
			ID:     counter + 1,
			Text:   "",
		}
		gp := &data.GossipPacket{
			Rumour: rm,
		}
		//Going to broadcast the message to everyone
		for peer := range peers {
			g.sendMessageToNeighbour(gp, peer)
		}
		time.Sleep(time.Duration(int64(rtTimer)) * time.Second)
	}
}
