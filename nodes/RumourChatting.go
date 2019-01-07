package nodes

import (
	"time"

	"github.com/BjornGudmundsson/Peerster/data"
)

// RumourChatting sends route rumours (rumour chat messages of sorts) to help
// update the routing table of every node.
func (g *Gossiper) RumourChatting(rtTimer int) {
	if rtTimer == 0 {
		return
	}
	for {
		counter := g.Counter.IncrementAndReturn()
		peers := g.Neighbours.GetAllNeighboursWithException("")
		rm := &data.RumourMessage{
			Origin: g.Name,
			ID:     counter,
			Text:   "",
		}
		g.RumourHolder.AddRumour(*rm)
		peer := data.GetRandomStringFromSlice(peers)
		g.SendRumourMessage(rm, peer)
		time.Sleep(time.Duration(rtTimer) * time.Second)
	}
}
