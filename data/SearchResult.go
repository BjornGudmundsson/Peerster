package data

//SearchResult is a struct
//that holds onto the results
//from a search reply.
type SearchResult struct {
	FileName     string
	MetafileHash []byte
	ChunkMap     []uint64
}
