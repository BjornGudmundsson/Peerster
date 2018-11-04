package nodes

import (
	"fmt"

	"github.com/BjornGudmundsson/Peerster/data"
)

func (g *Gossiper) delegateMessages(ch chan GossipAddress) {
	for {
		select {
		case msg := <-ch:
			g.Neighbours.PrintNeighbours()
			if msg.Msg.Simple != nil {
				go g.handleSimpleMessage(msg)
			}
			if msg.Msg.Rumour != nil {
				go g.handleRumourMessage(msg)
			}
			if msg.Msg.Status != nil {
				go g.handleStatusMessage(msg)
			}
			if msg.Msg.PrivateMessage != nil {
				go g.handlePrivateMessage(msg)
			}
			if msg.Msg.DataReply != nil {
				go g.handleDataReplyMessage(msg)
			}
			if msg.Msg.DataRequest != nil {
				go g.handleDataRequestMessage(msg)
			}
		}
	}
}

func (g *Gossiper) handleSimpleMessage(msg GossipAddress) {
	g.Neighbours.PrintNeighbours()
	sm := *msg.Msg.Simple
	fmt.Printf("SIMPLE MESSAGE origin %v from %v content %v \n", sm.OriginalName, sm.RelayPeerAddr, sm.Contents)
}

func (g *Gossiper) handleRumourMessage(msg GossipAddress) {
	gp := msg.Msg
	addr := msg.Addr
	rm := *gp.Rumour
	enPeer := g.enPeer.EntropyPeer
	if rm.Text != "" {
		fmt.Printf("\nRUMOR origin %v from %v ID %v contents %v \n", rm.Origin, addr, rm.ID, rm.Text)
	}
	isNew := g.Messages.CheckIfMsgIsNew(rm)
	if rm.Text == "" {
		if isNew {
			g.RoutingTable.UpdateRoutingTable(rm.Origin, addr)
		}
		return
	}
	if g.Status.IsMongering {
		if addr == g.Status.GetIP() {
			g.Mongering.Ch <- rm
		}
		return
	}
	if isNew {
		g.Messages.AddAMessage(rm)
		g.RoutingTable.UpdateRoutingTable(rm.Origin, addr)
		myMsgs := g.Messages.GetMessageVector()
		sp := data.GetStatusPacketFromVector(myMsgs)
		ngp := data.GossipPacket{
			Status: &sp,
		}
		if addr != enPeer {
			go g.sendMessageToNeighbour(&ngp, addr)
		}
		g.rumourMongering(&msg)
	}

}

//CheckIfUpToDate is a function that takes in a map that
//corresponds to a StatusPacket and returns if there are
//messages that this gossiper has not seen.
func (g *Gossiper) CheckIfUpToDate(m map[string]uint32) bool {
	messages := g.Messages.Messages
	for key, val := range m {
		msgs, ok := messages[key]
		if !ok {
			return false
		}
		n := len(msgs)
		if uint32(n) != val-1 {
			return false
		}
	}
	for key := range messages {
		_, ok := m[key]
		if !ok {
			return false
		}
	}
	return true
}

func (g *Gossiper) handleStatusMessage(msg GossipAddress) {
	PrintStatusPacket(msg)
	addr := msg.Addr
	m := msg.Msg
	smap := TurnStatusIntoMap(*m.Status)
	upToDate := g.CheckIfUpToDate(smap)
	if upToDate {
		fmt.Printf("\n IN SYNC WITH %v \n", addr)
		return
	}
	if addr == g.Status.GetIP() {
		g.Status.StatusChannel <- msg
		return
	}

	//Now I know that this status packet does not come from
	//someone I am mongering with. That means this status
	//packet is asking for messages
	g.SendMessageThatPeerNeeds(msg)
}

//PrintStatusPacket prints a StatusPacket
//in a very particular format.
func PrintStatusPacket(ga GossipAddress) {
	fmt.Printf("STATUS from %v ", ga.Addr)
	sp := ga.Msg.Status
	for _, ps := range sp.Want {
		fmt.Printf("peer %v NextID %v ", ps.Identifier, ps.NextID)
	}
}

//TurnStatusIntoMap takes a StatusPacket and returns a
//map that conveys the same information
func TurnStatusIntoMap(sp data.StatusPacket) map[string]uint32 {
	m := make(map[string]uint32)
	for _, ps := range sp.Want {
		m[ps.Identifier] = ps.NextID
	}
	return m
}
