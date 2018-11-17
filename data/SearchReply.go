package data

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
