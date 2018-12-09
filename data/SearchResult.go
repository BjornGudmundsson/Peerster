package data

//SearchResult is a struct
//that holds onto the results
//from a search reply.
type SearchResult struct {
	FileName     string
	MetafileHash []byte
	ChunkMap     []uint64
	ChunkCount   uint64
}

//CreateSearchResult is just an abstraction of creating a search result and getting the pointer.
func CreateSearchResult(fn string, metafilehash []byte, chunks []uint64, chunkCount uint64) *SearchResult {
	return &SearchResult{
		FileName:     fn,
		MetafileHash: metafilehash,
		ChunkMap:     chunks,
		ChunkCount:   chunkCount,
	}
}
