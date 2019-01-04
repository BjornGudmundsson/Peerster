package hashtable

//ChunkStoreRequest is a request for
//a peer to store a chunk.
type ChunkStoreRequest struct {
	Destination string
	Hash        string
	Data        []byte
	Src         string
}

//ChunkStoreReply is a reply
//to a store request. Has the reply
//and to which chunk it is replying to.
type ChunkStoreReply struct {
	Destination string
	Hash        string
	Reply       bool
	Src         string
}
