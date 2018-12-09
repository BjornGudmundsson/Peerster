package nodes

import (
	"fmt"
	"strings"
	"time"

	"github.com/BjornGudmundsson/Peerster/data"
)

//FileHandling is just way to handle dealing with file
//requests from the client by abstracting the functionality
//and make the client messages function maintainable
func (g *Gossiper) FileHandling(temp data.TextMessage) {
	if temp.Request == "" {
		g.HandleNewOSFile(temp.File)
		return
	}
	if temp.Request != "" {
		mf := temp.Request
		fn := temp.File
		md := data.MetaData{
			FileName:       fn,
			HashOfMetaFile: temp.Request,
			MetaFile:       nil,
			FileSize:       0,
		}
		g.Files[fn] = md
		dst := temp.Dst
		if dst != "" {
			g.ChunkToPeer.SetOwnerOfMetafileHash(mf, dst)
			g.ChunkToPeer.AddOwnerForMetafileHash(dst, mf)
		}
		go g.DownloadingFile(fn)
	}
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
	src := g.Name
	keywords := strings.Split(temp.Keywords, ",")
	fmt.Println(keywords)
	searchrequest := data.NewSearchRequest(src, budget, keywords)
	gp := &data.GossipPacket{
		SearchRequest: searchrequest,
	}
	if budget == 2 {
		//Do the logic if the budget is 2
		//This function will handle the logic behind
		//the constantly increasing messages.
		g.HandleIncrementedMessaging(src, keywords)
	} else {
		//This just broadcasts a search request once
		peers := g.Neighbours.ArrNeighbours
		//Just spread the request around  but only once
		for _, peer := range peers {
			g.sendMessageToNeighbour(gp, peer)
		}
		//Do the logic if the budget is not 2
	}
}

//HandleIncrementedMessaging handles the case if the budget is not specified
//or just happens to be 2. It will double the budget with each round and check
//if there are any full matches that match these keywords.
func (g *Gossiper) HandleIncrementedMessaging(src string, keywords []string) {
	budget := 2
	for {
		//Get all matches that I have for these keywords
		matches := g.FoundFileRepository.FindFullMatches(keywords)
		fullmatches := g.ChunkToPeer.FileMatches(matches)
		if fullmatches >= threshold {
			fmt.Println("SEARCH FINISHED")
			return
		}
		//Check if I have exceeded the budget set by the description.
		if budget > maxBudget {
			fmt.Println("SEARCH FINISHED")
			return
		}
		//Sending a data request to everyone
		searchrequest := data.NewSearchRequest(src, uint64(budget), keywords)
		gp := &data.GossipPacket{
			SearchRequest: searchrequest,
		}
		peers := g.Neighbours.ArrNeighbours
		//Just spread the request around  but only once
		for _, peer := range peers {
			g.sendMessageToNeighbour(gp, peer)
		}
		//Doubling the budget with each round
		budget = budget * 2
		time.Sleep(1 * time.Second)
	}
}
