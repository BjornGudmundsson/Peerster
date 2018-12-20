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
			if msg.Msg.SearchReply != nil {
				go g.HandleSearchReply(msg)
			}
			if msg.Msg.SearchRequest != nil {
				go g.HandleSearchRequest(msg)
			}

			if msg.Msg.BlockPublish != nil {
				go g.HandleBlockPublish(msg)
			}
			if msg.Msg.TxPublish != nil {
				go g.HandleTxPublish(msg)
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
	rm := gp.Rumour
	isNew := g.RumourHolder.IsNew(*rm)
	if isNew {
		g.RumourHolder.AddRumour(*rm)
		sp := g.RumourHolder.CreateStatusPacket()
		g.SendStatusPacket(sp, addr)
		go g.rumourMongering(rm, addr)
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
	sp := msg.Msg.Status
	hasEntry := g.StatusPeers.HasEntry(addr)
	if hasEntry {
		g.StatusPeers.PassPacketToProcess(addr, *m)
		return
	}
	fmt.Println("Not an entry")
	upToDate := g.RumourHolder.CheckIfUpToDate(m.Status)
	if upToDate {
		fmt.Println("In sync with", addr)
		return
	}
	peerNeeds := g.RumourHolder.GetRumoursPeerNeeds(sp)
	if len(peerNeeds) != 0 {
		//Sending a random rumour that the
		randomRumour := data.GetRandomRumourFromSlice(peerNeeds)
		fmt.Println("Random rumour: ", randomRumour)
		g.rumourMongering(&randomRumour, "")
	} else {
		//Get the messages that I need if there were no messages that
		//the other peer needs.
		INeed := g.RumourHolder.CheckIfNeedMessages(sp)
		IWant := INeed.Want
		if len(IWant) != 0 {
			g.SendStatusPacket(INeed, addr)
		}
	}

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
