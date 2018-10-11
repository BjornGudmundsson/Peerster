package nodes

import "github.com/BjornGudmundsson/Peerster/data"

//SendMessageThatPeerNeeds sends all the messages that a peer
//asked for in a StatusPacket
func (g *Gossiper) SendMessageThatPeerNeeds(ga GossipAddress) {
	messages := g.Messages.Messages
	sp := ga.Msg.Status
	addr := ga.Addr
	for _, ps := range sp.Want {
		id := ps.NextID
		rmsgs, ok := messages[ps.Identifier]
		if !ok {
			continue
		}
		n := uint32(len(rmsgs))
		if id > n {
			continue
		}
		for i := id - 1; i < n; i++ {
			rm := &rmsgs[i]
			gp := &data.GossipPacket{
				Rumour: rm,
			}
			g.sendMessageToNeighbour(gp, addr)
		}
	}
	msp := TurnStatusIntoMap(*sp)
	for key, val := range messages {
		if _, ok := msp[key]; !ok {
			for _, rm := range val {
				gp := &data.GossipPacket{
					Rumour: &rm,
				}
				g.sendMessageToNeighbour(gp, addr)
			}
		}
	}
}
