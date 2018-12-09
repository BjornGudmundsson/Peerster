package data

import (
	"encoding/hex"
	"fmt"
)

//SearchReply is a struct
//that holds onto the
//necessary variables to
//reply to a search request
//and forward it to the
//the node that requested it.
type SearchReply struct {
	Origin      string
	Destination string
	HopLimit    uint32
	Results     []*SearchResult
}

//NewSearchReply is an abstraction of creating a search reply and returning the pointer
func NewSearchReply(src, dst string, hl uint32, results []*SearchResult) *SearchReply {
	return &SearchReply{
		Origin:      src,
		Destination: dst,
		HopLimit:    hl,
		Results:     results,
	}
}

//Print prints a search repply
func (sr SearchReply) Print() {
	src := sr.Origin
	results := sr.Results
	for _, result := range results {
		fmt.Println("FOUND match ", result.FileName, " at ", src, " metafile=", hex.EncodeToString(result.MetafileHash), " chunks=", result.ChunkMap)
	}
}
