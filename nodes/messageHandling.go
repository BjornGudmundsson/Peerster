package nodes

import (
	"fmt"

	"github.com/BjornGudmundsson/Peerster/data"
)

func (g *Gossiper) delegateMessages(ch chan *GossipAddress) {
	for {
		select {
		case msg := <-ch:
			if msg.Msg.Simple != nil {
				go g.handleSimpleMessage(msg)
			}
			if msg.Msg.Rumour != nil {
				go g.handleRumourMessage(msg)
			}
			if msg.Msg.Status != nil {
				go g.handleStatusMessage(msg)
			}
		}
	}
}

func (g *Gossiper) handleSimpleMessage(msg *GossipAddress) {
	g.Neighbours.PrintNeighbours()
	sm := *msg.Msg.Simple
	fmt.Printf("SIMPLE MESSAGE origin %v from %v content %v \n", sm.OriginalName, sm.RelayPeerAddr, sm.Contents)
}

func (g *Gossiper) handleRumourMessage(msg *GossipAddress) {
	mh := g.Messages
	relay := msg.Addr
	newMsg := mh.AddAMessage(*msg.Msg.Rumour)
	//mh.PrintMessagesForOrigin(msg.Msg.Rumour.Origin)
	if g.Mongering.HasMongerer(relay) {
		return
	}
	if newMsg {
		msgVector := g.Messages.GetMessageVector()
		sp := data.CreateStatusPacketFromMessageVector(msgVector)
		fmt.Println(sp)
		gp := &data.GossipPacket{
			Status: sp,
		}
		go g.sendRumourMessageToNeighbour(gp, relay)
		fmt.Println("sending to this address ", relay)
		go g.rumourMongering(msg)
	}
}

func (g *Gossiper) handleStatusMessage(msg *GossipAddress) {
	fmt.Println("Status message")
}
