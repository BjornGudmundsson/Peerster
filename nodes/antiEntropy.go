package nodes

import (
	"time"

	"github.com/BjornGudmundsson/Peerster/data"
)

const antiEntropy time.Duration = 1 * time.Second

//AntiEntropy is an infinite loop that sends
//StatusPackets to a random peer at a pre-determined
//interval of all messages that this gossiper has as
//of sending that StatusPacket
func (g *Gossiper) AntiEntropy() {
	for {
		time.Sleep(antiEntropy)
		sp := g.RumourHolder.CreateStatusPacket()
		peers := g.Neighbours.GetAllNeighboursWithException("")
		randPeer := data.GetRandomStringFromSlice(peers)
		g.SendStatusPacket(sp, randPeer)
	}
}
