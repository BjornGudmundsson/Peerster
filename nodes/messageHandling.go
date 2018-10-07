package nodes

import (
	"fmt"

	"github.com/BjornGudmundsson/Peerster/data"
)

func (g *Gossiper) delegateMessages(ch chan *data.GossipPacket) {
	for {
		select {
		case msg := <-ch:
			if msg.Simple != nil {
				go g.handleSimpleMessage(msg)
			}
			if msg.Rumour != nil {
				fmt.Println("bjorn")
				go g.handleRumourMessage(msg)
			}
			if msg.Status != nil {
				go g.handleStatusMessage(msg)
			}
		}
	}
}

func (g *Gossiper) handleSimpleMessage(msg *data.GossipPacket) {
	g.Neighbours.PrintNeighbours()
	sm := *msg.Simple
	fmt.Printf("SIMPLE MESSAGE origin %v from %v content %v \n", sm.OriginalName, sm.RelayPeerAddr, sm.Contents)
}

func (g *Gossiper) handleRumourMessage(msg *data.GossipPacket) {
	mh := g.Messages
	mh.AddAMessage(*msg.Rumour)
	mh.PrintMessagesForOrigin(msg.Rumour.Origin)
	fmt.Println("Rumour messages")

}

func (g *Gossiper) handleStatusMessage(msg *data.GossipPacket) {
	fmt.Println("Status message")
}
