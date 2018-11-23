package nodes

import (
	"fmt"

	"github.com/BjornGudmundsson/Peerster/data"
)

//FileHandling is just way to handle dealing with file
//requests from the client by abstracting the functionality
//and make the client messages function maintainable
func (g *Gossiper) FileHandling(temp data.TextMessage) {
	/*if temp.Dst == "" {
		g.HandleNewOSFile(temp.File)
	}
	if temp.Dst != "" && temp.Request != "" {
		mf, e := hex.DecodeString(temp.Request)
		if e != nil {
			log.Fatal(e)
		}
		g.DownLoadAFile(temp.File, mf, temp.Dst)
	}*/
}

//ClientGossiperHandling is a function bound to a point to a Gossiper that
//handles gossiper messages. It abstracts the functionality and returns the
//gossip packet along with an address to send it to. Which just happens to be its own
//address. LoL
func (g *Gossiper) ClientGossiperHandling(temp data.TextMessage) *GossipAddress {
	id := g.Counter.IncrementAndReturn()
	fmt.Printf("CLIENT MESSAGE: %v", temp.Msg)
	rm := &data.RumourMessage{
		Origin: g.Name,
		ID:     id,
		Text:   temp.Msg,
	}
	g.Messages.AddAMessage(*rm)
	gp := &data.GossipPacket{
		Rumour: rm,
	}
	ga := &GossipAddress{
		Addr: g.address.String(),
		Msg:  gp,
	}
	return ga
}

//HandleClientSearchRequests is an abstraction of the functionality of
//search request send from the client
func (g *Gossiper) HandleClientSearchRequests(temp data.TextMessage) {
	budget := temp.Budget
	if budget == 2 {
		//Do the logic if the budget is 2
	} else {
		//Do the logic if the budget is not 2
	}
	fmt.Println("Handle client search request")
}
