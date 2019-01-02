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
	counter := g.Counter.IncrementAndReturn()
	peers := g.Neighbours.GetAllNeighboursWithException("")
	startingRumour := &data.RumourMessage{
		Origin: g.Name,
		ID:     counter,
		Text:   "",
	}
	g.RumourHolder.AddRumour(*startingRumour)
	p := data.GetRandomStringFromSlice(peers)
	g.SendRumourMessage(startingRumour, p)
	for {
		time.Sleep(time.Duration(int64(rtTimer)) * time.Second)
		//sendToEveryone
		counter = g.Counter.IncrementAndReturn()
		rm := &data.RumourMessage{
			Origin: g.Name,
			ID:     counter,
			Text:   "",
		}
		g.RumourHolder.AddRumour(*rm)
		peer := data.GetRandomStringFromSlice(peers)
		g.SendRumourMessage(rm, peer)
	}
}
