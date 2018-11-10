package nodes

import (
	"fmt"
	"log"
	"net"

	"github.com/BjornGudmundsson/Peerster/data"
	"github.com/dedis/protobuf"
)

func (g *Gossiper) handlePrivateMessage(msg GossipAddress) {
	priv := msg.Msg.PrivateMessage
	dst := priv.Destination
	origin := priv.Origin
	text := priv.Text
	privMessages := g.PrivateMessageStorage
	if dst == g.Name {
		privMessages.PutMessageFromOrigin(origin, text)
		fmt.Println("")
		fmt.Printf("PRIVATE origin %v hop-limit %v contents %v", origin, priv.HopLimit, text)
		fmt.Println("")
		return
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

func (g *Gossiper) SendPrivateMessageFromUser(dst, txt string) {
	priv := &data.PrivateMessage{
		Origin:      g.Name,
		ID:          0,
		Text:        txt,
		Destination: dst,
		HopLimit:    hoplimit,
	}
	gp := &data.GossipPacket{
		PrivateMessage: priv,
	}
	buf, _ := protobuf.Encode(gp)
	conn, e := net.Dial("udp", g.address.String())
	defer conn.Close()
	if e != nil {
		log.Fatal(e)
	}
	conn.Write(buf)
}
