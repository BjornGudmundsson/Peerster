package nodes

import (
	"sync"
	"time"

	"github.com/BjornGudmundsson/Peerster/data"
)

//EntropyPeer is actually an
//awful name for this struct.
//Just could not think of anything better.
type EntropyPeer struct {
	EntropyPeer string
	mux         sync.Mutex
}

//SetEntropyPeer set the entropyPeer to a new address
//and sets the current goroutine to sleep for 1 second
func (ep *EntropyPeer) SetEntropyPeer(peer string) {
	ep.mux.Lock()
	ep.EntropyPeer = peer
	ep.mux.Unlock()
}

//ResetEntropyPeer resets the entropyPeer
//back to the empty string
func (ep *EntropyPeer) ResetEntropyPeer() {
	ep.mux.Lock()
	ep.EntropyPeer = ""
	ep.mux.Unlock()
}

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
