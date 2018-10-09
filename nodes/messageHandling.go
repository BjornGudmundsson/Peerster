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
				go g.SendMessageVector()
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
	//fmt.Println("Rumour from ", relay)
	mh.PrintMessagesForOrigin(msg.Msg.Rumour.Origin)

	//Ignore all messages that are not from the person
	//I am mongering with
	if g.Status.IsMongering {
		fmt.Println("I am rumor mongering")
		return
	}
	fmt.Println("new msg", newMsg, "and is mongering ", g.Status.IsMongering)
	if newMsg && !g.Status.IsMongering {
		msgVector := g.Messages.GetMessageVector()
		sp := data.CreateStatusPacketFromMessageVector(msgVector)
		gp := &data.GossipPacket{
			Status: sp,
		}
		go g.sendRumourMessageToNeighbour(gp, relay)
		go g.rumourMongering(msg)
	}
}

func (g *Gossiper) handleStatusMessage(msg *GossipAddress) {
	addr := msg.Addr
	mongerAddr := g.Status.GetIP()
	if g.Status.IsMongering && addr == mongerAddr {
		g.Status.StatusChannel <- msg
		return
	}
	fmt.Println("got a random status message ", msg.Addr)
}

//SendMessageVector is a function that waits for a
//statuspacket and sends it
func (g *Gossiper) SendMessageVector() {
	for {
		select {
		case vector := <-g.Status.StatusChannel:
			gaddr := vector.Msg
			addr := vector.Addr
			needMsgs := g.Messages.NeedMsgs(*gaddr.Status)
			fmt.Println("Messages I need", needMsgs)
			gp := &data.GossipPacket{
				Status: &needMsgs,
			}
			g.sendRumourMessageToNeighbour(gp, addr)

		}
	}
}
