package data

//SearchRequest is a struct holding
//all the necessary variables to perform a
//search for a file over the peerster
//network.
type SearchRequest struct {
	Origin   string
	Budget   uint64
	Keywords []string
}

//NewSearchRequest creates a pointer to a new SearchRequest
func NewSearchRequest(src string, budget uint64, keywords []string) *SearchRequest {
	return &SearchRequest{
		Origin:   src,
		Budget:   budget,
		Keywords: keywords,
	}
}
