package nodes

func (g *Gossiper) handlePrivateMessage(msg GossipAddress) {
	priv := msg.Msg.PrivateMessage
	dst := priv.Destination
	origin := priv.Origin
	text := priv.Text
	privMessages := g.PrivateMessageStorage
	if dst == g.Name {
		privMessages.PutMessageFromOrigin(origin, text)
	}
	hLimit := priv.HopLimit
	nxtHop, ok := g.RoutingTable.Table[dst]
	//If I don't know the next hop, discard the message
	if !ok {
		return
	}
	nxtLimit := hLimit - 1
	if nxtLimit == 0 {
		return
	}
	priv.HopLimit = nxtLimit
	g.sendMessageToNeighbour(msg.Msg, nxtHop)
}
