package nodes

import (
	"strings"

	"github.com/BjornGudmundsson/Peerster/data"
)

//HandleSearchRequest handles all incoming search request
//messages coming to the gossiper.
func (g *Gossiper) HandleSearchRequest(msg GossipAddress) {
	addr := msg.Addr
	req := msg.Msg.SearchRequest
	src := req.Origin
	keywords := req.Keywords
	matchingFiles := g.HasFileThatMatchesKeyWords(keywords)
	if len(matchingFiles) != 0 {
		//Handle sending a reply back to origin.
	} else {
		//Start sending the request to all my neighbours except
		//for the one I just got the request from.
		peers := g.Neighbours.GetAllNeighboursWithException(addr)
		n := uint64(len(peers))
		budget := uint64(req.Budget - 1)
		if budget == uint64(0) {
			return
		}
		budgetDistribution := data.CreateBudgetList(n, budget)
		for i, b := range budgetDistribution {
			newReq := &data.SearchRequest{
				Origin:   src,
				Budget:   b,
				Keywords: keywords,
			}
			gp := &data.GossipPacket{
				SearchRequest: newReq,
			}
			g.sendMessageToNeighbour(gp, peers[i])
		}
	}
}

//HasFileThatMatchesKeyWords is a function on a Gossiper that takes in a
//slice of keywords and compares them the list of known files by the gossiper.
//and returns a slice of strings that are the filenames that matched the  keywords
func (g *Gossiper) HasFileThatMatchesKeyWords(keywords []string) []string {
	temp := make([]string, 0)
	files := g.Files
	for fileName := range files {
		for _, keyword := range keywords {
			containsKeyword := strings.Contains(fileName, keyword)
			if containsKeyword {
				temp = append(temp, fileName)
				break
			}
		}
	}
	return temp
}
