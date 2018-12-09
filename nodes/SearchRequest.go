package nodes

import (
	"encoding/hex"
	"log"
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
	if _, ok := g.RoutingTable.Table[src]; !ok {
		//I don't want to add a route to myself.
		g.RoutingTable.UpdateRoutingTable(src, addr)
	}
	isValidRequest := g.RecentRequest.AddSearchRequest(req)
	if !isValidRequest || g.Name == src {
		//Dropping a packet that has the same origin and keywords
		//and arrived in the last half second.
		//Or if this is the node that sent the request
		return
	}
	matchingFiles := g.HasFileThatMatchesKeyWords(keywords)
	if len(matchingFiles) != 0 {
		g.HandleMatchingFiles(matchingFiles, src)
	} else {
		g.HandleForwardingSearchRequest(src, addr, keywords, req.Budget)
	}
}

//HandleForwardingSearchRequest is an abstraction of forwarding a searhc request if there is no match
func (g *Gossiper) HandleForwardingSearchRequest(src, addr string, keywords []string, budget uint64) {
	//Start sending the request to all my neighbours except
	//for the one I just got the request from.
	peers := g.Neighbours.GetAllNeighboursWithException(addr)
	n := uint64(len(peers))
	budget = uint64(budget - 1)
	if budget <= uint64(1) {
		//The budget has expired
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

//HandleMatchingFiles is an abstraction of how to reply to a
//search request if there is atleat one matching file.
func (g *Gossiper) HandleMatchingFiles(matchingFiles []string, src string) {
	//Handle sending a reply back to origin.
	searchResults := make([]*data.SearchResult, len(matchingFiles))
	for i, fn := range matchingFiles {
		metafilehash, chunkmap, chunkCount := g.GetMetaFileHashAndChunkMap(fn)
		searchResult := data.CreateSearchResult(fn, metafilehash, chunkmap, chunkCount)
		searchResults[i] = searchResult
	}
	searchreply := data.NewSearchReply(g.Name, src, hoplimit, searchResults)
	g.SendSearchReply(searchreply)
}

//GetMetaFileHashAndChunkMap is function that returns the metafilehash and the
//chunkmap of all the chunks this node has. If it does not have the metafile
//then the second return value will be nil. If it does not have the hash of the
//metafile then the first return value will be nil. Kind of weird that that will
//happen but whatever.
func (g *Gossiper) GetMetaFileHashAndChunkMap(fn string) ([]byte, []uint64, uint64) {
	metadata := g.GetMetadataByFilename(fn)
	hashOfMetafile, _ := hex.DecodeString(metadata.HashOfMetaFile)
	metafile := metadata.MetaFile
	if metafile == nil {
		log.Fatal("Metafile was nil")
		return hashOfMetafile, nil, 0
	}
	chunkmap := g.CreateChunkmap(metafile)
	chunkCount := uint64(len(metafile) / 32)
	return hashOfMetafile, chunkmap, chunkCount
}

//HasFileThatMatchesKeyWords is a function on a Gossiper that takes in a
//slice of keywords and compares them the list of known files by the gossiper.
//and returns a slice of strings that are the filenames that matched the  keywords
func (g *Gossiper) HasFileThatMatchesKeyWords(keywords []string) []string {
	temp := make([]string, 0)
	files := g.Files
	for _, metadata := range files {
		fileName := metadata.FileName
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

//GetMetadataByFilename returns the metadata that has the corresponding filename.
//This is not a very elegant solution
func (g *Gossiper) GetMetadataByFilename(filename string) *data.MetaData {
	for _, metadata := range g.Files {
		if metadata.FileName == filename {
			return &metadata
		}
	}
	return nil
}
